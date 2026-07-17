// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

import (
	"regexp"
	"time"
)

const (
	SystemSkillsDir  = `clawbot/skills_system/<biz_type>`
	PublicSkillsDir  = `clawbot/skills_public`
	UserSkillsDir    = `clawbot/skills_user/<admin_user_id>`
	PrivateSkillsDir = `clawbot/skills_robot/<robot_key>`
	PrivateFileDir   = `clawbot/skills_robot/<robot_key>/query-local-docs/references`
	PrivateWorkDir   = `clawbot/working_dir/<robot_key>`

	// BookToSkillTaskDir book-to-skill task working directories (require <robot_key> + <task_id> replacement)
	BookToSkillTaskDir = `clawbot/working_dir/<robot_key>/book_to_skill_<task_id>` //book to skill dir
	WorkInputDir       = BookToSkillTaskDir + `/input`
	WorkOutputDir      = BookToSkillTaskDir + `/skill`
	WorkScriptsDir     = BookToSkillTaskDir + `/scripts`
	WorkChunksDir      = BookToSkillTaskDir + `/chunks`
	WorkTemplatesDir   = BookToSkillTaskDir + `/templates`
	WorkSummariesDir   = BookToSkillTaskDir + `/summaries`

	// WebToSkillWorkDir web-to-skill task working directories
	WebToSkillWorkDir = `clawbot/working_dir/web-to-skill/<task_batch>`
)

// skill management constants
const (
	SkillMdFileName      = `SKILL.md`
	MaxSkillZipSize      = 10 * 1024 * 1024
	MaxSkillUnzipSize    = 50 * 1024 * 1024
	MaxSkillZipEntries   = 2000
	SkillReservedName    = `query-local-docs`
	SkillTmpDir          = `.skill_tmp`
	SkillUploadKeyExpire = 10 * time.Minute
)

// skill source type, maps to chat_ai_clawbot_skill.source_type
const (
	SkillSourceTypeUpload = 1
)

var SkillNameRegexp = regexp.MustCompile(`^[A-Za-z0-9_-]{1,50}$`)

var SkillUploadKeyRegexp = regexp.MustCompile(`^skup_[0-9]{14}_[A-Za-z0-9]{16}$`)

var ClawbotLocalDocAllowExt = []string{`docx`, `xlsx`, `md`, `txt`, `pdf`}

type ClawbotLocalDocInfo struct {
	Name string    `json:"name"`
	Size int64     `json:"size"`
	Time time.Time `json:"time"`
	Ext  string    `json:"ext"`
}

type ClawbotLocalDocIndexItem struct {
	File        string   `yaml:"file"`
	Type        string   `yaml:"type"`
	Description string   `yaml:"description"`
	Keywords    []string `yaml:"keywords"`
	Updated     string   `yaml:"updated"`
}
