#!/usr/bin/env python3
"""Emit a bounded, source-proportional metadata outline from a web index."""

from __future__ import annotations

import argparse
import json
import sys
from collections import OrderedDict
from pathlib import Path
from typing import Any
from urllib.parse import urlsplit


if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")


DEFAULT_MAX_ITEMS = 60
MAX_ITEMS = 60
TITLE_MAX_CHARS = 160
DESCRIPTION_MAX_CHARS = 400
KEYWORD_MAX_ITEMS = 12
KEYWORD_MAX_CHARS = 80
URL_MAX_CHARS = 2048


def compact_text(value: object, limit: int) -> str:
    text = " ".join(str(value or "").split()).strip()
    if len(text) <= limit:
        return text
    return text[: limit - 1].rstrip() + "…"


def source_group(url: str) -> str:
    parsed = urlsplit(url)
    if parsed.scheme.lower() not in {"http", "https"} or not parsed.netloc:
        raise ValueError(f"index contains an invalid public URL: {url}")
    return f"{parsed.scheme.lower()}://{parsed.netloc.casefold()}"


def load_pages(path: Path) -> list[dict[str, Any]]:
    pages: list[dict[str, Any]] = []
    seen_urls: set[str] = set()
    for line_number, raw_line in enumerate(
        path.read_text(encoding="utf-8-sig").splitlines(),
        start=1,
    ):
        if not raw_line.strip():
            continue
        try:
            page = json.loads(raw_line)
        except json.JSONDecodeError as exc:
            raise ValueError(f"invalid JSONL at {path}:{line_number}: {exc}") from exc
        if not isinstance(page, dict):
            raise ValueError(f"index record at {path}:{line_number} must be an object")
        for field in ("url", "title", "description", "keywords", "html_path"):
            if field not in page:
                raise ValueError(f"index record at {path}:{line_number} is missing {field}")
        url = str(page["url"]).strip()
        source_group(url)
        if not url or url in seen_urls:
            raise ValueError(
                f"index contains an empty or duplicate URL at {path}:{line_number}: {url}"
            )
        if len(url) > URL_MAX_CHARS:
            raise ValueError(
                f"index URL exceeds {URL_MAX_CHARS} characters at {path}:{line_number}"
            )
        if not isinstance(page["keywords"], list):
            raise ValueError(
                f"index record keywords at {path}:{line_number} must be an array"
            )
        seen_urls.add(url)
        pages.append(page)
    if not pages:
        raise ValueError(f"index has no page records: {path}")
    return pages


def allocate_quotas(group_sizes: list[int], limit: int) -> list[int]:
    total = sum(group_sizes)
    if total <= limit:
        return list(group_sizes)

    quotas = [limit * size // total for size in group_sizes]
    remainders = [limit * size % total for size in group_sizes]
    remaining = limit - sum(quotas)
    order = sorted(range(len(group_sizes)), key=lambda item: (-remainders[item], item))
    for index in order[:remaining]:
        quotas[index] += 1

    if len(group_sizes) <= limit:
        for index, quota in enumerate(quotas):
            if quota:
                continue
            donors = [
                candidate
                for candidate, value in enumerate(quotas)
                if value > 1
            ]
            if not donors:
                raise AssertionError("cannot preserve source coverage within the outline limit")
            donor = max(
                donors,
                key=lambda candidate: (
                    quotas[candidate],
                    group_sizes[candidate],
                    -candidate,
                ),
            )
            quotas[donor] -= 1
            quotas[index] = 1
    return quotas


def evenly_sample(
    values: list[dict[str, Any]],
    quota: int,
) -> list[dict[str, Any]]:
    if quota <= 0:
        return []
    if quota >= len(values):
        return list(values)
    if quota == 1:
        return [values[(len(values) - 1) // 2]]
    indexes = [index * (len(values) - 1) // (quota - 1) for index in range(quota)]
    return [values[index] for index in indexes]


def compact_page(page: dict[str, Any], group: str) -> dict[str, Any]:
    keywords: list[str] = []
    seen: set[str] = set()
    for value in page["keywords"]:
        text = compact_text(value, KEYWORD_MAX_CHARS)
        key = text.casefold()
        if text and key not in seen:
            seen.add(key)
            keywords.append(text)
        if len(keywords) >= KEYWORD_MAX_ITEMS:
            break
    return {
        "source": group,
        "url": str(page["url"]).strip(),
        "title": compact_text(page["title"], TITLE_MAX_CHARS),
        "description": compact_text(page["description"], DESCRIPTION_MAX_CHARS),
        "keywords": keywords,
    }


def proportional_outline(
    pages: list[dict[str, Any]],
    limit: int,
) -> tuple[list[dict[str, Any]], list[dict[str, Any]]]:
    grouped: OrderedDict[str, list[dict[str, Any]]] = OrderedDict()
    for page in pages:
        grouped.setdefault(source_group(str(page["url"])), []).append(page)

    groups = list(grouped.items())
    quotas = allocate_quotas([len(group_pages) for _, group_pages in groups], limit)
    outline: list[dict[str, Any]] = []
    source_summary: list[dict[str, Any]] = []
    for (group, group_pages), quota in zip(groups, quotas):
        if quota <= 0:
            continue
        selected = evenly_sample(group_pages, quota)
        outline.extend(compact_page(page, group) for page in selected)
        source_summary.append({
            "source": group,
            "pages": len(group_pages),
            "outline_items": len(selected),
        })
    if len(outline) != min(len(pages), limit):
        raise AssertionError("proportional metadata outline did not fill the requested limit")
    return outline, source_summary


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Emit a bounded metadata outline from a WebToSkill JSONL index."
    )
    parser.add_argument("--index", required=True, type=Path)
    parser.add_argument("--max-items", type=int, default=DEFAULT_MAX_ITEMS)
    args = parser.parse_args()
    if not 1 <= args.max_items <= MAX_ITEMS:
        parser.error(f"--max-items must be between 1 and {MAX_ITEMS}")
    return args


def main() -> int:
    args = parse_args()
    pages = load_pages(args.index)
    outline, sources = proportional_outline(pages, args.max_items)
    print(json.dumps({
        "status": "complete",
        "pages": len(pages),
        "source_groups": len({
            source_group(str(page["url"]))
            for page in pages
        }),
        "represented_source_groups": len(sources),
        "outline_items": len(outline),
        "outline_truncated": len(outline) < len(pages),
        "sources": sources,
        "outline": outline,
    }, ensure_ascii=False, separators=(",", ":")))
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:
        print(f"ERROR: {exc}", file=sys.stderr, flush=True)
        raise SystemExit(1)
