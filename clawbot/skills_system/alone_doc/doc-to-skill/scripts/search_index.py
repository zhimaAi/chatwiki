#!/usr/bin/env python3
"""Search a Doc-to-Skill JSONL index with bounded lexical scoring."""

from __future__ import annotations

import argparse
import json
import math
import re
import sys
from pathlib import Path
from typing import Any


if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")


CHINESE_STOPWORDS = {"如何", "怎么", "怎样", "什么", "是否", "进行", "以及", "一个", "可以"}
ENGLISH_STOPWORDS = {
    "a", "an", "are", "can", "do", "does", "for", "how", "in", "is", "of", "on", "please", "the",
    "to", "use", "using", "what", "when", "where", "which", "who", "why",
}


def add_chinese_tokens(output: set[str], term: str) -> None:
    segments = [term]
    for stopword in sorted(CHINESE_STOPWORDS, key=len, reverse=True):
        segments = [piece for segment in segments for piece in segment.replace(stopword, " ").split()]
    for segment in segments:
        if len(segment) == 1:
            output.add(segment)
            continue
        for size in (2, 3):
            if len(segment) >= size:
                output.update(segment[index:index + size] for index in range(len(segment) - size + 1))


def tokens(value: str) -> set[str]:
    value = value.lower()
    output = {
        token
        for token in re.findall(r"[a-z0-9][a-z0-9_-]*", value)
        if token not in ENGLISH_STOPWORDS
    }
    for term in re.findall(r"[\u3400-\u9fff]+", value):
        add_chinese_tokens(output, term)
    return output


def required_matches(query_tokens: set[str]) -> int:
    if len(query_tokens) <= 2:
        return 1
    return max(2, math.ceil(len(query_tokens) * 0.3))


def score_row(row: dict[str, Any], query_tokens: set[str]) -> tuple[int, set[str]]:
    if not query_tokens:
        return 0, set()
    title = tokens(str(row.get("title", "")))
    keywords = tokens(" ".join(map(str, row.get("keywords", []))))
    questions = tokens(" ".join(map(str, row.get("questions", []))))
    body = tokens(str(row.get("summary", "")) + " " + str(row.get("content", "")))
    document_tokens = title | keywords | questions | body
    identifiers = {
        token
        for token in query_tokens
        if re.fullmatch(r"[a-z0-9][a-z0-9_-]*", token) and re.search(r"[0-9_-]", token)
    }
    if identifiers and not identifiers <= document_tokens:
        return 0, set()
    matched = query_tokens & document_tokens
    if not matched:
        return 0, set()
    coverage_bonus = round(10 * len(matched) / len(query_tokens))
    high_signal_matches = query_tokens & (title | keywords | questions)
    broad_match_bonus = 10 if len(high_signal_matches) >= required_matches(query_tokens) else 0
    score = (
        broad_match_bonus
        + coverage_bonus
        + 8 * len(query_tokens & title)
        + 6 * len(query_tokens & keywords)
        + 4 * len(query_tokens & questions)
        + len(query_tokens & body)
    )
    return score, matched


def score(row: dict[str, Any], query_tokens: set[str]) -> int:
    return score_row(row, query_tokens)[0]


def load_rows(index_path: Path) -> list[dict[str, Any]]:
    rows: list[dict[str, Any]] = []
    for line_number, raw_line in enumerate(index_path.read_text(encoding="utf-8-sig").splitlines(), start=1):
        if not raw_line.strip():
            continue
        try:
            row = json.loads(raw_line)
        except json.JSONDecodeError as exc:
            raise ValueError(f"invalid JSONL at {index_path}:{line_number}: {exc}") from exc
        if not isinstance(row, dict):
            raise ValueError(f"index record at {index_path}:{line_number} must be an object")
        rows.append(row)
    if not rows:
        raise ValueError(f"index has no knowledge records: {index_path}")
    return rows


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Search a generated document skill JSONL index with bounded output.")
    parser.add_argument("query")
    parser.add_argument("--index", default="references/doc-index.jsonl")
    parser.add_argument("--limit", type=int, default=5, help="Maximum candidates to return; clamped to 1-5.")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    skill_dir = Path(__file__).resolve().parents[1]
    index_path = Path(args.index)
    if not index_path.is_absolute():
        index_path = skill_dir / index_path
    if not index_path.is_file():
        print(f"index file not found: {index_path}", file=sys.stderr)
        return 1

    try:
        rows = load_rows(index_path)
    except (OSError, ValueError) as exc:
        print(str(exc), file=sys.stderr)
        return 1
    query_tokens = tokens(args.query)
    limit = min(max(args.limit, 1), 5)
    ranked: list[tuple[int, str, dict[str, Any], set[str]]] = []
    for row in rows:
        value, matched = score_row(row, query_tokens)
        if value > 0:
            ranked.append((value, str(row.get("id", "")), row, matched))
    ranked.sort(key=lambda item: (-item[0], item[1]))

    results: list[dict[str, Any]] = []
    for value, _, row, matched in ranked[:limit]:
        results.append({"score": value, "matched_tokens": sorted(matched), **row})
    print(
        json.dumps(
            {
                "query": args.query,
                "query_tokens": sorted(query_tokens),
                "result_count": len(results),
                "results": results,
            },
            ensure_ascii=False,
            indent=2,
        )
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
