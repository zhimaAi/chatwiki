// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

// BookToSkill task status constants
const (
	BookToSkillStatusPending = 0 // Pending
	BookToSkillStatusRunning = 1 // Running
	BookToSkillStatusSuccess = 2 // Success
	BookToSkillStatusFailed  = 3 // Failed
	BookToSkillStatusStopped = 4 // Stopped
)

// BookToSkill install status
const (
	BookToSkillInstallNone      = 0 // Not installed
	BookToSkillInstallInstalled = 1 // Installed
)

// BookToSkill file constraints
const (
	BookToSkillMaxFileCount      = 20  // Maximum number of uploaded files
	BookToSkillMaxFileName       = 20  // Maximum skill name characters (utf-8 chars)
	BookToSkillFileLimitMB       = 100 // Single file max MB
	BookToSkillFileLimitSize     = BookToSkillFileLimitMB * 1024 * 1024
	BookToSkillTaskTimeout       = 60 * 60 // Task timeout (seconds); NSQ async consumption doesn't block HTTP; large files need longer processing time
	BookToSkillDeepMaxIteration  = 300     // Deep agent single-step max iteration count (sufficient for several MB docs; too large leads to wasted rotations for large files; 0 falls back to default 20, cannot be 0)
	BookToSkillCompressThreshold = 120000  // Context compression threshold (rough token estimate, chars/3); reserves space for compressed content to avoid 128K context model overflow
	BookToSkillMaxToken          = 32768   // Single request max output tokens: steps 2/3 write_file need larger output, compatible with mainstream 32K+ output models
	// Single-agent batch processing constants (step D map-reduce)
	BookToSkillBatchSize = 5 // Step D chunks per batch
	// Split guidance constants (injected into step B prompt, AI references when generating split_chapters.py)
	BookToSkillSplitMaxBytes  = 16 * 1024 // 16KB, triggers recursive split when exceeded
	BookToSkillSplitMaxDepth  = 2         // Max recursive split depth
	BookToSkillSplitMaxChunks = 50        // Max chunks in split phase (log warning when exceeded, triggers AI re-split)
)

// BookToSkill allowed file extensions
var BookToSkillAllowExt = []string{
	`txt`, `docx`, `md`,
}

// BookToSkillSkillName is the skill package directory name under skills_public
const BookToSkillSkillName = `book-to-skill`

// BookToSkillTemplateDir is the source template directory (SKILL.md, scripts, templates)
const BookToSkillTemplateDir = `clawbot/skills_system/alone_book/` + BookToSkillSkillName

// BookToSkill default language when not specified
const BookToSkillDefaultLang = LangZhCn

// btsSingleAgentHeader_en is the header prompt in English
const btsSingleAgentHeader_en = "You are the Book-to-Skill single agent. Execute this phase task following the six-step workflow (A→A'→B→C→D→E) defined in SKILL.md.\n" +
	"Split params: single chunk ≤ %dKB, recursion ≤ %d levels\n\n"

// btsExecutePathRules_en path rules for execute tool (English)
const btsExecutePathRules_en = `
## ⚠️ execute Command Path Rules (every command must follow these, or execution will fail)
- The execute tool runs in the **current working directory**; you don't need cd and don't need to prepend the working directory path
- The working directory is the value of the "Working Directory" field above; all scripts and directories are children of this directory
- **Command arguments must use paths relative to the current directory**; workspace-relative or absolute paths are forbidden

Correct examples:
  python3 scripts/convert_to_md.py input/              ← script is in scripts/ subdirectory
  python3 scripts/split_chapters.py input_md/ chunks/ --pattern '^# .+'
  python3 scripts/write_resource.py copy --src chunks/chunk_010.md --dst skill/xxx/resources/Ch1/1.md
  ls chunks/

Wrong examples (the following formats will cause execution failure, strictly forbidden):
  python3 clawbot/working_dir/.../scripts/xxx.py       ← WRONG! Duplicated path prefix
  python3 /workspace/clawbot/working_dir/.../xxx.py     ← WRONG! Absolute path
  cd clawbot/working_dir/... && python3 scripts/xxx.py  ← WRONG! cd is ignored
  find / -name "xxx.py" 2>/dev/null                     ← WRONG! Redirect > forbidden

- When encountering "No such file or directory" error, check if the path is relative to the current working directory
- If you don't know what path to use, first run ls to see what files/subdirectories exist in the current directory
`

// compressSystemPrompt_en compression system prompt (English)
const compressSystemPrompt_en = "You are a summarization assistant. Compress single-agent conversation history into a structured progress summary, strictly following the format."

// compressUserPrompt_en compression user prompt template (English)
const compressUserPrompt_en = `Please compress the following processed batch conversation history into a progress summary not exceeding 2KB.
Keep: chapter path ranges per batch, titles, content summary highlights, source chunk ranges.
Discard: chunk raw text, tool call details, duplicate information.
Output format (one line per batch):
Batch 1 [chunk_001~050]《Chapter 1》: resources/Ch1/1.1.md,1.2.md | Character background overview
Batch 2 [chunk_051~100]《Chapter 2》: resources/Ch2/2.1.md | World-building setting
...
【Current Progress】Completed batches=N/Total, Current batch=M Processed chunk_XXX
【History to Compress】
`

// btsPhasePromptSplit_en split phase prompt (English)
const btsPhasePromptSplit_en = "【Current Phase】Steps A+A'+B+C (see SKILL.md for details)\n" +
	"Working Directory: %s\n" +
	"⚠️ Split target: single chunk ≤ %dKB, recursion ≤ %d levels, total chunks ≤ %d (if exceeded, switch to coarser --pattern and re-split)\n" +
	"⚠️ This phase forbids reading any files under chunks/\n" +
	"⚠️ Step C: only re-split chunks > %dKB, do not read one by one\n" +
	"Finish when done."

// btsPhasePromptBatch_en batch phase prompt (English)
const btsPhasePromptBatch_en = "【Current Phase】Step D Batch %d/%d (see SKILL.md for details)\n" +
	"Working Directory: %s | Skill Name: %s\n" +
	"This batch has %d chunks: %s\n" +
	"⚠️ Hard rule: after ≤3 consecutive reads, you must execute write_resource.py + write summary to disk\n" +
	"Finish after processing all chunks in this batch."

// btsPhasePromptMerge_en merge phase prompt (English)
const btsPhasePromptMerge_en = "【Current Phase】Step D' (see SKILL.md for details)\n" +
	"Working Directory: %s\n" +
	"Execute merge_summaries.py then finish."

// btsPhasePromptGenerate_en generate phase prompt (English)
const btsPhasePromptGenerate_en = "【Current Phase】Step E (see SKILL.md for details)\n" +
	"Working Directory: %s | Skill Name: %s\n" +
	"Generate SKILL.md following SKILL.md instructions and template structure, then finish."

// btsPhasePromptUnknown_en unknown phase fallback (English)
const btsPhasePromptUnknown_en = "Unknown phase: %s"

// BtsPromptKey identifies a BTS prompt
type BtsPromptKey string

const (
	BtsPromptHeader           BtsPromptKey = "header"
	BtsPromptExecutePathRules BtsPromptKey = "execute_path_rules"
	BtsPromptCompressSystem   BtsPromptKey = "compress_system"
	BtsPromptCompressUser     BtsPromptKey = "compress_user"
)

// GetBtsPrompt returns the bilingual prompt for the given key and language.
// Returns empty string for unknown keys. Only zh-CN and en-US are supported;
// fallback to zh-CN for unknown languages.
func GetBtsPrompt(lang string, key BtsPromptKey) string {
	switch key {
	case BtsPromptHeader:
		if lang == LangEnUs {
			return btsSingleAgentHeader_en
		}
		return btsSingleAgentHeader_en
	case BtsPromptExecutePathRules:
		if lang == LangEnUs {
			return btsExecutePathRules_en
		}
		return btsExecutePathRules_en
	case BtsPromptCompressSystem:
		if lang == LangEnUs {
			return compressSystemPrompt_en
		}
		return compressSystemPrompt_en
	case BtsPromptCompressUser:
		if lang == LangEnUs {
			return compressUserPrompt_en
		}
		return compressUserPrompt_en
	default:
		return ""
	}
}

// GetBtsPhasePrompt returns the bilingual phase-specific prompt body.
// lang: "zh-CN" or "en-US"; phase: "split", "batch", "merge", "generate".
// Returns a format string with positional args for fmt.Sprintf.
func GetBtsPhasePrompt(lang string, phase string) string {
	if lang == LangEnUs {
		switch phase {
		case "split":
			return btsPhasePromptSplit_en
		case "batch":
			return btsPhasePromptBatch_en
		case "merge":
			return btsPhasePromptMerge_en
		case "generate":
			return btsPhasePromptGenerate_en
		default:
			return btsPhasePromptUnknown_en
		}
	}
	switch phase {
	case "split":
		return btsPhasePromptSplit_en
	case "batch":
		return btsPhasePromptBatch_en
	case "merge":
		return btsPhasePromptMerge_en
	case "generate":
		return btsPhasePromptGenerate_en
	default:
		return btsPhasePromptUnknown_en
	}
}
