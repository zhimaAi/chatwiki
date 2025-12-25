// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func TriggerList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	list, err := work_flow.TriggerList(adminUserId, common.GetLang(c))
	if err != nil {
		common.FmtError(c, err.Error())
		return
	}
	common.FmtOk(c, list)
}

func TriggerSwitch(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	switchStatus := cast.ToInt(c.PostForm("switch_status"))
	if !tool.InArrayInt(switchStatus, []int{1, 0}) {
		common.FmtError(c, `param_err`, `switch_status`)
		return
	}
	id := cast.ToInt(c.PostForm("id"))
	if id == 0 {
		common.FmtError(c, `param_err`, `id`)
		return
	}
	info, err := msql.Model(`trigger_config`, define.Postgres).Where(`id`, cast.ToString(id)).Find()
	if err != nil {
		common.FmtError(c, err.Error())
		return
	}
	_, err = msql.Model(`trigger_config`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(id)).Update(msql.Datas{
		"switch_status": switchStatus,
		"update_time":   time.Now().Unix(),
	})
	if err != nil {
		common.FmtError(c, err.Error())
		return
	}
	lib_redis.DelCacheData(define.Redis, common.TriggerConfigCacheBuildHandler{
		AdminUserId: adminUserId,
		TriggerType: info[`trigger_type`],
	})
	common.FmtOk(c, nil)
}

func GetTriggerOfficialMessage(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	list := []map[string]any{
		{
			`msg_type`: define.TriggerOfficialMessage,
			`desc`:     `私信消息`,
			`fields`:   getMessage(),
		},
		{
			`msg_type`: define.TriggerOfficialSubscribeUnScribe,
			`desc`:     `关注/取关事件`,
			`fields`:   getSubscribeUnsubscribe(),
		},
		{
			`msg_type`: define.TriggerOfficialQrCodeScan,
			`desc`:     `二维码扫描事件`,
			`fields`:   getQrcodeScan(),
		},
		{
			`msg_type`: define.TriggerOfficialMenuClick,
			`desc`:     `菜单点击事件`,
			`fields`:   getMenuClick(),
		},
	}
	robotKey := cast.ToString(c.Query(`robot_key`))
	if !common.CheckRobotKey(robotKey) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `robot_key`))))
		return
	}
	robot, err := common.GetRobotInfo(robotKey)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(robot) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	appList, err := msql.Model(`chat_ai_wechat_app`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`app_type`, lib_define.AppOfficeAccount).
		Order(`id desc`).Field(`app_name,app_id,app_avatar,account_customer_type`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	for idx, appInfo := range appList {
		accountIsVerify := lib_define.WechatAccountIsVerify(appInfo[`account_customer_type`])
		appList[idx][`account_is_verify`] = cast.ToString(accountIsVerify)
	}
	common.FmtOk(c, map[string]any{
		`messages`: list,
		`apps`:     appList,
	})
	return
}
func getMessage() []work_flow.TriggerOutputParam {
	fields := make([]work_flow.TriggerOutputParam, 0)
	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "appid",
			Typ:      common.TypString,
			Required: true,
			Desc:     "公众号id",
		},
		Variable: "global.appid",
	})
	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "ToUserName",
			Typ:      common.TypString,
			Required: true,
			Desc:     "开发者微信号",
		},
		Variable: "global.ToUserName",
	})
	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "FromUserName",
			Typ:      common.TypString,
			Required: true,
			Desc:     "发送方账号（一个OpenID）",
		},
		Variable: "global.FromUserName",
	})
	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "CreateTime",
			Typ:      common.TypString,
			Required: true,
			Desc:     "消息创建时间（整型）",
		},
		Variable: "global.CreateTime",
	})
	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "MsgType",
			Typ:      common.TypString,
			Required: true,
			Desc:     "消息类型，文本为text，图片为image",
		},
		Variable: "global.MsgType",
	})
	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "Content",
			Typ:      common.TypString,
			Required: false,
			Desc:     "文本消息内容，MsgType=text时存在",
		},
		Variable: "global.Content",
	})
	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "PicUrl",
			Typ:      common.TypString,
			Required: false,
			Desc:     "图片链接（由系统生成），MsgType=image时存在",
		},
		Variable: "global.PicUrl",
	})
	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "MediaId",
			Typ:      common.TypString,
			Required: false,
			Desc:     "图片消息媒体id，可以调用获取临时素材接口拉取数据",
		},
		Variable: "global.MediaId",
	})
	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "MsgId",
			Typ:      common.TypString,
			Required: true,
			Desc:     "消息id，64位整型",
		},
		Variable: "global.MsgId",
	})
	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "MsgDataId",
			Typ:      common.TypString,
			Required: false,
			Desc:     "消息的数据ID（消息如果来自文章时才有）",
		},
		Variable: "global.MsgDataId",
	})
	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "Idx",
			Typ:      common.TypString,
			Required: false,
			Desc:     "多图文时第几篇文章，从1开始（消息如果来自文章时才有）",
		},
		Variable: "global.Idx",
	})
	return fields
}
func getMenuClick() []work_flow.TriggerOutputParam {
	fields := make([]work_flow.TriggerOutputParam, 0)

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "appid",
			Typ:      common.TypString,
			Required: true,
			Desc:     "公众号id",
		},
		Variable: "global.appid",
	})

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "ToUserName",
			Typ:      common.TypString,
			Required: true,
			Desc:     "开发者微信号",
		},
		Variable: "global.ToUserName",
	})

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "FromUserName",
			Typ:      common.TypString,
			Required: true,
			Desc:     "发送方账号（一个OpenID）",
		},
		Variable: "global.FromUserName",
	})

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "CreateTime",
			Typ:      common.TypString,
			Required: true,
			Desc:     "消息创建时间（整型）",
		},
		Variable: "global.CreateTime",
	})

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "MsgType",
			Typ:      common.TypString,
			Required: true,
			Desc:     "消息类型，event",
		},
		Variable: "global.MsgType",
	})

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "Event",
			Typ:      common.TypString,
			Required: true,
			Desc:     "事件类型，CLICK",
		},
		Variable: "global.Event",
	})

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "EventKey",
			Typ:      common.TypString,
			Required: false,
			Desc:     "事件KEY值，与自定义菜单接口中KEY值对应",
		},
		Variable: "global.EventKey",
	})

	return fields
}
func getSubscribeUnsubscribe() []work_flow.TriggerOutputParam {
	fields := make([]work_flow.TriggerOutputParam, 0)

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "appid",
			Typ:      common.TypString,
			Required: true,
			Desc:     "公众号id",
		},
		Variable: "global.appid",
	})

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "ToUserName",
			Typ:      common.TypString,
			Required: true,
			Desc:     "开发者微信号",
		},
		Variable: "global.ToUserName",
	})

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "FromUserName",
			Typ:      common.TypString,
			Required: true,
			Desc:     "发送方账号（一个OpenID）",
		},
		Variable: "global.FromUserName",
	})

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "CreateTime",
			Typ:      common.TypString,
			Required: true,
			Desc:     "消息创建时间（整型）",
		},
		Variable: "global.CreateTime",
	})

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "MsgType",
			Typ:      common.TypString,
			Required: true,
			Desc:     "消息类型，event",
		},
		Variable: "global.MsgType",
	})

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "Event",
			Typ:      common.TypString,
			Required: true,
			Desc:     "事件类型，subscribe(订阅)、unsubscribe(取消订阅)",
		},
		Variable: "global.Event",
	})

	// EventKey 对于关注/取消关注事件通常是空的，但保留字段
	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "EventKey",
			Typ:      common.TypString,
			Required: false,
			Desc:     "事件KEY值，对于关注/取消关注事件通常为空",
		},
		Variable: "global.EventKey",
	})

	return fields
}
func getQrcodeScan() []work_flow.TriggerOutputParam {
	fields := make([]work_flow.TriggerOutputParam, 0)

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "appid",
			Typ:      common.TypString,
			Required: true,
			Desc:     "公众号id",
		},
		Variable: "global.appid",
	})

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "ToUserName",
			Typ:      common.TypString,
			Required: true,
			Desc:     "开发者微信号",
		},
		Variable: "global.ToUserName",
	})

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "FromUserName",
			Typ:      common.TypString,
			Required: true,
			Desc:     "发送方账号（一个OpenID）",
		},
		Variable: "global.FromUserName",
	})

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "CreateTime",
			Typ:      common.TypString,
			Required: true,
			Desc:     "消息创建时间 （整型）",
		},
		Variable: "global.CreateTime",
	})

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "MsgType",
			Typ:      common.TypString,
			Required: true,
			Desc:     "消息类型，event",
		},
		Variable: "global.MsgType",
	})

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "Event",
			Typ:      common.TypString,
			Required: true,
			Desc:     "事件类型，用户未关注时为subscribe， 用户已关注时为SCAN",
		},
		Variable: "global.Event",
	})

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "EventKey",
			Typ:      common.TypString,
			Required: true,
			Desc:     "事件KEY值，qrscene_为前缀，后面为二维码的场景值ID",
		},
		Variable: "global.EventKey",
	})

	fields = append(fields, work_flow.TriggerOutputParam{
		StartNodeParam: work_flow.StartNodeParam{
			Key:      "Ticket",
			Typ:      common.TypString,
			Required: true,
			Desc:     "二维码的ticket，可用来换取二维码图片",
		},
		Variable: "global.Ticket",
	})

	return fields
}
