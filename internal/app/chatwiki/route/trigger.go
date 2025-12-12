// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package route

import (
	"chatwiki/internal/app/chatwiki/business/manage"
	"net/http"
)

func RegTriggerRoute() {
	Route[http.MethodPost][`/manage/triggerList`] = manage.TriggerList
	Route[http.MethodPost][`/manage/triggerSwitch`] = manage.TriggerSwitch
}
