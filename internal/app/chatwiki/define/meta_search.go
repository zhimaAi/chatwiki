// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

// 元数据过滤（机器人配置）
const (
	MetaSearchSwitchOff = 0
	MetaSearchSwitchOn  = 1
)

// 过滤条件组合类型：1 且 2 或
const (
	MetaSearchTypeAnd = 1
	MetaSearchTypeOr  = 2
)

// 最多条件数
const MetaSearchMaxConditions = 10

// 操作符枚举（存库/前后端传递数字）
// string字段：是、不是、内容包含、内容不包含、为空、不为空
// number/time字段：是、不是、为空、不为空、大于、等于、小于、大于等于、小于等于
const (
	MetaOpIs          = 1  // 是
	MetaOpIsNot       = 2  // 不是
	MetaOpContains    = 3  // 内容包含（string）
	MetaOpNotContains = 4  // 内容不包含（string）
	MetaOpEmpty       = 5  // 为空
	MetaOpNotEmpty    = 6  // 不为空
	MetaOpGt          = 7  // 大于
	MetaOpEq          = 8  // 等于
	MetaOpLt          = 9  // 小于
	MetaOpGte         = 10 // 大于等于
	MetaOpLte         = 11 // 小于等于
)

