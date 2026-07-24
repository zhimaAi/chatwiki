# Skill metadata contract

Create one UTF-8 JSON object:

```json
{
	"name": "example-docs",
	"title": "Example Documents",
	"description": "Answer questions using the indexed Example documents.",
	"source_summary": "The source set covers product setup and operating procedures.",
	"topic_groups": [
		{
			"name": "Product usage",
			"summary": "Setup, configuration, and common workflows.",
			"query_terms": [
				"product name",
				"configuration"
			],
			"document_ids": [
				"001-example.pdf"
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
  `query_terms` (maximum 80 each) and 5 `document_ids` (maximum 255 each).
- Every `document_ids` value must exactly match a document ID returned by the completed batch iterator.
- `aliases`: at most 20 objects. Each requires `canonical` (maximum 80) and at most 10 aliases (maximum 80 each).
- `coverage_notes`: at most 10 values, maximum 300 characters each.
- The completed iterator proportionally samples its bounded outline across documents and includes deterministic
  preservation-only descriptions for image content that was not analyzed.
- Keep every field grounded in the completed iterator's document list and outline. Omit uncertain aliases or coverage
  claims rather than inferring them.
- Add a coverage note only when the returned outline explicitly states that boundary. Missing topics, versions, or
  examples do not prove that the source set excludes them; use an empty array when no explicit boundary is available.
