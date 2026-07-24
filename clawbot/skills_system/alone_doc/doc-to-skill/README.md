# Doc to Skill

Doc to Skill turns one or more TXT, Markdown, DOCX, or PDF files into a portable skill ZIP. The result contains a
grounded JSONL knowledge index, normalized source Markdown, extracted assets, generated skill instructions, agent
metadata, and a local retrieval helper.

It is a standalone skill template. Its workflow is implemented by bundled Python scripts and does not depend on Go
business logic or a particular agent framework.

## Processing model

1. Convert all source documents to Markdown and preserve DOCX/PDF assets.
2. Split oversized sections at 5,000 Unicode characters and pack adjacent small sections into bounded chunks.
3. Group textual chunks into bounded model batches, normally up to 15,000 characters and 8 chunks.
4. Let the model write bounded line-based knowledge-point metadata and cite compact source-unit IDs.
5. Validate batch structure and hard limits, restore exact source units as final content, and merge the JSONL index
   deterministically.
6. Package the index, sources, assets, and retrieval script into an installable skill ZIP.

PDF pages without usable text are rendered as one image per page. Image-only DOCX or Markdown sections are handled by
the same preservation-only path instead of being sent to the model as text. Doc to Skill does not call OCR or a vision
model. Extracted DOCX/PDF images are preserved under `references/assets/`; original Markdown image references remain in
the normalized source so the generated skill can inspect the available source files later.

## Why batching reduces Token usage

The iterator returns several chunks in one call instead of one chunk per model call. Adjacent small sections are also
packed before batching. Source evidence units receive IDs such as `u001` and are capped at 1,000 Unicode characters,
preferring original line boundaries. Model output cites at most eight IDs per point instead of copying evidence or
authoring long content, and the merge script restores the exact excerpts locally. Each chunk is limited to twelve points
and each part to 64,000 characters. The line-based format avoids nested JSON escaping failures. Completed batches stay
on disk; their structure and evidence references are validated during resume without being returned to model context
again. Workflow state and chunk files must not be read through model file tools. After invalid output, the iterator
removes the rejected content, creates a small retry placeholder, and returns a fresh pending payload containing the
bounded errors and source units needed for one whole-file rewrite. No separate rejected path is exposed. An accidental
read of the target therefore returns the harmless placeholder instead of failing with a missing-file error. If the issue
list was truncated, the next validation pass exposes any remaining issues.

The final metadata outline is capped independently of indexing. When the index exceeds that cap, outline slots are
allocated in proportion to each document's knowledge-point count and sampled across the document instead of taking only
the first points. When the cap permits, every document with indexed content retains at least one outline item.

For a 22-chunk input, the default limits typically reduce indexing from 22 model turns to a small number of batch turns.
The exact number depends on chunk sizes; the prepared state reports it explicitly.

## Standalone usage

Assume this directory is the current skill directory and `workspace/input` contains staged input files such as
`001-guide.pdf` and `002-notes.md`.

Prepare conversion, chunks, and batch state:

```bash
python3 scripts/prepare_workflow.py \
  --input workspace/input \
  --markdown workspace/markdown \
  --assets workspace/assets \
  --chunks workspace/chunks \
  --state workspace/index-state.json \
  --log workspace/doc.log
```

Request the next batch:

```bash
python3 scripts/batch_index.py --next \
  --manifest workspace/chunks/chunks.jsonl \
  --state workspace/index-state.json \
  --parts workspace/index-parts \
  --log workspace/doc.log
```

Write one complete line-based text part using the returned `part_path` and the contract in
[`references/indexing.md`](references/indexing.md). Use only a whole-file write operation; do not use inline Python,
heredocs, shell redirection, or substring editing. When `retry_after_invalid` is true, overwrite the returned retry
placeholder from the fresh pending payload. Repeat until the iterator returns `status: complete`.

Only the bounded `.txt` batch protocol is accepted; older batch formats are intentionally ignored. Retain the completed
iterator response's `documents`, `outline`, and `outline_truncated` values for skill metadata. Do not read the complete
index for metadata generation, and do not infer unsupported topics from their absence in the bounded outline.

Merge validated batches:

```bash
python3 scripts/merge_index.py \
  --manifest workspace/chunks/chunks.jsonl \
  --state workspace/index-state.json \
  --parts workspace/index-parts \
  --output workspace/doc-index.jsonl \
  --log workspace/doc.log
```

If merge reports an invalid, missing, or retry-placeholder batch part, rerun the batch iterator and rewrite only that
part. For an invalid workflow state, a manifest mismatch, a missing source artifact, or an I/O error, return the
reported error instead of rewriting batch parts. Do not hand-edit converted Markdown, chunk manifests, workflow state,
or the merged index.

Create `workspace/skill-metadata.json` using only the retained completed-iterator values and the contract in
[`references/metadata.md`](references/metadata.md), then build:

```bash
python3 scripts/build_skill.py \
  --index workspace/doc-index.jsonl \
  --metadata workspace/skill-metadata.json \
  --markdown workspace/markdown \
  --assets workspace/assets \
  --zip-out workspace/generate_skill/example-docs.zip \
  --log workspace/doc.log
```

If build validation rejects `workspace/skill-metadata.json`, correct only that file and rerun the build command. Do not
repeat document preparation, batch indexing, or index merging when their validated outputs are already present.

## Input and naming rules

- Supported formats: `.txt`, `.md`, `.docx`, `.pdf`.
- Backend-staged filenames use `001-original.ext`; converted Markdown is `001-original.md`.
- The sequence prefix prevents same-stem, different-extension inputs from colliding.
- Running the converter without unique staged names rejects Markdown filename collisions rather than overwriting data.
- Every source must convert successfully; partial conversion is a task failure.

## Output structure

```text
example-docs/
|-- SKILL.md
|-- agents/openai.yaml
|-- references/doc-index.jsonl
|-- references/markdown/*.md
|-- references/assets/**
`-- scripts/search_index.py
```

The generated skill searches a bounded number of index rows before opening source Markdown or assets, keeping runtime
context bounded while retaining inspectable original material.

## Runtime requirements

- Python 3.10+
- `python-docx`, `pypdf`, and Pillow
- Poppler `pdftoppm`

The Docker image may include OCR tooling for future use, but the current workflow does not invoke it.

## Project structure

```text
.
|-- SKILL.md
|-- README.md
|-- agents/openai.yaml
|-- references/
|   |-- indexing.md
|   `-- metadata.md
`-- scripts/
    |-- prepare_workflow.py
    |-- convert_documents.py
    |-- split_markdown.py
    |-- batch_index.py
    |-- merge_index.py
    |-- build_skill.py
    `-- search_index.py
```
