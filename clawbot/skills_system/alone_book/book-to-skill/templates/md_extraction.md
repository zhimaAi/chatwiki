# Book-to-Skill Output Template

## summary_index.txt Format Definition (Markdown Table)

Each row corresponds to one resource file, four columns:

| Chapter Path | Title | Content Summary | Source File Path |
|---------|------|---------|-----------|
| resources/第1章/1.1 Overview.md | 1.1 Overview | Introduces character backgrounds, world-building settings | chunks/chunk_001.md |
| resources/第1章/1.2 Worldview.md | 1.2 Worldview | Describes continent layout, faction distribution | chunks/chunk_003.md |

Constraints:
- Chapter path is relative to skill/<skill_name>/resources/
- Title uses original chapter heading text from the source
- Content summary kept within 50 characters
- Source file path points to the original chunk filename under chunks/

## Split Level Description

Three-level structure: Chapter > Section > Sub-section

| Original Structure | Output Path |
|---------|---------|
| Has chapters and sections | resources/<chapter_name>/<section_title>.md |
| Has chapters, no sections | resources/<chapter_name>/<chapter_name>.md |
| No chapters or sections | resources/<topic_category>/<paragraph_title>.md (split by content paragraphs) |
| Recursive split | Single file > 16KB: re-split at next level, max 2 levels |

⚠️ Chapter naming rules:
- Chapter name MUST be taken from the file's `#` first-level heading original text (e.g. `# Product Introduction` → chapter name is `Product Introduction`)
- Never use the source file's directory name as chapter name
- If a single chunk contains multiple `###` sub-headings, must split into separate files: each `###` section = one resources/<chapter_name>/<section_title>.md

⚠️ Table of Contents tree rules:
- SKILL.md's complete TOC tree must be strictly built following file-level # > ## > ### heading hierarchy
- # heading = level-1 TOC node, ## heading = level-2 child node, ### heading = level-3 child node
- Never create a directory level higher than # heading (e.g. using directory path name as top-level heading)

---

## SKILL.md Output Structure

---
name: <book title>
description: <brief description based on book's preface or introduction>
---

# <Full Book Name>

## Overview
<If the book has a preface or abstract, copy directly; otherwise summarize from the first few paragraphs>

## 📚 Complete Table of Contents

- [第1章: Title](resources/第1章/1.1 Section.md)
    - [1.1 Section: Title](resources/第1章/1.1 Section.md)
    - [1.2 Section: Title](resources/第1章/1.2 Section.md)
        - [1.2.1 Subsection](resources/第1章/1.2.1 Subsection.md)
- [第2章: Title](resources/第2章/2.1 Section.md)
  ...

(Continue as above, fully reflecting the original book hierarchy; each node is clickable linking to its resource file)

## 🔍 Quick Reference Guide

- If user asks about **<keyword A>** → refer to [<chapter path>](resources/...)
- If user asks about **<keyword B>** → refer to [<chapter path>](resources/...)
- (Generate high-frequency keyword mappings from summary_index.txt)

## Usage Instructions

When the user asks a question, this skill will:
1. Find the most relevant chapter path based on keyword index;
2. Read the original text content from the corresponding resource file;
3. Use the original text as context to provide a comprehensive answer (citing original paragraphs when necessary).