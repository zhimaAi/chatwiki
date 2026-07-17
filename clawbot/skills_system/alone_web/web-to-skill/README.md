# Web To Skill

This mother skill turns one public website URL or an explicit URL batch into a reusable skill backed by rendered HTML
snapshots and a JSONL retrieval index.

## Workflow

Run all commands with Python 3 from this directory.

### 1. Prepare URLs

```bash
python3 scripts/prepare_urls.py \
  --out <task-dir>/crawl/url-list.txt \
  "https://example.com/docs"
```

One supplied URL triggers directory discovery. Two or more supplied URLs skip discovery and are only normalized and
deduplicated. The result is always a UTF-8 text file with one URL per line.

Directory discovery uses a 60-second navigation timeout, except Yuque, which uses 120 seconds. It retries one failed
navigation once. Directory scroll passes, convergence, and truncation reasons are appended to
`<task-dir>/crawl/crawl.log`.

ChatWiki Docs (`help.chatwiki.com`) uses its Docusaurus sitemap instead of the currently visible sidebar. The sitemap's
placeholder host is replaced with `help.chatwiki.com`, and only the language selected by the supplied start URL is
included.

KanCloud books use the complete summary tree embedded in `application/payload+json`. If that known directory cannot
produce more than the supplied page, URL preparation fails instead of silently building a one-page skill.

### 2. Crawl URLs

```bash
python3 scripts/crawl_urls.py \
  --url-list <task-dir>/crawl/url-list.txt \
  --out-dir <task-dir>/crawl
```

The crawler processes URLs sequentially with a fixed 60-second timeout and one retry for transient failures, including
HTTP 429 and 5xx responses. After four consecutive final timeouts, it records the stop condition and skips the remaining
URLs. Feishu pages retain the longest stable body snapshot when the final body becomes empty or shorter.

Outputs:

- `crawl/url-list.txt`: normalized crawl inputs.
- `crawl/html/*.html`: cleaned rendered snapshots.
- `crawl/index.jsonl`: one successful page record per line.
- `crawl/crawl.log`: directory preparation, page progress, retries, failures, and stop conditions.

Use `--debug` to process only the first five URLs:

```bash
python3 scripts/crawl_urls.py \
  --url-list <task-dir>/crawl/url-list.txt \
  --out-dir <task-dir>/crawl-debug \
  --debug
```

### 3. Build a skill

Create a metadata JSON file with model-provided `name`, `title`, `description`, `source_summary`, `topic_groups`,
optional `aliases`, and optional `coverage_notes`. Then run:

```bash
python3 scripts/build_skill.py \
  --index <task-dir>/crawl/index.jsonl \
  --metadata <task-dir>/skill-metadata.json \
  --zip-out <task-dir>/generate_skill/<skill-name>.zip
```

The generated skill keeps `scripts/search_index.py` for bounded index retrieval and `scripts/fetch_rendered_html.py` for
fetching the latest rendered version of one page.

See `SKILL.md` for the full metadata schema, grounding requirements, generated layout, and validation checklist.
