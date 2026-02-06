// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

// LibraryMetaType metadata type (numbers passed between frontend and backend)
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

// Built-in metadata keys (globally unique)
const (
	BuiltinMetaKeySource     = "source"
	BuiltinMetaKeyUpdateTime = "update_time"
	BuiltinMetaKeyCreateTime = "create_time"
	BuiltinMetaKeyGroup      = "group"
)

type BuiltinMetaSchema struct {
	Name string
	Key  string
	Type int // Same as LibraryMetaType*
}

func IsBuiltinMetaKey(key string) bool {
	switch key {
	case BuiltinMetaKeySource, BuiltinMetaKeyUpdateTime, BuiltinMetaKeyCreateTime, BuiltinMetaKeyGroup:
		return true
	default:
		return false
	}
}
