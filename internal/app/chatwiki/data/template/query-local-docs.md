---
name: query-local-docs
description: Query local reference documents stored in the skill's references directory. Use when the user asks about uploaded files, local documents, PDFs, DOCX files, Markdown files, product manuals, policy documents, README files, or any knowledge stored in the references folder.
---

# Query Local Docs

Use this skill to answer questions from local reference documents. All documents are stored in the `references/`
directory under this skill. The document set is dynamic: files can be added, removed, or replaced at any time.

## Directory Structure

```
query-local-docs/
├── SKILL.md              # 本 Skill 定义文件
└── references/           # 文档存放目录（只读）
    ├── index.yaml        # 文档索引清单（YAML 格式，便于代码动态更新）
    └── ...               # 用户上传的文档文件
```

## Document Index

`references/index.yaml` 是文档索引文件，采用 YAML 格式以便代码动态增删改查。格式如下：

```yaml
- file: product-manual.md
  type: Markdown
  description: 产品使用手册 V2.1
  keywords: [ 安装, 配置, 故障排除, 规格参数 ]
  updated: 2024-03-15

- file: pricing-policy.md
  type: Markdown
  description: 定价策略及折扣规则
  keywords: [ 价格, 折扣, 套餐, 阶梯计费 ]
  updated: 2024-04-01
```

字段说明：

- `file`：文件名（必填）
- `type`：文件类型，如 Markdown, PDF, DOCX 等
- `description`：文件内容简要描述
- `keywords`：关键词数组，用于快速匹配用户问题
- `updated`：最后更新日期（YYYY-MM-DD）

## Workflow — 分层检索

采用从粗到细的分层检索策略，最大化效率和精度：

### Step 1: 索引定位

**必须**优先读取 `references/index.yaml` 作为文件定位的第一步。禁止使用 `ls`/`glob` 代替 index.yaml 进行文件发现。

根据用户问题的语义和关键词，匹配 index 中的 `description` 和 `keywords` 字段，筛选出可能相关的候选文件（通常 1-3 个）。

- 当且仅当 `index.yaml` 文件不存在或读取失败时，才允许 fallback 到 `ls`/`glob` 扫描目录
- 当多个文件涉及同一主题且内容冲突时，优先采信「更新时间」较新的文件

### Step 2: 关键词搜索

对 Step 1 筛选出的候选文件使用 `grep` 搜索用户问题中的关键词及其近义词。

- grep 的 path 参数**必须**是具体文件路径，**禁止**传入目录路径
    - ✅ 正确：`grep pattern references/product-manual.md`
    - ❌ 错误：`grep pattern references/`
- 使用用户问题中的核心术语
- 同时尝试近义词和相关表达
- 获取匹配行前后的上下文（通常 3-5 行）

### Step 3: 按需精读

如果 grep 片段不足以完整回答问题，读取匹配段落的更大上下文（整个章节或前后 N 行）。

### Step 4: 兜底全读

仅当以上步骤均未找到足够信息时，才全文读取文件。

## 强制规则

- **禁止**在未经 grep 搜索的情况下直接全文读取任何文件。
- 全文读取（Step 4）仅在以下条件**全部满足**时允许：
    1. 已对候选文件执行 grep 搜索但未找到足够匹配
    2. 已尝试至少 2 组不同的关键词/近义词
    3. 文件总行数 ≤ 200 行
- 如果文件超过 200 行且 grep 无结果，必须调整关键词重新 grep，**不可**直接全读。

### 正确示例

✅ 用户问「安装步骤是什么」：

1. 从 index.yaml 定位到 product-manual.md
2. `grep "安装"` → 获取匹配片段
3. 精读匹配段落的上下文 → 回答

✅ 用户问「退款政策」：

1. 从 index.yaml 定位到 pricing-policy.md
2. `grep "退款"` → 无结果
3. `grep "退货|售后|取消"` → 获取匹配片段
4. 精读匹配段落 → 回答

### 反面示例（禁止）

❌ 用户问「安装步骤是什么」→ 直接 `read product-manual.md` 全文
❌ 用户问「价格是多少」→ 跳过 grep，直接全文读取 pricing-policy.md
❌ grep 一次无结果 → 立即全文读取（应先换关键词重试）
❌ grep 时 path 传入目录而非具体文件（如 `grep pattern references/`）

## Query Strategy

- 对于「有哪些文件」「列出文档」等清单类请求，直接返回 `index.yaml` 的内容或扫描 `references/` 目录。
- 对于具体问题，严格遵循分层检索流程：index → grep → 精读 → 全读。
- 对于政策、产品、排障、价格、参数、流程类问题，比对所有匹配文件，不要在第一个命中处停止。
- 如果多个文件内容冲突，报告冲突并标注每条信息的来源文件。不要静默合并矛盾内容。
- 如果文件是二进制或不可直接读取，告知用户该文件无法在当前文本视图中检查，建议上传可读的文本/Markdown 版本。

## Boundaries

- 仅读取 `references/` 目录下的文件，不可写入或修改。
- 不可读取 `references/` 目录之外的任何文件。
- 除非用户明确要求文件转换或提取且环境允许，否则不执行任何命令。
- 不可仅凭文件名臆测内容。文件名可引导搜索方向，但答案必须来自可读的文件内容。
- 回答时必须引用来源文件名，并尽可能包含章节标题、页码或其他定位标记。
- 如果文档中没有足够证据回答问题，明确告知用户「当前文档中未找到相关信息」。
