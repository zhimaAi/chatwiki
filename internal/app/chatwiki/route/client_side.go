// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package route

import (
	"chatwiki/internal/app/chatwiki/business"
	"chatwiki/internal/app/chatwiki/business/manage"
	"net/http"
)

func RegClientSideRoute() {
	/*admin config API*/
	Route[http.MethodGet][`/manage/getClientSideLoginSwitch`] = manage.GetClientSideLoginSwitch
	Route[http.MethodPost][`/manage/setClientSideLoginSwitch`] = manage.SetClientSideLoginSwitch
	Route[http.MethodPost][`/manage/clientSideDownload`] = manage.ClientSideDownload
	/*client side API*/
	noAuthFuns(Route[http.MethodGet], `/manage/clientSide/getCompany`, business.GetCompany)
	noAuthFuns(Route[http.MethodGet], `/manage/clientSide/getRobotList`, business.ClientSideGetRobotList)
	noAuthFuns(Route[http.MethodPost], `/manage/clientSide/login`, business.ClientSideLogin)
}
