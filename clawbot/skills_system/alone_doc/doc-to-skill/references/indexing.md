# Batch indexing contract

For each `pending` response, write one complete UTF-8 text part to the exact returned `part_path`. This is a line-based
format, not JSON. It deliberately requires no quote escaping and never includes copied source content.

```text
DOC_TO_SKILL_BATCH_V1
BATCH batch-000001
CHUNK chunk-000001
POINT
TITLE Concise knowledge-point title
SUMMARY Short factual summary
KEYWORD keyword one
KEYWORD keyword two
QUESTION A natural question answered by this point
EVIDENCE u002
EVIDENCE u003
END_POINT
END_CHUNK
END_BATCH
```

Structure and limits:

- Preserve the returned `batch_id` and include every returned `chunk_id` exactly once, in the same order.
- Start with exactly `DOC_TO_SKILL_BATCH_V1`. End every point, chunk, and batch with its matching marker.
- Put `TITLE` and `SUMMARY` first in every point, then zero or more `KEYWORD`, `QUESTION`, and `EVIDENCE` lines.
- Every value occupies one physical line. Ordinary punctuation and quotes are allowed in values.
- Produce 1 to 12 useful knowledge points per textual chunk. Split unrelated facts and avoid near-duplicates.
- Each point requires a title of at most 160 characters, a summary of at most 500 characters, 1 to 12 unique keywords of
  at most 80 characters each, 0 to 8 unique questions of at most 300 characters each, and 1 to 8 unique evidence IDs.
- The complete part must not exceed 64,000 Unicode characters.
- Do not write `content`. The merge script deterministically restores final `content` by joining the exact selected
  evidence units, so model-authored prose cannot replace or expand the source.

Grounding rules:

- An `EVIDENCE` value must be one exact ID from that chunk's returned `allowed_evidence_ids`. Never use an ID from
  another chunk, invent an ID, expand a range, continue a sequence, or exceed the returned `last_evidence_id`.
- Select the smallest set of at most eight units that directly supports the title and summary. Every returned unit is
  already long enough to serve as evidence; do not quote partial units.
- Units preserve source text and are capped at 1,000 Unicode characters, preferring original line boundaries.
- Do not add facts, prerequisites, relationships, procedures, ordering, totals, or limits that require guessing or
  calculation. Preserve material versions, identifiers, filenames, commands, numbers, units, and limit wording.
- Images are preservation-only during generation. Do not use vision, OCR, filenames, link text, or nearby prose to infer
  image content.
- Retrieval questions may use natural wording but must remain answerable from the selected evidence.
- Do not add source metadata. Trusted manifests attach document IDs, source files, heading paths, page ranges, chunk
  IDs, and image references during merge.

Writing and recovery rules:

- Use one whole-file `write_file` call for the returned `part_path`. Do not use `execute`, inline Python, heredocs,
  shell redirection, append operations, substring editing, or any equivalent method to create or modify a batch part.
- Never call `read_file`, `ls`, `edit_file`, or an equivalent operation on a batch part. The iterator is the only
  validator.
- After invalid output, the iterator removes the rejected content, writes a small retry placeholder at the new `.txt`
  `part_path`, and returns another `pending` response with `retry_after_invalid: true`, bounded `errors`, and fresh
  source units. It does not expose a separate rejected file path.
- Overwrite the retry placeholder as one whole file from that current response. Do not inspect or patch it. If an
  accidental read occurs despite this rule, the placeholder exists and only instructs the caller to overwrite it.
- Continue one part at a time until the iterator reports `complete`.

The scripts validate the text structure, hard limits, batch/chunk identity and order, and evidence membership. They
restore exact source evidence and final content locally. Faithful semantic extraction within those structural bounds
remains the model's responsibility.
