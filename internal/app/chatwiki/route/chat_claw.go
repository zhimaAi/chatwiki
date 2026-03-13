// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package route

import (
	"chatwiki/internal/app/chatwiki/business"
	"chatwiki/internal/app/chatwiki/business/manage"
	"net/http"
)

func RegChatClawRoute() {
	noAuthFuns(Route[http.MethodPost], `/manage/chatclaw/login`, manage.ChatClawLogin)
	noAuthFuns(Route[http.MethodPost], `/manage/chatclaw/chat/completions`, business.ChatClawChat)
	Route[http.MethodGet][`/manage/chatclaw/getRobotList`] = manage.ChatClawGetRobotList
	Route[http.MethodGet][`/manage/chatclaw/getLibraryList`] = manage.ChatClawGetLibraryList
	Route[http.MethodGet][`/manage/chatclaw/getLibFileList`] = manage.ChatClawGetLibFileList
	Route[http.MethodGet][`/manage/chatclaw/getParagraphList`] = manage.ChatClawGetParagraphList
	Route[http.MethodGet][`/manage/chatclaw/getLibraryGroup`] = manage.ChatClawGetLibraryGroup
	Route[http.MethodGet][`/manage/chatclaw/getEnabledRobotList`] = manage.ChatClawGetEnabledRobotList
	Route[http.MethodGet][`/manage/chatclaw/getEnabledLibraryList`] = manage.ChatClawGetEnabledLibraryList
	Route[http.MethodPost][`/manage/chatclaw/updateRobotSwitchStatus`] = manage.ChatClawUpdateRobotSwitchStatus
	Route[http.MethodPost][`/manage/chatclaw/updateLibrarySwitchStatus`] = manage.ChatClawUpdateLibrarySwitchStatus
	Route[http.MethodGet][`/manage/chatclaw/tokenLogList`] = manage.ChatClawTokenLogList
	Route[http.MethodPost][`/manage/chatclaw/tokenForceOffline`] = manage.ChatClawTokenForceOffline
	Route[http.MethodPost][`/manage/chatclaw/refreshToken`] = manage.ChatClawRefreshToken
	Route[http.MethodPost][`/manage/chatclaw/libraryRecall`] = manage.ChatClawLibraryRecall
}
