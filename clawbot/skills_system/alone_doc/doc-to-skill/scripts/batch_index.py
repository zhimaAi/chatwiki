#!/usr/bin/env python3
"""Plan and iterate compact evidence-ID batches for DocToSkill indexing."""

from __future__ import annotations

import argparse
import datetime as dt
import hashlib
import json
import re
import sys
from pathlib import Path

from merge_index import (
    BATCH_PART_HEADER,
    BATCH_RETRY_HEADER,
    MAX_BATCH_PART_CHARS,
    MAX_EVIDENCE_IDS_PER_POINT,
    MAX_KEYWORD_CHARS,
    MAX_KEYWORDS_PER_POINT,
    MAX_KNOWLEDGE_POINTS_PER_CHUNK,
    MAX_QUESTION_CHARS,
    MAX_QUESTIONS_PER_POINT,
    MAX_SUMMARY_CHARS,
    MAX_TITLE_CHARS,
    ValidationIssues,
    image_only_points,
    load_jsonl,
    load_validated_batch_part,
    new_batch_part_path,
    source_units,
)


if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")


STATE_VERSION = 1
BATCH_ID_RE = re.compile(r"batch-\d{6}")
MAX_REPORTED_VALIDATION_ERRORS = 50


class WorkflowLogger:
    def __init__(self, path: Path) -> None:
        self.path = path

    def emit(self, event: str, **fields: object) -> None:
        self.path.parent.mkdir(parents=True, exist_ok=True)
        record = {
            "time": dt.datetime.now().astimezone().isoformat(timespec="seconds"),
            "stage": "index",
            "event": event,
            **fields,
        }
        with self.path.open("a", encoding="utf-8", newline="\n") as handle:
            handle.write(json.dumps(record, ensure_ascii=False, default=str, separators=(",", ":")) + "\n")


def manifest_sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def chunk_character_count(chunk: dict) -> int:
    value = chunk.get("char_count")
    if isinstance(value, int) and value >= 0:
        return value
    chunk_path = Path(str(chunk.get("chunk_path", "")))
    if not chunk_path.is_file():
        raise ValueError(f"chunk file is missing: {chunk_path}")
    return len(chunk_path.read_text(encoding="utf-8"))


def build_batches(manifest: list[dict], max_chars: int, max_chunks: int) -> list[dict]:
    batches: list[dict] = []
    current_ids: list[str] = []
    current_chars = 0

    def flush() -> None:
        nonlocal current_ids, current_chars
        if not current_ids:
            return
        batches.append({
            "batch_id": f"batch-{len(batches) + 1:06d}",
            "chunk_ids": current_ids,
            "char_count": current_chars,
        })
        current_ids = []
        current_chars = 0

    seen: set[str] = set()
    for chunk in manifest:
        chunk_id = str(chunk.get("chunk_id", "")).strip()
        if not chunk_id or chunk_id in seen:
            raise ValueError(f"manifest contains an invalid or duplicate chunk_id: {chunk_id}")
        seen.add(chunk_id)
        if bool(chunk.get("image_only")):
            continue
        char_count = chunk_character_count(chunk)
        if current_ids and (len(current_ids) >= max_chunks or current_chars + char_count > max_chars):
            flush()
        current_ids.append(chunk_id)
        current_chars += char_count
    flush()
    return batches


def load_state(path: Path, manifest_path: Path) -> dict:
    try:
        state = json.loads(path.read_text(encoding="utf-8"))
    except Exception as exc:
        raise ValueError(f"invalid workflow state {path}: {exc}") from exc
    if not isinstance(state, dict) or state.get("version") != STATE_VERSION:
        raise ValueError(f"unsupported workflow state: {path}")
    if state.get("manifest_sha256") != manifest_sha256(manifest_path):
        raise ValueError("workflow state does not match the current chunk manifest")
    manifest = load_jsonl(manifest_path)
    expected_ids = [
        str(chunk.get("chunk_id", "")).strip()
        for chunk in manifest
        if not bool(chunk.get("image_only"))
    ]
    batches = state.get("batches")
    if not isinstance(batches, list):
        raise ValueError("workflow state batches must be an array")
    batch_ids: set[str] = set()
    actual_ids: list[str] = []
    for batch in batches:
        if not isinstance(batch, dict):
            raise ValueError("workflow state contains an invalid batch")
        batch_id = str(batch.get("batch_id", "")).strip()
        chunk_ids = batch.get("chunk_ids")
        if (
            BATCH_ID_RE.fullmatch(batch_id) is None
            or batch_id in batch_ids
            or not isinstance(chunk_ids, list)
            or not chunk_ids
        ):
            raise ValueError("workflow state contains an invalid or duplicate batch")
        batch_ids.add(batch_id)
        actual_ids.extend(str(value).strip() for value in chunk_ids)
    if actual_ids != expected_ids:
        raise ValueError("workflow state batch chunks do not match the current manifest")
    expected_documents = len({str(chunk.get("document_id", "")) for chunk in manifest})
    expected_image_only = sum(1 for chunk in manifest if bool(chunk.get("image_only")))
    if (
        state.get("documents") != expected_documents
        or state.get("chunks") != len(manifest)
        or state.get("image_only_chunks") != expected_image_only
    ):
        raise ValueError("workflow state summary does not match the current manifest")
    return state


def plan_batches(args: argparse.Namespace, logger: WorkflowLogger) -> int:
    manifest = load_jsonl(args.manifest)
    if not manifest:
        raise ValueError("chunk manifest is empty")
    if args.state.is_file():
        state = load_state(args.state, args.manifest)
        resumed = True
    else:
        batches = build_batches(manifest, args.max_batch_chars, args.max_batch_chunks)
        state = {
            "version": STATE_VERSION,
            "manifest": args.manifest.as_posix(),
            "manifest_sha256": manifest_sha256(args.manifest),
            "documents": len({str(chunk.get("document_id", "")) for chunk in manifest}),
            "chunks": len(manifest),
            "image_only_chunks": sum(1 for chunk in manifest if bool(chunk.get("image_only"))),
            "max_batch_chars": args.max_batch_chars,
            "max_batch_chunks": args.max_batch_chunks,
            "batches": batches,
        }
        args.state.parent.mkdir(parents=True, exist_ok=True)
        args.state.write_text(
            json.dumps(state, ensure_ascii=False, indent=2) + "\n",
            encoding="utf-8",
            newline="\n",
        )
        resumed = False
    logger.emit(
        "batch.plan",
        documents=state["documents"],
        chunks=state["chunks"],
        batches=len(state["batches"]),
        resumed=resumed,
    )
    print(json.dumps({
        "status": "prepared",
        "documents": state["documents"],
        "chunks": state["chunks"],
        "image_only_chunks": state["image_only_chunks"],
        "batches": len(state["batches"]),
        "resumed": resumed,
    }, ensure_ascii=False, separators=(",", ":")))
    return 0


def batch_payload(batch: dict, chunks_by_id: dict[str, dict], part_path: Path) -> dict:
    chunks: list[dict] = []
    for chunk_id in batch["chunk_ids"]:
        chunk = chunks_by_id.get(str(chunk_id))
        if chunk is None:
            raise ValueError(f"batch references an unknown chunk: {chunk_id}")
        chunk_path = Path(str(chunk.get("chunk_path", "")))
        if not chunk_path.is_file():
            raise ValueError(f"chunk file is missing: {chunk_path}")
        units = source_units(chunk_path.read_text(encoding="utf-8"))
        chunks.append({
            "chunk_id": chunk_id,
            "document_id": chunk.get("document_id", ""),
            "source_name": chunk.get("source_name", ""),
            "section_paths": chunk.get("section_paths", [chunk.get("section_path", [])]),
            "allowed_evidence_ids": [unit["id"] for unit in units],
            "last_evidence_id": units[-1]["id"],
            "units": units,
        })
    return {
        "status": "pending",
        "action": "write_entire_part",
        "batch_id": batch["batch_id"],
        "part_path": part_path.as_posix(),
        "part_format": BATCH_PART_HEADER,
        "constraints": {
            "max_part_chars": MAX_BATCH_PART_CHARS,
            "max_knowledge_points_per_chunk": MAX_KNOWLEDGE_POINTS_PER_CHUNK,
            "max_evidence_ids_per_point": MAX_EVIDENCE_IDS_PER_POINT,
            "max_title_chars": MAX_TITLE_CHARS,
            "max_summary_chars": MAX_SUMMARY_CHARS,
            "max_keywords_per_point": MAX_KEYWORDS_PER_POINT,
            "max_keyword_chars": MAX_KEYWORD_CHARS,
            "max_questions_per_point": MAX_QUESTIONS_PER_POINT,
            "max_question_chars": MAX_QUESTION_CHARS,
            "content_is_derived_from_evidence": True,
        },
        "chunks": chunks,
    }


def retry_placeholder_text(batch_id: str) -> str:
    return "\n".join([
        BATCH_RETRY_HEADER,
        f"BATCH {batch_id}",
        "ACTION Overwrite this entire file from the iterator pending payload. Do not inspect or edit it.",
        "",
    ])


def write_retry_placeholder(path: Path, batch_id: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(retry_placeholder_text(batch_id), encoding="utf-8", newline="\n")


def is_retry_placeholder(path: Path, batch_id: str) -> bool:
    try:
        lines = path.read_text(encoding="utf-8").splitlines()
    except OSError:
        return False
    return len(lines) >= 2 and lines[0] == BATCH_RETRY_HEADER and lines[1] == f"BATCH {batch_id}"


def pending_retry_payload(
    batch: dict,
    chunks_by_id: dict[str, dict],
    part_path: Path,
    completed_chunks: int,
    total_chunks: int,
) -> dict:
    payload = batch_payload(batch, chunks_by_id, part_path)
    payload.update({
        "retry_after_invalid": True,
        "part_placeholder": True,
        "part_exists": True,
        "read_part": False,
        "completed_chunks": completed_chunks,
        "total_chunks": total_chunks,
    })
    return payload


def public_validation_error(value: str, part_path: Path) -> str:
    output = value
    for path_value in {str(part_path), part_path.as_posix()}:
        if path_value:
            output = output.replace(f"invalid batch part {path_value}", "invalid batch part")
            output = output.replace(path_value, "batch part")
    return output


def proportional_outline(
    items_by_document: dict[str, list[dict]],
    documents: list[str],
    limit: int,
) -> list[dict]:
    """Allocate outline slots proportionally, then sample each document evenly."""
    eligible = [document_id for document_id in documents if items_by_document.get(document_id)]
    total_items = sum(len(items_by_document[document_id]) for document_id in eligible)
    if total_items <= limit:
        return [item for document_id in eligible for item in items_by_document[document_id]]

    quotas: dict[str, int] = {}
    remainders: list[tuple[int, int, str]] = []
    allocated = 0
    for position, document_id in enumerate(eligible):
        numerator = limit * len(items_by_document[document_id])
        quota = numerator // total_items
        quotas[document_id] = quota
        allocated += quota
        remainders.append((numerator % total_items, -position, document_id))
    for _, _, document_id in sorted(remainders, reverse=True)[:limit - allocated]:
        quotas[document_id] += 1

    # Pure proportional allocation can round a small document down to zero.
    # Preserve coverage when the bounded outline has at least one slot per document,
    # taking those slots from the documents with the largest allocations.
    if len(eligible) <= limit:
        positions = {document_id: position for position, document_id in enumerate(eligible)}
        for document_id in eligible:
            if quotas[document_id] > 0:
                continue
            donors = [candidate for candidate in eligible if quotas[candidate] > 1]
            if not donors:
                raise AssertionError("cannot preserve document coverage within the outline limit")
            donor = max(
                donors,
                key=lambda candidate: (
                    quotas[candidate],
                    len(items_by_document[candidate]),
                    -positions[candidate],
                ),
            )
            quotas[donor] -= 1
            quotas[document_id] = 1

    outline: list[dict] = []
    for document_id in eligible:
        items = items_by_document[document_id]
        quota = quotas[document_id]
        if quota >= len(items):
            outline.extend(items)
        elif quota == 1:
            outline.append(items[(len(items) - 1) // 2])
        elif quota > 1:
            indexes = [index * (len(items) - 1) // (quota - 1) for index in range(quota)]
            outline.extend(items[index] for index in indexes)
    if len(outline) != limit:
        raise AssertionError("proportional outline allocation did not fill the requested limit")
    return outline


def iterate_batches(args: argparse.Namespace, logger: WorkflowLogger) -> int:
    manifest = load_jsonl(args.manifest)
    state = load_state(args.state, args.manifest)
    chunks_by_id = {str(chunk.get("chunk_id", "")): chunk for chunk in manifest}
    args.parts.mkdir(parents=True, exist_ok=True)
    completed_chunks = state["image_only_chunks"]
    validated_by_chunk: dict[str, list[dict]] = {}
    for batch in state["batches"]:
        batch_id = str(batch["batch_id"])
        part_path = new_batch_part_path(args.parts, batch_id)
        if not part_path.is_file():
            payload = batch_payload(batch, chunks_by_id, part_path)
            payload.update({"completed_chunks": completed_chunks, "total_chunks": state["chunks"]})
            logger.emit("batch.pending", batch_id=batch_id, completed_chunks=completed_chunks, total_chunks=state["chunks"])
            print(json.dumps(payload, ensure_ascii=False, separators=(",", ":")))
            return 0
        if is_retry_placeholder(part_path, batch_id):
            payload = pending_retry_payload(
                batch,
                chunks_by_id,
                part_path,
                completed_chunks,
                state["chunks"],
            )
            logger.emit(
                "batch.retry_pending",
                batch_id=batch_id,
                completed_chunks=completed_chunks,
                total_chunks=state["chunks"],
            )
            print(json.dumps(payload, ensure_ascii=False, separators=(",", ":")))
            return 0
        try:
            validated = load_validated_batch_part(part_path, batch, chunks_by_id)
        except Exception as exc:
            all_errors = list(exc.errors) if isinstance(exc, ValidationIssues) else [str(exc)]
            if not all_errors:
                all_errors = [str(exc) or type(exc).__name__]
            raw_errors = all_errors[:MAX_REPORTED_VALIDATION_ERRORS]
            errors = [public_validation_error(value, part_path) for value in raw_errors]
            errors_truncated = len(all_errors) > len(errors)
            try:
                part_path.unlink()
            except FileNotFoundError:
                pass
            except OSError as cleanup_exc:
                raise RuntimeError(f"failed to remove invalid batch part {part_path}: {cleanup_exc}") from cleanup_exc
            replacement_path = new_batch_part_path(args.parts, batch_id)
            write_retry_placeholder(replacement_path, batch_id)
            payload = pending_retry_payload(
                batch,
                chunks_by_id,
                replacement_path,
                completed_chunks,
                state["chunks"],
            )
            payload.update({
                "invalid_part_removed": True,
                "error": errors[0],
                "errors": errors,
                "error_count": len(all_errors),
                "errors_truncated": errors_truncated,
                "rewrite_policy": {
                    "whole_file": True,
                    "read_existing_part": False,
                    "substring_edits": False,
                    "write_file_only": True,
                },
            })
            logger.emit(
                "batch.invalid",
                batch_id=batch_id,
                error=raw_errors[0],
                errors=raw_errors,
                error_count=len(all_errors),
                errors_truncated=errors_truncated,
                part_removed=True,
            )
            print(json.dumps(payload, ensure_ascii=False, separators=(",", ":")))
            return 0
        completed_chunks += len(batch["chunk_ids"])
        validated_by_chunk.update(validated)

    documents: list[str] = []
    for chunk in manifest:
        document_id = str(chunk.get("document_id", ""))
        if document_id and document_id not in documents:
            documents.append(document_id)

    items_by_document: dict[str, list[dict]] = {document_id: [] for document_id in documents}
    for chunk in manifest:
        document_id = str(chunk.get("document_id", ""))
        chunk_id = str(chunk.get("chunk_id", ""))
        points = (
            image_only_points(chunk)
            if bool(chunk.get("image_only"))
            else validated_by_chunk.get(chunk_id, [])
        )
        for point in points:
            items_by_document.setdefault(document_id, []).append({
                "document_id": document_id,
                "title": point["title"],
                "summary": point["summary"],
                "keywords": point["keywords"],
            })

    knowledge_point_count = sum(len(items) for items in items_by_document.values())
    outline = proportional_outline(items_by_document, documents, args.max_outline_items)
    logger.emit(
        "batch.complete",
        completed_chunks=completed_chunks,
        total_chunks=state["chunks"],
        knowledge_points=knowledge_point_count,
        outline_items=len(outline),
        outline_documents=len({item["document_id"] for item in outline}),
    )
    print(json.dumps({
        "status": "complete",
        "completed_chunks": completed_chunks,
        "total_chunks": state["chunks"],
        "documents": documents,
        "outline": outline,
        "outline_truncated": knowledge_point_count > len(outline),
    }, ensure_ascii=False, separators=(",", ":")))
    return 0


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Plan or iterate DocToSkill evidence-ID batches.")
    mode = parser.add_mutually_exclusive_group(required=True)
    mode.add_argument("--plan", action="store_true")
    mode.add_argument("--next", action="store_true", dest="emit_next")
    parser.add_argument("--manifest", required=True, type=Path)
    parser.add_argument("--state", required=True, type=Path)
    parser.add_argument("--parts", type=Path)
    parser.add_argument("--max-batch-chars", type=int, default=15000)
    parser.add_argument("--max-batch-chunks", type=int, default=8)
    parser.add_argument("--max-outline-items", type=int, default=60)
    parser.add_argument("--log", required=True, type=Path)
    args = parser.parse_args()
    if args.max_batch_chars <= 0 or args.max_batch_chunks <= 0 or args.max_outline_items <= 0:
        parser.error("batch limits must be positive")
    if args.emit_next and args.parts is None:
        parser.error("--parts is required with --next")
    return args


def main() -> int:
    args = parse_args()
    logger = WorkflowLogger(args.log)
    return plan_batches(args, logger) if args.plan else iterate_batches(args, logger)


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:
        print(f"ERROR: {exc}", file=sys.stderr, flush=True)
        raise SystemExit(1)
