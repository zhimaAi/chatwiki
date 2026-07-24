#!/usr/bin/env python3
"""Validate crawl artifacts and emit a bounded completion summary."""

from __future__ import annotations

import argparse
import json
import re
import sys
from collections import Counter
from pathlib import Path
from typing import Any

from build_skill import load_crawl_coverage, load_index, resolve_html_path


if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")


MAX_FAILURE_EXAMPLES = 5
MAX_ERROR_CHARS = 240
FIELD_RE = re.compile(
    r'(?P<key>[A-Za-z_][A-Za-z0-9_]*)='
    r'(?P<value>"(?:\\.|[^"\\])*"|true|false|null|-?\d+(?:\.\d+)?)'
)


def compact_text(value: object, limit: int) -> str:
    text = " ".join(str(value or "").split()).strip()
    if len(text) <= limit:
        return text
    return text[: limit - 1].rstrip() + "…"


def load_url_count(path: Path) -> int:
    if not path.is_file():
        raise ValueError(f"URL list not found: {path}")
    urls = [line.strip() for line in path.read_text(encoding="utf-8-sig").splitlines()]
    urls = [url for url in urls if url and not url.startswith("#")]
    if not urls:
        raise ValueError(f"URL list is empty: {path}")
    if len(set(urls)) != len(urls):
        raise ValueError(f"URL list contains duplicate entries: {path}")
    return len(urls)


def parse_fields(line: str) -> dict[str, Any]:
    fields: dict[str, Any] = {}
    for match in FIELD_RE.finditer(line):
        fields[match.group("key")] = json.loads(match.group("value"))
    return fields


def load_failure_summary(path: Path) -> dict[str, Any]:
    lines = path.read_text(encoding="utf-8-sig").splitlines()
    starts = [
        index
        for index, line in enumerate(lines)
        if "[crawl_urls]" in line and " run.start " in f" {line} "
    ]
    completions = [
        index
        for index, line in enumerate(lines)
        if "[crawl_urls]" in line and " run.done " in f" {line} "
    ]
    if not starts or not completions or starts[-1] >= completions[-1]:
        raise ValueError(f"crawl log has no complete latest crawl run: {path}")

    failures: list[dict[str, str]] = []
    by_error_type: Counter[str] = Counter()
    for line in lines[starts[-1] : completions[-1] + 1]:
        if "[crawl_urls]" not in line or " page.error " not in f" {line} ":
            continue
        fields = parse_fields(line)
        error_type = compact_text(fields.get("error_type", "UnknownError"), 80)
        by_error_type[error_type] += 1
        if len(failures) < MAX_FAILURE_EXAMPLES:
            failures.append({
                "url": compact_text(fields.get("url", ""), 2048),
                "error_type": error_type,
                "error": compact_text(fields.get("error", ""), MAX_ERROR_CHARS),
            })
    count = sum(by_error_type.values())
    return {
        "count": count,
        "by_error_type": dict(sorted(by_error_type.items())),
        "examples": failures,
        "examples_truncated": count > len(failures),
    }


def load_duplicate_summary(path: Path) -> dict[str, Any]:
    lines = path.read_text(encoding="utf-8-sig").splitlines()
    starts = [
        index
        for index, line in enumerate(lines)
        if "[crawl_urls]" in line and " run.start " in f" {line} "
    ]
    completions = [
        index
        for index, line in enumerate(lines)
        if "[crawl_urls]" in line and " run.done " in f" {line} "
    ]
    if not starts or not completions or starts[-1] >= completions[-1]:
        raise ValueError(f"crawl log has no complete latest crawl run: {path}")

    examples: list[dict[str, str]] = []
    count = 0
    for line in lines[starts[-1] : completions[-1] + 1]:
        if "[crawl_urls]" not in line or " page.duplicate " not in f" {line} ":
            continue
        count += 1
        if len(examples) < MAX_FAILURE_EXAMPLES:
            fields = parse_fields(line)
            examples.append({
                "url": compact_text(fields.get("url", ""), 2048),
                "final_url": compact_text(fields.get("final_url", ""), 2048),
            })
    return {
        "count": count,
        "examples": examples,
        "examples_truncated": count > len(examples),
    }


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Validate a WebToSkill crawl without loading its full log into model context."
    )
    parser.add_argument(
        "--index",
        required=True,
        type=Path,
        help="Crawl index.jsonl; url-list.txt and crawl.log are read from the same directory.",
    )
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    index_path = args.index.resolve()
    crawl_dir = index_path.parent
    pages = load_index(index_path)
    url_count = load_url_count(crawl_dir / "url-list.txt")
    coverage = load_crawl_coverage(crawl_dir / "crawl.log", len(pages))
    if int(coverage["requested"]) != url_count:
        raise ValueError("crawl coverage mismatch: requested must equal URL-list rows")

    for page in pages:
        resolve_html_path(index_path, str(page["html_path"]))

    failures = load_failure_summary(crawl_dir / "crawl.log")
    if failures["count"] != int(coverage["failed"]):
        raise ValueError("crawl failure summary mismatch: page.error count must equal failed")
    duplicates = load_duplicate_summary(crawl_dir / "crawl.log")
    if duplicates["count"] != int(coverage["duplicate_final_urls"]):
        raise ValueError(
            "crawl duplicate summary mismatch: page.duplicate count must equal duplicate_final_urls"
        )

    print(json.dumps({
        "status": "complete",
        "url_count": url_count,
        "index_rows": len(pages),
        "html_paths_checked": len(pages),
        "run_done": coverage,
        "failure_summary": failures,
        "duplicate_summary": duplicates,
    }, ensure_ascii=False, separators=(",", ":")))
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:
        print(f"ERROR: {exc}", file=sys.stderr, flush=True)
        raise SystemExit(1)
