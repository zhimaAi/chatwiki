// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

type BridgeGetModelConfigOptionReq struct {
	ModelType string `form:"model_type"`
}

func BridgeGetModelConfigOption(adminUserId, userId int, lang string, req *BridgeGetModelConfigOptionReq) ([]map[string]any, int, error) {
	modelType := strings.TrimSpace(req.ModelType)
	if len(modelType) == 0 {
		return nil, -1, errors.New(i18n.Show(lang, `param_lack`))
	}
	configs, err := msql.Model(`chat_ai_model_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).Order(`id desc`).Select()
	if err != nil {
		logs.Error(err.Error())
		return nil, -1, errors.New(i18n.Show(lang, `sys_err`))
	}
	list := make([]map[string]any, 0)
	for _, config := range configs {
		if tool.InArrayString(modelType, strings.Split(config[`model_types`], `,`)) {
			modelInfo, ok := common.GetModelInfoByDefine(config[`model_define`])
			if !ok {
				continue
			}
			choosableThinking := make(map[string]bool)
			if cast.ToInt(config[`thinking_type`]) == 2 { //可选设置在配置上的
				for _, modelName := range modelInfo.LlmModelList {
					choosableThinking[fmt.Sprintf(`%s#%s`, config[`id`], modelName)] = true
				}
				choosableThinking[fmt.Sprintf(`%s#%s`, config[`id`], config[`show_model_name`])] = true //兼容
			} else { //可选设置在模型上的
				for _, modelName := range modelInfo.ChoosableThinkingModels {
					choosableThinking[fmt.Sprintf(`%s#%s`, config[`id`], modelName)] = true
				}
			}
			list = append(list, map[string]any{`model_config`: config, `model_info`: modelInfo, `choosable_thinking`: choosableThinking})
		}
	}
	return list, 0, nil
}
