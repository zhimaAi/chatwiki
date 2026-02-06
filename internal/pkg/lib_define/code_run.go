// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package lib_define

const CodeRunKey = `QwXKBkHZRJ4StIYr06hvUOLVu9AemFa2`

type CodeRunBody struct {
	MainFunc string         `json:"main_func"`
	Params   map[string]any `json:"params"`
}

const (
	LanguageJavaScript = `javaScript`
	LanguagePython     = `python`
)
