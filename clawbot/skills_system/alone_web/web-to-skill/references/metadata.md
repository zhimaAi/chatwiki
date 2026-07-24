# Skill metadata contract

Create one UTF-8 JSON object from the completed metadata outline:

```json
{
	"name": "example-docs",
	"title": "Example Docs",
	"description": "Answer questions using the indexed Example Docs pages and rendered snapshots.",
	"source_summary": "The sampled page metadata covers Example product setup and operating guidance.",
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
	"coverage_notes": []
}
```

Validation limits:

- `name`, `title`, `description`, `source_summary`, and `topic_groups` are required.
- `name`: lowercase letters, digits, and single hyphens only; maximum 50 characters.
- `title`: maximum 80 characters. `description`: maximum 500. `source_summary`: maximum 1,200.
- `topic_groups`: 1-12 objects. Each requires `name` (maximum 80) and `summary` (maximum 300), with at most 10
  `query_terms` (maximum 80 each) and 3 `page_urls` (maximum 2,048 each).
- Every `page_urls` value must exactly match a URL in the completed outline. Do not reopen `index.jsonl` to find more
  representative URLs.
- `aliases`: at most 20 objects. Each requires `canonical` (maximum 80) and at most 10 aliases (maximum 80 each).
- `coverage_notes`: at most 10 values, maximum 300 characters each.
- Ground every field in the completed outline. When `outline_truncated` is true, describe only represented themes and do
  not claim that the outline exhaustively covers the index.
- Do not infer a missing feature, product boundary, or unsupported topic merely because it is absent from the outline.
  Add a coverage note only when the sampled page metadata states that boundary explicitly; otherwise use an empty array.
- Crawl success, failure, and timeout coverage is attached deterministically by the build script. Do not duplicate or
  reinterpret those counts in model-authored metadata.
