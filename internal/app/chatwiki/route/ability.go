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
	Route[http.MethodGet][`/manage/ability/getSpecifyAbilityConfig`] = manage.GetSpecifyAbilityConfig

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

	/*robot smart menu API*/
	Route[http.MethodPost][`/manage/ability/saveSmartMenu`] = manage.SaveSmartMenu
	Route[http.MethodPost][`/manage/ability/deleteSmartMenu`] = manage.DeleteSmartMenu
	Route[http.MethodGet][`/manage/ability/getSmartMenu`] = manage.GetSmartMenu
	Route[http.MethodGet][`/manage/ability/getSmartMenuList`] = manage.GetSmartMenuList

	/*official custom menu API*/
	Route[http.MethodPost][`/manage/ability/saveCustomMenu`] = manage.SaveCustomMenu
	Route[http.MethodGet][`/manage/ability/getCustomMenuList`] = manage.GetCustomMenuList
	Route[http.MethodGet][`/manage/ability/syncWxMenuToShow`] = manage.SyncWxMenuToShow
	Route[http.MethodPost][`/manage/ability/closeWxMenu`] = manage.CloseWxMenu

	/* official account batch send API */
	Route[http.MethodGet][`/manage/officialAccount/getSyncDraftList`] = manage.SyncDraftList
	Route[http.MethodGet][`/manage/officialAccount/draftList`] = manage.DraftList
	Route[http.MethodGet][`/manage/officialAccount/draftGroupList`] = manage.DraftGroupList
	Route[http.MethodPost][`/manage/officialAccount/saveDraftGroup`] = manage.SaveDraftGroup
	Route[http.MethodPost][`/manage/officialAccount/deleteDraftGroup`] = manage.DeleteDraftGroup
	Route[http.MethodPost][`/manage/officialAccount/moveDraftGroup`] = manage.MoveDraftGroup
	Route[http.MethodPost][`/manage/officialAccount/createBatchSendTask`] = manage.CreateBathSendTask
	Route[http.MethodGet][`/manage/officialAccount/batchSendTaskList`] = manage.BathSendTaskList
	Route[http.MethodPost][`/manage/officialAccount/setBatchSendTaskTopStatus`] = manage.SetBatchSendTaskTopStatus
	Route[http.MethodPost][`/manage/officialAccount/deleteBatchSendTask`] = manage.DeleteBatchSendTask
	Route[http.MethodPost][`/manage/officialAccount/setBatchSendTaskOpenStatus`] = manage.SetBatchSendTaskOpenStatus
	Route[http.MethodPost][`/manage/officialAccount/changeBatchTaskCommentRule`] = manage.ChangeBatchTaskCommentRule
	Route[http.MethodPost][`/manage/officialAccount/changeBatchTaskCommentRuleStatus`] = manage.ChangeBatchTaskCommentRuleStatus

	/* official account comment check API */
	Route[http.MethodPost][`/manage/officialAccount/changeCommentStatus`] = manage.ChangeCommentStatus
	Route[http.MethodPost][`/manage/officialAccount/markElect`] = manage.MarkElect
	Route[http.MethodPost][`/manage/officialAccount/unMarkElect`] = manage.UnMarkElect
	Route[http.MethodPost][`/manage/officialAccount/replyComment`] = manage.ReplyComment
	Route[http.MethodPost][`/manage/officialAccount/deleteComment`] = manage.DeleteComment
	Route[http.MethodPost][`/manage/officialAccount/deleteCommentReply`] = manage.DeleteCommentReply
	Route[http.MethodPost][`/manage/officialAccount/saveCommentRule`] = manage.SaveCommentRule
	Route[http.MethodPost][`/manage/officialAccount/deleteCommentRule`] = manage.DeleteCommentRule
	Route[http.MethodPost][`/manage/officialAccount/changeCommentRuleStatus`] = manage.ChangeCommentRuleStatus
	Route[http.MethodGet][`/manage/officialAccount/getCommentRuleList`] = manage.GetCommentRuleList
	Route[http.MethodGet][`/manage/officialAccount/getCommentRuleInfo`] = manage.GetCommentRuleInfo
	Route[http.MethodGet][`/manage/officialAccount/getCommentList`] = manage.GetCommentList

}
