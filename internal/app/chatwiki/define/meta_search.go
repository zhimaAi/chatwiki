// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

// Metadata filtering (robot configuration)
const (
	MetaSearchSwitchOff = 0
	MetaSearchSwitchOn  = 1
)

// Filter condition combination type: 1 AND 2 OR
const (
	MetaSearchTypeAnd = 1
	MetaSearchTypeOr  = 2
)

// Maximum number of conditions
const MetaSearchMaxConditions = 10

// Operator enumeration (stored in database/numbers passed between frontend and backend)
// string fields: is, is not, contains, does not contain, is empty, is not empty
// number/time fields: is, is not, is empty, is not empty, greater than, equal to, less than, greater than or equal to, less than or equal to
const (
	MetaOpIs          = 1  // is
	MetaOpIsNot       = 2  // is not
	MetaOpContains    = 3  // contains content (string)
	MetaOpNotContains = 4  // does not contain content (string)
	MetaOpEmpty       = 5  // is empty
	MetaOpNotEmpty    = 6  // is not empty
	MetaOpGt          = 7  // greater than
	MetaOpEq          = 8  // equal to
	MetaOpLt          = 9  // less than
	MetaOpGte         = 10 // greater than or equal to
	MetaOpLte         = 11 // less than or equal to
)
