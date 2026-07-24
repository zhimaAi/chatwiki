#!/usr/bin/env python3
"""Run deterministic DocToSkill preparation stages with compact output."""

from __future__ import annotations

import argparse
import json
import subprocess
import sys
from pathlib import Path

from batch_index import load_state


if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")


SCRIPT_DIR = Path(__file__).resolve().parent


def run_stage(script: str, arguments: list[str]) -> None:
    command = [sys.executable, str(SCRIPT_DIR / script), *arguments]
    completed = subprocess.run(
        command,
        check=False,
        capture_output=True,
        text=True,
        encoding="utf-8",
        errors="replace",
    )
    if completed.returncode == 0:
        return
    details = (completed.stderr or completed.stdout or "unknown error").strip()
    if len(details) > 2000:
        details = details[-2000:]
    raise RuntimeError(f"{script} failed: {details}")


def compact_state(state: dict, resumed: bool) -> dict:
    return {
        "status": "prepared",
        "next_action": "run_batch_iterator",
        "read_policy": "Do not read chunk files, chunks.jsonl, or workflow state directly.",
        "documents": state["documents"],
        "chunks": state["chunks"],
        "image_only_chunks": state["image_only_chunks"],
        "batches": len(state["batches"]),
        "resumed": resumed,
    }


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Prepare a standalone DocToSkill workflow.")
    parser.add_argument("--input", required=True, type=Path)
    parser.add_argument("--markdown", required=True, type=Path)
    parser.add_argument("--assets", required=True, type=Path)
    parser.add_argument("--chunks", required=True, type=Path)
    parser.add_argument("--state", required=True, type=Path)
    parser.add_argument("--log", required=True, type=Path)
    parser.add_argument("--max-chars", type=int, default=5000)
    parser.add_argument("--max-batch-chars", type=int, default=15000)
    parser.add_argument("--max-batch-chunks", type=int, default=8)
    args = parser.parse_args()
    if args.max_chars <= 0 or args.max_batch_chars <= 0 or args.max_batch_chunks <= 0:
        parser.error("character and batch limits must be positive")
    return args


def main() -> int:
    args = parse_args()
    manifest = args.chunks / "chunks.jsonl"
    if args.state.is_file():
        if not manifest.is_file():
            raise ValueError(f"workflow state exists but the chunk manifest is missing: {manifest}")
        state = load_state(args.state, manifest)
        print(json.dumps(compact_state(state, resumed=True), ensure_ascii=False, separators=(",", ":")))
        return 0

    run_stage("convert_documents.py", [
        "--input", str(args.input),
        "--output", str(args.markdown),
        "--assets", str(args.assets),
        "--log", str(args.log),
    ])
    run_stage("split_markdown.py", [
        "--input", str(args.markdown),
        "--output", str(args.chunks),
        "--max-chars", str(args.max_chars),
        "--log", str(args.log),
    ])
    run_stage("batch_index.py", [
        "--plan",
        "--manifest", str(manifest),
        "--state", str(args.state),
        "--max-batch-chars", str(args.max_batch_chars),
        "--max-batch-chunks", str(args.max_batch_chunks),
        "--log", str(args.log),
    ])
    state = load_state(args.state, manifest)
    print(json.dumps(compact_state(state, resumed=False), ensure_ascii=False, separators=(",", ":")))
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:
        print(f"ERROR: {exc}", file=sys.stderr, flush=True)
        raise SystemExit(1)
