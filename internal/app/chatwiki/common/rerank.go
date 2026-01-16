// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"sort"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func RerankData(adminUserId int, openid, appType string, robot msql.Params, req *adaptor.ZhimaRerankReq) ([]msql.Params, error) {
	modelConfigId, useModel := cast.ToInt(robot[`rerank_model_config_id`]), robot[`rerank_use_model`]
	handler, err := GetModelCallHandler(adminUserId, modelConfigId, useModel, robot)
	if err != nil {
		return nil, err
	}
	res, err := handler.RequestRerank(adminUserId, openid, appType, robot, req)
	if err != nil {
		return nil, err
	}
	if handler.modelInfo != nil && handler.modelInfo.TokenUseReport != nil { //token use report
		handler.modelInfo.TokenUseReport(handler.config, useModel, res.InputToken, res.OutputToken, robot, 0)
	}
	sort.Slice(res.Data, func(i, j int) bool {
		return res.Data[i].RelevanceScore > res.Data[j].RelevanceScore
	})
	rerankList := make([]msql.Params, 0)
	for _, item := range res.Data {
		if len(req.Data) > item.Index {
			rerankList = append(rerankList, req.Data[item.Index])
		}
	}
	return rerankList, nil
}
