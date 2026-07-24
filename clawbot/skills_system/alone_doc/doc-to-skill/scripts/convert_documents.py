#!/usr/bin/env python3
"""Convert TXT, Markdown, DOCX, and PDF inputs to Markdown with local image assets."""

from __future__ import annotations

import argparse
import datetime as dt
import hashlib
import json
import os
import re
import shutil
import statistics
import subprocess
import sys
from pathlib import Path


if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")


SUPPORTED = {".txt", ".md", ".docx", ".pdf"}
PDF_MIN_TEXT_CHARS = 20
PDF_RENDER_DPI = 180


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


def safe_name(value: str) -> str:
    value = re.sub(r"[\\/:*?\"<>|\x00-\x1f]", "-", value).strip(" .-")
    return value or "document"


def read_text(path: Path) -> str:
    data = path.read_bytes()
    for encoding in ("utf-8-sig", "utf-8", "gb18030"):
        try:
            return data.decode(encoding).replace("\r\n", "\n").replace("\r", "\n")
        except UnicodeDecodeError:
            continue
    raise ValueError(f"cannot decode text file: {path.name}")


def unique_asset_path(asset_dir: Path, stem: str, extension: str, payload: bytes) -> Path:
    digest = hashlib.sha256(payload).hexdigest()[:12]
    extension = re.sub(r"[^a-zA-Z0-9]", "", extension.lower()) or "bin"
    return asset_dir / f"{safe_name(stem)}-{digest}.{extension}"


def docx_heading_level(paragraph) -> int | None:
    style = (getattr(getattr(paragraph, "style", None), "name", "") or "").lower().replace(" ", "")
    match = re.match(r"heading([1-6])$", style)
    if match:
        return int(match.group(1))
    p_pr = paragraph._element.pPr
    if p_pr is not None and p_pr.outlineLvl is not None:
        return min(int(p_pr.outlineLvl.val) + 1, 6)
    return None


def docx_paragraph_markdown(paragraph, asset_dir: Path, markdown_path: Path, image_no: list[int]) -> str:
    from docx.oxml.ns import qn

    pieces: list[str] = []
    contents = paragraph.iter_inner_content() if hasattr(paragraph, "iter_inner_content") else paragraph.runs
    for run in contents:
        if hasattr(run, "url"):
            text = (run.text or "").strip()
            url = (run.url or "").strip()
            if text:
                pieces.append(f"[{text}]({url})" if url else text)
            continue
        text = run.text or ""
        if text:
            if run.bold and run.italic:
                text = f"***{text}***"
            elif run.bold:
                text = f"**{text}**"
            elif run.italic:
                text = f"*{text}*"
            pieces.append(text)
        for node in run._element.iter():
            if node.tag != qn("a:blip"):
                continue
            rel_id = node.get(qn("r:embed"))
            if not rel_id or rel_id not in paragraph.part.related_parts:
                continue
            part = paragraph.part.related_parts[rel_id]
            payload = part.blob
            ext = Path(str(part.partname)).suffix.lstrip(".") or "png"
            image_no[0] += 1
            target = unique_asset_path(asset_dir, f"image-{image_no[0]:04d}", ext, payload)
            target.parent.mkdir(parents=True, exist_ok=True)
            target.write_bytes(payload)
            relative = Path("..") / "assets" / asset_dir.name / target.name
            pieces.append(f"![{paragraph.text.strip() or target.stem}]({relative.as_posix()})")
    text = "".join(pieces).strip()
    if not text:
        return ""
    level = docx_heading_level(paragraph)
    if level:
        return f"{'#' * level} {text}"
    p_pr = paragraph._element.pPr
    if p_pr is not None and p_pr.numPr is not None:
        indent = "  " * int(p_pr.numPr.ilvl.val if p_pr.numPr.ilvl is not None else 0)
        return f"{indent}- {text}"
    return text


def docx_table_markdown(table, asset_dir: Path, markdown_path: Path, image_no: list[int]) -> str:
    rows: list[list[str]] = []
    for row in table.rows:
        values: list[str] = []
        for cell in row.cells:
            paragraphs = [docx_paragraph_markdown(item, asset_dir, markdown_path, image_no) for item in cell.paragraphs]
            values.append("<br>".join(value.replace("\n", "<br>") for value in paragraphs if value))
        rows.append(values)
    if not rows:
        return ""
    width = max(len(row) for row in rows)
    rows = [row + [""] * (width - len(row)) for row in rows]
    output = ["| " + " | ".join(rows[0]) + " |", "| " + " | ".join(["---"] * width) + " |"]
    output.extend("| " + " | ".join(row) + " |" for row in rows[1:])
    return "\n".join(output)


def convert_docx(path: Path, markdown_path: Path, asset_dir: Path) -> str:
    try:
        from docx import Document
        from docx.table import Table
        from docx.text.paragraph import Paragraph
    except ModuleNotFoundError as exc:
        raise RuntimeError("python-docx is required for DOCX conversion") from exc

    document = Document(path)
    blocks: list[str] = []
    image_no = [0]
    for child in document.element.body.iterchildren():
        if child.tag.endswith("}p"):
            value = docx_paragraph_markdown(Paragraph(child, document), asset_dir, markdown_path, image_no)
        elif child.tag.endswith("}tbl"):
            value = docx_table_markdown(Table(child, document), asset_dir, markdown_path, image_no)
        else:
            continue
        if value:
            blocks.append(value)
    return "\n\n".join(blocks).strip() + "\n"


def compact_pdf_text(value: str) -> str:
    return re.sub(r"[ \t]+", " ", value.replace("\r", "")).strip()


def pdf_heading_level(text: str, font_size: float, body_size: float) -> int | None:
    if not text or len(text) > 120:
        return None
    numeric = re.match(r"^(\d+(?:\.\d+)*)\s*([、.．])\s*\S+", text)
    if numeric:
        base_level = numeric.group(1).count(".") + 1
        if font_size >= body_size * 1.35 or (font_size <= 0 and numeric.group(2) == "、"):
            return min(base_level, 6)
        if font_size > 0 and numeric.group(2) == "、":
            return min(base_level + 1, 6)
        return None
    if font_size <= 0:
        return None
    if font_size >= body_size * 1.6:
        return 1
    if font_size >= body_size * 1.35:
        return 2
    if font_size >= body_size * 1.15:
        return 3
    return None


def pdf_text_markdown(page) -> tuple[str, int]:
    fragments: list[tuple[str, float]] = []

    def visitor(text, _cm, tm, _font, font_size):
        value = compact_pdf_text(str(text))
        if value:
            fragments.append((value, float(font_size or 0)))

    page.extract_text(visitor_text=visitor)
    fallback = page.extract_text(extraction_mode="layout") or ""
    text_count = len(re.sub(r"\s+", "", fallback))
    if text_count == 0:
        return "", 0
    sizes = [item[1] for item in fragments if item[1] > 0]
    body_size = statistics.median(sizes) if sizes else 10.0
    output: list[str] = []
    previous_blank = True
    for raw_line in fallback.splitlines():
        text = compact_pdf_text(raw_line)
        if not text:
            if output and not previous_blank:
                output.append("")
            previous_blank = True
            continue
        normalized = compact_pdf_text(text)
        numeric_line = bool(re.match(r"^\d+(?:\.\d+)*\s*[、.．]\s*\S+", normalized))
        matching_sizes = [
            size
            for fragment, size in fragments
            if fragment == normalized or (numeric_line and len(fragment) >= 4 and fragment in normalized)
        ]
        level = pdf_heading_level(text, max(matching_sizes, default=0), body_size)
        if level:
            text = f"{'#' * level} {text}"
        output.append(text)
        previous_blank = False
    value = "\n".join(output).strip()
    return value, text_count


def render_pdf_page(path: Path, page_number: int, asset_dir: Path) -> Path:
    executable = os.environ.get("PDFTOPPM_BIN") or shutil.which("pdftoppm")
    if not executable:
        raise RuntimeError("pdftoppm is required to render scanned PDF pages")
    asset_dir.mkdir(parents=True, exist_ok=True)
    prefix = asset_dir / f".__render-page-{page_number:04d}"
    rendered = Path(str(prefix) + ".png")
    rendered.unlink(missing_ok=True)
    command = [
        executable,
        "-f", str(page_number),
        "-l", str(page_number),
        "-singlefile",
        "-r", str(PDF_RENDER_DPI),
        "-png",
        str(path),
        str(prefix),
    ]
    result = subprocess.run(command, capture_output=True, text=True, timeout=120, check=False)
    if result.returncode != 0 or not rendered.is_file():
        rendered.unlink(missing_ok=True)
        detail = (result.stderr or result.stdout or "unknown pdftoppm error").strip()
        raise RuntimeError(f"cannot render PDF page {page_number}: {detail}")
    payload = rendered.read_bytes()
    rendered.unlink(missing_ok=True)
    target = unique_asset_path(asset_dir, f"page-{page_number}-scan", "png", payload)
    target.write_bytes(payload)
    return target


def pdf_page_render_markdown(path: Path, page_number: int, asset_dir: Path) -> str:
    target = render_pdf_page(path, page_number, asset_dir)
    relative = Path("..") / "assets" / asset_dir.name / target.name
    return f"![PDF page {page_number} render]({relative.as_posix()})"


def extract_pdf_images(page) -> list[tuple[str, bytes]]:
    images: list[tuple[str, bytes]] = []
    for image in page.images:
        payload = image.data
        extension = Path(image.name).suffix.lstrip(".") or "bin"
        if payload:
            images.append((extension, payload))
    return images


def convert_pdf(path: Path, markdown_path: Path, asset_dir: Path) -> str:
    try:
        from pypdf import PdfReader
    except ModuleNotFoundError as exc:
        raise RuntimeError("pypdf is required for PDF conversion") from exc

    document = PdfReader(str(path))
    output: list[str] = []
    image_no = 0
    for page_index, page in enumerate(document.pages):
        page_number = page_index + 1
        page_blocks: list[str] = [f"<!-- source-page: {page_number} -->"]
        page_text, text_count = pdf_text_markdown(page)
        if text_count == 0:
            page_blocks.extend([
                f"# PDF page {page_number} (image)",
                "<!-- source-page-mode: image -->",
                pdf_page_render_markdown(path, page_number, asset_dir),
            ])
            output.append("\n\n".join(page_blocks))
            continue
        if page_text:
            page_blocks.append(page_text)
        if text_count < PDF_MIN_TEXT_CHARS:
            page_blocks.append(pdf_page_render_markdown(path, page_number, asset_dir))
            output.append("\n\n".join(page_blocks))
            continue
        try:
            images = extract_pdf_images(page)
        except Exception as exc:
            print(f"WARN: cannot enumerate images on page {page_number}: {exc}", file=sys.stderr, flush=True)
            images = []
        for extension, payload in images:
            image_no += 1
            target = unique_asset_path(asset_dir, f"page-{page_number}-image-{image_no:04d}", extension, payload)
            target.parent.mkdir(parents=True, exist_ok=True)
            target.write_bytes(payload)
            relative = Path("..") / "assets" / asset_dir.name / target.name
            page_blocks.append(f"![Page {page_number} image {image_no}]({relative.as_posix()})")
        output.append("\n\n".join(page_blocks))
    if not output:
        raise ValueError(f"PDF contains no pages: {path.name}")
    return "\n\n".join(output).strip() + "\n"


def convert_one(path: Path, output_dir: Path, assets_root: Path) -> dict:
    document_id = path.name
    markdown_path = output_dir / path.with_suffix(".md").name
    asset_dir = assets_root / path.stem
    asset_dir.mkdir(parents=True, exist_ok=True)
    if path.suffix.lower() in {".txt", ".md"}:
        content = read_text(path)
    elif path.suffix.lower() == ".docx":
        content = convert_docx(path, markdown_path, asset_dir)
    elif path.suffix.lower() == ".pdf":
        content = convert_pdf(path, markdown_path, asset_dir)
    else:
        raise ValueError(f"unsupported file type: {path.name}")
    markdown_path.write_text(content, encoding="utf-8", newline="\n")
    asset_count = sum(1 for item in asset_dir.rglob("*") if item.is_file())
    if asset_count == 0:
        shutil.rmtree(asset_dir)
    return {
        "document_id": document_id,
        "source_name": path.name,
        "source_ext": path.suffix.lower().lstrip("."),
        "markdown_path": markdown_path.as_posix(),
        "char_count": len(content),
        "asset_count": asset_count,
    }


def main() -> int:
    global LOGGER
    parser = argparse.ArgumentParser()
    parser.add_argument("--input", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    parser.add_argument("--assets", required=True, type=Path)
    parser.add_argument("--log", required=True, type=Path)
    args = parser.parse_args()
    LOGGER = WorkflowLogger(args.log)
    LOGGER.emit("convert", "start", input=args.input.as_posix(), output=args.output.as_posix(), assets=args.assets.as_posix())
    if not args.input.is_dir():
        message = f"input directory does not exist: {args.input}"
        LOGGER.emit("convert", "error", error_type="ArgumentError", error=message)
        parser.error(message)
    files = sorted(item for item in args.input.iterdir() if item.is_file() and not item.name.startswith("."))
    unsupported = [item.name for item in files if item.suffix.lower() not in SUPPORTED]
    if unsupported:
        raise ValueError("unsupported input files: " + ", ".join(unsupported))
    if not files:
        raise ValueError("no input documents found")
    converted_names: dict[str, str] = {}
    for path in files:
        markdown_name = path.with_suffix(".md").name
        key = markdown_name.casefold()
        if key in converted_names:
            raise ValueError(
                f"source files map to the same Markdown filename: {converted_names[key]}, {path.name} -> {markdown_name}"
            )
        converted_names[key] = path.name
    args.output.mkdir(parents=True, exist_ok=True)
    args.assets.mkdir(parents=True, exist_ok=True)
    manifest = args.output / "documents.jsonl"
    records: list[dict] = []
    errors: list[str] = []
    for path in files:
        try:
            record = convert_one(path, args.output, args.assets)
            records.append(record)
            LOGGER.emit("convert", "document.succeed", **record)
            print(json.dumps({"stage": "convert", "status": "succeed", **record}, ensure_ascii=False), flush=True)
        except Exception as exc:
            errors.append(f"{path.name}: {exc}")
            LOGGER.emit("convert", "document.failed", source_name=path.name, error_type=type(exc).__name__, error=str(exc))
            print(json.dumps({"stage": "convert", "status": "failed", "source_name": path.name, "error": str(exc)}, ensure_ascii=False), flush=True)
    with manifest.open("w", encoding="utf-8", newline="\n") as handle:
        for record in records:
            handle.write(json.dumps(record, ensure_ascii=False) + "\n")
    if errors:
        raise RuntimeError("; ".join(errors))
    LOGGER.emit("convert", "summary", documents=len(records), manifest=manifest.as_posix())
    print(json.dumps({"stage": "convert.summary", "documents": len(records), "manifest": manifest.as_posix()}, ensure_ascii=False), flush=True)
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:
        LOGGER.emit("convert", "error", error_type=type(exc).__name__, error=str(exc))
        print(f"ERROR: {exc}", file=sys.stderr, flush=True)
        raise SystemExit(1)
