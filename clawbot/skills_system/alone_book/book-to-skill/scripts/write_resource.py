#!/usr/bin/env python3
"""Write chunk content to resource files without LLM token overhead.

AI reads a chunk → identifies its structure → selects a subcommand:
- copy:   single chunk → single resource file
- merge:  multiple chunks → single resource file (concatenated)
- split:  single chunk with multiple ##/### headings → split into separate files

Usage:
  python3 scripts/write_resource.py copy --src chunks/chunk_010.md \
      --dst skill/tskill/resources/第1章/概述.md
  python3 scripts/write_resource.py merge --src chunks/chunk_010.md chunks/chunk_011.md \
      --dst skill/tskill/resources/第1章/概述.md
  python3 scripts/write_resource.py split --src chunks/chunk_005.md \
      --pattern '^## .+' --dst-dir skill/tskill/resources/第1章/
"""

import os
import re
import sys
import argparse
from pathlib import Path

# Safe character set for filenames: Chinese, English, digits, hyphen, underscore.
# Spaces are NOT retained — they are replaced with hyphens to prevent broken
# Markdown link navigation (e.g. [text](path/file name.md) won't render).
# Everything else (commas, semicolons, question marks, full-width punctuation, etc.) is stripped.
_SAFE_CHARS_RE = re.compile(r'[^\u4e00-\u9fa5a-zA-Z0-9\-_]')

# Section number prefix patterns to strip from filenames:
#   Chinese: 一 、, 二、, 三.
#   Arabic:  1 、, 1. , 1) , (1) , 1.1 , 1.1.1
#   Order matters: multi-level Arabic (1.1) must precede single-level (1.)
_SECTION_NUM_RE = re.compile(
    r'^(?:'
    r'[一二三四五六七八九十]+\s*[、．.]?\s*'      # 一 、, 二、, 三.
    r'|\(\d+\)\s*'                              # (1)
    r'|\d+(?:[\.\-]\d+)+\s*'                    # 1.1, 1.1.1 (multi-level, no delimiter needed)
    r'|\d+\s*[、．.)]\s*'                         # 1 、, 1. , 1)  (single-level with delimiter)
    r')'
)


def _clean_filename(text):
    """Extract a safe filename from heading text.

    Strips leading # markers, section-number prefixes, and all unsafe characters.
    Only Chinese, English, digits, hyphens, and underscores are retained.
    Spaces are replaced with hyphens to prevent broken Markdown links.
    Example: '## 一、Overview' → 'Overview'
    Example: '### 1.1 Setup Guide' → 'Setup-Guide'
    Example: 'Hello World, test' → 'Hello-World-test'
    Example: 'No auto-reply but message auto-replies' → 'No-auto-reply-but-message-auto-replies'
    """
    text = text.strip()
    # Strip leading # markers (e.g. '### 1.1 Overview' → '1.1 Overview')
    text = re.sub(r'^#+\s*', '', text)
    # Strip section number prefixes (一 、, 1. , (1) , 1.1 , etc.)
    text = _SECTION_NUM_RE.sub('', text)
    # Strip all unsafe characters (only keep Chinese, English, digits, hyphen, underscore)
    text = _SAFE_CHARS_RE.sub('', text)
    # Normalize whitespace: collapse multiple spaces → hyphen, trim
    # Hyphens prevent broken Markdown links like [text](path/file name.md)
    text = re.sub(r'\s+', '-', text).strip()
    # Strip trailing separators that may result from stripping
    text = text.strip('.-_ ')
    # Fallback: if the result is empty, use a safe default
    if not text:
        text = 'untitled'
    return text


def _resolve_duplicate_name(filepath):
    """If filepath exists, return filepath with (1), (2) suffix appended.

    Does NOT modify disk — caller writes after resolving the name.
    """
    if not filepath.exists():
        return filepath

    stem = filepath.stem
    suffix = filepath.suffix
    parent = filepath.parent
    counter = 1
    while True:
        new_name = f"{stem} ({counter}){suffix}"
        new_path = parent / new_name
        if not new_path.exists():
            return new_path
        counter += 1


def _ensure_parent(path):
    """Create parent directories if they don't exist."""
    path.parent.mkdir(parents=True, exist_ok=True)


def _read_file(filepath):
    """Read a file with UTF-8 encoding. Returns content or sys.exit(1)."""
    if not os.path.isfile(filepath):
        print(f"Error: file not found: {filepath}", file=sys.stderr)
        sys.exit(1)
    try:
        with open(filepath, "r", encoding="utf-8") as f:
            return f.read()
    except Exception as e:
        print(f"Error: failed to read {filepath}: {e}", file=sys.stderr)
        sys.exit(1)


def _write_file(filepath, content):
    """Write content to filepath with UTF-8, creating parent dirs."""
    _ensure_parent(filepath)
    with open(filepath, "w", encoding="utf-8") as f:
        f.write(content)


# ========== Subcommand: copy ==========


def cmd_copy(args):
    """Copy a single chunk file to a resource file."""
    src = Path(args.src)
    dst = Path(args.dst)

    if not src.is_file():
        print(f"Error: source file not found: {src}", file=sys.stderr)
        sys.exit(1)

    content = _read_file(src)
    _write_file(dst, content)
    print(f"Copied {src.name} → {dst}")


# ========== Subcommand: merge ==========


def cmd_merge(args):
    """Merge multiple chunk files into a single resource file."""
    src_paths = [Path(s) for s in args.src]
    dst = Path(args.dst)

    # Validate all sources exist BEFORE writing (atomicity)
    missing = [s for s in src_paths if not s.is_file()]
    if missing:
        names = ", ".join(str(s) for s in missing)
        print(f"Error: source file(s) not found: {names}", file=sys.stderr)
        sys.exit(1)

    # Warn if exceeding --max-src
    if len(src_paths) > args.max_src:
        print(
            f"Warning: merging {len(src_paths)} chunks, exceeds recommended max of {args.max_src}. "
            f"Consider splitting or reducing merge scope."
        )

    # Read and concatenate
    parts = []
    for sp in src_paths:
        parts.append(_read_file(sp))

    merged = "\n\n".join(parts)
    _write_file(dst, merged)
    names = ", ".join(s.name for s in src_paths)
    print(f"Merged [{names}] → {dst}")


# ========== Subcommand: split ==========


def cmd_split(args):
    """Split a single chunk into multiple resource files by heading pattern."""
    src = Path(args.src)
    dst_dir = Path(args.dst_dir)
    pattern = args.pattern
    chapter_name = args.chapter_name

    if not src.is_file():
        print(f"Error: source file not found: {src}", file=sys.stderr)
        sys.exit(1)

    content = _read_file(src)

    # Compile regex safely
    try:
        compiled = re.compile(pattern, re.MULTILINE)
    except re.error as e:
        print(f"Error: invalid regex pattern '{pattern}': {e}", file=sys.stderr)
        sys.exit(1)

    matches = list(re.finditer(pattern, content, re.MULTILINE))

    if not matches:
        print(f"No headings matching '{pattern}' found in {src.name}. Nothing to split.")
        return

    _ensure_parent(dst_dir)

    written = 0

    # 1. Preamble: content before the first heading
    preamble_start = 0
    preamble_end = matches[0].start()
    preamble = content[preamble_start:preamble_end].strip()

    if preamble:
        # Extract chapter name from # heading in preamble
        if not chapter_name:
            h1_match = re.search(r'^#\s+(.+)', preamble, re.MULTILINE)
            if h1_match:
                chapter_name = h1_match.group(1).strip()
            elif not chapter_name:
                # Fallback: use last component of dst-dir
                chapter_name = dst_dir.name

        # Write preamble (including # heading line) as chapter file
        ch_filename = _clean_filename(chapter_name) if chapter_name else "chapter"
        ch_path = _resolve_duplicate_name(dst_dir / f"{ch_filename}.md")
        _write_file(ch_path, preamble)
        print(f"  {ch_path}")
        written += 1

    # Fallback chapter name for sub-sections if not set yet
    if not chapter_name:
        chapter_name = dst_dir.name

    # 2. Sub-sections: each heading + body → a file
    for i, m in enumerate(matches):
        # Extract full heading line
        line_start = m.start()
        line_end = content.find('\n', m.start())
        if line_end < 0:
            line_end = len(content)
        heading_line = content[line_start:line_end].strip()
        heading_text = _clean_filename(heading_line)

        # Body from end of heading line to start of next heading (or EOF)
        body_start = line_end + 1
        body_end = matches[i + 1].start() if i + 1 < len(matches) else len(content)
        body = content[body_start:body_end].strip()

        section_content = f"{heading_line}\n\n{body}" if body else heading_line

        filename = heading_text if heading_text else f"section_{i + 1:02d}"
        sec_path = _resolve_duplicate_name(dst_dir / f"{filename}.md")
        _write_file(sec_path, section_content)
        print(f"  {sec_path}")
        written += 1

    if written > 0:
        print(f"Split {src.name} → {written} file(s) in {dst_dir}/")


# ========== Main ==========


def main():
    parser = argparse.ArgumentParser(
        description="Write chunk content to resource files without LLM token overhead"
    )
    sub = parser.add_subparsers(dest="command", required=True, help="Subcommand")

    # copy
    p_copy = sub.add_parser("copy", help="Copy single chunk → single resource")
    p_copy.add_argument("--src", required=True, help="Source chunk file (e.g. chunks/chunk_001.md)")
    p_copy.add_argument("--dst", required=True, help="Destination resource file (e.g. skill/t/resources/Ch1.md)")

    # merge
    p_merge = sub.add_parser("merge", help="Merge multiple chunks → single resource")
    p_merge.add_argument("--src", nargs="+", required=True, help="Source chunk files (1 or more)")
    p_merge.add_argument("--dst", required=True, help="Destination resource file")
    p_merge.add_argument("--max-src", type=int, default=3,
                         help="Max source files before warning (default: 3)")

    # split
    p_split = sub.add_parser("split", help="Split chunk by headings → multiple resource files")
    p_split.add_argument("--src", required=True, help="Source chunk file")
    p_split.add_argument("--pattern", required=True, help="Regex pattern for sub-headings (e.g. '^## .+')")
    p_split.add_argument("--dst-dir", required=True, help="Destination directory for resource files")
    p_split.add_argument("--chapter-name", default=None,
                         help="Chapter name (extracted from # heading if omitted)")

    args = parser.parse_args()

    if args.command == "copy":
        cmd_copy(args)
    elif args.command == "merge":
        cmd_merge(args)
    elif args.command == "split":
        cmd_split(args)


if __name__ == "__main__":
    main()
