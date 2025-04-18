// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func RerankData(adminUserId int, openid, appType string, robot msql.Params, req *adaptor.ZhimaRerankReq) ([]msql.Params, error) {
	modelConfigId, useModel := cast.ToInt(robot[`rerank_model_config_id`]), robot[`rerank_use_model`]
	handler, err := GetModelCallHandler(modelConfigId, useModel)
	if err != nil {
		return nil, err
	}
	res, err := handler.RequestRerank(adminUserId, openid, appType, robot, req)
	if err != nil {
		return nil, err
	}
	// todo token统计
	if handler.modelInfo != nil && handler.modelInfo.TokenUseReport != nil { //token use report
		handler.modelInfo.TokenUseReport(handler.config, useModel, res.InputToken, res.OutputToken)
	}
	for key, item := range req.Data {
		item[`similarity`] = cast.ToString(res.Data[key].RelevanceScore)
	}
	return req.Data, nil
}
