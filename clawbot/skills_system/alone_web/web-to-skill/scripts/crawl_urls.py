#!/usr/bin/env python3
from __future__ import annotations

import argparse
import asyncio
import datetime as dt
import hashlib
import json
import os
import re
import sys
import time
from collections import Counter
from pathlib import Path
from typing import Any, TextIO
from urllib.parse import urlparse

from bs4 import BeautifulSoup
from jieba import analyse as jieba_analyse

from fetch_rendered_html import (
    URLSafetyError,
    URLSafetyGuard,
    USER_AGENT,
    ensure_nonempty_rendered_body,
    read_body_text,
    read_page_metadata,
    rendered_html_snapshot,
    rule_for_url,
    wait_for_rendered_body,
)


if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")


TIMEOUT_MS = 60_000
CONSECUTIVE_TIMEOUT_LIMIT = 4
KEYWORD_LIMIT = 20
MAX_INDEX_KEYWORDS = 12
MAX_KEYWORD_DOCUMENT_FREQUENCY = 0.30
MIN_PAGES_FOR_KEYWORD_FREQUENCY_FILTER = 4
ASCII_KEYWORD_RE = re.compile(r"^[a-z0-9][a-z0-9._+-]*$", re.IGNORECASE)
SKIPPED_RESOURCE_TYPES = {"font", "image", "media"}


class DuplicateFinalURL(Exception):
    def __init__(self, requested_url: str, final_url: str, first_requested_url: str):
        super().__init__(f"{requested_url} redirects to the already captured URL {final_url}")
        self.requested_url = requested_url
        self.final_url = final_url
        self.first_requested_url = first_requested_url


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Sequentially crawl every URL in a URL-list file into HTML snapshots and a JSONL index.")
    parser.add_argument("--url-list", required=True, help="UTF-8 text file with one URL per line.")
    parser.add_argument("--out-dir", required=True, help="Output directory for index.jsonl, html/, and crawl.log.")
    parser.add_argument("--debug", action="store_true", help="Crawl only the first five URLs.")
    return parser.parse_args()


def read_urls(path: Path) -> list[str]:
    seen: set[str] = set()
    urls: list[str] = []
    for line in path.read_text(encoding="utf-8-sig").splitlines():
        url = line.strip()
        if not url or url.startswith("#") or url in seen:
            continue
        if not url.startswith(("http://", "https://")):
            raise ValueError(f"invalid HTTP(S) URL in {path}: {url}")
        seen.add(url)
        urls.append(url)
    if not urls:
        raise ValueError(f"URL list is empty: {path}")
    return urls


def launch_args() -> list[str]:
    args = ["--disable-blink-features=AutomationControlled", "--disable-dev-shm-usage"]
    if os.environ.get("PLAYWRIGHT_NO_SANDBOX", "1").lower() not in {"0", "false", "no"}:
        args.extend(["--no-sandbox", "--disable-setuid-sandbox"])
    return args


def page_id(url: str) -> str:
    return hashlib.sha256(url.encode("utf-8")).hexdigest()[:16]


def split_keywords(raw: list[str]) -> list[str]:
    seen: set[str] = set()
    result: list[str] = []
    for item in raw:
        value = " ".join(str(item).split()).strip()
        key = value.casefold()
        if value and key not in seen:
            seen.add(key)
            result.append(value)
    return result


def extract_keywords(title: str, description: str, body_text: str, original: list[str]) -> list[str]:
    source = "\n".join(part for part in (title, description, body_text) if part.strip())
    extracted = jieba_analyse.extract_tags(source, topK=KEYWORD_LIMIT) if source else []
    return split_keywords([*original, *extracted])[:KEYWORD_LIMIT]


def keyword_in_context(keyword: str, context: str) -> bool:
    if ASCII_KEYWORD_RE.fullmatch(keyword):
        pattern = rf"(?<![a-z0-9]){re.escape(keyword)}(?![a-z0-9])"
        return re.search(pattern, context, re.IGNORECASE) is not None
    return keyword.casefold() in context.casefold()


def filter_common_keywords(records: list[dict[str, Any]]) -> None:
    document_frequency: Counter[str] = Counter()
    normalized: list[list[str]] = []
    for record in records:
        keywords = [item for item in split_keywords(record.get("keywords", [])) if len(item) >= 2]
        normalized.append(keywords)
        document_frequency.update({keyword.casefold() for keyword in keywords})

    page_count = max(len(records), 1)
    for record, keywords in zip(records, normalized):
        context = f"{record.get('title', '')} {record.get('description', '')}"
        filtered: list[str] = []
        for keyword in keywords:
            frequency = document_frequency[keyword.casefold()] / page_count
            if (
                page_count >= MIN_PAGES_FOR_KEYWORD_FREQUENCY_FILTER
                and frequency > MAX_KEYWORD_DOCUMENT_FREQUENCY
                and not keyword_in_context(keyword, context)
            ):
                continue
            filtered.append(keyword)
            if len(filtered) >= MAX_INDEX_KEYWORDS:
                break
        record["keywords"] = filtered


def retryable_error(exc: Exception) -> bool:
    if isinstance(exc, URLSafetyError):
        return False
    message = str(exc).casefold()
    return (
        type(exc).__name__ == "TimeoutError"
        or "net::err_" in message
        or "rendered body is empty" in message
        or re.search(r"\bhttp status (?:429|5\d\d)\b", message) is not None
    )


def is_timeout_error(exc: Exception) -> bool:
    return type(exc).__name__ == "TimeoutError"


def clean_body_text(value: str) -> str:
    value = value.replace("\r\n", "\n").replace("\r", "\n")
    lines = [" ".join(line.split()) for line in value.split("\n")]
    return "\n".join(line for line in lines if line)


def prefers_stable_body(url: str) -> bool:
    rule = rule_for_url(url)
    host = (urlparse(url).hostname or "").lower()
    return bool(rule and (host == "feishu.cn" or host.endswith(".feishu.cn")))


def choose_browser_body_text(current: str, stable: str, prefer_stable: bool) -> tuple[str, bool]:
    current_clean = clean_body_text(current)
    stable_clean = clean_body_text(stable)
    if stable_clean and (not current_clean or (prefer_stable and len(stable_clean) > len(current_clean))):
        return stable, True
    return current, False


def snapshot_body_text(snapshot: str) -> str:
    soup = BeautifulSoup(snapshot, "html.parser")
    main = soup.select_one("main[data-rendered-snapshot='true']") or soup.body or soup
    return "\n".join(line.strip() for line in main.get_text("\n", strip=True).splitlines() if line.strip())


def update_snapshot_keywords(snapshot: str, keywords: list[str]) -> str:
    soup = BeautifulSoup(snapshot, "html.parser")
    head = soup.head
    if head is not None:
        node = head.select_one("meta[name='keywords']")
        if node is None:
            node = soup.new_tag("meta")
            node["name"] = "keywords"
            head.append(node)
        node["content"] = ", ".join(keywords)
    return soup.prettify(formatter="minimal") + "\n"


class CrawlLogger:
    def __init__(self, path: Path) -> None:
        path.parent.mkdir(parents=True, exist_ok=True)
        self.file: TextIO = path.open("a", encoding="utf-8", buffering=1)

    def close(self) -> None:
        self.file.close()

    def emit(self, event: str, *, progress: str = "", **fields: Any) -> None:
        timestamp = dt.datetime.now().astimezone().isoformat(timespec="seconds")
        data = {"progress": progress, **fields} if progress else fields
        details = " ".join(
            f"{key}={json.dumps(value, ensure_ascii=False, default=str)}"
            for key, value in data.items()
            if value is not None
        )
        line = f"[crawl_urls] {timestamp} {event}" + (f" {details}" if details else "")
        print(line, file=sys.stderr, flush=True)
        self.file.write(line + "\n")


async def wait_for_stable_body_snapshot(
    page: Any,
    logger: CrawlLogger,
    index: int,
    total: int,
) -> str:
    started = asyncio.get_running_loop().time()
    deadline = started + min(TIMEOUT_MS / 1000, 30)
    next_progress = started + 5
    last_text = ""
    best_text = ""
    stable_count = 0
    while asyncio.get_running_loop().time() < deadline:
        text = clean_body_text(await read_body_text(page))
        if len(text) > len(best_text):
            best_text = text
        if text == last_text and text:
            stable_count += 1
        else:
            stable_count = 0
        last_text = text
        if stable_count >= 2:
            logger.emit(
                "page.body_stable",
                progress=f"{index}/{total}-step2/5",
                stable=True,
                chars=len(text),
                url=page.url,
            )
            return text
        now = asyncio.get_running_loop().time()
        if now >= next_progress:
            logger.emit(
                "page.body_stable.wait",
                progress=f"{index}/{total}-step2/5",
                elapsed_seconds=round(now - started, 3),
                chars=len(text),
                url=page.url,
            )
            next_progress = now + 5
        await page.wait_for_timeout(1000)
    logger.emit(
        "page.body_stable",
        progress=f"{index}/{total}-step2/5",
        stable=False,
        chars=len(last_text),
        retained_chars=len(best_text),
        url=page.url,
    )
    return best_text


async def crawl_one(
    page: Any,
    guard: URLSafetyGuard,
    url: str,
    html_dir: Path,
    logger: CrawlLogger,
    index: int,
    total: int,
    captured_final_urls: dict[str, str],
) -> dict[str, Any]:
    blocked_navigation_error: URLSafetyError | None = None
    route_requests = 0
    route_blocked = 0
    route_skipped = 0

    async def route_request(route: Any, request: Any) -> None:
        nonlocal blocked_navigation_error, route_blocked, route_requests, route_skipped
        route_requests += 1
        if not request.is_navigation_request() and request.resource_type in SKIPPED_RESOURCE_TYPES:
            route_skipped += 1
            await route.abort("blockedbyclient")
            return
        try:
            await guard.validate(request.url)
        except URLSafetyError as exc:
            route_blocked += 1
            if request.is_navigation_request():
                blocked_navigation_error = exc
            await route.abort("blockedbyclient")
            return
        await route.continue_()

    await page.route("**/*", route_request)
    started = time.monotonic()
    try:
        await guard.validate(url)
        logger.emit("page.goto.start", progress=f"{index}/{total}-step1/5", url=url, timeout_ms=TIMEOUT_MS)
        try:
            response = await page.goto(url, timeout=TIMEOUT_MS, wait_until="domcontentloaded")
        except Exception as exc:
            if blocked_navigation_error is not None:
                raise blocked_navigation_error from exc
            raise
        if response is not None and response.status >= 400:
            raise RuntimeError(f"HTTP status {response.status}")
        final_url = page.url
        await guard.validate(final_url)
        if final_url in captured_final_urls:
            raise DuplicateFinalURL(url, final_url, captured_final_urls[final_url])

        logger.emit("page.render.start", progress=f"{index}/{total}-step2/5", final_url=final_url)
        await wait_for_rendered_body(page, final_url, None, TIMEOUT_MS)
        stable_body_text = await wait_for_stable_body_snapshot(page, logger, index, total)

        logger.emit("page.read.start", progress=f"{index}/{total}-step3/5", final_url=final_url)
        browser_body_text = await read_body_text(page)
        current_body_chars = len(clean_body_text(browser_body_text))
        stable_body_chars = len(clean_body_text(stable_body_text))
        browser_body_text, retained_stable_body = choose_browser_body_text(
            browser_body_text,
            stable_body_text,
            prefers_stable_body(final_url),
        )
        if retained_stable_body:
            logger.emit(
                "page.body_text.retained",
                progress=f"{index}/{total}-step3/5",
                url=final_url,
                current_chars=current_body_chars,
                stable_chars=stable_body_chars,
            )
        ensure_nonempty_rendered_body(browser_body_text)
        title, description, original_keywords = await read_page_metadata(page)
        html_text = await page.content()

        logger.emit("page.snapshot.start", progress=f"{index}/{total}-step4/5", final_url=final_url)
        snapshot = rendered_html_snapshot(
            html_text,
            final_url,
            browser_body_text,
            title,
            description,
            original_keywords,
        )
        body_text = snapshot_body_text(snapshot)
        ensure_nonempty_rendered_body(body_text)
        keywords = extract_keywords(title, description, body_text, original_keywords)
        snapshot = update_snapshot_keywords(snapshot, keywords)
        html_dir.mkdir(parents=True, exist_ok=True)
        html_path = html_dir / f"{page_id(final_url)}.html"
        html_path.write_text(snapshot, encoding="utf-8")

        logger.emit("page.index.ready", progress=f"{index}/{total}-step5/5", final_url=final_url, html_path=str(html_path))
        return {
            "url": final_url,
            "title": title,
            "description": description,
            "keywords": keywords,
            "html_path": html_path.relative_to(html_dir.parent).as_posix(),
        }
    finally:
        logger.emit(
            "page.routes.summary",
            progress=f"{index}/{total}-step5/5",
            url=url,
            requests=route_requests,
            blocked=route_blocked,
            skipped=route_skipped,
            elapsed_seconds=round(time.monotonic() - started, 3),
        )


async def run(args: argparse.Namespace) -> tuple[int, int]:
    urls = read_urls(Path(args.url_list))
    if args.debug:
        urls = urls[:5]
    out_dir = Path(args.out_dir)
    out_dir.mkdir(parents=True, exist_ok=True)
    effective_url_list = out_dir / "url-list.txt"
    effective_url_list.write_text("\n".join(urls) + "\n", encoding="utf-8", newline="\n")
    html_dir = out_dir / "html"
    index_path = out_dir / "index.jsonl"
    logger = CrawlLogger(out_dir / "crawl.log")
    successes = 0
    records: list[dict[str, Any]] = []
    captured_final_urls: dict[str, str] = {}
    duplicate_final_urls = 0
    consecutive_timeouts = 0
    max_consecutive_timeouts = 0
    stopped_by_consecutive_timeouts = False
    skipped_pages_after_timeout_stop = 0
    attempted = 0
    try:
        logger.emit(
            "run.start",
            urls=len(urls),
            debug=args.debug,
            timeout_ms=TIMEOUT_MS,
            consecutive_timeout_limit=CONSECUTIVE_TIMEOUT_LIMIT,
            index_path=str(index_path),
        )
        from playwright.async_api import async_playwright

        guard = URLSafetyGuard()
        async with async_playwright() as playwright:
            browser = await playwright.chromium.launch(headless=True, args=launch_args())
            context = await browser.new_context(
                ignore_https_errors=True,
                locale="zh-CN",
                user_agent=USER_AGENT,
                viewport={"width": 1440, "height": 1000},
            )
            try:
                with index_path.open("w", encoding="utf-8", buffering=1) as index_file:
                    for position, url in enumerate(urls, start=1):
                        attempted = position
                        record: dict[str, Any] | None = None
                        final_error: Exception | None = None
                        for attempt in range(2):
                            page = await context.new_page()
                            try:
                                record = await crawl_one(
                                    page,
                                    guard,
                                    url,
                                    html_dir,
                                    logger,
                                    position,
                                    len(urls),
                                    captured_final_urls,
                                )
                                final_error = None
                                break
                            except Exception as exc:
                                final_error = exc
                                if attempt == 0 and retryable_error(exc):
                                    logger.emit(
                                        "page.retry",
                                        progress=f"{position}/{len(urls)}-step1/5",
                                        url=url,
                                        error_type=type(exc).__name__,
                                        error=str(exc),
                                    )
                                    await page.wait_for_timeout(3000)
                                    continue
                                break
                            finally:
                                await page.close()
                        if record is not None:
                            index_file.write(json.dumps(record, ensure_ascii=False, separators=(",", ":")) + "\n")
                            records.append(record)
                            captured_final_urls[str(record["url"])] = url
                            successes += 1
                            logger.emit("page.done", progress=f"{position}/{len(urls)}-step5/5", url=record["url"])
                            if consecutive_timeouts:
                                logger.emit(
                                    "crawl.timeout_streak.reset",
                                    previous=consecutive_timeouts,
                                    outcome="success",
                                    url=record["url"],
                                )
                            consecutive_timeouts = 0
                        elif isinstance(final_error, DuplicateFinalURL):
                            duplicate_final_urls += 1
                            logger.emit(
                                "page.duplicate",
                                progress=f"{position}/{len(urls)}-step5/5",
                                url=final_error.requested_url,
                                final_url=final_error.final_url,
                                first_requested_url=final_error.first_requested_url,
                            )
                            if consecutive_timeouts:
                                logger.emit(
                                    "crawl.timeout_streak.reset",
                                    previous=consecutive_timeouts,
                                    outcome="duplicate",
                                    url=final_error.requested_url,
                                )
                            consecutive_timeouts = 0
                        elif final_error is not None:
                            logger.emit(
                                "page.error",
                                progress=f"{position}/{len(urls)}-step5/5",
                                url=url,
                                error_type=type(final_error).__name__,
                                error=str(final_error),
                            )
                            if is_timeout_error(final_error):
                                consecutive_timeouts += 1
                                max_consecutive_timeouts = max(max_consecutive_timeouts, consecutive_timeouts)
                                logger.emit(
                                    "crawl.timeout_streak",
                                    url=url,
                                    consecutive_timeouts=consecutive_timeouts,
                                    limit=CONSECUTIVE_TIMEOUT_LIMIT,
                                )
                                if consecutive_timeouts >= CONSECUTIVE_TIMEOUT_LIMIT:
                                    stopped_by_consecutive_timeouts = True
                                    skipped_pages_after_timeout_stop = len(urls) - position
                                    logger.emit(
                                        "crawl.timeout_limit.reached",
                                        url=url,
                                        consecutive_timeouts=consecutive_timeouts,
                                        limit=CONSECUTIVE_TIMEOUT_LIMIT,
                                        skipped=skipped_pages_after_timeout_stop,
                                    )
                                    break
                            else:
                                if consecutive_timeouts:
                                    logger.emit(
                                        "crawl.timeout_streak.reset",
                                        previous=consecutive_timeouts,
                                        outcome="error",
                                        error_type=type(final_error).__name__,
                                        url=url,
                                    )
                                consecutive_timeouts = 0
            finally:
                await context.close()
                await browser.close()
        filter_common_keywords(records)
        index_path.write_text(
            "\n".join(json.dumps(record, ensure_ascii=False, separators=(",", ":")) for record in records) + "\n",
            encoding="utf-8",
        )
        logger.emit(
            "run.done",
            requested=len(urls),
            attempted=attempted,
            succeeded=successes,
            failed=attempted - successes - duplicate_final_urls,
            duplicate_final_urls=duplicate_final_urls,
            max_consecutive_timeouts=max_consecutive_timeouts,
            stopped_by_consecutive_timeouts=stopped_by_consecutive_timeouts,
            skipped_pages_after_timeout_stop=skipped_pages_after_timeout_stop,
        )
    finally:
        logger.close()
    return successes, len(urls)


def main() -> int:
    args = parse_args()
    try:
        successes, total = asyncio.run(run(args))
    except Exception as exc:
        print(str(exc), file=sys.stderr)
        return 1
    if successes == 0:
        print(f"No pages were crawled successfully out of {total} URLs.", file=sys.stderr)
        return 1
    print(Path(args.out_dir) / "index.jsonl")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
