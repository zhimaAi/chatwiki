---
name: web-to-skill
description: "Convert one public website URL or an explicit batch of public URLs into a reusable skill zip backed by rendered HTML snapshots and a bounded JSONL retrieval index. Use to discover a documentation directory from one URL, crawl a supplied URL set sequentially, generate a source profile, or package indexed web content as a specialized skill."
---

# Web To Skill

Use three deterministic stages. Keep every agent-created intermediate file and the final zip under the writable task
directory provided by the system prompt. Run all scripts with Python 3. Do not install or upgrade packages or browser
binaries at runtime.

Let `<skill-dir>` be the skill base directory supplied by the skill loader immediately above these instructions. Use
that directory exactly as supplied. Execute bundled scripts in place with `python3 <skill-dir>/scripts/...`; never copy
them into the task directory or a temporary directory.

Let `<task-dir>` be the writable task directory supplied by the system prompt. Use it exactly as supplied for every
intermediate artifact, model-authored JSON file, and final zip.

## Stage 1: Prepare the URL list

Always run `<skill-dir>/scripts/prepare_urls.py` first.

- When the user supplies exactly one URL, the script opens that page with Playwright, detects its directory or
  navigation tree, and writes the start page plus all discovered in-scope page URLs.
- When the user supplies two or more URLs, the script skips all network discovery. It only validates, normalizes,
  deduplicates, and writes the supplied URLs.
- The output is always a UTF-8 text file with one URL per line.

```bash
python3 <skill-dir>/scripts/prepare_urls.py \
  --out <task-dir>/crawl/url-list.txt \
  "https://example.com/docs"
```

```bash
python3 <skill-dir>/scripts/prepare_urls.py \
  --out <task-dir>/crawl/url-list.txt \
  "https://example.com/page-a" \
  "https://example.com/page-b"
```

Do not manually create the URL-list file and do not run directory discovery for a batch supplied by the user.

Directory navigation uses a fixed 60-second timeout except for Yuque, which uses 120 seconds. A failed navigation is
retried once. Directory scroll passes, convergence, and any truncation reason are appended to
`<task-dir>/crawl/crawl.log`; they are not added to the URL-list file or the retrieval index.

For ChatWiki Docs (`help.chatwiki.com`), use the Docusaurus sitemap rather than the currently visible sidebar. Rebuild
sitemap paths on the supplied origin because the sitemap publishes a placeholder host, and preserve the language
selected by the supplied start URL.

For KanCloud books, read the complete summary tree from `application/payload+json` before using rendered catalog links.
Treat a result containing only the supplied page as a preparation failure; do not continue to build a one-page skill.

## Stage 2: Crawl the prepared URLs

Run `<skill-dir>/scripts/crawl_urls.py` with the URL-list file from Stage 1.

```bash
python3 <skill-dir>/scripts/crawl_urls.py \
  --url-list <task-dir>/crawl/url-list.txt \
  --out-dir <task-dir>/crawl
```

The crawler has intentionally fixed behavior:

- Crawl URLs sequentially with no concurrency.
- Use Playwright and a fixed 60-second page timeout.
- Retry one time after a timeout, browser network error, HTTP 429 or 5xx response, or empty rendered body.
- Stop after four consecutive final timeouts and skip the remaining URLs. A success or a non-timeout failure resets the
  timeout streak.
- Apply built-in body selectors for ChatWiki Docs, Yuque, Feishu, OpenClaw Docs, Alibaba Cloud Help, KanCloud, and
  WeChat Official Account articles; use the rendered page body as the fallback.
- For Feishu, retain the longest stable body snapshot when the final body is empty or shorter.
- Save cleaned rendered HTML under `<task-dir>/crawl/html/`.
- When different prepared URLs redirect to the same final URL, capture and index that page once. Record the other
  prepared URLs as redirect duplicates in crawl coverage; they are not crawl failures.
- Extract keywords with jieba from the title, description, and selected body, then merge them with the page's original
  metadata keywords.
- When at least four pages succeed, remove cross-page high-frequency noise terms unless they occur in the page title or
  description, then keep at most 12 keywords per page.
- Append one successful page object per line to `<task-dir>/crawl/index.jsonl`.
- Append immediate structured progress, retry, timeout-stop, and error events to the same `<task-dir>/crawl/crawl.log`
  created during Stage 1.

Use debug mode only when validating the workflow. It processes at most the first five URLs:

```bash
python3 <skill-dir>/scripts/crawl_urls.py \
  --url-list <task-dir>/crawl/url-list.txt \
  --out-dir <task-dir>/crawl-debug \
  --debug
```

Do not pass browser, concurrency, wait, retry, depth, link-scope, or timeout options. Those controls are not part of
this workflow.

Never edit or delete the generated URL list, JSONL index, crawl log, or rendered HTML snapshots by hand. After crawl
validation succeeds, never rerun URL preparation or crawling because a downstream stage fails. If crawl validation
fails, rerun the crawler at most once without deleting artifacts; if the same error repeats, return that error instead
of restarting the workflow.

Validate the completed crawl with the bundled helper:

```bash
python3 <skill-dir>/scripts/validate_crawl.py \
  --index <task-dir>/crawl/index.jsonl
```

The helper resolves every crawl-index `html_path` relative to the directory containing `index.jsonl`, never relative to
the process working directory. It returns a bounded completion result containing the latest `crawl_urls run.done`
counts, a bounded failure summary, and a bounded redirect-duplicate summary. Do not open, list, print, `cat`, or
otherwise load `crawl.log` through model file tools or ad hoc commands. Use only the helper output when evaluating crawl
completion, failures, or redirect duplicates.

## Stage 3: Build the specialized skill

Run the deterministic metadata-outline helper. Do not open, list, or read `index.jsonl` through model file tools.

```bash
python3 <skill-dir>/scripts/metadata_outline.py \
  --index <task-dir>/crawl/index.jsonl
```

The helper returns at most 60 compact page records. It allocates slots in proportion to each source site's successful
page count, preserves source-site coverage when the limit permits, and samples evenly within each site instead of taking
only its first pages. Use only the returned `outline` and exact URLs to create one UTF-8 JSON metadata file at
`<task-dir>/skill-metadata.json`.

Read [references/metadata.md](references/metadata.md) for the schema and limits. The model must provide the skill
identity and source profile; the build script does not infer them. If `outline_truncated` is true, describe represented
themes without claiming exhaustive index coverage. Do not infer negative coverage from topics missing in the outline.

If build validation rejects the metadata file, correct only `<task-dir>/skill-metadata.json` and rerun the build stage.
Do not repeat URL preparation or crawling when their validated outputs are already present.

The build stage also reads the last `crawl_urls run.done` event from `<task-dir>/crawl/crawl.log`, validates its counts
against `index.jsonl`, and writes a deterministic succeeded, failed, redirect-duplicate, and timeout-skipped coverage
note into the generated Skill. Model-authored `coverage_notes` supplement this deterministic crawl boundary and cannot
replace or hide it.

Use the metadata `name` value as `<skill-name>` and build the zip:

```bash
python3 <skill-dir>/scripts/build_skill.py \
  --index <task-dir>/crawl/index.jsonl \
  --metadata <task-dir>/skill-metadata.json \
  --zip-out <task-dir>/generate_skill/<skill-name>.zip
```

The generated skill contains:

```text
<skill-name>/
|-- SKILL.md
|-- agents/
|   `-- openai.yaml
|-- references/
|   |-- web-index.jsonl
|   `-- html/
|       `-- *.html
`-- scripts/
    |-- search_index.py
    `-- fetch_rendered_html.py
```

`search_index.py` is the bounded local retrieval helper. `fetch_rendered_html.py` fetches one current rendered page only
when the saved snapshots are insufficient or the user explicitly requests current content.

## Validation and completion

Before returning success:

1. Confirm the crawl-validation helper reports `status: complete` with nonzero `url_count` and `index_rows`. Require
   `run_done.requested == url_count`, `run_done.succeeded == index_rows`, `html_paths_checked == index_rows`,
   `failure_summary.count == run_done.failed`, and `duplicate_summary.count == run_done.duplicate_final_urls`.
2. Confirm the metadata-outline helper reports `status: complete`, at least one outline item, and no more than 60.
3. Confirm the build command succeeds. Its deterministic input validation and packaging create the required Skill
   contents at `<task-dir>/generate_skill/<skill-name>.zip`.
4. Immediately return only the final zip path required by the system prompt. Do not call `ls` or reopen the zip, index,
   HTML snapshots, or `crawl.log` after a successful build. Inspect them only to diagnose a reported script error.
