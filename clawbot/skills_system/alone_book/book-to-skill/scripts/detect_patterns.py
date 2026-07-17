#!/usr/bin/env python3
"""Detect section-delimiter patterns in markdown files.

Scans all .md files and tests a set of common chapter/section patterns,
reporting match counts, chunk-size estimates, and sample headers.
The AI can then instantly pick the best top-level splitter without
iterating through patterns one by one.

Usage:
  python3 scripts/detect_patterns.py input_md/
  python3 scripts/detect_patterns.py input_md/ --sample-count 8
  python3 scripts/detect_patterns.py input_md/ --extra-pattern '^第[0-9]+节.*'
"""

import os
import re
import sys
import argparse
from pathlib import Path

# ── Predefined common section patterns ──────────────────────────────────
# Format: (description, regex)
# Order matters: more specific / higher-priority patterns first so the AI
# can see the most likely candidate at a glance.
PREDEFINED_PATTERNS = [
    ("Markdown H1",             r"^# .+"),
    ("Markdown H2",             r"^## .+"),
    ("Markdown H3",             r"^### .+"),
    ("Markdown H4",             r"^#### .+"),
    ('Chinese "第X章/节/篇/部"',     r"^第[一二三四五六七八九十百千\d]+[章节篇部].*"),
    ('Chinese "第X卷/集"',          r"^第[一二三四五六七八九十百千\d]+[卷集].*"),
    ("Numeric 1. / 1、/ 1 ",   r"^\d+[\.\、\s]\s*.+"),
    ("Chinese ordinals 一、/ 一.",      r"^[一二三四五六七八九十]+[、．\.]\s*.+"),
]


def count_matches(content, pattern):
    """Return list of (start_pos, matched_text) for each pattern match."""
    try:
        compiled = re.compile(pattern, re.MULTILINE)
    except re.error as e:
        return None, str(e)
    matches = [(m.start(), m.group().strip()) for m in compiled.finditer(content)]
    return matches, None


def estimate_chunks(content, pattern):
    """Simulate split_chapters.py chunk-size distribution for a pattern.

    Returns (chunk_sizes_bytes, sample_headers) where chunk_sizes_bytes is
    a list of byte-lengths for chunks that would be produced, and
    sample_headers is the first N matched header texts.
    """
    try:
        compiled = re.compile(pattern, re.MULTILINE)
        ms = list(compiled.finditer(content))
    except re.error:
        return [], []

    if not ms:
        return [], []

    samples = []
    sizes = []

    for i, m in enumerate(ms):
        header = m.group().strip()
        if len(samples) < 10:
            samples.append(header)

        # content from after this header line to before next header
        line_end = content.find("\n", m.start())
        body_start = (line_end + 1) if line_end >= 0 else len(content)
        body_end = ms[i + 1].start() if i + 1 < len(ms) else len(content)
        body = content[body_start:body_end].strip()

        chunk_text = f"{header}\n\n{body}" if body else header
        sizes.append(len(chunk_text.encode("utf-8")))

    # Prepend preamble to first chunk size
    if ms[0].start() > 0:
        preamble = content[:ms[0].start()].strip()
        if preamble and sizes:
            sizes[0] += len(("\n\n" + preamble).encode("utf-8"))

    return sizes, samples


def get_size_stats(sizes):
    """Return (min_kb, avg_kb, max_kb) from a list of byte sizes."""
    if not sizes:
        return 0.0, 0.0, 0.0
    kb_list = [s / 1024.0 for s in sizes]
    return min(kb_list), sum(kb_list) / len(kb_list), max(kb_list)


def scan_directory(input_dir, extra_patterns, sample_count):
    """Core: scan all .md files, test every pattern, build report."""
    md_files = sorted(Path(input_dir).glob("*.md"))
    if not md_files:
        print(f"No .md files found in {input_dir}")
        return

    # Read all files into memory (aggregate for cross-file patterns like # headers)
    file_contents = {}
    total_bytes = 0
    for fp in md_files:
        try:
            content = fp.read_text(encoding="utf-8")
            file_contents[fp.name] = content
            total_bytes += len(content.encode("utf-8"))
        except Exception as e:
            print(f"Warning: skip {fp.name}: {e}", file=sys.stderr)

    if not file_contents:
        print("No readable .md files.")
        return

    # Build pattern list: predefined + extra
    patterns_to_test = list(PREDEFINED_PATTERNS)
    for ep in extra_patterns:
        patterns_to_test.append((f"Custom: {ep}", ep))

    # ── Collect results per pattern ──────────────────────────────────
    results = []

    for desc, regex in patterns_to_test:
        total_matches = 0
        matched_files = 0
        all_sizes = []
        all_samples = []
        per_file_info = []  # (filename, matches, sizes)

        for fname, content in file_contents.items():
            sizes, samples = estimate_chunks(content, regex)
            if sizes:
                matched_files += 1
                total_matches += len(sizes)
                all_sizes.extend(sizes)
                all_samples.extend(samples)
                per_file_info.append((fname, len(sizes), sizes))
            else:
                per_file_info.append((fname, 0, []))

        min_kb, avg_kb, max_kb = get_size_stats(all_sizes)
        results.append({
            "desc": desc,
            "regex": regex,
            "matches": total_matches,
            "files": matched_files,
            "min_kb": min_kb,
            "avg_kb": avg_kb,
            "max_kb": max_kb,
            "samples": all_samples[:sample_count],
            "per_file": per_file_info,
        })

    # ── Output report ─────────────────────────────────────────────────
    border = "=" * 80
    print(border)
    print(f"Pattern Detection Report")
    print(f"Source: {input_dir}  ({len(md_files)} file(s), {total_bytes / 1024:.1f} KB total)")
    print(border)
    print()

    # Sort: primary by file coverage (desc), secondary by match count (desc)
    results.sort(key=lambda r: (-r["files"], -r["matches"]))

    # Table header
    print(f"{'Pattern':<36} {'Matches':>7} {'Files':>5} {'MinKB':>7} {'AvgKB':>7} {'MaxKB':>8}")
    print("-" * 76)

    for r in results:
        if r["matches"] == 0:
            print(f"{r['desc']:<36} {'—':>7} {'—':>5} {'—':>7} {'—':>7} {'—':>8}")
        else:
            print(f"{r['desc']:<36} {r['matches']:>7} {r['files']:>5} "
                  f"{r['min_kb']:>6.1f} {r['avg_kb']:>6.1f} {r['max_kb']:>7.1f}")

    print()

    # ── Recommendation ───────────────────────────────────────────────
    # Heuristic: pick the pattern with the highest % of chunks in 4-20KB range
    # and the highest file coverage. If nothing fits, report closest.
    best = None
    best_score = -1
    for r in results:
        if r["matches"] == 0:
            continue
        # Score = file_coverage * in_range_ratio
        file_cov = r["files"] / len(file_contents)
        # Count chunks in 4-20KB target range
        in_range = 0
        all_sz = []
        for _, _, sizes in r["per_file"]:
            all_sz.extend(sizes)
        for sz in all_sz:
            kb = sz / 1024.0
            if 4 <= kb <= 20:
                in_range += 1
        in_range_ratio = in_range / len(all_sz) if all_sz else 0
        score = file_cov * 0.4 + in_range_ratio * 0.6
        if score > best_score:
            best_score = score
            best = r

    if best:
        print("── Recommendation ──")
        print(f"  Top-level splitter: --pattern '{best['regex']}'")
        print(f"  Reason: {best['matches']} sections in {best['files']} file(s), "
              f"avg chunk {best['avg_kb']:.1f} KB (target: 5-16KB)")
        print()

    # ── Per-pattern details (samples + per-file breakdown) ───────────
    for r in results:
        if r["matches"] == 0:
            continue
        print(f"── Pattern: {r['desc']}  ({r['regex']}) ──")
        print(f"   Sections: {r['matches']}  "
              f"Files: {r['files']}/{len(file_contents)}  "
              f"Size: {r['min_kb']:.1f}~{r['avg_kb']:.1f}~{r['max_kb']:.1f} KB")

        # Per-file breakdown (only when multi-file)
        if len(r["per_file"]) > 1:
            for fname, cnt, sizes in r["per_file"]:
                if cnt == 0:
                    print(f"   {fname}: no match")
                else:
                    fmin, favg, fmax = get_size_stats(sizes)
                    print(f"   {fname}: {cnt} sections  {fmin:.1f}~{favg:.1f}~{fmax:.1f} KB")
        print()

        # Sample headers
        if r["samples"]:
            print(f"   Sample headers (first {len(r['samples'])}):")
            for s in r["samples"]:
                print(f"     {s}")
            print()

    print(border)
    print("Done. Use the recommendation above or pick a pattern based on samples.")
    print(border)


def main():
    parser = argparse.ArgumentParser(
        description="Detect section-delimiter patterns in markdown files"
    )
    parser.add_argument(
        "input_dir",
        help="Directory containing .md files (typically input_md/)",
    )
    parser.add_argument(
        "--sample-count", "-n",
        type=int,
        default=5,
        help="Number of sample headers to show per pattern (default: 5)",
    )
    parser.add_argument(
        "--extra-pattern", "-e",
        action="append",
        default=[],
        help="Extra regex pattern to test (repeatable). Use shell-safe quoting.",
    )

    args = parser.parse_args()

    if not os.path.isdir(args.input_dir):
        print(f"Error: {args.input_dir} is not a directory", file=sys.stderr)
        sys.exit(1)

    scan_directory(args.input_dir, args.extra_pattern, args.sample_count)


if __name__ == "__main__":
    main()
