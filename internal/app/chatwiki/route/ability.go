// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package route

import (
	"chatwiki/internal/app/chatwiki/business/manage"
	"net/http"
)

func RegAbilityRoute() {
	/*ability API*/
	Route[http.MethodPost][`/manage/ability/saveUserAbility`] = manage.SaveUserAbility
	Route[http.MethodGet][`/manage/ability/getAbilityList`] = manage.GetAbilityList

	/*robot ability API*/
	Route[http.MethodPost][`/manage/ability/saveRobotAbility`] = manage.SaveRobotAbility
	Route[http.MethodPost][`/manage/ability/saveRobotAbilitySwitchStatus`] = manage.SaveRobotAbilitySwitchStatus
	Route[http.MethodPost][`/manage/ability/saveRobotAbilityFixedMenu`] = manage.SaveRobotAbilityFixedMenu
	Route[http.MethodPost][`/manage/ability/saveRobotAbilityAiReplyStatus`] = manage.SaveRobotAbilityAiReplyStatus
	Route[http.MethodGet][`/manage/ability/getRobotAbilityList`] = manage.GetRobotAbilityList
	Route[http.MethodGet][`/manage/ability/getRobotSpecifyAbilityConfig`] = manage.GetRobotSpecifyAbilityConfig

	/*robot keyword reply API*/
	Route[http.MethodPost][`/manage/ability/saveRobotKeywordReply`] = manage.SaveRobotKeywordReply
	Route[http.MethodPost][`/manage/ability/checkKeyWordRepeat`] = manage.CheckKeyWordRepeat
	Route[http.MethodPost][`/manage/ability/deleteRobotKeywordReply`] = manage.DeleteRobotKeywordReply
	Route[http.MethodGet][`/manage/ability/getRobotKeywordReply`] = manage.GetRobotKeywordReply
	Route[http.MethodGet][`/manage/ability/getRobotKeywordReplyList`] = manage.GetRobotKeywordReplyList
	Route[http.MethodPost][`/manage/ability/updateRobotKeywordReplySwitchStatus`] = manage.UpdateRobotKeywordReplySwitchStatus

	/*robot received message reply API*/
	Route[http.MethodPost][`/manage/ability/saveRobotReceivedMessageReply`] = manage.SaveRobotReceivedMessageReply
	Route[http.MethodPost][`/manage/ability/deleteRobotReceivedMessageReply`] = manage.DeleteRobotReceivedMessageReply
	Route[http.MethodGet][`/manage/ability/getRobotReceivedMessageReply`] = manage.GetRobotReceivedMessageReply
	Route[http.MethodGet][`/manage/ability/getRobotReceivedMessageReplyList`] = manage.GetRobotReceivedMessageReplyList
	Route[http.MethodPost][`/manage/ability/updateRobotReceivedMessageReplySwitchStatus`] = manage.UpdateRobotReceivedMessageReplySwitchStatus
	Route[http.MethodPost][`/manage/ability/updateRobotReceivedMessageReplyPriorityNum`] = manage.UpdateRobotReceivedMessageReplyPriorityNum

	/*robot subscribe reply API*/
	Route[http.MethodPost][`/manage/ability/saveRobotSubscribeReply`] = manage.SaveRobotSubscribeReply
	Route[http.MethodPost][`/manage/ability/deleteRobotSubscribeReply`] = manage.DeleteRobotSubscribeReply
	Route[http.MethodGet][`/manage/ability/getRobotSubscribeReply`] = manage.GetRobotSubscribeReply
	Route[http.MethodGet][`/manage/ability/getRobotSubscribeReplyList`] = manage.GetRobotSubscribeReplyList
	Route[http.MethodPost][`/manage/ability/updateRobotSubscribeReplySwitchStatus`] = manage.UpdateRobotSubscribeReplySwitchStatus
	Route[http.MethodPost][`/manage/ability/updateRobotSubscribeReplyPriorityNum`] = manage.UpdateRobotSubscribeReplyPriorityNum
}
