---
name: web-to-skill
description: "Convert one public website URL or an explicit batch of public URLs into a reusable skill zip backed by rendered HTML snapshots and a bounded JSONL retrieval index. Use to discover a documentation directory from one URL, crawl a supplied URL set sequentially, generate a source profile, or package indexed web content as a specialized skill."
---

# Web To Skill

Use three deterministic stages. Keep every intermediate file and the final zip under the writable task directory
provided by the system prompt. Run all scripts with Python 3.

## Stage 1: Prepare the URL list

Always run `scripts/prepare_urls.py` first.

- When the user supplies exactly one URL, the script opens that page with Playwright, detects its directory or
  navigation tree, and writes the start page plus all discovered in-scope page URLs.
- When the user supplies two or more URLs, the script skips all network discovery. It only validates, normalizes,
  deduplicates, and writes the supplied URLs.
- The output is always a UTF-8 text file with one URL per line.

```bash
python3 scripts/prepare_urls.py \
  --out <task-dir>/crawl/url-list.txt \
  "https://example.com/docs"
```

```bash
python3 scripts/prepare_urls.py \
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

Run `scripts/crawl_urls.py` with the URL-list file from Stage 1.

```bash
python3 scripts/crawl_urls.py \
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
- Extract keywords with jieba from the title, description, and selected body, then merge them with the page's original
  metadata keywords.
- When at least four pages succeed, remove cross-page high-frequency noise terms unless they occur in the page title or
  description, then keep at most 12 keywords per page.
- Append one successful page object per line to `<task-dir>/crawl/index.jsonl`.
- Append immediate structured progress, retry, timeout-stop, and error events to the same `<task-dir>/crawl/crawl.log`
  created during Stage 1.

Use debug mode only when validating the workflow. It processes at most the first five URLs:

```bash
python3 scripts/crawl_urls.py \
  --url-list <task-dir>/crawl/url-list.txt \
  --out-dir <task-dir>/crawl-debug \
  --debug
```

Do not pass browser, concurrency, wait, retry, depth, link-scope, or timeout options. Those controls are not part of
this workflow.

## Stage 3: Build the specialized skill

Inspect the bounded page metadata in `index.jsonl` and create exactly one UTF-8 JSON metadata file. The model must
provide the skill identity and source profile; the build script does not infer them.

Use this schema:

```json
{
	"name": "example-docs",
	"title": "Example Docs",
	"description": "Answer questions using the indexed Example Docs pages and rendered snapshots.",
	"source_summary": "The indexed pages cover the Example product documentation and operational guidance.",
	"topic_groups": [
		{
			"name": "Product usage",
			"summary": "Setup, configuration, and common product workflows.",
			"query_terms": [
				"product name",
				"configuration"
			],
			"page_urls": [
				"https://example.com/docs/setup"
			]
		}
	],
	"aliases": [
		{
			"canonical": "Example Product",
			"aliases": [
				"Example"
			]
		}
	],
	"coverage_notes": [
		"The index does not cover account billing."
	]
}
```

Requirements:

- `name`, `title`, `description`, `source_summary`, and `topic_groups` are required.
- `name` must contain only lowercase letters, digits, and single hyphens, and must not exceed 50 characters.
- `page_urls` must use exact URLs present in `index.jsonl`.
- Keep profile statements grounded in the indexed page metadata. Do not invent product behavior or coverage.
- `aliases` and `coverage_notes` may be empty arrays.

Build the zip:

```bash
python3 scripts/build_skill.py \
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

1. Confirm `crawl/url-list.txt` contains at least one URL.
2. Confirm `index.jsonl` contains at least one valid JSON object and every `html_path` exists.
3. Confirm the final zip exists and contains `<skill-name>/SKILL.md`, the JSONL index, at least one HTML snapshot, and
   both helper scripts.
4. Output only the final zip path required by the system prompt.
