#!/usr/bin/env python3
"""Build a validated Doc-to-Skill zip."""

from __future__ import annotations

import argparse
import datetime as dt
import json
import re
import shutil
import sys
import tempfile
import zipfile
from pathlib import Path
from typing import Any

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")


NAME_RE = re.compile(r"^[a-z0-9]+(?:-[a-z0-9]+)*$")
REQUIRED_METADATA_FIELDS = ("name", "title", "description", "source_summary", "topic_groups")


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


def require_text(value: Any, label: str, *, max_length: int) -> str:
    if not isinstance(value, str):
        raise ValueError(f"{label} must be a string")
    text = " ".join(value.split()).strip()
    if not text:
        raise ValueError(f"{label} must not be empty")
    if len(text) > max_length:
        raise ValueError(f"{label} must not exceed {max_length} characters")
    return text


def optional_text_list(value: Any, label: str, *, max_items: int, max_length: int) -> list[str]:
    if value is None:
        return []
    if not isinstance(value, list) or len(value) > max_items:
        raise ValueError(f"{label} must be an array with at most {max_items} items")
    output: list[str] = []
    seen: set[str] = set()
    for index, item in enumerate(value):
        text = require_text(item, f"{label}[{index}]", max_length=max_length)
        key = text.casefold()
        if key not in seen:
            seen.add(key)
            output.append(text)
    return output


def resolve_reference_path(root: Path, stored_path: str, prefix: str, label: str) -> Path:
    if not stored_path.startswith(prefix):
        raise ValueError(f"{label} must start with {prefix}: {stored_path}")
    relative = stored_path[len(prefix):]
    if not relative:
        raise ValueError(f"{label} is empty after {prefix}")
    root = root.resolve()
    candidate = (root / Path(relative)).resolve()
    try:
        candidate.relative_to(root)
    except ValueError as exc:
        raise ValueError(f"{label} escapes its source directory: {stored_path}") from exc
    if not candidate.is_file():
        raise ValueError(f"{label} references a missing file: {stored_path}")
    return candidate


EXTRACTED_IMAGE_SOURCE_EXTENSIONS = {"docx", "pdf"}


def validate_image_reference(assets_dir: Path, stored_path: str, source_ext: str, label: str) -> None:
    if source_ext not in EXTRACTED_IMAGE_SOURCE_EXTENSIONS:
        return
    stored_path = stored_path.strip()
    if not stored_path:
        raise ValueError(f"{label} must not be empty")
    if not stored_path.startswith("references/assets/"):
        raise ValueError(f"{label} must reference a packaged extracted image: {stored_path}")
    resolve_reference_path(assets_dir, stored_path, "references/assets/", label)


def load_index(path: Path, markdown_dir: Path, assets_dir: Path) -> tuple[list[dict[str, Any]], set[str]]:
    rows: list[dict[str, Any]] = []
    document_ids: set[str] = set()
    for line_number, raw_line in enumerate(path.read_text(encoding="utf-8-sig").splitlines(), start=1):
        if not raw_line.strip():
            continue
        try:
            row = json.loads(raw_line)
        except json.JSONDecodeError as exc:
            raise ValueError(f"invalid JSONL at {path}:{line_number}: {exc}") from exc
        if not isinstance(row, dict):
            raise ValueError(f"index record at {path}:{line_number} must be an object")
        for field in ("id", "title", "summary", "keywords", "questions", "content", "evidence", "source"):
            if field not in row:
                raise ValueError(f"index record at {path}:{line_number} is missing {field}")
        if not isinstance(row["keywords"], list) or not isinstance(row["questions"], list):
            raise ValueError(f"index record arrays at {path}:{line_number} are invalid")
        if (
            not isinstance(row["evidence"], list)
            or not row["evidence"]
            or any(not isinstance(item, str) or not item.strip() for item in row["evidence"])
        ):
            raise ValueError(f"index evidence at {path}:{line_number} must be a non-empty string array")
        source = row["source"]
        if not isinstance(source, dict):
            raise ValueError(f"index source at {path}:{line_number} must be an object")
        document_id = require_text(source.get("document_id"), f"index source document_id at {path}:{line_number}", max_length=255)
        source_name = require_text(source.get("source_name"), f"index source source_name at {path}:{line_number}", max_length=255)
        source_ext = Path(source_name).suffix.lower().lstrip(".")
        if document_id != source_name:
            raise ValueError(f"index source document_id at {path}:{line_number} must equal source_name")
        source_markdown = str(source.get("source_markdown", ""))
        resolve_reference_path(markdown_dir, source_markdown, "references/markdown/", f"index source_markdown at {path}:{line_number}")
        expected_markdown_name = Path(source_name).with_suffix(".md").name
        if Path(source_markdown).name != expected_markdown_name:
            raise ValueError(
                f"index source_markdown at {path}:{line_number} does not preserve source_name: {source_name}"
            )
        images = source.get("images", [])
        if not isinstance(images, list):
            raise ValueError(f"index source images at {path}:{line_number} must be an array")
        for image in images:
            validate_image_reference(assets_dir, str(image), source_ext, f"index image at {path}:{line_number}")
        document_ids.add(document_id)
        rows.append(row)
    if not rows:
        raise ValueError(f"knowledge index has no records: {path}")
    return rows, document_ids


def load_metadata(path: Path, document_ids: set[str]) -> dict[str, Any]:
    try:
        raw = json.loads(path.read_text(encoding="utf-8-sig"))
    except json.JSONDecodeError as exc:
        raise ValueError(f"invalid metadata JSON at {path}: {exc}") from exc
    if not isinstance(raw, dict):
        raise ValueError("metadata root must be an object")
    missing = [field for field in REQUIRED_METADATA_FIELDS if field not in raw]
    if missing:
        raise ValueError(f"metadata is missing required fields: {', '.join(missing)}")

    name = require_text(raw["name"], "metadata.name", max_length=50)
    if not NAME_RE.fullmatch(name):
        raise ValueError("metadata.name must use lowercase letters, digits, and single hyphens only")
    title = require_text(raw["title"], "metadata.title", max_length=80)
    description = require_text(raw["description"], "metadata.description", max_length=500)
    source_summary = require_text(raw["source_summary"], "metadata.source_summary", max_length=1200)

    groups = raw["topic_groups"]
    if not isinstance(groups, list) or not 1 <= len(groups) <= 12:
        raise ValueError("metadata.topic_groups must contain 1-12 items")
    topic_groups: list[dict[str, Any]] = []
    for index, item in enumerate(groups):
        label = f"metadata.topic_groups[{index}]"
        if not isinstance(item, dict):
            raise ValueError(f"{label} must be an object")
        group_document_ids = optional_text_list(
            item.get("document_ids"), f"{label}.document_ids", max_items=5, max_length=255
        )
        unknown_ids = [document_id for document_id in group_document_ids if document_id not in document_ids]
        if unknown_ids:
            raise ValueError(f"{label}.document_ids contains an ID that is not present in the index: {unknown_ids[0]}")
        topic_groups.append(
            {
                "name": require_text(item.get("name"), f"{label}.name", max_length=80),
                "summary": require_text(item.get("summary"), f"{label}.summary", max_length=300),
                "query_terms": optional_text_list(item.get("query_terms"), f"{label}.query_terms", max_items=10, max_length=80),
                "document_ids": group_document_ids,
            }
        )

    aliases_raw = raw.get("aliases", [])
    if not isinstance(aliases_raw, list) or len(aliases_raw) > 20:
        raise ValueError("metadata.aliases must be an array with at most 20 items")
    aliases: list[dict[str, Any]] = []
    for index, item in enumerate(aliases_raw):
        label = f"metadata.aliases[{index}]"
        if not isinstance(item, dict):
            raise ValueError(f"{label} must be an object")
        aliases.append(
            {
                "canonical": require_text(item.get("canonical"), f"{label}.canonical", max_length=80),
                "aliases": optional_text_list(item.get("aliases"), f"{label}.aliases", max_items=10, max_length=80),
            }
        )

    return {
        "name": name,
        "title": title,
        "description": description,
        "source_summary": source_summary,
        "topic_groups": topic_groups,
        "aliases": aliases,
        "coverage_notes": optional_text_list(raw.get("coverage_notes"), "metadata.coverage_notes", max_items=10, max_length=300),
    }


def yaml_quote(value: str) -> str:
    return json.dumps(value, ensure_ascii=False)


def markdown_text(value: str) -> str:
    return value.replace("\r", " ").replace("\n", " ").strip()


def profile_markdown(metadata: dict[str, Any]) -> str:
    lines = ["## Source profile", "", markdown_text(metadata["source_summary"]), "", "### Topic guide", ""]
    for group in metadata["topic_groups"]:
        lines.append(f"- **{markdown_text(group['name'])}**: {markdown_text(group['summary'])}")
        if group["query_terms"]:
            lines.append(f"  - Query terms: {', '.join(f'`{term}`' for term in group['query_terms'])}")
        if group["document_ids"]:
            lines.append(f"  - Representative documents: {', '.join(f'`{item}`' for item in group['document_ids'])}")
    if metadata["aliases"]:
        lines.extend(["", "### Terminology aliases", ""])
        for item in metadata["aliases"]:
            lines.append(f"- `{item['canonical']}`: {', '.join(f'`{alias}`' for alias in item['aliases'])}")
    if metadata["coverage_notes"]:
        lines.extend(["", "### Coverage boundaries", ""])
        lines.extend(f"- {markdown_text(note)}" for note in metadata["coverage_notes"])
    return "\n".join(lines)


def skill_markdown(metadata: dict[str, Any]) -> str:
    return f'''---
name: {metadata["name"]}
description: {yaml_quote(metadata["description"])}
---

# {markdown_text(metadata["title"])}

Use this skill to answer questions grounded in the indexed source documents.

{profile_markdown(metadata)}

## Retrieval workflow

1. Use the source profile and terminology aliases to choose two to four discriminating query terms from the user's question.
   Separate multiple query terms with a literal `|`, for example `product|feature|error`; do not use spaces as the
   separator.
2. Run the search command, for example `python3 scripts/search_index.py "product|feature|error" --limit 5`.
   When searching the bundled skill content, do not pass `--index`; the script locates
   `references/doc-index.jsonl` from its own location.
3. Read the returned knowledge points, their exact `evidence` excerpts, and their `source` metadata.
4. When more context is needed, open the referenced file under `references/markdown/` and inspect associated images
   under `references/assets/` through file operations.
5. Answer only from claims supported by `evidence` or by the referenced source Markdown. Preserve material warnings,
   commands, tables, links, image references, and relevant details found in the packaged source images.
6. If the indexed documents do not contain the answer, state that the supplied sources do not cover it.
'''


def bounded_short_description(description: str) -> str:
    def is_word_character(character: str) -> bool:
        return (character.isascii() and character.isalnum()) or character in "-'"

    compact = " ".join(description.split())
    if len(compact) < 25:
        compact = f"Use indexed source documents. {compact}".strip()
    if len(compact) <= 64:
        return compact
    raw_candidate = compact[:64]
    candidate = raw_candidate.rstrip()
    next_character = compact[64]
    if (
        candidate
        and is_word_character(raw_candidate[-1])
        and is_word_character(next_character)
    ):
        boundary = len(candidate)
        while boundary > 0:
            character = candidate[boundary - 1]
            if is_word_character(character):
                boundary -= 1
                continue
            break
        whole_words = candidate[:boundary].rstrip(" -")
        candidate = whole_words if len(whole_words) >= 25 else "Use indexed source documents"
    return candidate


def build_openai_yaml(metadata: dict[str, Any]) -> str:
    prompt = f"Use ${metadata['name']} to search the bounded document index before answering from the saved source documents."
    return (
        "interface:\n"
        f"  display_name: {yaml_quote(metadata['title'])}\n"
        f"  short_description: {yaml_quote(bounded_short_description(metadata['description']))}\n"
        f"  default_prompt: {yaml_quote(prompt)}\n"
    )


def zip_directory(source_dir: Path, zip_path: Path) -> None:
    zip_path.parent.mkdir(parents=True, exist_ok=True)
    temporary_zip = zip_path.with_suffix(zip_path.suffix + ".tmp")
    if temporary_zip.exists():
        temporary_zip.unlink()
    try:
        with zipfile.ZipFile(temporary_zip, "w", compression=zipfile.ZIP_DEFLATED) as archive:
            for path in sorted(source_dir.rglob("*")):
                if path.is_file():
                    archive.write(path, path.relative_to(source_dir.parent).as_posix())
        temporary_zip.replace(zip_path)
    finally:
        if temporary_zip.exists():
            temporary_zip.unlink()


def build(args: argparse.Namespace) -> tuple[Path, int]:
    index_path = Path(args.index).resolve()
    metadata_path = Path(args.metadata).resolve()
    markdown_dir = Path(args.markdown).resolve()
    assets_dir = Path(args.assets).resolve()
    zip_path = Path(args.zip_out).resolve()
    if zip_path.suffix.lower() != ".zip":
        raise ValueError("--zip-out must end with .zip")

    rows, document_ids = load_index(index_path, markdown_dir, assets_dir)
    metadata = load_metadata(metadata_path, document_ids)
    markdown_files = sorted(markdown_dir.glob("*.md"))
    if not markdown_files:
        raise ValueError("source Markdown directory is empty")
    search_script = Path(__file__).with_name("search_index.py")
    if not search_script.is_file():
        raise ValueError("search_index.py is missing")

    with tempfile.TemporaryDirectory(prefix="doc-to-skill-") as temporary_root:
        root = Path(temporary_root) / metadata["name"]
        (root / "agents").mkdir(parents=True)
        (root / "references" / "markdown").mkdir(parents=True)
        (root / "references" / "assets").mkdir(parents=True)
        (root / "scripts").mkdir(parents=True)
        (root / "SKILL.md").write_text(skill_markdown(metadata), encoding="utf-8", newline="\n")
        (root / "agents" / "openai.yaml").write_text(build_openai_yaml(metadata), encoding="utf-8", newline="\n")
        shutil.copy2(index_path, root / "references" / "doc-index.jsonl")
        for path in markdown_files:
            shutil.copy2(path, root / "references" / "markdown" / path.name)
        if assets_dir.is_dir():
            for path in assets_dir.rglob("*"):
                if path.is_file():
                    target = root / "references" / "assets" / path.relative_to(assets_dir)
                    target.parent.mkdir(parents=True, exist_ok=True)
                    shutil.copy2(path, target)
        shutil.copy2(search_script, root / "scripts" / "search_index.py")
        zip_directory(root, zip_path)
    return zip_path, len(rows)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Build a reusable skill zip from a grounded document index.")
    parser.add_argument("--index", required=True)
    parser.add_argument("--metadata", required=True)
    parser.add_argument("--markdown", required=True)
    parser.add_argument("--assets", required=True)
    parser.add_argument("--zip-out", required=True)
    parser.add_argument("--log", required=True)
    return parser.parse_args()


def main() -> int:
    global LOGGER
    args = parse_args()
    reported_zip_path = Path(args.zip_out).as_posix()
    LOGGER = WorkflowLogger(args.log)
    LOGGER.emit("build", "start", index=args.index, metadata=args.metadata, zip_out=args.zip_out)
    try:
        zip_path, knowledge_count = build(args)
    except Exception as exc:
        LOGGER.emit("build", "error", error_type=type(exc).__name__, error=str(exc))
        print(str(exc), file=sys.stderr, flush=True)
        return 1
    LOGGER.emit(
        "build",
        "summary",
        skill_name=zip_path.stem,
        knowledge_points=knowledge_count,
        zip_path=reported_zip_path,
    )
    print(json.dumps({
        "status": "complete",
        "stage": "build.summary",
        "action": "return_zip_path",
        "next_action": "return_zip_path_without_listing",
        "skill_name": zip_path.stem,
        "knowledge_points": knowledge_count,
        "zip_path": reported_zip_path,
    }, ensure_ascii=False), flush=True)
    print(reported_zip_path, flush=True)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
