// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/i18n"
	"errors"

	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

type RrfWeightConfig struct {
	Vector uint `json:"vector"`
	Search uint `json:"search"`
	Graph  uint `json:"graph"`
}

// GetDefaultRrfWeight 获取默认的RRF算法权重配置
func GetDefaultRrfWeight(adminUserId int) RrfWeightConfig {
	if GetNeo4jStatus(adminUserId) {
		return RrfWeightConfig{Vector: 50, Search: 30, Graph: 20}
	}
	return RrfWeightConfig{Vector: 70, Search: 30, Graph: 0}
}

// CheckRrfWeight 校验前端配置的RRF算法权重参数
func CheckRrfWeight(rrfWeight, lang string) error {
	if len(rrfWeight) == 0 {
		return nil //空值不校验,直接使用默认值
	}
	config := RrfWeightConfig{}
	if err := tool.JsonDecodeUseNumber(rrfWeight, &config); err != nil {
		return errors.New(i18n.Show(lang, `param_invalid`, `rrf_weight`))
	}
	if config.Vector+config.Search+config.Graph != 100 {
		return errors.New(i18n.Show(lang, `rrf_weight_config_sum_invalid`))
	}
	return nil
}

// ParseRrfWeight 解析RRF算法权重配置参数
func ParseRrfWeight(adminUserId int, rrfWeight string) RrfWeightConfig {
	if len(rrfWeight) == 0 { //未配置,使用默认值
		return GetDefaultRrfWeight(adminUserId)
	}
	config := RrfWeightConfig{}
	if err := tool.JsonDecodeUseNumber(rrfWeight, &config); err != nil {
		logs.Error(`admin_user_id:%d,rrf_weight:%s,err:%v`, adminUserId, rrfWeight, err)
		return GetDefaultRrfWeight(adminUserId)
	}
	return config
}
