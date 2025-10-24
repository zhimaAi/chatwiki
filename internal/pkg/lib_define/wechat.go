// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package lib_define

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

const AppYunH5 = `yun_h5`
const AppYunPc = `yun_pc`
const AppOpenApi = `yun_open_api`

var AppTypeList = []string{
	AppOfficeAccount,
	AppMini,
	AppWechatKefu,
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
	MsgTypeEvent              = `event`
	EventEnterSession         = `enter_session`
	EventUserEnterTempsession = `user_enter_tempsession`
)
