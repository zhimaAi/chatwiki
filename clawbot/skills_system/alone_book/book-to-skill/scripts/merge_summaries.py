#!/usr/bin/env python3
"""Merge individual summary files in summaries/ into a single summary_index.txt.

Each file in summaries/ contains one Markdown table row (no header).
This script sorts them by filename, concatenates with a header row,
and writes the result to summary_index.txt.

Usage: python3 scripts/merge_summaries.py summaries/ summary_index.txt
"""
import os
import sys


def merge_summaries(summaries_dir, output_path):
    if not os.path.isdir(summaries_dir):
        print(f"Error: {summaries_dir} does not exist", file=sys.stderr)
        sys.exit(1)

    files = sorted(f for f in os.listdir(summaries_dir) if f.endswith(".txt"))
    if not files:
        print(f"Error: no .txt files found in {summaries_dir}", file=sys.stderr)
        sys.exit(1)

    header = "| Chapter Path | Title | Content Summary | Source File Path |\n|---|---|---|---|\n"
    rows = []
    skipped = 0
    for fname in files:
        fpath = os.path.join(summaries_dir, fname)
        try:
            with open(fpath, "r", encoding="utf-8") as f:
                content = f.read().strip()
        except Exception as e:
            print(f"Warning: skip {fname}: {e}", file=sys.stderr)
            skipped += 1
            continue
        if content:
            rows.append(content)
        else:
            skipped += 1

    # If the output file already exists (e.g. pre-created by write_file with
    # wrong permissions), try to remove it so the container user (uid 10001)
    # can recreate it with correct ownership.
    if os.path.exists(output_path):
        try:
            os.remove(output_path)
        except (PermissionError, OSError) as e:
            print(f"Error: cannot overwrite {output_path}: {e}", file=sys.stderr)
            sys.exit(1)

    with open(output_path, "w", encoding="utf-8") as out:
        out.write(header)
        for row in rows:
            if not row.endswith("\n"):
                row += "\n"
            out.write(row)

    print(f"Merged {len(rows)} summary rows into {output_path} (skipped {skipped} empty/error files)")


if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python3 scripts/merge_summaries.py <summaries_dir> <output_path>")
        sys.exit(1)
    merge_summaries(sys.argv[1], sys.argv[2])
