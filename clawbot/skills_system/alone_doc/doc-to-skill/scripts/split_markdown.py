#!/usr/bin/env python3
"""Split Markdown into heading-aware chunks bounded by Unicode character count."""

from __future__ import annotations

import argparse
import datetime as dt
import json
import re
import sys
from pathlib import Path

from merge_index import has_usable_source_text


if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")


HEADING_RE = re.compile(r"^(#{1,6})\s+(.+?)\s*$")
PAGE_RE = re.compile(r"<!--\s*source-page:\s*(\d+)\s*-->")
SCAN_PAGE_RE = re.compile(r"<!--\s*source-page-mode:\s*image\s*-->")
IMAGE_RE = re.compile(r"!\[[^\]]*\]\(([^)]+)\)")
SEPARATORS = ("\n\n", "\n", "。", "！", "？", ". ", "! ", "? ", "; ", "；", "，", ", ", " ")


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


def load_jsonl(path: Path) -> list[dict]:
    if not path.is_file():
        raise ValueError(f"document manifest is missing: {path}")
    rows: list[dict] = []
    for line_number, raw_line in enumerate(path.read_text(encoding="utf-8-sig").splitlines(), 1):
        if not raw_line.strip():
            continue
        try:
            row = json.loads(raw_line)
        except json.JSONDecodeError as exc:
            raise ValueError(f"invalid JSONL at {path}:{line_number}: {exc}") from exc
        if not isinstance(row, dict):
            raise ValueError(f"document manifest record at {path}:{line_number} must be an object")
        rows.append(row)
    if not rows:
        raise ValueError(f"document manifest is empty: {path}")
    return rows


def split_exact(text: str, limit: int) -> list[str]:
    return [text[start:start + limit] for start in range(0, len(text), limit)]


def recursive_split(text: str, limit: int, separators: tuple[str, ...] = SEPARATORS) -> list[str]:
    text = text.strip()
    if not text:
        return []
    if len(text) <= limit:
        return [text]
    if not separators:
        return split_exact(text, limit)
    separator = separators[0]
    pieces = text.split(separator)
    if len(pieces) == 1:
        return recursive_split(text, limit, separators[1:])
    output: list[str] = []
    current = ""
    for index, piece in enumerate(pieces):
        candidate_piece = piece + (separator if index < len(pieces) - 1 else "")
        if len(candidate_piece) > limit:
            if current.strip():
                output.append(current.strip())
                current = ""
            output.extend(recursive_split(candidate_piece, limit, separators[1:]))
            continue
        if current and len(current) + len(candidate_piece) > limit:
            output.append(current.strip())
            current = candidate_piece
        else:
            current += candidate_piece
    if current.strip():
        output.append(current.strip())
    return output


def parse_sections(text: str) -> list[tuple[list[str], str, list[int]]]:
    hierarchy: list[tuple[int, str]] = []
    sections: list[tuple[list[str], list[str], list[int]]] = []
    current_path: list[str] = []
    current_lines: list[str] = []
    current_pages: list[int] = []
    current_page: int | None = None
    empty_heading_pending = False

    def flush_current() -> None:
        nonlocal current_lines, current_pages, empty_heading_pending
        body = "\n".join(current_lines).strip()
        if PAGE_RE.sub("", body).strip():
            sections.append((list(current_path), list(current_lines), list(current_pages)))
        elif empty_heading_pending and any(IMAGE_RE.search(value) for value in current_path):
            sections.append((list(current_path), [], list(current_pages)))
        current_lines = []
        current_pages = []
        empty_heading_pending = False

    for line in text.splitlines():
        page_match = PAGE_RE.search(line)
        if page_match:
            previous_was_scan = any(SCAN_PAGE_RE.search(value) for value in current_lines)
            flush_current()
            # The generated scan heading must not become the heading path of an
            # unheaded textual page that follows it.
            if previous_was_scan:
                hierarchy = []
                current_path = []
            current_page = int(page_match.group(1))
        match = HEADING_RE.match(line)
        if match:
            flush_current()
            level = len(match.group(1))
            title = match.group(2).strip()
            hierarchy = [item for item in hierarchy if item[0] < level]
            hierarchy.append((level, title))
            current_path = [item[1] for item in hierarchy]
            current_lines = []
            current_pages = [current_page] if current_page is not None else []
            empty_heading_pending = True
        else:
            current_lines.append(line)
            if current_page is not None and current_page not in current_pages:
                current_pages.append(current_page)
    flush_current()
    return [(path, "\n".join(lines).strip(), pages) for path, lines, pages in sections]


def heading_prefix(path: list[str]) -> str:
    return "\n\n".join(f"{'#' * min(index, 6)} {title}" for index, title in enumerate(path, 1))


def page_range(content: str) -> list[int]:
    values = sorted({int(value) for value in PAGE_RE.findall(content)})
    if not values:
        return []
    return [values[0], values[-1]]


def common_section_path(paths: list[list[str]]) -> list[str]:
    if not paths:
        return []
    common = list(paths[0])
    for path in paths[1:]:
        length = 0
        while length < len(common) and length < len(path) and common[length] == path[length]:
            length += 1
        common = common[:length]
        if not common:
            break
    return common


def build_section_parts(markdown_path: Path, max_chars: int) -> list[dict]:
    text = markdown_path.read_text(encoding="utf-8")
    output: list[dict] = []
    for section_path, body, section_pages in parse_sections(text):
        prefix = heading_prefix(section_path)
        available = max_chars - len(prefix) - (2 if prefix and body else 0) - 1
        if available <= 0:
            raise ValueError(f"heading path exceeds max chars in {markdown_path.name}: {' > '.join(section_path)}")
        parts = recursive_split(body, available) or [""]
        for part_no, part in enumerate(parts, 1):
            content = "\n\n".join(value for value in (prefix, part) if value).strip() + "\n"
            if len(content) > max_chars:
                raise AssertionError(f"section part exceeds {max_chars} characters")
            pages = [min(section_pages), max(section_pages)] if section_pages else page_range(content)
            images = IMAGE_RE.findall(content)
            output.append({
                "content": content,
                "section_path": section_path,
                "section_part": part_no,
                "section_part_total": len(parts),
                "pages": pages,
                "images": images,
                "image_only": bool(SCAN_PAGE_RE.search(content)) or (
                    bool(images) and not has_usable_source_text(content)
                ),
            })
    return output


def pack_section_parts(parts: list[dict], max_chars: int) -> list[dict]:
    packed: list[dict] = []
    current: list[dict] = []

    def flush() -> None:
        nonlocal current
        if not current:
            return
        content = "\n\n".join(str(item["content"]).strip() for item in current).strip() + "\n"
        if len(content) > max_chars:
            raise AssertionError(f"packed chunk exceeds {max_chars} characters")
        page_values = [page for item in current for page in item["pages"]]
        images: list[str] = []
        for item in current:
            for image in item["images"]:
                if image not in images:
                    images.append(image)
        section_paths = [list(item["section_path"]) for item in current]
        packed.append({
            "content": content,
            "section_path": common_section_path(section_paths),
            "section_paths": section_paths,
            "part": current[0]["section_part"] if len(current) == 1 else 1,
            "part_total": current[0]["section_part_total"] if len(current) == 1 else 1,
            "pages": [min(page_values), max(page_values)] if page_values else [],
            "images": images,
            "image_only": False,
        })
        current = []

    for item in parts:
        if item["image_only"]:
            flush()
            packed.append({
                "content": item["content"],
                "section_path": list(item["section_path"]),
                "section_paths": [list(item["section_path"])],
                "part": item["section_part"],
                "part_total": item["section_part_total"],
                "pages": list(item["pages"]),
                "images": list(item["images"]),
                "image_only": True,
            })
            continue
        candidate = "\n\n".join(
            [*(str(value["content"]).strip() for value in current), str(item["content"]).strip()]
        ).strip() + "\n"
        if current and len(candidate) > max_chars:
            flush()
        current.append(item)
    flush()
    return packed


def split_document(
    markdown_path: Path,
    output_dir: Path,
    max_chars: int,
    sequence_start: int,
    document_id: str,
    source_name: str,
) -> tuple[list[dict], int]:
    records: list[dict] = []
    sequence = sequence_start
    chunks = pack_section_parts(build_section_parts(markdown_path, max_chars), max_chars)
    for chunk in chunks:
        content = str(chunk["content"])
        sequence += 1
        chunk_id = f"chunk-{sequence:06d}"
        chunk_path = output_dir / f"{chunk_id}.md"
        chunk_path.write_text(content, encoding="utf-8", newline="\n")
        records.append({
            "chunk_id": chunk_id,
            "document_id": document_id,
            "source_name": source_name,
            "source_markdown": markdown_path.as_posix(),
            "section_path": chunk["section_path"],
            "section_paths": chunk["section_paths"],
            "part": chunk["part"],
            "part_total": chunk["part_total"],
            "chunk_path": chunk_path.as_posix(),
            "char_count": len(content),
            "pages": chunk["pages"],
            "images": chunk["images"],
            "image_only": chunk["image_only"],
        })
    return records, sequence


def main() -> int:
    global LOGGER
    parser = argparse.ArgumentParser()
    parser.add_argument("--input", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--max-chars", type=int, default=5000)
    parser.add_argument("--log", required=True, type=Path)
    args = parser.parse_args()
    LOGGER = WorkflowLogger(args.log)
    LOGGER.emit("split", "start", input=args.input.as_posix(), output=args.output.as_posix(), max_chars=args.max_chars)
    if args.max_chars <= 0:
        message = "max chars must be positive"
        LOGGER.emit("split", "error", error_type="ArgumentError", error=message)
        parser.error(message)
    files = sorted(item for item in args.input.glob("*.md") if item.name != "documents.jsonl")
    if not files:
        raise ValueError("no converted Markdown documents found")
    document_manifest = load_jsonl(args.input / "documents.jsonl")
    documents_by_markdown: dict[Path, dict] = {}
    input_dir = args.input.resolve()
    for record in document_manifest:
        stored_markdown_path = str(record.get("markdown_path", "")).strip()
        document_id = str(record.get("document_id", "")).strip()
        source_name = str(record.get("source_name", "")).strip()
        if not stored_markdown_path or not document_id or not source_name:
            raise ValueError("document manifest record is missing markdown_path, document_id, or source_name")
        markdown_path = Path(stored_markdown_path).resolve()
        try:
            markdown_path.relative_to(input_dir)
        except ValueError as exc:
            raise ValueError(f"document manifest markdown_path is outside the input directory: {stored_markdown_path}") from exc
        if markdown_path in documents_by_markdown:
            raise ValueError(f"duplicate markdown_path in document manifest: {markdown_path}")
        documents_by_markdown[markdown_path] = record
    files_by_path = {path.resolve(): path for path in files}
    unlisted = [path.name for resolved, path in files_by_path.items() if resolved not in documents_by_markdown]
    if unlisted:
        raise ValueError("converted Markdown is absent from document manifest: " + ", ".join(unlisted))
    missing = [
        Path(str(record["markdown_path"])).name
        for resolved, record in documents_by_markdown.items()
        if resolved not in files_by_path
    ]
    if missing:
        raise ValueError("document manifest references missing Markdown: " + ", ".join(missing))
    args.output.mkdir(parents=True, exist_ok=True)
    records: list[dict] = []
    sequence = 0
    for path in files:
        document = documents_by_markdown[path.resolve()]
        document_records, sequence = split_document(
            path,
            args.output,
            args.max_chars,
            sequence,
            str(document["document_id"]),
            str(document["source_name"]),
        )
        records.extend(document_records)
        LOGGER.emit("split", "document.done", document=path.name, chunks=len(document_records))
        print(json.dumps({"stage": "split", "document": path.name, "chunks": len(document_records)}, ensure_ascii=False), flush=True)
    manifest = args.output / "chunks.jsonl"
    with manifest.open("w", encoding="utf-8", newline="\n") as handle:
        for record in records:
            handle.write(json.dumps(record, ensure_ascii=False) + "\n")
    if not records:
        raise ValueError("document split produced no chunks")
    LOGGER.emit("split", "summary", documents=len(files), chunks=len(records), manifest=manifest.as_posix())
    print(json.dumps({"stage": "split.summary", "chunks": len(records), "manifest": manifest.as_posix()}, ensure_ascii=False), flush=True)
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:
        LOGGER.emit("split", "error", error_type=type(exc).__name__, error=str(exc))
        print(f"ERROR: {exc}", file=sys.stderr, flush=True)
        raise SystemExit(1)
