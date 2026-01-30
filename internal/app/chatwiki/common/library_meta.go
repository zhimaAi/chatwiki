// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
)

func GetBuiltinMetaSchemaList(lang string) []define.BuiltinMetaSchema {
	return []define.BuiltinMetaSchema{
		{Name: i18n.Show(lang, `meta_source`), Key: define.BuiltinMetaKeySource, Type: define.LibraryMetaTypeString},
		{Name: i18n.Show(lang, `meta_update_time`), Key: define.BuiltinMetaKeyUpdateTime, Type: define.LibraryMetaTypeTime},
		{Name: i18n.Show(lang, `meta_create_time`), Key: define.BuiltinMetaKeyCreateTime, Type: define.LibraryMetaTypeTime},
		{Name: i18n.Show(lang, `meta_group`), Key: define.BuiltinMetaKeyGroup, Type: define.LibraryMetaTypeString},
	}
}
