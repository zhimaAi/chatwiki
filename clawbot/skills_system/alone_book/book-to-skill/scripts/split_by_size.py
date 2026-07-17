#!/usr/bin/env python3
"""Fallback splitter: cut a markdown file into fixed-size chunks by UTF-8 byte count.

Used when no clear section delimiter pattern is available (detect_patterns.py
reports all patterns producing chunks >20KB or <3KB, or total matches ≤ 3).

Strategy:
  - Target chunk size: 10 KB (10240 bytes) by default
  - Break at paragraph boundaries (\n\n) nearest to the target size
  - Fall back to line boundaries (\n) for extra-long paragraphs
  - Output naming follows existing convention: chunk_NNN.md (3-digit zero-padded)

Usage:
  python3 scripts/split_by_size.py --input input_md/article.md --output chunks/
  python3 scripts/split_by_size.py --input input_md/article.md --output chunks/ --max-size 10
"""

import os
import sys
import argparse
from pathlib import Path

DEFAULT_MAX_KB = 10


def _utf8_len(text):
    """Return UTF-8 byte length of text."""
    return len(text.encode("utf-8"))


def _split_long_para(para, max_bytes):
    """Split an oversized paragraph at line boundaries.

    Falls back to byte-level split for single lines exceeding max_bytes.
    Returns list of chunk strings.
    """
    lines = para.split("\n")
    result = []
    line_chunk = ""

    for line in lines:
        if not line.strip():
            continue

        if line_chunk:
            test_line = line_chunk + "\n" + line
            if _utf8_len(test_line) > max_bytes:
                result.append(line_chunk)
                line_chunk = line
            else:
                line_chunk = test_line
        else:
            if _utf8_len(line) > max_bytes:
                # Single line exceeds max_bytes — split at safe UTF-8 byte boundary
                raw = line.encode("utf-8")
                pos = 0
                while pos < len(raw):
                    segment_end = min(pos + max_bytes, len(raw))
                    # Walk back to a valid UTF-8 boundary to avoid corrupting multi-byte chars
                    while segment_end > pos:
                        try:
                            raw[pos:segment_end].decode("utf-8")
                            break
                        except UnicodeDecodeError:
                            segment_end -= 1
                    if segment_end == pos:
                        # Cannot find valid boundary, skip this segment
                        break
                    segment = raw[pos:segment_end].decode("utf-8")
                    result.append(segment)
                    pos = segment_end
                line_chunk = ""
            else:
                line_chunk = line

    if line_chunk:
        result.append(line_chunk)

    return result


def split_content(content, max_bytes):
    """Split text into chunks, breaking at paragraph boundaries.

    Each chunk is approximately max_bytes in UTF-8 size.
    Paragraphs exceeding max_bytes are further split at line boundaries.
    Returns list of chunk strings.
    """
    paragraphs = content.split("\n\n")
    chunks = []
    current = ""

    for para in paragraphs:
        if not para.strip():
            continue

        if _utf8_len(para) > max_bytes:
            # Oversized paragraph — split at line boundaries first
            sub_parts = _split_long_para(para, max_bytes)
            for part in sub_parts:
                if current:
                    test = current + "\n\n" + part
                    if _utf8_len(test) > max_bytes:
                        chunks.append(current)
                        current = part
                    else:
                        current = test
                else:
                    current = part
        else:
            if current:
                test = current + "\n\n" + para
                if _utf8_len(test) > max_bytes:
                    chunks.append(current)
                    current = para
                else:
                    current = test
            else:
                current = para

    if current:
        chunks.append(current)

    return chunks if chunks else [content.strip()]


def process_file(filepath, output_dir, max_bytes):
    """Process a single .md file and write chunks."""
    try:
        with open(filepath, "r", encoding="utf-8") as f:
            content = f.read()
    except Exception as e:
        print(f"Error reading {filepath}: {e}", file=sys.stderr)
        return 0

    if not content.strip():
        return 0

    chunks = split_content(content, max_bytes)

    Path(output_dir).mkdir(parents=True, exist_ok=True)

    # Determine starting chunk number from existing files
    existing = list(Path(output_dir).glob("chunk_*.md"))
    chunk_num = len(existing) + 1

    written = 0
    for chunk in chunks:
        if not chunk.strip():
            continue
        out_path = Path(output_dir) / f"chunk_{chunk_num:03d}.md"
        with open(out_path, "w", encoding="utf-8") as f:
            f.write(chunk.strip())
        chunk_num += 1
        written += 1

    return written


def main():
    parser = argparse.ArgumentParser(
        description="Fallback splitter: cut markdown into fixed-size chunks by UTF-8 byte count"
    )
    parser.add_argument(
        "--input", "-i",
        required=True,
        help="Input .md file to split",
    )
    parser.add_argument(
        "--output", "-o",
        required=True,
        help="Output directory for chunk_NNN.md files (typically chunks/)",
    )
    parser.add_argument(
        "--max-size", "-s",
        type=int,
        default=DEFAULT_MAX_KB,
        help=f"Target chunk size in KB (default: {DEFAULT_MAX_KB})",
    )

    args = parser.parse_args()

    if not os.path.isfile(args.input):
        print(f"Error: {args.input} is not a file", file=sys.stderr)
        sys.exit(1)

    max_bytes = args.max_size * 1024

    fname = os.path.basename(args.input)
    total = process_file(args.input, args.output, max_bytes)

    if total == 0:
        print(f"No content in {fname}, nothing written")
    else:
        print(f"Done: {fname} → {total} chunk(s) in {args.output}/")

    # Also report chunk size distribution
    existing = sorted(Path(args.output).glob("chunk_*.md"))
    if existing:
        print(f"Total chunks in {args.output}/: {len(existing)}")


if __name__ == "__main__":
    main()
