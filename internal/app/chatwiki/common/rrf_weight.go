// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"math"

	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

type RrfWeightConfig struct {
	Vector uint `json:"vector"`
	Search uint `json:"search"`
}

// GetDefaultRrfWeight Get the default RRF algorithm weight configuration
func GetDefaultRrfWeight() RrfWeightConfig {
	return RrfWeightConfig{Vector: 70, Search: 30}
}

// NormalizeRrfWeight drops the retired graph weight and rescales vector/search to 100.
func NormalizeRrfWeight(config RrfWeightConfig) RrfWeightConfig {
	sum := config.Vector + config.Search
	if sum == 0 || sum < config.Vector { // Empty or overflowed sum.
		return GetDefaultRrfWeight()
	}
	if sum == 100 {
		return RrfWeightConfig{Vector: config.Vector, Search: config.Search}
	}
	vector := min(uint(math.Round(float64(config.Vector)*100/float64(sum))), 100)
	return RrfWeightConfig{Vector: vector, Search: 100 - vector}
}

// ParseRrfWeight Parse the RRF algorithm weight configuration parameters
func ParseRrfWeight(rrfWeight string) RrfWeightConfig {
	if len(rrfWeight) == 0 { // Not configured, use default value
		return GetDefaultRrfWeight()
	}
	config := RrfWeightConfig{}
	if err := tool.JsonDecodeUseNumber(rrfWeight, &config); err != nil {
		logs.Error(`parse rrf_weight error: %v`, err)
		return GetDefaultRrfWeight()
	}
	return NormalizeRrfWeight(config)
}
