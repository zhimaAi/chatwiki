#!/usr/bin/env python3
"""Split markdown by chapter headers into size-constrained chunks.

Uses `re.finditer` (NOT re.split with capturing groups) for reliable
header detection and content extraction. AI only needs to pass --pattern.

Usage:
  python3 scripts/split_chapters.py input_md/ chunks/ --pattern '^# .+'
  python3 scripts/split_chapters.py input_md/ chunks/ --pattern '^## .+'
  python3 scripts/split_chapters.py input_md/ chunks/ --pattern '^第[一二三四五六七八九十百千\\d]+[章节篇]'
"""

import os
import re
import sys
import argparse
from pathlib import Path

DEFAULT_MIN_KB = 5
DEFAULT_MAX_KB = 16


def find_headers_and_split(content, pattern):
    """Find all headers matching `pattern` and split content into chunks.

    Uses re.finditer (position-based slicing) instead of re.split with
    capturing groups, which is fragile and error-prone.

    Returns list of chunk strings. Each chunk starts with its header line.
    Content before the first header (preamble) is prepended to chunk[0].
    """
    try:
        matches = list(re.finditer(pattern, content, re.MULTILINE))
    except re.error as e:
        print(f"Error: invalid regex pattern '{pattern}': {e}", file=sys.stderr)
        sys.exit(1)

    if not matches:
        return [content.strip()] if content.strip() else []

    chunks = []

    for i, m in enumerate(matches):
        # Extend header to cover the FULL line (up to newline), not just the
        # regex match. This prevents heading text from leaking into body when
        # the AI's pattern only matches the marker (e.g. '^# ' vs '^# .+').
        line_start = m.start()
        line_end = content.find('\n', m.start())
        if line_end < 0:
            line_end = len(content)
        header = content[line_start:line_end].strip()
        body_start = line_end + 1  # skip newline after the header line
        body_end = matches[i + 1].start() if i + 1 < len(matches) else len(content)
        body = content[body_start:body_end].strip()

        chunks.append(f"{header}\n\n{body}" if body else header)

    # Prepend preamble (content before the first header) to first chunk
    preamble_start = 0
    preamble_end = matches[0].start()
    if preamble_end > 0:
        preamble = content[preamble_start:preamble_end].strip()
        if preamble:
            chunks[0] = f"{preamble}\n\n{chunks[0]}"

    return chunks


def merge_small_chunks(chunks, min_bytes, max_bytes):
    """Merge consecutive chunks < min_bytes, keeping merged result ≤ max_bytes.

    Uses index-based while loop (NOT for-each with list mutation) so the
    merged chunk is re-checked and no elements are silently skipped.
    """
    i = 0
    while i < len(chunks):
        size = len(chunks[i].encode("utf-8"))
        if size < min_bytes and i + 1 < len(chunks):
            merged = chunks[i] + "\n\n" + chunks[i + 1]
            if len(merged.encode("utf-8")) <= max_bytes:
                chunks[i] = merged
                chunks.pop(i + 1)
                continue  # re-check merged chunk at the same index
        i += 1
    return chunks


def split_large_chunk(chunk, max_bytes):
    """Split a single chunk at paragraph boundaries (\\n\\n).

    Each sub-chunk is ≤ max_bytes. Empty paragraphs are skipped.
    Returns list of sub-chunks.
    """
    paragraphs = [p.strip() for p in chunk.split("\n\n") if p.strip()]
    if not paragraphs:
        return [chunk]

    sub_chunks = []
    current = ""

    for para in paragraphs:
        test = (current + "\n\n" + para) if current else para
        if len(test.encode("utf-8")) <= max_bytes:
            current = test
        else:
            if current:
                sub_chunks.append(current)
            current = para

    if current:
        sub_chunks.append(current)

    return sub_chunks if sub_chunks else [chunk]


def process_file(filepath, output_dir, pattern, min_bytes, max_bytes):
    """Process a single .md file: find headers → split → merge → write."""
    with open(filepath, "r", encoding="utf-8") as f:
        content = f.read()

    if not content.strip():
        return 0

    chunks = find_headers_and_split(content, pattern)
    chunks = merge_small_chunks(chunks, min_bytes, max_bytes)

    Path(output_dir).mkdir(parents=True, exist_ok=True)

    # Determine starting chunk number from existing files
    existing = list(Path(output_dir).glob("chunk_*.md"))
    chunk_num = len(existing) + 1

    for chunk in chunks:
        sub_chunks = split_large_chunk(chunk, max_bytes) if len(chunk.encode("utf-8")) > max_bytes else [chunk]
        for sub in sub_chunks:
            if not sub.strip():
                continue
            out_path = Path(output_dir) / f"chunk_{chunk_num:03d}.md"
            with open(out_path, "w", encoding="utf-8") as f:
                f.write(sub.strip())
            chunk_num += 1

    return chunk_num - len(existing) - 1


def main():
    parser = argparse.ArgumentParser(
        description="Split markdown by chapter headers into size-constrained chunks"
    )
    parser.add_argument("input_dir", help="Directory containing .md files")
    parser.add_argument("output_dir", help="Output directory for chunk_*.md files")
    parser.add_argument(
        "--pattern", "-p",
        default=r"^# .+",
        help="Regex pattern for chapter headers (must match full header line, e.g. '^# .+' or '^第[\\d]+章')",
    )
    parser.add_argument(
        "--min-size", type=int, default=DEFAULT_MIN_KB,
        help=f"Minimum chunk size in KB (default: {DEFAULT_MIN_KB})",
    )
    parser.add_argument(
        "--max-size", type=int, default=DEFAULT_MAX_KB,
        help=f"Maximum chunk size in KB (default: {DEFAULT_MAX_KB})",
    )

    args = parser.parse_args()

    if not os.path.isdir(args.input_dir):
        print(f"Error: {args.input_dir} is not a directory", file=sys.stderr)
        sys.exit(1)

    min_bytes = args.min_size * 1024
    max_bytes = args.max_size * 1024

    md_files = sorted(Path(args.input_dir).glob("*.md"))
    if not md_files:
        print(f"No .md files found in {args.input_dir}")
        sys.exit(0)

    total = 0
    for f in md_files:
        n = process_file(f, args.output_dir, args.pattern, min_bytes, max_bytes)
        total += n

    print(f"Done: {total} chunk(s) written to {args.output_dir}/")


if __name__ == "__main__":
    main()
