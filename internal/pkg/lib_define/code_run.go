// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package lib_define

const CodeRunKey = `QwXKBkHZRJ4StIYr06hvUOLVu9AemFa2`

type CodeRunBody struct {
	MainFunc string         `json:"main_func"`
	Params   map[string]any `json:"params"`
}
