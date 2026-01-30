// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

// LibraryMetaType 元数据类型（前后端传递数字）
const (
	LibraryMetaTypeString = 0
	LibraryMetaTypeTime   = 1
	LibraryMetaTypeNumber = 2
)

var LibraryMetaTypeList = [...]int{
	LibraryMetaTypeString,
	LibraryMetaTypeTime,
	LibraryMetaTypeNumber,
}

func IsLibraryMetaTypeValid(t int) bool {
	for _, v := range LibraryMetaTypeList {
		if v == t {
			return true
		}
	}
	return false
}

// 内置元数据 key（全局唯一）
const (
	BuiltinMetaKeySource     = "source"
	BuiltinMetaKeyUpdateTime = "update_time"
	BuiltinMetaKeyCreateTime = "create_time"
	BuiltinMetaKeyGroup      = "group"
)

type BuiltinMetaSchema struct {
	Name string
	Key  string
	Type int // 同 LibraryMetaType*
}

func IsBuiltinMetaKey(key string) bool {
	switch key {
	case BuiltinMetaKeySource, BuiltinMetaKeyUpdateTime, BuiltinMetaKeyCreateTime, BuiltinMetaKeyGroup:
		return true
	default:
		return false
	}
}
