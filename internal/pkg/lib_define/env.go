// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package lib_define

import (
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

const ElectronPath = `front-end/chat-ai-electron/`

func GetElectronVersion() (version string) {
	version = `0.0.0`
	content, err := tool.ReadFile(ElectronPath + `package.json`)
	if err != nil {
		logs.Error(err.Error())
		return
	}
	params := make(map[string]any)
	if err := tool.JsonDecodeUseNumber(content, &params); err != nil {
		logs.Error(err.Error())
		return
	}
	if v := cast.ToString(params[`version`]); len(v) > 0 {
		return v
	}
	return
}
