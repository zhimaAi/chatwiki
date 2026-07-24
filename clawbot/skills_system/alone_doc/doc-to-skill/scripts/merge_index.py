#!/usr/bin/env python3
"""Validate model-authored batch parts and merge a grounded JSONL index."""

from __future__ import annotations

import argparse
import datetime as dt
import json
import re
import sys
from pathlib import Path


if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")


class WorkflowLogger:
    def __init__(self, path: str | Path | None) -> None:
        self.path = Path(path) if path else None

    def emit(self, stage: str, event: str, **fields: object) -> None:
        if self.path is None:
            return
        self.path.parent.mkdir(parents=True, exist_ok=True)
        record = {
            "time": dt.datetime.now().astimezone().isoformat(timespec="seconds"),
            "stage": stage,
            "event": event,
            **fields,
        }
        with self.path.open("a", encoding="utf-8", newline="\n") as handle:
            handle.write(json.dumps(record, ensure_ascii=False, default=str, separators=(",", ":")) + "\n")


LOGGER = WorkflowLogger(None)


class ValidationIssues(ValueError):
    """Carry every independently actionable validation issue in one response."""

    def __init__(self, errors: list[str]) -> None:
        self.errors = errors
        super().__init__("; ".join(errors))


MARKDOWN_IMAGE_RE = re.compile(r"!\[[^\]]*\]\([^)]*\)")
IMAGE_ONLY_MARKER = "<!-- source-page-mode: image -->"
EXTRACTED_IMAGE_SOURCE_EXTENSIONS = {"docx", "pdf"}
BATCH_ID_RE = re.compile(r"batch-\d{6}")
SOURCE_UNIT_MAX_CHARS = 1000
BATCH_PART_HEADER = "DOC_TO_SKILL_BATCH_V1"
BATCH_RETRY_HEADER = "DOC_TO_SKILL_RETRY_V1"
MAX_BATCH_PART_CHARS = 64000
MAX_KNOWLEDGE_POINTS_PER_CHUNK = 12
MAX_EVIDENCE_IDS_PER_POINT = 8
MAX_TITLE_CHARS = 160
MAX_SUMMARY_CHARS = 500
MAX_KEYWORDS_PER_POINT = 12
MAX_KEYWORD_CHARS = 80
MAX_QUESTIONS_PER_POINT = 8
MAX_QUESTION_CHARS = 300


def package_image_reference(value, source_ext: str):
    if not isinstance(value, str) or source_ext not in EXTRACTED_IMAGE_SOURCE_EXTENSIONS:
        return value
    prefix = "../assets/"
    if value.startswith(prefix):
        return "references/assets/" + value[len(prefix):]
    return value


def load_jsonl(path: Path) -> list[dict]:
    rows: list[dict] = []
    for line_no, line in enumerate(path.read_text(encoding="utf-8").splitlines(), 1):
        if not line.strip():
            continue
        try:
            value = json.loads(line)
        except json.JSONDecodeError as exc:
            raise ValueError(f"{path}:{line_no}: invalid JSON: {exc}") from exc
        if not isinstance(value, dict):
            raise ValueError(f"{path}:{line_no}: JSON object required")
        rows.append(value)
    return rows


def clean_strings(values: list[str]) -> list[str]:
    output: list[str] = []
    for value in values:
        value = value.strip()
        if value and value not in output:
            output.append(value)
    return output


def split_source_block(value: str, limit: int = SOURCE_UNIT_MAX_CHARS) -> list[str]:
    output: list[str] = []
    current: list[str] = []
    current_chars = 0

    def flush() -> None:
        nonlocal current, current_chars
        if current:
            output.append("\n".join(current).strip())
        current = []
        current_chars = 0

    for line in value.splitlines() or [value]:
        if len(line) > limit:
            flush()
            output.extend(line[start:start + limit] for start in range(0, len(line), limit))
            continue
        separator_chars = 1 if current else 0
        if current and current_chars + separator_chars + len(line) > limit:
            flush()
            separator_chars = 0
        current.append(line)
        current_chars += separator_chars + len(line)
    flush()
    return [item for item in output if item]


def source_units(chunk_content: str) -> list[dict]:
    units: list[dict] = []
    for block in re.split(r"\n\s*\n", chunk_content):
        text = block.strip()
        if not text or re.fullmatch(r"<!--.*?-->", text, flags=re.DOTALL):
            continue
        if MARKDOWN_IMAGE_RE.fullmatch(text):
            continue
        for unit in split_source_block(text):
            if len(re.sub(r"\s+", "", unit)) < 4:
                continue
            units.append({"id": f"u{len(units) + 1:03d}", "text": unit})
    if not units:
        raise ValueError("textual chunk contains no usable evidence units")
    return units


def has_usable_source_text(chunk_content: str) -> bool:
    visible = re.sub(r"<!--.*?-->", "", chunk_content, flags=re.DOTALL)
    visible = MARKDOWN_IMAGE_RE.sub("", visible)
    return len(re.sub(r"[\W_]+", "", visible)) >= 4


def validate_evidence(
    path: Path,
    point_no: int | str,
    evidence: list[str],
    chunk_content: str,
) -> None:
    normalized_chunk = " ".join(chunk_content.split())
    evidence_errors: list[str] = []
    for evidence_no, quote in enumerate(evidence, 1):
        visible_quote = re.sub(r"<!--.*?-->", "", quote, flags=re.DOTALL)
        marker_only = quote.strip() == IMAGE_ONLY_MARKER and IMAGE_ONLY_MARKER in chunk_content
        if not marker_only and len(re.sub(r"\s+", "", visible_quote)) < 4:
            evidence_errors.append(f"{path}: knowledge point {point_no} evidence {evidence_no} is too short")
        if " ".join(quote.split()) not in normalized_chunk:
            evidence_errors.append(
                f"{path}: knowledge point {point_no} evidence {evidence_no} is absent from the chunk"
            )

    if evidence_errors:
        raise ValidationIssues(evidence_errors)


def new_batch_part_path(parts_dir: Path, batch_id: str) -> Path:
    return parts_dir / f"{batch_id}.txt"


def parse_text_batch_part(path: Path, text: str) -> dict:
    lines = [(line_no, line.strip()) for line_no, line in enumerate(text.splitlines(), 1) if line.strip()]
    if not lines or lines[0][1] != BATCH_PART_HEADER:
        raise ValueError(f"{path}: first non-empty line must be {BATCH_PART_HEADER}")
    position = 1

    def take(prefix: str | None = None, exact: str | None = None) -> tuple[int, str]:
        nonlocal position
        if position >= len(lines):
            expected = exact or prefix or "another line"
            raise ValueError(f"{path}: expected {expected} before end of file")
        line_no, value = lines[position]
        position += 1
        if exact is not None and value != exact:
            raise ValueError(f"{path}:{line_no}: expected {exact}")
        if prefix is not None and not value.startswith(prefix):
            raise ValueError(f"{path}:{line_no}: expected {prefix}<value>")
        return line_no, value

    _, batch_line = take(prefix="BATCH ")
    batch_id = batch_line[len("BATCH "):].strip()
    chunks: list[dict] = []
    while position < len(lines) and lines[position][1] != "END_BATCH":
        _, chunk_line = take(prefix="CHUNK ")
        chunk_id = chunk_line[len("CHUNK "):].strip()
        points: list[dict] = []
        while position < len(lines) and lines[position][1] != "END_CHUNK":
            take(exact="POINT")
            _, title_line = take(prefix="TITLE ")
            _, summary_line = take(prefix="SUMMARY ")
            keywords: list[str] = []
            questions: list[str] = []
            evidence_ids: list[str] = []
            while position < len(lines) and lines[position][1] != "END_POINT":
                line_no, value = lines[position]
                position += 1
                if value.startswith("KEYWORD "):
                    keywords.append(value[len("KEYWORD "):].strip())
                elif value.startswith("QUESTION "):
                    questions.append(value[len("QUESTION "):].strip())
                elif value.startswith("EVIDENCE "):
                    evidence_ids.append(value[len("EVIDENCE "):].strip())
                else:
                    raise ValueError(
                        f"{path}:{line_no}: expected KEYWORD, QUESTION, EVIDENCE, or END_POINT"
                    )
            take(exact="END_POINT")
            points.append({
                "title": title_line[len("TITLE "):].strip(),
                "summary": summary_line[len("SUMMARY "):].strip(),
                "keywords": keywords,
                "questions": questions,
                "evidence_ids": evidence_ids,
            })
        take(exact="END_CHUNK")
        chunks.append({"chunk_id": chunk_id, "knowledge_points": points})
    take(exact="END_BATCH")
    if position != len(lines):
        line_no, _ = lines[position]
        raise ValueError(f"{path}:{line_no}: content after END_BATCH is not allowed")
    return {"batch_id": batch_id, "chunks": chunks}


def load_batch_part(path: Path) -> dict:
    text = path.read_text(encoding="utf-8")
    if len(text) > MAX_BATCH_PART_CHARS:
        raise ValueError(
            f"{path}: batch part has {len(text)} characters; maximum is {MAX_BATCH_PART_CHARS}"
        )
    if text.lstrip().startswith(BATCH_RETRY_HEADER):
        raise ValueError(f"{path}: batch part is a retry placeholder and must be overwritten")
    if not text.lstrip().startswith(BATCH_PART_HEADER):
        raise ValueError(f"{path}: batch part must start with {BATCH_PART_HEADER}")
    return parse_text_batch_part(path, text)


def validate_bounded_strings(
    path: Path,
    label: str,
    values: list[str],
    maximum_items: int,
    maximum_chars: int,
    required: bool,
) -> list[str]:
    cleaned = clean_strings(values)
    errors: list[str] = []
    if required and not cleaned:
        errors.append(f"{path}: {label} must contain at least one value")
    if len(cleaned) > maximum_items:
        errors.append(f"{path}: {label} has {len(cleaned)} values; maximum is {maximum_items}")
    for item_no, value in enumerate(cleaned, 1):
        if len(value) > maximum_chars:
            errors.append(
                f"{path}: {label} value {item_no} has {len(value)} characters; maximum is {maximum_chars}"
            )
    if errors:
        raise ValidationIssues(errors)
    return cleaned


def load_validated_batch_part(path: Path, batch: dict, chunks_by_id: dict[str, dict]) -> dict[str, list[dict]]:
    raw = load_batch_part(path)
    batch_id = str(batch.get("batch_id", "")).strip()
    if (
        not isinstance(raw, dict)
        or not isinstance(raw.get("batch_id"), str)
        or raw["batch_id"].strip() != batch_id
    ):
        raise ValueError(f"{path}: batch_id must equal {batch_id}")
    chunk_outputs = raw.get("chunks")
    if not isinstance(chunk_outputs, list):
        raise ValueError(f"{path}: chunks must be an array")
    expected_ids = [str(value) for value in batch.get("chunk_ids", [])]
    actual_ids: list[str] = []
    output: dict[str, list[dict]] = {}
    validation_errors: list[str] = []
    for chunk_no, chunk_output in enumerate(chunk_outputs, 1):
        if not isinstance(chunk_output, dict):
            raise ValueError(f"{path}: chunk result {chunk_no} must be an object")
        raw_chunk_id = chunk_output.get("chunk_id")
        if not isinstance(raw_chunk_id, str):
            raise ValueError(f"{path}: chunk result {chunk_no} chunk_id must be a string")
        chunk_id = raw_chunk_id.strip()
        if chunk_id not in expected_ids or chunk_id in output:
            raise ValueError(f"{path}: unexpected or duplicate chunk_id: {chunk_id}")
        actual_ids.append(chunk_id)
        chunk = chunks_by_id.get(chunk_id)
        if chunk is None:
            raise ValueError(f"{path}: chunk is absent from manifest: {chunk_id}")
        chunk_path = Path(str(chunk.get("chunk_path", "")))
        if not chunk_path.is_file():
            raise ValueError(f"chunk file is missing: {chunk_path}")
        chunk_content = chunk_path.read_text(encoding="utf-8")
        units = {item["id"]: item["text"] for item in source_units(chunk_content)}
        points = chunk_output.get("knowledge_points")
        if not isinstance(points, list) or not points:
            raise ValueError(f"{path}: chunk {chunk_id} has no knowledge points")
        if len(points) > MAX_KNOWLEDGE_POINTS_PER_CHUNK:
            validation_errors.append(
                f"{path}: chunk {chunk_id} has {len(points)} knowledge points; "
                f"maximum is {MAX_KNOWLEDGE_POINTS_PER_CHUNK}"
            )
        validated: list[dict] = []
        for point_no, point in enumerate(points, 1):
            if not isinstance(point, dict):
                raise ValueError(f"{path}: chunk {chunk_id} knowledge point {point_no} must be an object")
            required_text_fields = ("title", "summary")
            text_fields = {field: point.get(field) for field in required_text_fields}
            if any(not isinstance(value, str) for value in text_fields.values()):
                validation_errors.append(
                    f"{path}: chunk {chunk_id} knowledge point {point_no} has invalid text fields"
                )
                continue
            raw_keywords = point.get("keywords")
            raw_questions = point.get("questions", [])
            if (
                not isinstance(raw_keywords, list)
                or any(not isinstance(item, str) for item in raw_keywords)
                or not isinstance(raw_questions, list)
                or any(not isinstance(item, str) for item in raw_questions)
            ):
                validation_errors.append(
                    f"{path}: chunk {chunk_id} knowledge point {point_no} keywords and questions must be string arrays"
                )
                continue
            title = text_fields["title"].strip()
            summary = text_fields["summary"].strip()
            if not title or len(title) > MAX_TITLE_CHARS:
                validation_errors.append(
                    f"{path}: chunk {chunk_id} knowledge point {point_no} title must contain 1-{MAX_TITLE_CHARS} characters"
                )
            if not summary or len(summary) > MAX_SUMMARY_CHARS:
                validation_errors.append(
                    f"{path}: chunk {chunk_id} knowledge point {point_no} summary must contain 1-{MAX_SUMMARY_CHARS} characters"
                )
            try:
                keywords = validate_bounded_strings(
                    path,
                    f"chunk {chunk_id} knowledge point {point_no} keywords",
                    raw_keywords,
                    MAX_KEYWORDS_PER_POINT,
                    MAX_KEYWORD_CHARS,
                    True,
                )
                questions = validate_bounded_strings(
                    path,
                    f"chunk {chunk_id} knowledge point {point_no} questions",
                    raw_questions,
                    MAX_QUESTIONS_PER_POINT,
                    MAX_QUESTION_CHARS,
                    False,
                )
            except ValidationIssues as exc:
                validation_errors.extend(exc.errors)
                continue
            raw_ids = point.get("evidence_ids")
            if (
                not isinstance(raw_ids, list)
                or not raw_ids
                or any(not isinstance(item, str) or not item.strip() for item in raw_ids)
            ):
                validation_errors.append(
                    f"{path}: chunk {chunk_id} knowledge point {point_no} evidence_ids must be a non-empty string array"
                )
                continue
            evidence_ids = clean_strings(raw_ids)
            if len(evidence_ids) > MAX_EVIDENCE_IDS_PER_POINT:
                validation_errors.append(
                    f"{path}: chunk {chunk_id} knowledge point {point_no} has {len(evidence_ids)} evidence IDs; "
                    f"maximum is {MAX_EVIDENCE_IDS_PER_POINT}"
                )
                continue
            unknown_ids = [value for value in evidence_ids if value not in units]
            if unknown_ids:
                validation_errors.append(
                    f"{path}: chunk {chunk_id} knowledge point {point_no} has an unknown evidence_id: {unknown_ids[0]}"
                )
                continue
            evidence = [units[value] for value in evidence_ids]
            content = "\n\n".join(evidence)
            if not title or not summary or not content or not keywords:
                validation_errors.append(
                    f"{path}: chunk {chunk_id} knowledge point {point_no} misses title, summary, content, or keywords"
                )
                continue
            try:
                validate_evidence(
                    path,
                    f"{chunk_id}/{point_no}",
                    evidence,
                    chunk_content,
                )
            except ValidationIssues as exc:
                validation_errors.extend(exc.errors)
                continue
            except ValueError as exc:
                validation_errors.append(str(exc))
                continue
            validated.append({
                "title": title,
                "summary": summary,
                "keywords": keywords,
                "questions": questions,
                "content": content,
                "evidence": evidence,
            })
        output[chunk_id] = validated
    if actual_ids != expected_ids:
        validation_errors.append(f"{path}: chunks must include every expected chunk_id once and in order")
    if validation_errors:
        raise ValidationIssues(validation_errors)
    return output


def image_only_points(chunk: dict) -> list[dict]:
    pages = chunk.get("pages", [])
    page_number = pages[0] if isinstance(pages, list) and pages else None
    document = str(chunk.get("document_id") or chunk.get("source_name") or "document").strip()
    chunk_path = Path(str(chunk.get("chunk_path", "")))
    if not chunk_path.is_file():
        raise ValueError(f"image-only chunk file is missing: {chunk_path}")
    chunk_content = chunk_path.read_text(encoding="utf-8")
    if IMAGE_ONLY_MARKER in chunk_content:
        evidence = [IMAGE_ONLY_MARKER]
    else:
        image_references = [match.group(0) for match in MARKDOWN_IMAGE_RE.finditer(chunk_content)]
        if not image_references:
            raise ValueError(f"image-only chunk contains no image reference: {chunk_path}")
        evidence = [image_references[0]]

    if page_number is not None:
        title = f"{document} - page {page_number} (image only)"
        summary = f"Page {page_number} has no usable extractable text and is preserved as an image."
        keywords = [document, "image-only page", f"page {page_number}"]
        content = (
            f"No OCR or vision analysis was performed for page {page_number}. "
            "Open the associated source image to inspect its original content."
        )
    else:
        title = f"{document} (image only)"
        summary = "This source section has no usable extractable text and is preserved as an image reference."
        keywords = [document, "image-only content"]
        content = (
            "No OCR or vision analysis was performed for this source section. "
            "Open the associated source image or image reference to inspect its original content."
        )
    return [{
        "title": title,
        "summary": summary,
        "keywords": keywords,
        "questions": [],
        "content": content,
        "evidence": evidence,
    }]


def index_rows_for_chunk(chunk: dict, points: list[dict]) -> list[dict]:
    chunk_id = str(chunk.get("chunk_id", "")).strip()
    source_name = str(chunk.get("source_name", "")).strip()
    source_ext = Path(source_name).suffix.lower().lstrip(".")
    rows: list[dict] = []
    for point_no, point in enumerate(points, 1):
        rows.append({
            "id": f"{chunk_id}-kp-{point_no:03d}",
            "title": point["title"],
            "summary": point["summary"],
            "keywords": point["keywords"],
            "questions": point["questions"],
            "content": point["content"],
            "evidence": point["evidence"],
            "source": {
                "document_id": chunk.get("document_id", ""),
                "source_name": source_name,
                "source_markdown": "references/markdown/" + Path(str(chunk.get("source_markdown", ""))).name,
                "section_path": chunk.get("section_path", []),
                "section_paths": chunk.get("section_paths", [chunk.get("section_path", [])]),
                "chunk_id": chunk_id,
                "pages": chunk.get("pages", []),
                "images": [
                    package_image_reference(value, source_ext)
                    for value in chunk.get("images", [])
                ],
            },
        })
    return rows


def merge_batch_parts(manifest: list[dict], state_path: Path, parts_dir: Path) -> list[dict]:
    try:
        state = json.loads(state_path.read_text(encoding="utf-8"))
    except Exception as exc:
        raise ValueError(f"invalid workflow state {state_path}: {exc}") from exc
    batches = state.get("batches") if isinstance(state, dict) else None
    if not isinstance(batches, list):
        raise ValueError(f"workflow state has no batches: {state_path}")
    chunks_by_id = {str(chunk.get("chunk_id", "")): chunk for chunk in manifest}
    points_by_chunk: dict[str, list[dict]] = {}
    expected_ids = [
        str(chunk.get("chunk_id", "")).strip()
        for chunk in manifest
        if not bool(chunk.get("image_only"))
    ]
    state_ids: list[str] = []
    batch_ids: set[str] = set()
    for batch in batches:
        batch_id = str(batch.get("batch_id", "")).strip() if isinstance(batch, dict) else ""
        chunk_ids = batch.get("chunk_ids") if isinstance(batch, dict) else None
        if (
            BATCH_ID_RE.fullmatch(batch_id) is None
            or batch_id in batch_ids
            or not isinstance(chunk_ids, list)
            or not chunk_ids
        ):
            raise ValueError(f"workflow state contains an invalid batch: {state_path}")
        batch_ids.add(batch_id)
        state_ids.extend(str(value).strip() for value in chunk_ids)
        part_path = new_batch_part_path(parts_dir, batch_id)
        if not part_path.is_file():
            raise ValueError(f"missing batch part: {part_path}")
        points_by_chunk.update(load_validated_batch_part(part_path, batch, chunks_by_id))
    if state_ids != expected_ids:
        raise ValueError("workflow state batch chunks do not match the current manifest")
    output: list[dict] = []
    for chunk in manifest:
        chunk_id = str(chunk.get("chunk_id", "")).strip()
        points = image_only_points(chunk) if bool(chunk.get("image_only")) else points_by_chunk.get(chunk_id)
        if not points:
            raise ValueError(f"no validated knowledge points for chunk: {chunk_id}")
        output.extend(index_rows_for_chunk(chunk, points))
    return output


def main() -> int:
    global LOGGER
    parser = argparse.ArgumentParser()
    parser.add_argument("--manifest", required=True, type=Path)
    parser.add_argument("--parts", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--state", required=True, type=Path, help="Workflow state for batch evidence-ID parts.")
    parser.add_argument("--log", required=True, type=Path)
    args = parser.parse_args()
    LOGGER = WorkflowLogger(args.log)
    LOGGER.emit("index", "merge.start", manifest=args.manifest.as_posix(), parts=args.parts.as_posix())
    manifest = load_jsonl(args.manifest)
    if not manifest:
        raise ValueError("chunk manifest is empty")
    output = merge_batch_parts(manifest, args.state, args.parts)
    args.output.parent.mkdir(parents=True, exist_ok=True)
    with args.output.open("w", encoding="utf-8", newline="\n") as handle:
        for row in output:
            handle.write(json.dumps(row, ensure_ascii=False) + "\n")
    LOGGER.emit("index", "summary", chunks=len(manifest), knowledge_points=len(output), output=args.output.as_posix())
    print(json.dumps({"stage": "index.summary", "chunks": len(manifest), "knowledge_points": len(output), "output": args.output.as_posix()}, ensure_ascii=False), flush=True)
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:
        LOGGER.emit("index", "error", error_type=type(exc).__name__, error=str(exc))
        print(f"ERROR: {exc}", file=sys.stderr, flush=True)
        raise SystemExit(1)
