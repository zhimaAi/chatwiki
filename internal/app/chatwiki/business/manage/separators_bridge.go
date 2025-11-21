// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"

	"github.com/spf13/cast"
)

func BridgeGetSeparatorsList(adminUserId, loginUserId int, lang string) ([]map[string]any, int, error) {
	list := make([]map[string]any, 0)
	for _, item := range define.SeparatorsList {
		name := i18n.Show(lang, cast.ToString(item[`name`]))
		list = append(list, map[string]any{`no`: item[`no`], `name`: name})
	}
	return list, 0, nil
}
