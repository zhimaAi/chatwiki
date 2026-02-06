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

// GetDefaultRrfWeight Get the default RRF algorithm weight configuration
func GetDefaultRrfWeight(adminUserId int) RrfWeightConfig {
	if GetNeo4jStatus(adminUserId) {
		return RrfWeightConfig{Vector: 50, Search: 30, Graph: 20}
	}
	return RrfWeightConfig{Vector: 70, Search: 30, Graph: 0}
}

// CheckRrfWeight Validate the RRF algorithm weight parameters configured by the frontend
func CheckRrfWeight(rrfWeight, lang string) error {
	if len(rrfWeight) == 0 {
		return nil // Skip validation for empty values, use default values directly
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

// ParseRrfWeight Parse the RRF algorithm weight configuration parameters
func ParseRrfWeight(adminUserId int, rrfWeight string) RrfWeightConfig {
	if len(rrfWeight) == 0 { // Not configured, use default value
		return GetDefaultRrfWeight(adminUserId)
	}
	config := RrfWeightConfig{}
	if err := tool.JsonDecodeUseNumber(rrfWeight, &config); err != nil {
		logs.Error(`admin_user_id:%d,rrf_weight:%s,err:%v`, adminUserId, rrfWeight, err)
		return GetDefaultRrfWeight(adminUserId)
	}
	return config
}
