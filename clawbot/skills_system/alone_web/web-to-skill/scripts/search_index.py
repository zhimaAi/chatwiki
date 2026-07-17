#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path
from typing import Any


if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")


ASCII_TERM_RE = re.compile(r"^[a-z0-9][a-z0-9._+-]*$", re.IGNORECASE)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Search a generated web skill JSONL index with bounded output.")
    parser.add_argument("--query", required=True, help="Two to four discriminating terms separated with |.")
    parser.add_argument("--index", default="references/web-index.jsonl", help="JSONL index path relative to the skill.")
    parser.add_argument("--limit", default=3, type=int, help="Maximum candidates to return; clamped to 1-5.")
    return parser.parse_args()


def parse_terms(value: str) -> list[str]:
    terms: list[str] = []
    seen: set[str] = set()
    for raw in re.split(r"[|,，;；\n]+", value):
        term = re.sub(r"\s+", " ", raw).strip().casefold()
        if len(term) < 2 or term in seen:
            continue
        seen.add(term)
        terms.append(term)
    if not terms:
        raise SystemExit("--query must contain at least one term with 2 or more characters")
    return terms[:6]


def load_pages(index_path: Path) -> list[dict[str, Any]]:
    pages: list[dict[str, Any]] = []
    for line_number, raw in enumerate(index_path.read_text(encoding="utf-8-sig").splitlines(), start=1):
        if not raw.strip():
            continue
        try:
            page = json.loads(raw)
        except json.JSONDecodeError as exc:
            raise SystemExit(f"invalid JSONL at {index_path}:{line_number}: {exc}") from exc
        if not isinstance(page, dict):
            raise SystemExit(f"invalid JSONL object at {index_path}:{line_number}")
        pages.append(page)
    return pages


def contains_term(text: str, term: str) -> bool:
    if ASCII_TERM_RE.fullmatch(term):
        pattern = rf"(?<![a-z0-9]){re.escape(term)}(?![a-z0-9])"
        return re.search(pattern, text, re.IGNORECASE) is not None
    return term.casefold() in text.casefold()


def score_page(page: dict[str, Any], terms: list[str]) -> tuple[int, list[str]]:
    title = str(page.get("title", "")).casefold()
    description = str(page.get("description", "")).casefold()
    url = str(page.get("url", "")).casefold()
    keywords = [str(item).casefold() for item in page.get("keywords", []) if str(item).strip()]
    score = 0
    matched: list[str] = []
    for term in terms:
        term_score = 0
        if contains_term(title, term):
            term_score = max(term_score, 24)
        if any(contains_term(keyword, term) for keyword in keywords):
            term_score = max(term_score, 14)
        if contains_term(description, term):
            term_score = max(term_score, 7)
        if contains_term(url, term):
            term_score = max(term_score, 3)
        if term_score:
            score += term_score
            matched.append(term)
    if len(matched) == len(terms):
        score += 12
    return score, matched


def main() -> int:
    args = parse_args()
    terms = parse_terms(args.query)
    limit = min(max(args.limit, 1), 5)
    skill_dir = Path(__file__).resolve().parents[1]
    index_path = Path(args.index)
    if not index_path.is_absolute():
        index_path = skill_dir / index_path
    if not index_path.is_file():
        raise SystemExit(f"index file not found: {index_path}")

    ranked: list[tuple[int, int, str, dict[str, Any], list[str]]] = []
    for page in load_pages(index_path):
        score, matched = score_page(page, terms)
        if not matched:
            continue
        ranked.append((len(matched), score, str(page.get("title", "")), page, matched))
    ranked.sort(key=lambda item: (-item[0], -item[1], item[2]))

    results: list[dict[str, Any]] = []
    for _, score, _, page, matched in ranked[:limit]:
        results.append(
            {
                "score": score,
                "matched_terms": matched,
                "title": page.get("title", ""),
                "description": page.get("description", ""),
                "keywords": page.get("keywords", []),
                "html_path": page.get("html_path", ""),
                "url": page.get("url", ""),
            }
        )

    print(json.dumps({"query_terms": terms, "result_count": len(results), "results": results}, ensure_ascii=False, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
