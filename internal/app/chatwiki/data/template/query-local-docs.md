---
name: query-local-docs
description: Query local reference documents stored in the skill's references directory. Use when the user asks about uploaded files, local documents, PDFs, DOCX files, Markdown files, product manuals, policy documents, README files, or any knowledge stored in the references folder.
---

# Query Local Docs

从 `references/` 目录下的本地文档中检索信息。

## 文档索引

```yaml
<index_yaml>
```

## 检索规则

根据索引中的 `description` 和 `keywords` 定位候选文件后，按以下流程执行：

1. **grep 搜索** — 对候选文件用关键词及近义词 grep，获取匹配行及上下文
2. **按需精读** — grep 片段不足时，读取匹配段落的更大上下文
3. **兜底全读** — 已 grep 至少 2 组关键词无结果且文件 ≤ 200 行时才允许

**禁止：**

- 未经 grep 直接 read 全文
- grep path 传入目录（必须为具体文件，如 `references/filename.ext`）
- 用 `ls`/`glob` 代替索引定位
- 多文件冲突时优先采信 `updated` 较新的文件
