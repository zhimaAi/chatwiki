#!/usr/bin/env python3
from __future__ import annotations

import argparse
import asyncio
import datetime as dt
import json
import os
import sys
import xml.etree.ElementTree as ET
from dataclasses import dataclass
from pathlib import Path
from typing import Any, TextIO
from urllib.parse import parse_qs, quote, urljoin, urlparse, urlunparse

from bs4 import BeautifulSoup

from fetch_rendered_html import (
    URLSafetyError,
    URLSafetyGuard,
    USER_AGENT,
    wait_for_rendered_body,
    wait_for_stable_body_text,
)


if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")


NAVIGATION_TIMEOUT_MS = 60_000
YUQUE_NAVIGATION_TIMEOUT_MS = 120_000
DIRECTORY_COLLECTION_TIMEOUT_MS = 120_000
SKIPPED_RESOURCE_TYPES = {"font", "image", "media"}
GENERIC_DIRECTORY_SELECTORS = (
    "nav a[href]",
    "aside a[href]",
    "[role='navigation'] a[href]",
    "[class*='sidebar'] a[href]",
    "[class*='sider'] a[href]",
    "[class*='toc'] a[href]",
    "[class*='catalog'] a[href]",
    "[class*='directory'] a[href]",
    "[class*='menu'] a[href]",
    "[class*='tree'] a[href]",
)


@dataclass(frozen=True)
class DirectoryRule:
    domains: tuple[str, ...]
    item_selectors: tuple[str, ...]
    scroll_selectors: tuple[str, ...] = ("body",)
    path_prefixes: tuple[str, ...] = ()
    linked_domains: tuple[str, ...] = ()
    linked_path_prefixes: tuple[str, ...] = ()
    token_attr: str = ""
    token_param: str = ""
    token_url_template: str = ""
    text_selector: str = ""
    sitemap_path: str = ""
    sitemap_path_prefixes: tuple[str, ...] = ()
    payload_summary: bool = False
    strict_directory: bool = False


DIRECTORY_RULES = (
    DirectoryRule(
        domains=("help.chatwiki.com",),
        item_selectors=(),
        sitemap_path="/sitemap.xml",
        sitemap_path_prefixes=("/docs/",),
        strict_directory=True,
    ),
    DirectoryRule(
        domains=("www.yuque.com",),
        item_selectors=(".ant-tabs-content-holder a[href]",),
        scroll_selectors=(".ant-tabs-content-holder",),
    ),
    DirectoryRule(
        domains=("feishu.cn",),
        item_selectors=('[data-node-uid*="wikiToken="]',),
        scroll_selectors=(".workspace-scroll-area",),
        path_prefixes=("/wiki/", "/docx/"),
        linked_domains=("feishu.cn",),
        linked_path_prefixes=("/wiki/", "/docx/"),
        token_attr="data-node-uid",
        token_param="wikiToken",
        token_url_template="https://my.feishu.cn/wiki/%s",
        text_selector=".workspace-tree-view-node-content",
    ),
    DirectoryRule(
        domains=("docs.openclaw.ai",),
        item_selectors=('aside.sidebar nav a[href^="/zh-CN/"], nav.tabs a[href^="/zh-CN/"]',),
    ),
    DirectoryRule(
        domains=("help.aliyun.com",),
        item_selectors=('#pc-markdown-container a[href*="/zh/"], main.aliyun-docs-view a[href*="/zh/"]',),
    ),
    DirectoryRule(
        domains=("kancloud.cn",),
        item_selectors=(".catalog a[href]",),
        payload_summary=True,
        strict_directory=True,
    ),
    DirectoryRule(
        domains=("mp.weixin.qq.com",),
        item_selectors=(),
    ),
)


COLLECT_LINKS_SCRIPT = r"""
async ({itemSelector, scrollSelector, tokenAttr, tokenParam, tokenURLTemplate, textSelector, maxDurationMs}) => {
  const sleep = ms => new Promise(resolve => setTimeout(resolve, ms));
  const startedAt = Date.now();
  const deadline = startedAt + maxDurationMs;
  const seen = new Map();
  const expanded = new Set();
  const stats = {
    roots: 0,
    passes: 0,
    scrollSteps: 0,
    expanded: 0,
    converged: true,
    truncated: false,
    truncationReason: "",
    elapsedMs: 0
  };
  const attr = (el, name) => (!el || !name) ? "" : (el.getAttribute(name) || "");
  const readURL = el => {
    if (tokenAttr && tokenURLTemplate) {
      const raw = attr(el, tokenAttr);
      let token = raw;
      if (raw && tokenParam) token = new URLSearchParams(raw).get(tokenParam) || "";
      if (token) return tokenURLTemplate.replace("%s", encodeURIComponent(token));
    }
    const anchor = el.matches("a") ? el : el.querySelector("a");
    return anchor ? (anchor.getAttribute("href") || anchor.href || "") : "";
  };
  const readText = el => {
    const node = textSelector ? el.querySelector(textSelector) : el;
    return ((node || el).innerText || (node || el).textContent || "").trim();
  };
  const collect = () => {
    for (const el of document.querySelectorAll(itemSelector)) {
      const url = readURL(el);
      if (url && !seen.has(url)) seen.set(url, readText(el));
    }
  };
  const expandVisible = scroller => {
    const selector = [
      '[class*="collapseIconWrapper"][class*="collapsed"]',
      '[class*="workspace-tree-view-node-expand-arrow--collapsed"]'
    ].join(",");
    const viewport = scroller.getBoundingClientRect();
    for (const icon of document.querySelectorAll(selector)) {
      if (Date.now() >= deadline) return;
      const rect = icon.getBoundingClientRect();
      if (rect.bottom < viewport.top || rect.top > viewport.bottom) continue;
      const owner = icon.closest('a, [data-node-uid], [data-id], [data-key], [role="treeitem"]') || icon;
      const key = ["href", "data-node-uid", "data-id", "data-key", "aria-label", "title"]
        .map(name => attr(owner, name)).filter(Boolean).join("|") || (owner.textContent || "").trim();
      if (expanded.has(key)) continue;
      expanded.add(key);
      stats.expanded += 1;
      const anchor = icon.closest("a");
      const href = anchor ? anchor.getAttribute("href") : null;
      if (anchor) anchor.removeAttribute("href");
      icon.click();
      if (anchor && href !== null) anchor.setAttribute("href", href);
    }
  };
  collect();
  const roots = Array.from(document.querySelectorAll(scrollSelector || "body"));
  stats.roots = roots.length;
  if (!roots.length) {
    stats.converged = false;
    stats.truncated = true;
    stats.truncationReason = "scroll_root_missing";
  }
  for (const root of roots) {
    let previousCount = -1;
    let previousHeight = -1;
    let rootConverged = false;
    for (let pass = 0; pass < 8 && Date.now() < deadline; pass++) {
      stats.passes += 1;
      const candidates = [root, ...root.querySelectorAll("*")]
        .filter(el => el.scrollHeight > el.clientHeight + 20)
        .sort((a, b) => (b.scrollHeight - b.clientHeight) - (a.scrollHeight - a.clientHeight));
      const scroller = candidates[0] || root;
      const step = Math.max(100, Math.floor((scroller.clientHeight || 600) * 0.75));
      const maxScroll = scroller.scrollHeight + scroller.clientHeight;
      for (let y = 0; y <= maxScroll && Date.now() < deadline; y += step) {
        stats.scrollSteps += 1;
        scroller.scrollTop = y;
        scroller.dispatchEvent(new Event("scroll", {bubbles: true}));
        await sleep(120);
        collect();
        expandVisible(scroller);
      }
      await sleep(200);
      collect();
      if (seen.size === previousCount && scroller.scrollHeight === previousHeight) {
        rootConverged = true;
        break;
      }
      previousCount = seen.size;
      previousHeight = scroller.scrollHeight;
    }
    if (!rootConverged) {
      stats.converged = false;
      stats.truncated = true;
      stats.truncationReason = Date.now() >= deadline ? "time_budget" : "pass_budget";
      break;
    }
  }
  stats.elapsedMs = Date.now() - startedAt;
  return {
    links: Array.from(seen.entries()).map(([url, text]) => ({url, text})),
    stats
  };
}
"""


class PrepareLogger:
    def __init__(self, path: Path) -> None:
        path.parent.mkdir(parents=True, exist_ok=True)
        self.file: TextIO = path.open("a", encoding="utf-8", buffering=1)

    def close(self) -> None:
        self.file.close()

    def emit(self, event: str, **fields: Any) -> None:
        timestamp = dt.datetime.now().astimezone().isoformat(timespec="seconds")
        details = " ".join(
            f"{key}={json.dumps(value, ensure_ascii=False, default=str)}"
            for key, value in fields.items()
            if value is not None
        )
        line = f"[prepare_urls] {timestamp} {event}" + (f" {details}" if details else "")
        print(line, file=sys.stderr, flush=True)
        self.file.write(line + "\n")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Create a normalized URL list. A single URL triggers directory discovery; multiple URLs are written directly."
    )
    parser.add_argument("--out", required=True, help="Output UTF-8 text file, one URL per line.")
    parser.add_argument("urls", nargs="+", help="One or more HTTP(S) URLs.")
    return parser.parse_args()


def normalize_url(raw_url: str, base_url: str | None = None) -> str:
    value = raw_url.strip()
    parsed = urlparse(urljoin(base_url or value, value))
    if parsed.scheme not in ("http", "https") or not parsed.hostname:
        raise ValueError(f"invalid HTTP(S) URL: {raw_url}")
    host = parsed.hostname.lower()
    netloc = host
    if parsed.port and not ((parsed.scheme == "http" and parsed.port == 80) or (parsed.scheme == "https" and parsed.port == 443)):
        netloc = f"{host}:{parsed.port}"
    return urlunparse((parsed.scheme, netloc, parsed.path or "/", "", parsed.query, ""))


def host_matches(host: str, domain: str) -> bool:
    normalized = domain.lower().lstrip(".")
    return host == normalized or host.endswith("." + normalized)


def rule_for_url(url: str) -> DirectoryRule | None:
    parsed = urlparse(url)
    host = (parsed.hostname or "").lower()
    for rule in DIRECTORY_RULES:
        if any(host_matches(host, domain) for domain in rule.domains) and (
            not rule.path_prefixes or any(parsed.path.startswith(prefix) for prefix in rule.path_prefixes)
        ):
            return rule
    return None


def navigation_timeout_ms(url: str) -> int:
    host = (urlparse(url).hostname or "").lower()
    if host_matches(host, "www.yuque.com"):
        return YUQUE_NAVIGATION_TIMEOUT_MS
    return NAVIGATION_TIMEOUT_MS


def in_scope(candidate: str, start_url: str, rule: DirectoryRule | None) -> bool:
    start = urlparse(start_url)
    target = urlparse(candidate)
    if target.hostname == start.hostname:
        return True
    if not rule:
        return False
    host = (target.hostname or "").lower()
    return any(host_matches(host, domain) for domain in rule.linked_domains) and (
        not rule.linked_path_prefixes or any(target.path.startswith(prefix) for prefix in rule.linked_path_prefixes)
    )


def static_directory_links(html_text: str, page_url: str, rule: DirectoryRule | None) -> list[dict[str, str]]:
    soup = BeautifulSoup(html_text, "html.parser")
    selectors = rule.item_selectors if rule else GENERIC_DIRECTORY_SELECTORS
    links: list[dict[str, str]] = []
    for selector in selectors:
        try:
            elements = soup.select(selector)
        except Exception:
            continue
        for element in elements:
            raw_url = ""
            if rule and rule.token_attr and rule.token_url_template:
                token_source = str(element.get(rule.token_attr) or "")
                token = parse_qs(token_source).get(rule.token_param, [""])[0] if rule.token_param else token_source
                if token:
                    raw_url = rule.token_url_template.replace("%s", quote(token, safe=""))
            if not raw_url:
                anchor = element if element.name == "a" else element.find("a")
                if anchor is not None:
                    raw_url = str(anchor.get("href") or anchor.get("data-href") or "")
            if raw_url:
                links.append({"url": raw_url, "text": element.get_text(" ", strip=True)})
    return links


def payload_directory_links(html_text: str, rule: DirectoryRule | None) -> list[dict[str, str]]:
    if not rule or not rule.payload_summary or not html_text:
        return []
    soup = BeautifulSoup(html_text, "html.parser")
    links: list[dict[str, str]] = []

    def append_nodes(nodes: Any) -> None:
        if not isinstance(nodes, list):
            return
        for node in nodes:
            if not isinstance(node, dict):
                continue
            path = str(node.get("path") or node.get("id") or "").strip()
            if path:
                links.append({"url": path, "text": str(node.get("title") or node.get("name") or "")})
            append_nodes(node.get("articles"))
            append_nodes(node.get("children"))

    for script in soup.select('script[type="application/payload+json"]'):
        try:
            payload = json.loads(script.get_text())
        except (TypeError, json.JSONDecodeError):
            continue
        if isinstance(payload, dict):
            append_nodes(payload.get("summary"))
    return links


def launch_args() -> list[str]:
    args = ["--disable-blink-features=AutomationControlled", "--disable-dev-shm-usage"]
    if os.environ.get("PLAYWRIGHT_NO_SANDBOX", "1").lower() not in {"0", "false", "no"}:
        args.extend(["--no-sandbox", "--disable-setuid-sandbox"])
    return args


def sitemap_document_urls(xml_bytes: bytes, start_url: str, rule: DirectoryRule) -> tuple[list[str], dict[str, Any]]:
    root = ET.fromstring(xml_bytes)
    start = urlparse(start_url)
    locale_prefix = "/zh" if start.path == "/zh" or start.path.startswith("/zh/") else ""
    urls = [start_url]
    seen = {start_url}
    entries = 0
    matched = 0
    rewritten_hosts = 0
    for node in root.findall(".//{*}loc"):
        value = (node.text or "").strip()
        if not value:
            continue
        entries += 1
        source = urlparse(value)
        source_path = source.path
        if source_path.startswith("/zh/docs/"):
            source_path = source_path[3:]
        if not any(source_path.startswith(prefix) for prefix in rule.sitemap_path_prefixes):
            continue
        matched += 1
        if (source.hostname or "").lower() != (start.hostname or "").lower():
            rewritten_hosts += 1
        candidate = urlunparse((start.scheme, start.netloc, locale_prefix + source_path, "", source.query, ""))
        candidate = normalize_url(candidate)
        if candidate not in seen:
            seen.add(candidate)
            urls.append(candidate)
    return urls, {
        "entries": entries,
        "matched": matched,
        "rewritten_hosts": rewritten_hosts,
        "locale_prefix": locale_prefix,
    }


async def discover_sitemap_urls(
    start_url: str,
    rule: DirectoryRule,
    guard: URLSafetyGuard,
    logger: PrepareLogger,
    playwright_factory: Any,
) -> list[str]:
    sitemap_url = normalize_url(rule.sitemap_path, start_url)
    await guard.validate(sitemap_url)
    timeout_ms = navigation_timeout_ms(start_url)
    logger.emit("directory.sitemap.start", url=sitemap_url, timeout_ms=timeout_ms)
    async with playwright_factory() as playwright:
        browser = await playwright.chromium.launch(headless=True, args=launch_args())
        context = await browser.new_context(ignore_https_errors=True, locale="zh-CN", user_agent=USER_AGENT)
        page = await context.new_page()
        blocked_navigation_error: URLSafetyError | None = None

        async def route_request(route: Any, request: Any) -> None:
            nonlocal blocked_navigation_error
            try:
                await guard.validate(request.url)
            except URLSafetyError as exc:
                if request.is_navigation_request():
                    blocked_navigation_error = exc
                await route.abort("blockedbyclient")
                return
            await route.continue_()

        await page.route("**/*", route_request)
        try:
            response = None
            for attempt in range(2):
                try:
                    logger.emit(
                        "directory.sitemap.navigation.start",
                        url=sitemap_url,
                        attempt=attempt + 1,
                        timeout_ms=timeout_ms,
                    )
                    response = await page.goto(sitemap_url, timeout=timeout_ms, wait_until="domcontentloaded")
                    break
                except Exception as exc:
                    if blocked_navigation_error is not None:
                        raise blocked_navigation_error from exc
                    if attempt == 1:
                        raise
                    logger.emit(
                        "directory.sitemap.navigation.retry",
                        url=sitemap_url,
                        attempt=attempt + 1,
                        error_type=type(exc).__name__,
                        error=str(exc),
                    )
                    await page.wait_for_timeout(3000)
            if response is None:
                raise RuntimeError("sitemap navigation returned no response")
            if response.status >= 400:
                raise RuntimeError(f"sitemap HTTP status {response.status}")
            await guard.validate(response.url)
            xml_bytes = await response.body()
            urls, stats = sitemap_document_urls(xml_bytes, start_url, rule)
            if len(urls) <= 1:
                raise RuntimeError("sitemap contains no matching documentation URLs")
            logger.emit("directory.sitemap.done", url=sitemap_url, urls=len(urls), **stats)
            return urls
        finally:
            await page.close()
            await context.close()
            await browser.close()


async def discover_urls(start_url: str, logger: PrepareLogger) -> list[str]:
    try:
        from playwright.async_api import async_playwright
    except ImportError as exc:
        raise RuntimeError("Playwright is required to discover directory URLs.") from exc

    guard = URLSafetyGuard()
    await guard.validate(start_url)
    rule = rule_for_url(start_url)
    if rule and rule.sitemap_path:
        return await discover_sitemap_urls(start_url, rule, guard, logger, async_playwright)
    if rule and not rule.item_selectors:
        logger.emit("directory.scan.skipped", url=start_url, reason="no_directory_rule")
        return [start_url]

    navigation_timeout = navigation_timeout_ms(start_url)
    logger.emit(
        "directory.scan.start",
        url=start_url,
        navigation_timeout_ms=navigation_timeout,
        collection_timeout_ms=DIRECTORY_COLLECTION_TIMEOUT_MS,
    )

    async with async_playwright() as playwright:
        browser = await playwright.chromium.launch(headless=True, args=launch_args())
        context = await browser.new_context(
            ignore_https_errors=True,
            locale="zh-CN",
            user_agent=USER_AGENT,
            viewport={"width": 1440, "height": 1000},
        )
        page = await context.new_page()
        blocked_navigation_error: URLSafetyError | None = None

        async def route_request(route: Any, request: Any) -> None:
            nonlocal blocked_navigation_error
            if not request.is_navigation_request() and request.resource_type in SKIPPED_RESOURCE_TYPES:
                await route.abort("blockedbyclient")
                return
            try:
                await guard.validate(request.url)
            except URLSafetyError as exc:
                if request.is_navigation_request():
                    blocked_navigation_error = exc
                await route.abort("blockedbyclient")
                return
            await route.continue_()

        await page.route("**/*", route_request)
        try:
            response = None
            for attempt in range(2):
                try:
                    logger.emit(
                        "directory.navigation.start",
                        url=start_url,
                        attempt=attempt + 1,
                        timeout_ms=navigation_timeout,
                    )
                    response = await page.goto(start_url, timeout=navigation_timeout, wait_until="domcontentloaded")
                    logger.emit("directory.navigation.done", url=start_url, attempt=attempt + 1, final_url=page.url)
                    break
                except Exception as exc:
                    if blocked_navigation_error is not None:
                        raise blocked_navigation_error from exc
                    if attempt == 1:
                        raise
                    logger.emit(
                        "directory.navigation.retry",
                        url=start_url,
                        attempt=attempt + 1,
                        error_type=type(exc).__name__,
                        error=str(exc),
                    )
                    await page.wait_for_timeout(3000)
            if response is not None and response.status >= 400:
                raise RuntimeError(f"HTTP status {response.status}")
            final_url = normalize_url(page.url)
            await guard.validate(final_url)
            raw_links: list[dict[str, str]] = []
            if response is not None and rule and rule.payload_summary:
                try:
                    response_html = (await response.body()).decode("utf-8", errors="replace")
                    payload_links = payload_directory_links(response_html, rule)
                except Exception as exc:
                    payload_links = []
                    logger.emit(
                        "directory.payload.error",
                        final_url=final_url,
                        error_type=type(exc).__name__,
                        error=str(exc),
                    )
                if payload_links:
                    raw_links.extend(payload_links)
                    logger.emit(
                        "directory.payload.done",
                        final_url=final_url,
                        links=len(payload_links),
                        payload_type="application/payload+json",
                    )

            if not raw_links:
                await wait_for_rendered_body(page, final_url, None, navigation_timeout)
                await wait_for_stable_body_text(page, navigation_timeout)
                await page.wait_for_timeout(1000)

                selectors = rule.item_selectors if rule else (", ".join(GENERIC_DIRECTORY_SELECTORS),)
                scroll_selectors = rule.scroll_selectors if rule else ("body",)
                try:
                    await page.wait_for_selector(", ".join(selectors), state="attached", timeout=15_000)
                except Exception:
                    pass
                for item_selector in selectors:
                    for scroll_selector in scroll_selectors:
                        result = await page.evaluate(
                            COLLECT_LINKS_SCRIPT,
                            {
                                "itemSelector": item_selector,
                                "scrollSelector": scroll_selector,
                                "tokenAttr": rule.token_attr if rule else "",
                                "tokenParam": rule.token_param if rule else "",
                                "tokenURLTemplate": rule.token_url_template if rule else "",
                                "textSelector": rule.text_selector if rule else "",
                                "maxDurationMs": DIRECTORY_COLLECTION_TIMEOUT_MS,
                            },
                        )
                        result_links = list(result.get("links", []))
                        raw_links.extend(result_links)
                        logger.emit(
                            "directory.scroll.done",
                            item_selector=item_selector,
                            scroll_selector=scroll_selector,
                            links=len(result_links),
                            **dict(result.get("stats", {})),
                        )
                static_links = static_directory_links(await page.content(), final_url, rule)
                raw_links.extend(static_links)
                logger.emit("directory.static.done", links=len(static_links), final_url=final_url)
        finally:
            await page.close()
            await context.close()
            await browser.close()

    urls = [final_url]
    seen = {final_url}
    for item in raw_links:
        try:
            candidate = normalize_url(str(item.get("url", "")), final_url)
        except ValueError:
            continue
        if candidate in seen or not in_scope(candidate, final_url, rule):
            continue
        seen.add(candidate)
        urls.append(candidate)
    if rule and rule.strict_directory and len(urls) <= 1:
        raise RuntimeError(f"directory discovery returned only the supplied URL for {final_url}")
    logger.emit("directory.scan.done", url=start_url, final_url=final_url, raw_links=len(raw_links), urls=len(urls))
    return urls


async def run(args: argparse.Namespace, logger: PrepareLogger) -> list[str]:
    normalized = []
    seen: set[str] = set()
    for raw_url in args.urls:
        url = normalize_url(raw_url)
        if url not in seen:
            seen.add(url)
            normalized.append(url)
    if len(normalized) == 1:
        try:
            return await discover_urls(normalized[0], logger)
        except URLSafetyError:
            raise
        except Exception as exc:
            rule = rule_for_url(normalized[0])
            if rule and rule.strict_directory:
                logger.emit(
                    "directory.scan.error",
                    url=normalized[0],
                    error_type=type(exc).__name__,
                    error=str(exc),
                    fallback="disabled",
                )
                raise
            logger.emit(
                "directory.scan.error",
                url=normalized[0],
                error_type=type(exc).__name__,
                error=str(exc),
                fallback="supplied_url",
            )
            print(f"directory discovery failed; keeping the supplied URL: {exc}", file=sys.stderr)
            return normalized
    logger.emit("batch.prepare", supplied=len(args.urls), unique_urls=len(normalized), discovery=False)
    return normalized


def main() -> int:
    args = parse_args()
    out_path = Path(args.out)
    logger = PrepareLogger(out_path.parent / "crawl.log")
    try:
        logger.emit("run.start", supplied_urls=len(args.urls), out_path=str(out_path))
        urls = asyncio.run(run(args, logger))
        out_path.parent.mkdir(parents=True, exist_ok=True)
        out_path.write_text("\n".join(urls) + "\n", encoding="utf-8")
        logger.emit("run.done", urls=len(urls), out_path=str(out_path))
    except Exception as exc:
        logger.emit("run.error", error_type=type(exc).__name__, error=str(exc))
        print(str(exc), file=sys.stderr)
        return 1
    finally:
        logger.close()
    print(out_path)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
