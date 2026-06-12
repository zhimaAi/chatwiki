// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

import "time"

const (
	PublicSkillsDir  = `clawbot/skills_public`
	PrivateSkillsDir = `clawbot/skills_robot/<robot_key>`
	PrivateFileDir   = `clawbot/skills_robot/<robot_key>/query-local-docs/references`
	PrivateWorkDir   = `clawbot/working_dir/<robot_key>`
)

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

const QueryLocalDocsIndexDesc = "# Document Index - Records the meta-information of all documents under the references/ directory\r\n"
