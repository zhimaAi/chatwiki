#!/usr/bin/env python3
from __future__ import annotations

import argparse
import asyncio
import html as html_escape
import ipaddress
import os
import socket
import sys
import time
from dataclasses import dataclass
from pathlib import Path
from typing import Any
from urllib.parse import urlparse

try:
    from bs4 import BeautifulSoup
except Exception:  # pragma: no cover - optional cleanup dependency
    BeautifulSoup = None

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")


USER_AGENT = (
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 "
    "(KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"
)
DEFAULT_DNS_TIMEOUT_MS = 1000
DEFAULT_DNS_CACHE_TTL_MS = 60000
DEFAULT_TIMEOUT_MS = 30000
DEFAULT_WAIT_MS = 1000


class URLSafetyError(RuntimeError):
    pass


def _blocked_address(address: ipaddress.IPv4Address | ipaddress.IPv6Address) -> bool:
    return (
        address.is_loopback
        or address.is_private
        or address.is_link_local
        or address.is_unspecified
        or address.is_multicast
    )


class URLSafetyGuard:
    def __init__(
        self,
        dns_timeout_ms: int = DEFAULT_DNS_TIMEOUT_MS,
        dns_cache_ttl_ms: int = DEFAULT_DNS_CACHE_TTL_MS,
    ) -> None:
        if dns_timeout_ms <= 0:
            raise ValueError("dns_timeout_ms must be positive")
        if dns_cache_ttl_ms <= 0:
            raise ValueError("dns_cache_ttl_ms must be positive")
        self.dns_timeout = dns_timeout_ms / 1000
        self.dns_cache_ttl = dns_cache_ttl_ms / 1000
        self._dns_cache: dict[tuple[str, int], tuple[float, tuple[str, ...]]] = {}
        self._dns_inflight: dict[tuple[str, int], asyncio.Task[tuple[str, ...]]] = {}

    async def validate(self, raw_url: str) -> str:
        parsed = urlparse(raw_url)
        if parsed.scheme not in ("http", "https"):
            raise URLSafetyError(f"unsupported URL scheme: {parsed.scheme or '<empty>'}")
        host = parsed.hostname
        if not host:
            raise URLSafetyError("URL host is required")
        try:
            port = parsed.port or (443 if parsed.scheme == "https" else 80)
        except ValueError as exc:
            raise URLSafetyError(f"invalid URL port: {exc}") from exc

        addresses = await self._resolve(host, port)
        for raw_address in addresses:
            address_text = raw_address.split("%", 1)[0]
            try:
                address = ipaddress.ip_address(address_text)
            except ValueError as exc:
                raise URLSafetyError(f"DNS returned an invalid address for {host}: {raw_address}") from exc
            if _blocked_address(address):
                raise URLSafetyError(f"blocked non-public address for {host}: {address}")
        return raw_url

    async def _resolve(self, host: str, port: int) -> tuple[str, ...]:
        try:
            literal = ipaddress.ip_address(host.split("%", 1)[0])
            return (str(literal),)
        except ValueError:
            pass

        key = (host.rstrip(".").lower(), port)
        now = time.monotonic()
        cached = self._dns_cache.get(key)
        if cached is not None:
            expires_at, addresses = cached
            if expires_at > now:
                return addresses
            self._dns_cache.pop(key, None)

        task = self._dns_inflight.get(key)
        if task is None:
            task = asyncio.create_task(self._resolve_uncached(host, port))
            self._dns_inflight[key] = task

            def clear_inflight(completed: asyncio.Task[tuple[str, ...]]) -> None:
                if self._dns_inflight.get(key) is completed:
                    self._dns_inflight.pop(key, None)

            task.add_done_callback(clear_inflight)

        addresses = await asyncio.shield(task)
        self._dns_cache[key] = (time.monotonic() + self.dns_cache_ttl, addresses)
        return addresses

    async def _resolve_uncached(self, host: str, port: int) -> tuple[str, ...]:

        loop = asyncio.get_running_loop()
        try:
            infos = await asyncio.wait_for(
                loop.getaddrinfo(host, port, type=socket.SOCK_STREAM),
                timeout=self.dns_timeout,
            )
        except asyncio.TimeoutError as exc:
            raise URLSafetyError(f"DNS resolution timed out after {self.dns_timeout:g}s: {host}") from exc
        except OSError as exc:
            raise URLSafetyError(f"DNS resolution failed for {host}: {exc}") from exc

        addresses = tuple(dict.fromkeys(info[4][0] for info in infos if info[4]))
        if not addresses:
            raise URLSafetyError(f"DNS returned no addresses for {host}")
        return addresses


@dataclass(frozen=True)
class DomainRule:
    domains: tuple[str, ...]
    body_selector: str
    match_path_prefixes: tuple[str, ...] = ()


DOMAIN_RULES = (
    DomainRule(
        domains=("help.chatwiki.com",),
        body_selector="article.theme-doc-markdown, .theme-doc-markdown, main article",
    ),
    DomainRule(
        domains=("www.yuque.com",),
        body_selector="article.article-content",
    ),
    DomainRule(
        domains=("feishu.cn",),
        body_selector=".docx-page-main, .wiki-doc-content, .suite-docx, main",
        match_path_prefixes=("/wiki/", "/docx/"),
    ),
    DomainRule(
        domains=("docs.openclaw.ai",),
        body_selector="article.article, main article",
    ),
    DomainRule(
        domains=("help.aliyun.com",),
        body_selector="#pc-markdown-container .markdown-body, main.aliyun-docs-view",
    ),
    DomainRule(
        domains=("kancloud.cn",),
        body_selector=".content",
    ),
    DomainRule(
        domains=("mp.weixin.qq.com",),
        body_selector="#js_content, .rich_media_content",
    ),
)


def normalize_domain(domain: str) -> str:
    domain = domain.strip().lower()
    if not domain:
        return ""
    if "://" not in domain:
        domain = "https://" + domain
    return (urlparse(domain).hostname or "").lower()


def rule_for_url(raw_url: str) -> DomainRule | None:
    parsed = urlparse(raw_url)
    host = (parsed.hostname or "").lower()
    for rule in DOMAIN_RULES:
        for domain in rule.domains:
            domain_host = normalize_domain(domain)
            domain_matches = host == domain_host or host.endswith("." + domain_host)
            path_matches = not rule.match_path_prefixes or any(
                parsed.path.startswith(prefix) for prefix in rule.match_path_prefixes
            )
            if domain_matches and path_matches:
                return rule
    return None


def nonnegative_int(value: str) -> int:
    try:
        parsed = int(value)
    except ValueError as exc:
        raise argparse.ArgumentTypeError(f"expected integer, got {value!r}") from exc
    if parsed < 0:
        raise argparse.ArgumentTypeError("expected a non-negative integer")
    return parsed


def positive_int(value: str) -> int:
    parsed = nonnegative_int(value)
    if parsed <= 0:
        raise argparse.ArgumentTypeError("expected a positive integer")
    return parsed


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Use headless Chromium to fetch the latest script-stripped rendered HTML snapshot for one URL.",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    parser.add_argument("url", help="URL to fetch.")
    parser.add_argument("--out", help=argparse.SUPPRESS)
    parser.add_argument("--timeout-ms", default=DEFAULT_TIMEOUT_MS, type=positive_int, help="Navigation timeout.")
    parser.add_argument("--wait-ms", default=DEFAULT_WAIT_MS, type=nonnegative_int, help="Extra wait before capturing HTML.")
    parser.add_argument("--selector", help="Optional CSS selector to wait for before capturing HTML.")
    parser.add_argument("--dns-timeout-ms", default=DEFAULT_DNS_TIMEOUT_MS, type=positive_int, help=argparse.SUPPRESS)
    return parser.parse_args()


async def read_body_text(page: Any) -> str:
    try:
        return await page.locator("body").inner_text(timeout=5000)
    except Exception:
        return ""


async def wait_for_stable_body_text(page: Any, timeout_ms: int) -> None:
    deadline = asyncio.get_running_loop().time() + min(timeout_ms / 1000, 30)
    last_len = -1
    stable_count = 0
    while asyncio.get_running_loop().time() < deadline:
        text = await read_body_text(page)
        current_len = len(" ".join(text.split()))
        if current_len == last_len and current_len > 0:
            stable_count += 1
        else:
            stable_count = 0
        last_len = current_len
        if stable_count >= 2:
            return
        await page.wait_for_timeout(1000)


async def wait_for_rendered_body(page: Any, page_url: str, explicit_selector: str | None, timeout_ms: int) -> None:
    rule = rule_for_url(page_url)
    selector = explicit_selector or (rule.body_selector if rule and rule.body_selector else "article, main")
    if not selector:
        return
    try:
        await page.wait_for_selector(
            selector,
            state="attached",
            timeout=timeout_ms if explicit_selector else min(timeout_ms, 5000),
        )
    except Exception:
        if explicit_selector:
            raise


def collapse_space(value: str) -> str:
    return " ".join((value or "").split()).strip()


def clean_body_text(value: str) -> str:
    value = value.replace("\r\n", "\n").replace("\r", "\n")
    lines = [collapse_space(line) for line in value.split("\n")]
    return "\n".join(line for line in lines if line)


def ensure_nonempty_rendered_body(body_text: str) -> None:
    if not clean_body_text(body_text):
        raise RuntimeError(
            "rendered body is empty after page load; the page scripts or content may have failed to load"
        )


def select_first_body_nodes_soup(soup: Any, selector_group: str) -> list[Any]:
    """Return nodes from the first selector that matches, treating commas as priority fallbacks."""
    for selector in (item.strip() for item in selector_group.split(",")):
        if not selector:
            continue
        try:
            nodes = list(soup.select(selector))
        except Exception:
            continue
        if nodes:
            return nodes
    return []


def split_keywords(raw: str) -> list[str]:
    if not raw.strip():
        return []
    seen: set[str] = set()
    result: list[str] = []
    for part in raw.replace("，", ",").replace("；", ";").replace(";", ",").split(","):
        for item in part.split():
            item = item.strip()
            key = item.lower()
            if item and key not in seen:
                seen.add(key)
                result.append(item)
    return result


def first_nonempty(*values: str) -> str:
    for value in values:
        value = collapse_space(value)
        if value:
            return value
    return ""


def append_text_section(snapshot: Any, parent: Any, body_text: str) -> None:
    section = snapshot.new_tag("section", attrs={"data-rendered-text": "true"})
    for paragraph in [item.strip() for item in body_text.splitlines() if item.strip()]:
        node = snapshot.new_tag("p")
        node.string = paragraph
        section.append(node)
    if section.contents:
        parent.append(section)


def remove_unwanted_soup_nodes(soup: Any) -> None:
    for node in soup.select("script, style, noscript, template, link[rel='stylesheet']"):
        node.decompose()
    for node in soup.find_all(True):
        attrs = {}
        for key, value in node.attrs.items():
            lowered = key.lower().strip()
            if lowered == "style" or lowered.startswith("on"):
                continue
            attrs[key] = value
        node.attrs = attrs


def meta_content_soup(soup: Any, attr_name: str, attr_value: str) -> str:
    for meta in soup.find_all("meta"):
        if collapse_space(str(meta.get(attr_name, ""))).lower() == attr_value.lower():
            return collapse_space(str(meta.get("content", "")))
    return ""


def build_snapshot_from_soup(
    page_url: str,
    title: str,
    description: str,
    keywords: list[str],
    body_nodes: list[Any],
    body_text: str,
) -> str:
    snapshot = BeautifulSoup("<!DOCTYPE html><html><head></head><body></body></html>", "html.parser")
    head = snapshot.head
    body = snapshot.body
    if head is None or body is None:
        return "<!DOCTYPE html>\n<html><body></body></html>\n"

    head.append(snapshot.new_tag("meta", attrs={"charset": "utf-8"}))
    head.append(snapshot.new_tag("meta", attrs={"name": "source-url", "content": page_url}))
    title_node = snapshot.new_tag("title")
    title_node.string = title
    head.append(title_node)
    if description:
        head.append(snapshot.new_tag("meta", attrs={"name": "description", "content": description}))
    if keywords:
        head.append(snapshot.new_tag("meta", attrs={"name": "keywords", "content": ", ".join(keywords)}))

    body["data-source-url"] = page_url
    main = snapshot.new_tag("main", attrs={"data-rendered-snapshot": "true"})
    if body_text:
        append_text_section(snapshot, main, body_text)

    rendered_html = snapshot.new_tag("section", attrs={"data-rendered-html": "true"})
    for node in body_nodes:
        if getattr(node, "name", "") == "body":
            html_fragment = "".join(str(child) for child in node.contents)
        else:
            html_fragment = str(node)
        fragment = BeautifulSoup(html_fragment, "html.parser")
        for child in list(fragment.contents):
            rendered_html.append(child)
    if rendered_html.contents:
        main.append(rendered_html)

    body.append(main)
    return snapshot.prettify(formatter="minimal") + "\n"


def build_snapshot_from_text(page_url: str, title: str, description: str, keywords: list[str], body_text: str) -> str:
    escaped_title = html_escape.escape(title)
    escaped_description = html_escape.escape(description, quote=True)
    escaped_keywords = html_escape.escape(", ".join(keywords), quote=True)
    escaped_url = html_escape.escape(page_url, quote=True)
    paragraphs = "\n".join(
        f"<p>{html_escape.escape(paragraph)}</p>"
        for paragraph in [item.strip() for item in body_text.splitlines() if item.strip()]
    )
    meta_description = f'<meta name="description" content="{escaped_description}">\n' if description else ""
    meta_keywords = f'<meta name="keywords" content="{escaped_keywords}">\n' if keywords else ""
    return (
        "<!DOCTYPE html>\n"
        "<html><head>\n"
        '<meta charset="utf-8">\n'
        f'<meta name="source-url" content="{escaped_url}">\n'
        f"<title>{escaped_title}</title>\n"
        f"{meta_description}{meta_keywords}"
        "</head>\n"
        f'<body data-source-url="{escaped_url}"><main data-rendered-snapshot="true">\n'
        f'<section data-rendered-text="true">\n{paragraphs}\n</section>\n'
        "</main></body></html>\n"
    )


async def meta_content_from_page(page: Any, selector: str) -> str:
    try:
        return collapse_space(await page.locator(selector).first().get_attribute("content", timeout=1000) or "")
    except Exception:
        return ""


async def read_page_metadata(page: Any) -> tuple[str, str, list[str]]:
    title = collapse_space(await page.title())
    description = await meta_content_from_page(
        page,
        "meta[name='description'], meta[property='og:description'], meta[name='twitter:description']",
    )
    keywords = split_keywords(await meta_content_from_page(page, "meta[name='keywords']"))
    return title, description, keywords


def rendered_html_snapshot(
    html_text: str,
    final_url: str,
    browser_body_text: str,
    fallback_title: str,
    fallback_description: str,
    fallback_keywords: list[str],
) -> str:
    rendered_text = clean_body_text(browser_body_text)
    if BeautifulSoup is None:
        return build_snapshot_from_text(final_url, fallback_title, fallback_description, fallback_keywords, rendered_text)

    soup = BeautifulSoup(html_text, "html.parser")
    remove_unwanted_soup_nodes(soup)
    head = soup.head or soup

    title = first_nonempty(
        meta_content_soup(head, "property", "og:title"),
        meta_content_soup(head, "name", "twitter:title"),
        head.title.get_text(" ", strip=True) if head.title else "",
        fallback_title,
    )
    description = first_nonempty(
        meta_content_soup(head, "name", "description"),
        meta_content_soup(head, "property", "og:description"),
        meta_content_soup(head, "name", "twitter:description"),
        fallback_description,
    )
    keywords = split_keywords(meta_content_soup(head, "name", "keywords")) or fallback_keywords

    rule = rule_for_url(final_url)
    body_nodes = select_first_body_nodes_soup(soup, rule.body_selector) if rule and rule.body_selector else []
    matched_rule_body = bool(body_nodes)
    if not body_nodes:
        body_nodes = [soup.body or soup]
    body_text = clean_body_text("\n\n".join(node.get_text("\n", strip=True) for node in body_nodes))
    if len(rendered_text) > len(body_text) and (not matched_rule_body or not body_text):
        body_text = rendered_text
    return build_snapshot_from_soup(final_url, title, description, keywords, body_nodes, body_text)


async def fetch_rendered_html(args: argparse.Namespace) -> tuple[str, str]:
    try:
        from playwright.async_api import async_playwright
    except ImportError as exc:
        raise RuntimeError(
            "Playwright is required. Install it with: "
            "python3 -m pip install playwright beautifulsoup4 && python3 -m playwright install chromium"
        ) from exc

    launch_args = [
        "--disable-blink-features=AutomationControlled",
        "--disable-dev-shm-usage",
    ]
    if os.environ.get("PLAYWRIGHT_NO_SANDBOX", "1").lower() not in {"0", "false", "no"}:
        launch_args.extend(["--no-sandbox", "--disable-setuid-sandbox"])

    url_guard = URLSafetyGuard(args.dns_timeout_ms)
    await url_guard.validate(args.url)

    async with async_playwright() as playwright:
        browser = await playwright.chromium.launch(headless=True, args=launch_args)
        context = await browser.new_context(
            ignore_https_errors=True,
            locale="zh-CN",
            user_agent=USER_AGENT,
            viewport={"width": 1440, "height": 1000},
        )
        page = await context.new_page()
        blocked_navigation_error: URLSafetyError | None = None

        async def validate_route(route: Any, request: Any) -> None:
            nonlocal blocked_navigation_error
            request_url = request.url
            if not request_url.startswith(("http://", "https://")):
                if request.is_navigation_request():
                    blocked_navigation_error = URLSafetyError(f"blocked navigation scheme: {urlparse(request_url).scheme}")
                await route.abort("blockedbyclient")
                return
            try:
                await url_guard.validate(request_url)
            except URLSafetyError as exc:
                if request.is_navigation_request():
                    blocked_navigation_error = exc
                await route.abort("blockedbyclient")
                return
            await route.continue_()

        await page.route("**/*", validate_route)
        try:
            try:
                response = await page.goto(args.url, timeout=args.timeout_ms, wait_until="domcontentloaded")
            except Exception as exc:
                if blocked_navigation_error is not None:
                    raise blocked_navigation_error from exc
                raise
            if response is not None and response.status >= 400:
                raise RuntimeError(f"http status {response.status}")
            try:
                await page.wait_for_load_state("networkidle", timeout=min(args.timeout_ms, 10000))
            except Exception:
                pass
            await wait_for_rendered_body(page, page.url, args.selector, args.timeout_ms)
            await wait_for_stable_body_text(page, args.timeout_ms)
            if args.wait_ms:
                await page.wait_for_timeout(args.wait_ms)
            final_url = page.url
            await url_guard.validate(final_url)
            browser_body_text = await read_body_text(page)
            ensure_nonempty_rendered_body(browser_body_text)
            title, description, keywords = await read_page_metadata(page)
            html_text = await page.content()
            return final_url, rendered_html_snapshot(html_text, final_url, browser_body_text, title, description, keywords)
        finally:
            await page.close()
            await context.close()
            await browser.close()


def write_output(html_text: str, out_path: str | None) -> None:
    if not out_path:
        sys.stdout.write(html_text)
        return
    path = Path(out_path)
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(html_text, encoding="utf-8")
    print(path)


def main() -> int:
    args = parse_args()
    try:
        _, html_text = asyncio.run(fetch_rendered_html(args))
    except Exception as exc:
        print(str(exc), file=sys.stderr)
        return 1
    write_output(html_text, args.out)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
