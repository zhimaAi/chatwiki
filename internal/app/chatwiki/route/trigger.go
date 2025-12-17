// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package route

import (
	"chatwiki/internal/app/chatwiki/business/manage"
	"net/http"
)

func RegTriggerRoute() {
	Route[http.MethodPost][`/manage/triggerList`] = manage.TriggerList
	Route[http.MethodPost][`/manage/triggerSwitch`] = manage.TriggerSwitch
}
