// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func RerankData(modelConfigId int, useModel string, req *adaptor.ZhimaRerankReq) ([]msql.Params, error) {
	handler, err := GetModelCallHandler(modelConfigId, useModel)
	if err != nil {
		return nil, err
	}
	return handler.RequestRerank(req)
}
