// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"context"
	"sort"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/llm_adaptor/v2/rerank"
)

func RerankData(ctx context.Context, lang string, adminUserId int, openid, appType string, robot msql.Params, req *rerank.CreateRequest, data []msql.Params) ([]msql.Params, error) {
	modelConfigId, useModel := cast.ToInt(robot[`rerank_model_config_id`]), robot[`rerank_use_model`]
	handler, err := GetModelCallHandler(lang, adminUserId, modelConfigId, useModel, robot)
	if err != nil {
		return nil, err
	}
	request := *req
	request.Model = handler.Model
	ctx, usageUUID := withModelUsageUUID(ctx)
	res, err := handler.RequestRerank(ctx, lang, adminUserId, openid, appType, robot, &request)
	if err != nil {
		return nil, err
	}
	if handler.modelInfo != nil && handler.modelInfo.TokenUseReport != nil { //token use report
		inputToken, outputToken := rerankTokens(res)
		handler.modelInfo.TokenUseReport(handler.config, useModel, inputToken, outputToken, robot, 0, usageUUID)
	}
	sort.Slice(res.Results, func(i, j int) bool {
		return res.Results[i].RelevanceScore > res.Results[j].RelevanceScore
	})
	rerankList := make([]msql.Params, 0)
	for _, item := range res.Results {
		if len(data) > item.Index {
			rerankList = append(rerankList, data[item.Index])
		}
	}
	return rerankList, nil
}
