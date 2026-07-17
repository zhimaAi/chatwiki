#!/usr/bin/env python3
"""Extract a line range from a file and write to a target file.

Used when AI reads a chunk, identifies logical sections by line numbers,
and needs to carve out each section into a separate resource file.
Also works as a fine-grained alternative to write_resource.py split/merge
in the normal flow.

Usage:
  python3 scripts/extract_lines.py chunks/chunk_010.md 10 50 \
      --dst skill/tskill/resources/概述/安装步骤.md

  python3 scripts/extract_lines.py chunks/chunk_010.md 10 50 \
      --dst skill/tskill/resources/概述/安装步骤.md --prepend '## 安装步骤\n\n'

Parameters:
  src          Source file path
  start_line   1-based start line (inclusive)
  end_line     1-based end line (inclusive)
  --dst        Destination file path (required)
  --prepend    Optional text to prepend before extracted content (e.g. heading)
"""

import os
import sys
import argparse
from pathlib import Path


def extract_lines(src_path, start, end):
    """Extract lines [start, end] (1-based, inclusive) from src_path.

    Returns the exact content of those lines.
    """
    if not os.path.isfile(src_path):
        print(f"Error: source file not found: {src_path}", file=sys.stderr)
        sys.exit(1)

    with open(src_path, "r", encoding="utf-8") as f:
        all_lines = f.readlines()

    total_lines = len(all_lines)
    if start < 1:
        print(f"Error: start line must be >= 1, got {start}", file=sys.stderr)
        sys.exit(1)
    if end > total_lines:
        print(f"Warning: end line {end} exceeds file length {total_lines}, "
              f"clamping to {total_lines}")
        end = total_lines
    if start > end:
        print(f"Error: start line {start} > end line {end}", file=sys.stderr)
        sys.exit(1)

    extracted = "".join(all_lines[start - 1:end])
    return extracted


def write_output(content, dst_path, prepend=None):
    """Write content to dst_path, optionally with prepended text.
    Creates parent directories automatically.

    The prepend text supports \\n escape sequences (converted to real newlines)
    to work reliably across different shell environments.
    """
    Path(dst_path).parent.mkdir(parents=True, exist_ok=True)

    output = content
    if prepend:
        # Convert literal \\n to real newlines for shell compatibility.
        # In bash single quotes, \n is literal; in double quotes it's expanded.
        # This normalizes both cases.
        prepend = prepend.replace("\\n", "\n")
        output = prepend + content

    with open(dst_path, "w", encoding="utf-8") as f:
        f.write(output)


def main():
    parser = argparse.ArgumentParser(
        description="Extract a line range from a file and write to target"
    )
    parser.add_argument(
        "src",
        help="Source file path (e.g. chunks/chunk_010.md)",
    )
    parser.add_argument(
        "start_line",
        type=int,
        help="1-based start line number (inclusive)",
    )
    parser.add_argument(
        "end_line",
        type=int,
        help="1-based end line number (inclusive)",
    )
    parser.add_argument(
        "--dst",
        required=True,
        help="Destination file path (e.g. skill/tskill/resources/Ch1/section.md)",
    )
    parser.add_argument(
        "--prepend",
        default=None,
        help="Optional text to prepend before extracted content (e.g. heading line)",
    )

    args = parser.parse_args()

    content = extract_lines(args.src, args.start_line, args.end_line)
    write_output(content, args.dst, args.prepend)

    # Report
    src_name = os.path.basename(args.src)
    dst_name = args.dst
    print(f"Extracted lines {args.start_line}-{args.end_line} "
          f"from {src_name} → {dst_name}")


if __name__ == "__main__":
    main()
