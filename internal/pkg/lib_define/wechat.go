// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package lib_define

import (
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/tool"
)

const SignToken = `chatwiki`

const AesKey = `zVL5R9LBMtdbRiBuUGpcvi86dOBVMVtmgJl8k7iExsL`

var OpenAppid string

const SUCCESS = `success`
const SuccessJson = `{
    "ErrCode": 0,
    "ErrMsg": "Success"
}`
const SuccessXml = `<xml>
	<ErrCode>0</ErrCode>
	<ErrMsg>Success</ErrMsg>
</xml>`

const AppOfficeAccount = `official_account`
const AppMini = `mini_program`
const AppWechatKefu = `wechat_kefu`
const FeiShuRobot = `feishu_robot`

const AppYunH5 = `yun_h5`
const AppYunPc = `yun_pc`
const AppOpenApi = `yun_open_api`

var AppTypeList = []string{
	AppOfficeAccount,
	AppMini,
	AppWechatKefu,
	FeiShuRobot,
}

const PwdSetType = 1
const AuthSetType = 2

const EventAuthorized = `authorized`
const EventUnauthorized = `unauthorized`
const EventUpdateAuthorized = `updateauthorized`
const EventComponentVerifyTicket = `component_verify_ticket`

var DiscardEvents = []string{
	`TEMPLATESENDJOBFINISH`,
	`ADD_GUIDE_BUYER_RELATION_EVENT`,
}

const (
	MsgTypeText               = `text`
	MsgTypeImage              = `image`
	MsgTypeVoice              = `voice`
	MsgTypeVideo              = `video`
	MsgTypeShortVideo         = `shortvideo`
	MsgTypeMinirogrampage     = `miniprogrampage`
	MsgTypeLocation           = `location`
	MsgTypeLink               = `link`
	MsgTypeEvent              = `event`
	EventEnterSession         = `enter_session`
	EventUserEnterTempsession = `user_enter_tempsession`
	EventSubscribe            = `subscribe`
)

const (
	FeShuMsgTypeText        = `text`
	FeShuMsgTypePost        = `post`
	FeShuMsgTypeImage       = `image`
	FeShuMsgTypeFile        = `file`
	FeShuMsgTypeAudio       = `audio`
	FeShuMsgTypeMedia       = `media`
	FeShuMsgTypeSticker     = `sticker`
	FeShuMsgTypeInteractive = `interactive`
	FeShuMsgTypeShareChat   = `share_chat`
	FeShuMsgTypeShareUser   = `share_user`
	FeShuMsgTypeSystem      = `system`
)

var MsgTypeNameMap = map[string]string{
	MsgTypeText:           `文本`,
	MsgTypeImage:          `图片`,
	MsgTypeVoice:          `语音`,
	MsgTypeVideo:          `视频`,
	MsgTypeShortVideo:     `小视频`,
	MsgTypeMinirogrampage: `小程序`,
	MsgTypeLocation:       `地理位置`,
	MsgTypeLink:           `链接`,

	FeShuMsgTypeFile:        `文件`,
	FeShuMsgTypeAudio:       `音频`,
	FeShuMsgTypeMedia:       `视频`,
	FeShuMsgTypeSticker:     `表情包`,
	FeShuMsgTypeInteractive: `卡片`,
	FeShuMsgTypeShareChat:   `分享群名片`,
	FeShuMsgTypeShareUser:   `分享个人名片`,
	FeShuMsgTypeSystem:      `系统消息`,
	FeShuMsgTypePost:        `富文本`,

	MsgTypeEvent:              `事件`,
	EventEnterSession:         `进入会话`,
	EventUserEnterTempsession: `临时会话`,
	EventSubscribe:            `关注`,
}

/*
*
ADD_SCENE_SEARCH 公众号搜索
ADD_SCENE_ACCOUNT_MIGRATION 公众号迁移
ADD_SCENE_PROFILE_CARD 名片分享
ADD_SCENE_QR_CODE 扫描二维码
ADD_SCENE_PROFILE_LINK 图文页内名称点击
ADD_SCENE_PROFILE_ITEM 图文页右上角菜单
ADD_SCENE_PAID 支付后关注
ADD_SCENE_WECHAT_ADVERTISEMENT 微信广告
ADD_SCENE_REPRINT 他人转载
ADD_SCENE_LIVESTREAM 视频号直播
ADD_SCENE_CHANNELS 视频号
ADD_SCENE_WXA 小程序关注
ADD_SCENE_OTHERS 其他
*/
const (
	SubscribeSceneAddSceneSearch              = `ADD_SCENE_SEARCH`
	SubscribeSceneAddSceneAccountMigration    = `ADD_SCENE_ACCOUNT_MIGRATION`
	SubscribeSceneAddSceneProfileCard         = `ADD_SCENE_PROFILE_CARD`
	SubscribeSceneAddSceneQrCode              = `ADD_SCENE_QR_CODE`
	SubscribeSceneAddSceneProfileLink         = `ADD_SCENE_PROFILE_LINK`
	SubscribeSceneAddSceneProfileItem         = `ADD_SCENE_PROFILE_ITEM`
	SubscribeSceneAddScenePaid                = `ADD_SCENE_PAID`
	SubscribeSceneAddSceneWeChatAdvertisement = `ADD_SCENE_WECHAT_ADVERTISEMENT`
	SubscribeSceneAddSceneReprint             = `ADD_SCENE_REPRINT`
	SubscribeSceneAddSceneLivestream          = `ADD_SCENE_LIVESTREAM`
	SubscribeSceneAddSceneChannels            = `ADD_SCENE_CHANNELS`
	SubscribeSceneAddSceneWxa                 = `ADD_SCENE_WXA`
	SubscribeSceneAddSceneOthers              = `ADD_SCENE_OTHERS`
)

// SubscribeSceneNameMap 场景名称
var SubscribeSceneNameMap = map[string]string{
	SubscribeSceneAddSceneSearch:              `公众号搜索`,
	SubscribeSceneAddSceneAccountMigration:    `公众号迁移`,
	SubscribeSceneAddSceneProfileCard:         `名片分享`,
	SubscribeSceneAddSceneQrCode:              `扫描二维码`,
	SubscribeSceneAddSceneProfileLink:         `图文页内名称点击`,
	SubscribeSceneAddSceneProfileItem:         `图文页右上角菜单`,
	SubscribeSceneAddScenePaid:                `支付后关注`,
	SubscribeSceneAddSceneWeChatAdvertisement: `微信广告`,
	SubscribeSceneAddSceneReprint:             `他人转载`,
	SubscribeSceneAddSceneLivestream:          `视频号直播`,
	SubscribeSceneAddSceneChannels:            `视频号`,
	SubscribeSceneAddSceneWxa:                 `小程序关注`,
	SubscribeSceneAddSceneOthers:              `其他`,
}

// WechatAccountIsVerify 根据认证类型判断是否认证
func WechatAccountIsVerify(accountCustomerType any) bool {
	//https://developers.weixin.qq.com/doc/oplatform/openApi/miniprogram-management/basic-info-management/api_getaccountbasicinfo.html#Enum_Res__customer_type
	customerType := cast.ToInt(accountCustomerType)
	return tool.InArrayInt(customerType, []int{-1, 1, 2, 3, 4, 5, 6, 8, 9, 11, 12})
}
