---
name: doc-to-skill
description: "Convert one or more TXT, Markdown, DOCX, or PDF documents into a reusable skill zip backed by normalized Markdown, extracted images, and a grounded JSONL knowledge index. Invoke this skill before inspecting task files, then execute its workflow directly without listing directories."
---

# Doc To Skill

Convert every supplied document into one portable, indexed skill. Keep all generated files under the writable task
directory supplied by the system prompt. Do not install packages, call OCR or vision models, or depend on application
framework code.

Let `<skill-dir>` be the skill directory supplied by the skill loader. Use it exactly as supplied. Run bundled scripts
in place with `python3 <skill-dir>/scripts/...`; never copy scripts, guess an absolute path, prepend another directory,
or add `cd`.

Let `<task-dir>` be the writable task directory supplied by the system prompt. Inputs are in `<task-dir>/input`. Store
all workflow artifacts and the final zip below `<task-dir>`. Every command appends diagnostics to
`<task-dir>/doc.log`.

Proceed directly with the commands below. Do not inventory the skill or task directory, pre-create output directories,
or read bundled scripts. The scripts create their own directories and report the bounded data needed for each step.

## 1. Prepare documents and batches

Run once. The command is resumable after the workflow state exists.

```bash
python3 <skill-dir>/scripts/prepare_workflow.py \
  --input <task-dir>/input \
  --markdown <task-dir>/markdown \
  --assets <task-dir>/assets \
  --chunks <task-dir>/chunks \
  --state <task-dir>/index-state.json \
  --log <task-dir>/doc.log
```

The converter accepts `.txt`, `.md`, `.docx`, and `.pdf`. It preserves DOCX/PDF images, extracts usable PDF text, and
renders PDF pages without usable text as one image per page. It never performs OCR or image interpretation.

After `status: prepared`, immediately run the batch iterator in step 2. Never open, list, or read `chunks/`,
`chunks.jsonl`, `index-state.json`, or individual chunk files through file tools. They are private workflow artifacts;
the iterator is the only source of model-visible document text.

Source files staged as `001-original.ext` retain that exact value as their document ID. Their normalized Markdown is
`001-original.md`; no script adds a `doc-` prefix. A source conversion failure fails the whole task.

Chunks preserve heading context and contain at most 5,000 Unicode characters. Each Chinese character, English letter,
digit, punctuation mark, whitespace, and newline counts as one character. Adjacent small sections are packed together;
oversized sections are split by semantic boundaries before exact character boundaries. Sections that contain images but
no usable evidence text remain standalone and are excluded from model batches, including scanned PDF pages and
image-only DOCX or Markdown content.

## 2. Extract grounded knowledge in batches

Run the iterator:

```bash
python3 <skill-dir>/scripts/batch_index.py --next \
  --manifest <task-dir>/chunks/chunks.jsonl \
  --state <task-dir>/index-state.json \
  --parts <task-dir>/index-parts \
  --log <task-dir>/doc.log
```

When `status` is `pending`, read [references/indexing.md](references/indexing.md), generate the complete line-based text
part, write the whole returned `part_path` once with `write_file`, and rerun the same command. A batch contains multiple
chunks and source units identified as `u001`, `u002`, and so on. Cite only the allowed IDs; do not copy source text or
write model-authored `content`. The merge script restores final content exactly from the selected units.

When `retry_after_invalid` is true, the same `pending` response includes fresh source units and bounded validation
errors. Its `part_path` contains only a small retry placeholder so accidental reads do not fail. Do not read or edit the
placeholder; overwrite the whole file once with `write_file` using the current response. Never use `execute`, inline
Python, heredocs, shell redirection, append operations, `read_file`, `ls`, `edit_file`, or an equivalent operation to
create, inspect, or alter a batch part. Continue until `status` is `complete`. Never retain multiple pending batches
before writing the current part.

Image-only pages are added deterministically during merge. Do not infer their contents from filenames, links, or
surrounding text. Images are preserved for the generated skill to inspect later through file operations; this generation
workflow does not analyze them.

## 3. Merge the validated index

After the iterator reports `complete`, run:

```bash
python3 <skill-dir>/scripts/merge_index.py \
  --manifest <task-dir>/chunks/chunks.jsonl \
  --state <task-dir>/index-state.json \
  --parts <task-dir>/index-parts \
  --output <task-dir>/doc-index.jsonl \
  --log <task-dir>/doc.log
```

If merge reports an invalid, missing, or retry-placeholder batch part, return to step 2 and rewrite only that batch
part. For an invalid workflow state, a manifest mismatch, a missing source artifact, or an I/O error, return the
reported error instead of rewriting batch parts. Do not hand-edit converted Markdown, chunk manifests, workflow state,
or the merged index.

## 4. Create skill metadata

Use the iterator's final `documents`, bounded `outline`, and `outline_truncated` values to create
`<task-dir>/skill-metadata.json`. Read [references/metadata.md](references/metadata.md) for the schema and limits. The
outline is sampled proportionally by each document's knowledge-point count, preserves document coverage when the limit
permits, and includes deterministic descriptions for image-only content. When `outline_truncated` is true, describe
represented themes without claiming exhaustive coverage. Keep all descriptions, topics, aliases, and coverage notes
grounded in the final outline and document IDs. Add a coverage note only for a boundary explicitly stated by the
outline; absence of a topic or version is not evidence that the documents exclude it.

Do not reopen all source chunks or load the complete index into model context. If the bounded outline is insufficient
for a metadata statement, omit that statement instead of guessing.

## 5. Build the skill

Use metadata `name` as `<skill-name>`:

```bash
python3 <skill-dir>/scripts/build_skill.py \
  --index <task-dir>/doc-index.jsonl \
  --metadata <task-dir>/skill-metadata.json \
  --markdown <task-dir>/markdown \
  --assets <task-dir>/assets \
  --zip-out <task-dir>/generate_skill/<skill-name>.zip \
  --log <task-dir>/doc.log
```

If build validation rejects `<task-dir>/skill-metadata.json`, correct only that file and rerun the build stage. Do not
repeat document preparation, batch indexing, or index merging when their validated outputs are already present.

When `--zip-out` is relative, the successful command result preserves that relative path instead of resolving it to an
absolute path. Return the reported `zip_path` unchanged.

The generated archive contains:

```text
<skill-name>/
|-- SKILL.md
|-- agents/openai.yaml
|-- references/doc-index.jsonl
|-- references/markdown/*.md
|-- references/assets/** (when extracted assets exist)
`-- scripts/search_index.py
```

The generated skill reads its bounded JSONL index and may inspect packaged Markdown or images through file operations.
It has no dependency on DocToSkill's host application.

## Completion checks

Before success, confirm from command results:

1. The prepare command reports all expected documents and at least one chunk.
2. The iterator reports `complete`, with every text part and hard limit validated, every selected evidence ID and final
   content restored from its own chunk, and image-only chunks counted automatically.
3. Merge produces a non-empty `<task-dir>/doc-index.jsonl`.
4. Build reports `status: complete`, `action: return_zip_path`, and the final `zip_path`; its deterministic validation
   guarantees the zip contains `SKILL.md`, `agents/openai.yaml`, the index, source Markdown, retrieval script, and every
   referenced extracted asset.
5. Return only the reported `zip_path` immediately. Do not call `ls` or reopen the zip, index, source chunks, or
   `doc.log` after a successful build. Inspect them only to diagnose a reported script error.
