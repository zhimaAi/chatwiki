---
name: book-to-skill
description: Meta-skill that converts uploaded documents into reusable skills. Single agent, six steps, Read template → Smart split → Recursive check → Per-file process → Generate package.
---

# Book to Skill (Single Agent Six Steps A→E)

You are a single agent. Execute in order: A→A'→B→C→D→D'→E. Go side batches scheduling with automatic context reset after each batch. No resume support.

## Command Execution Constraints (MUST follow)
- **Allowed command names** only: `cat` `find` `grep` `head` `jq` `ls` `node` `npm` `pwd` `python3` `rg` `tail` `wc`
- Python and Node direct code or module execution is allowed, including `python3 -c`, `python3 -m`, `node -e`, and `node -p`
- Use relative paths for commands (e.g. `python3 scripts/split_chapters.py input/xxx chunks/`), no need to cd to working directory
- Forbidden shell control characters: `>` `<` `|` `;` `&` (except `&&`). Scripts use command-line args for output paths, no redirection
- Forbidden: `pip install` (all dependencies pre-installed in container)
- Forbidden: writing helper scripts to check/install Python packages (e.g. `check_deps.py`). If a built-in script reports `ModuleNotFoundError`, report the error and skip the file — do NOT attempt to fix environment issues.
- Forbidden: `touch`, `rm`, `chmod` and other filesystem commands
- If a command is blocked, check for forbidden syntax above, fix and retry

## Path Conventions (MUST follow)
- `write_file` / `read_file`: paths relative to workspace root, e.g. `clawbot/working_dir/<key>/<id>/chunks/chunk_001.md`
- `execute`: commands run in the **current working directory**, command args must use paths relative to current directory
- Correct: `execute python3 scripts/split_chapters.py input/ chunks/`
- Wrong: `execute python3 clawbot/working_dir/.../scripts/xxx.py` ← path would be duplicated!
- When writing Python scripts, all paths inside the script must also be relative to current working directory
- Correct: `os.makedirs('skill/tskill/resources/第10章')`
- Wrong: `os.makedirs('clawbot/working_dir/.../skill/...')` ← would create in wrong location!

## Directory Structure

```
input/           ← Original uploaded files
input_md/        ← Step A' converted Markdown output (*.md)
chunks/          ← Step B split output (chunk_001.md, chunk_002.md...)
templates/       ← Output templates (md_extraction.md, defines summary_index format and SKILL.md structure)
scripts/         ← Conversion + check + merge scripts (convert_to_md.py, find_large_chunks.py, merge_summaries.py)
summaries/       ← Step D writes (one .txt per chunk, single table row)
skill/           ← Step D writes resources/ dir (chapter original text) → Step E writes SKILL.md
```

## Built-in Scripts (execute directly, no need to read)

| Script | Purpose | Usage |
|------|------|------|
| `scripts/convert_to_md.py` | Convert txt/md/docx to Markdown | `python3 scripts/convert_to_md.py input/` → `input_md/*.md` |
| `scripts/detect_patterns.py` | One-click scan all common section symbols, output stats + samples | `python3 scripts/detect_patterns.py input_md/` → shows match count/size distribution/samples per pattern |
| `scripts/find_large_chunks.py` | List chunks exceeding specified size (Step C) | `python3 scripts/find_large_chunks.py chunks/ 16` (>16KB descending) |
| `scripts/merge_summaries.py` | Merge all single-line files in summaries/ into a table | `python3 scripts/merge_summaries.py summaries/ summary_index.txt` |
| `scripts/split_chapters.py` | Split input_md/*.md by chapter headings → chunks/ | `python3 scripts/split_chapters.py input_md/ chunks/ --pattern '<regex>'` |
| `scripts/split_by_size.py` | Fallback split: no-paragraph docs split by 10KB fixed size | `python3 scripts/split_by_size.py --input input_md/<file>.md --output chunks/` |
| `scripts/extract_lines.py` | Extract chunk lines by line number → resources/ (fine-grained/fallback) | `python3 scripts/extract_lines.py <src> <start_line> <end_line> --dst <target_path>` |
| `scripts/write_resource.py` | Write chunk content to resources/ (copy/merge/split) | See Step D |

> ⚠️ Execute all scripts directly, **do not read their content**. Run `detect_patterns.py` as the first step in Step B. Use `--pattern` with `split_chapters.py` to specify the chapter heading regex.

## Step A: Read Template

1. Read `templates/md_extraction.md` to get output format definitions
2. Clarify: summary_index.txt table has four columns (chapter path | title | content summary | source file path), SKILL.md structure, split level rules
3. Constraint: Only read the template once in this step, then proceed to Step A'

## Step A': Document Format Conversion

1. `ls input/` to view all uploaded files
2. If there are txt/md/docx files under `input/`, execute:
   `python3 scripts/convert_to_md.py input/`
   → Convert each file to Markdown, output to `input_md/<original_name>.md`
3. All subsequent Step B sampling/splitting is based on `input_md/*.md`
4. ⚠️ Only run the conversion script once, do not repeat
5. Proceed to Step B when done

## Step B: Smart Split

1. `ls input_md/` to view converted md files
2. **First step** execute `python3 scripts/detect_patterns.py input_md/` → one-click scan all common section symbols
   - Auto-detection: Markdown #/##/###/####, Chinese 第X章/节/篇/部, numeric numbering, Chinese ordinals, etc.
   - Output: total matches per symbol type, chunk size distribution (min/avg/max), top 5 sample headings
   - Directly determine the best `--pattern` from the report, no need for per-pattern grep trials
3. Split level selection (strictly highest level first) & fallback judgment:
   **First check if fallback is needed**: If detect_patterns.py report satisfies any of the following, the document has no effective section structure, **skip steps 3-4 and go directly to fallback**:
   - All patterns have avg chunk size > 20KB (can't split, continuous text)
   - All patterns have avg chunk size < 3KB (too fragmented, meaningless)
   - All patterns have total matches ≤ 3 (very few separators overall)
   
   **Fallback flow**: For each .md file under `input_md/`, execute:
   ```
   python3 scripts/split_by_size.py --input input_md/<filename>.md --output chunks/
   ```
   Each chunk ~10KB, breaks at paragraph boundaries. After execution, jump directly to Step C.
   
   **Normal flow** (when effective section structure exists):
   - Markdown H1 matches > 0 → use `^# .+` to split (even if only 1 #)
   - No H1 but H2 matches > 0 → use `^## .+`
   - No H1/H2 but H3 matches > 0 → use `^### .+`
   - No Markdown headings at all → pick pattern with most matches and avg chunk in 5-16KB range
   - ⚠️ Never match multiple heading levels simultaneously (e.g. matching both # and ##)
4. Execute the split script with the selected pattern:
   ```
   python3 scripts/split_chapters.py input_md/ chunks/ --pattern '<chapter regex>'
   ```
   Common pattern examples (must match the full heading line):
   - Markdown H1: `--pattern '^# .+'`
   - Markdown H2: `--pattern '^## .+'`
   - Markdown H3: `--pattern '^### .+'`
   - Chinese "第X章": `--pattern '^第[一二三四五六七八九十百千\d]+[章节篇].*'`
   - Numeric numbering: `--pattern '^\d+[\.\、\s].+'`
   - Chinese ordinals: `--pattern '^[一二三四五六七八九十]+[、．].+'`
5. Produces chunks/chunk_001.md, chunk_002.md... (zero-padded 3-digit naming, 5-16KB)
6. ⚠️ No need to read chunks/ before splitting is complete, go directly to Step C verification
7. ⚠️ Sub-headings (##/###) splitting will be handled by Step D. Step B only does coarse split by top-level headings
8. Naturally proceed to Step C when done

## Step C: Recursive Check

1. Execute `python3 scripts/find_large_chunks.py chunks/ 16` to list chunks > 16KB
2. If no results → Step C complete, proceed to Step D
3. For each oversized chunk, recursively re-split (max 2 levels): write `scripts/re_split.py` (note: do NOT overwrite built-in `split_chapters.py`) → execute `python3 scripts/re_split.py <chunk_file>`
4. Constraint: Only re-split oversized files, do not read chunk contents one by one

## Step D: Per-File Processing (batch map-reduce)

Go side schedules by batch, 50 chunks per batch. Your tasks:

⚠️ **Context control (hard rule)**: After 3 consecutive chunk reads without writing, you MUST execute write_resource to persist and write summary. Never read 3+ chunks consecutively without any write operation. Violation causes context overflow.

1. Read chunk files one by one (in this batch's chunk list order)
2. Write complete chunk original text to resources/ by chapter structure:
   ⚠️ **Never use `write_file` to output chunk original text! Must use `execute write_resource.py`.**
   - Read chunk → identify chapter boundaries → select subcommand → execute script
   - Three subcommands and their use cases:
     · `copy`: single chunk = single resource (entire chapter in one chunk, no split or merge needed)
       `python3 scripts/write_resource.py copy --src chunks/chunk_010.md --dst skill/<skill_name>/resources/<chapter>/<filename>.md`
     · `merge`: multiple consecutive chunks belong to same chapter → merge into one resource file
       `python3 scripts/write_resource.py merge --src chunks/chunk_010.md chunks/chunk_011.md --dst skill/<skill_name>/resources/<chapter>/<filename>.md`
       ⚠️ Max 3 consecutive chunks (including current), force write if still incomplete after exceeding
     · `split`: single chunk contains multiple `##` or `###` sub-headings → split into multiple resource files
       `python3 scripts/write_resource.py split --src chunks/chunk_005.md --pattern '^## .+' --dst-dir skill/<skill_name>/resources/<chapter>/`
       ⚠️ Chapter name comes from the `#` heading text in the file (script extracts automatically), never use directory path name as chapter name
   - Script auto-creates parent directories, no manual mkdir needed
   - ⚠️ **Write original text directly, no summarizing, no abbreviating, no rewriting**. Each resource file = the complete original text of the corresponding chapter from the chunk
   - When there are clear sub-sections: `<section_title>.md` (e.g. `1.1 Overview.md`)
   - When no clear sub-sections: `<chapter_name>.md` (e.g. `第27章 配灵药.md`), one file for the entire chapter, do not force split
   
   **Fine-grained/fallback extraction** (`extract_lines.py`): When the chunk has no `##`/`###` headings matching the split command, or when finer granularity extraction is needed:
   ```
   python3 scripts/extract_lines.py <chunk_file> <start_line> <end_line> --dst <target_path>
   ```
   - AI reads the chunk first, determines start/end line numbers for each logical segment, then extracts segment by segment into resources/
   - If a heading needs to be added to extracted content, use `--prepend` arg: `--prepend '## Installation Steps\n\n'`
   - Use cases: fallback-split documents without headings, normal flow scenarios needing fine-grained line-range extraction
   - Script auto-creates parent directories, line numbers start from 1
3. Then write `summaries/<chunk_filename>.txt` (single table row, e.g. chunk_001.md → summaries/chunk_001.txt)
   - ⚠️ summaries/ and resources/ are two independent outputs: summaries store short summaries for indexing, resources store complete original text for Q&A
   - ⚠️ Summary still uses `write_file` (short content, no token issues)
   - Format (one line per file, no header):
   ```
   | resources/<chapter>/<filename>.md | <section title> | <50-char summary> | chunks/chunk_NNN.md |
   ```
   - ⚠️ Each chunk's summary file is written only once, do not rewrite
4. ⚠️ Pure-title chunks (only "第X章" with no body): skip — don't write resource or summary, wait to merge when the body chunk is read. Skipping more than 5 consecutive will trigger system interception
5. ⚠️ When one chapter spans multiple chunks, use the `merge` subcommand to combine into one resource file, do not write one file per chunk
6. End when all chunks in this batch are processed. Do not process chunks outside this batch, do not execute Step E

## Step D': Merge Index (Go side scheduling)

1. Execute `python3 scripts/merge_summaries.py summaries/ summary_index.txt`
   - ⚠️ **Output path MUST be `summary_index.txt` at working directory root**, NOT inside `skill/<name>/` subdirectory. Go side validates `summary_index.txt` at root and will fail if it's written elsewhere.
2. Read `summary_index.txt` to verify merged results
3. Constraint: Only execute the merge script, do not reprocess chunks

## Step E: Generate Skill Package

1. Execute `find skill/<skill_name>/resources/ -name '*.md' | sort` to get the list of all resource files
2. For each resource file, `head -5` to get #/##/### heading lines, build a complete table of contents tree (must cover all .md files)
3. Read `summary_index.txt` to supplement content summaries for each file (auxiliary only, actual file structure from steps 1-2 takes priority)
4. Write `skill/<skill_name>/SKILL.md` following the template structure, including:
   - Frontmatter (name/description)
   - Full book overview
   - Complete table of contents tree (each node clickable linking to resources/ files)
   - Quick reference guide (keyword → chapter path)
   - Usage instructions
5. ⚠️ Chapter files under resources/ were already written in Step D, no need to rewrite
6. ⚠️ summary_index.txt is only an auxiliary source for content summaries; TOC completeness must be based on actual files under resources/
7. ⚠️ The TOC tree must be strictly built following file-level # > ## > ### heading hierarchy. Never use directory path names as top-level headings. # headings must never be child nodes of other entries
8. Done

## Constraints
- No iteration round limit, but push forward efficiently, avoid re-reading the same chunk
- All intermediate outputs must be persisted to disk, do not rely on memory
- Write only to chunks/ skill/ summaries/ and working directory root
