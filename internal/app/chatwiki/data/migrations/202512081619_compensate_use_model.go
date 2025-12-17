// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package migrations

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func init() {
	goose.AddMigrationNoTxContext(func(ctx context.Context, db *sql.DB) error {
		return CompensateUseModel()
	}, nil)
}

func CompensateUseModel() error {
	var maxId int
	modelConfig := msql.Model(`chat_ai_model_config`, define.Postgres)
	if maxIdStr, err := modelConfig.Max(`id`); err != nil {
		logs.Error(`sql:%s,err:%s`, modelConfig.GetLastSql(), err.Error())
		return err
	} else {
		maxId = cast.ToInt(maxIdStr)
	}
	logs.Other(`compensate`, `获取最大ID:%d`, maxId)
	var size = 100 //每一批次数
	for i := 0; ; i++ {
		start, end := i*size, (i+1)*size
		logs.Other(`compensate`, `第%d轮:%d~%d`, i+1, start, end)
		list, err := modelConfig.Where(`id`, `>`, cast.ToString(start)).
			Where(`id`, `<=`, cast.ToString(end)).Select()
		if err != nil {
			logs.Error(`sql:%s,err:%s`, modelConfig.GetLastSql(), err.Error())
			return err
		}
		for _, config := range list { //逐个处理
			if err = CompensateUseModelOne(config); err != nil {
				logs.Error(`config:%s,err:%s`, tool.JsonEncodeNoError(config), err.Error())
				return err
			}
		}
		if end >= maxId {
			break //处理完毕,结束循环
		}
	}
	return nil
}

func CompensateUseModelOne(config msql.Params) error {
	switch config[`model_define`] {
	case common.ModelChatWiki:
		//nothing to do
	case common.ModelOpenAIAgent, common.ModelAzureOpenAI, common.ModelOllama, common.ModelXnference, common.ModelDoubao:
		useModel := common.UseModelConfig{
			ModelType:     config[`model_types`],
			UseModelName:  config[`deployment_name`],
			ShowModelName: config[`show_model_name`],
			ThinkingType:  cast.ToUint(config[`thinking_type`]),
		}
		switch config[`model_types`] {
		case common.Llm:
			useModel.ModelInputSupport = common.ModelInputSupport{InputText: 1}
			useModel.ModelOutputSupport = common.ModelOutputSupport{OutputText: 1}
		case common.TextEmbedding:
			useModel.ModelInputSupport = common.ModelInputSupport{InputText: 1}
		case common.Rerank:
			//nothing to do
		default:
			return fmt.Errorf(`特殊模型的model_types参数错误:%s`, config[`model_types`])
		}
		if len(useModel.ShowModelName) == 0 { //填充空值
			useModel.ShowModelName = useModel.UseModelName
		}
		if err := useModel.ToSave(cast.ToInt(config[`admin_user_id`]), cast.ToInt(config[`id`])); err != nil {
			return err
		}
	default:
		common.AutoAddDefaultUseModel(cast.ToInt(config[`admin_user_id`]), cast.ToInt(config[`id`]), config[`model_define`])
	}
	return nil
}
