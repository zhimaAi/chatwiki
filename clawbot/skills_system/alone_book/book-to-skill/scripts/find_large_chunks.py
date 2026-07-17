#!/usr/bin/env python3
"""Find chunks larger than a given threshold.

Usage: python3 scripts/find_large_chunks.py chunks/ 16
       (lists chunks/ files > 16 KB, sorted by size descending)
       python3 scripts/find_large_chunks.py chunks/ 0
       (lists all chunks with sizes)
"""

import os
import sys


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 scripts/find_large_chunks.py <chunks_dir> [size_kb_threshold]")
        sys.exit(1)

    chunks_dir = sys.argv[1]
    threshold_kb = int(sys.argv[2]) if len(sys.argv) > 2 else 16

    if not os.path.isdir(chunks_dir):
        print(f"Error: {chunks_dir} not found")
        sys.exit(1)

    files = []
    for name in sorted(os.listdir(chunks_dir)):
        if name.endswith(".md") or name.endswith(".txt"):
            path = os.path.join(chunks_dir, name)
            size = os.path.getsize(path)
            kb = size / 1024.0
            if threshold_kb == 0 or kb > threshold_kb:
                files.append((name, size, kb))

    files.sort(key=lambda x: -x[1])  # sort by size descending

    if not files:
        if threshold_kb > 0:
            print(f"No chunks > {threshold_kb}KB in {chunks_dir}")
        else:
            print(f"No chunks found in {chunks_dir}")
        return

    # table format: chunk | size(KB)
    print(f"{'chunk':<22} {'size(KB)':>9}  (threshold={threshold_kb}KB)")
    print("-" * 42)
    for name, size, kb in files:
        print(f"{name:<22} {kb:>9.1f}")


if __name__ == "__main__":
    main()
