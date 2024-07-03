// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package lib_web

import (
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func FmtJson(data interface{}, err error) string {
	type response struct {
		Res  int         `json:"res"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}
	if data == nil {
		data = ``
	}
	var obj response
	if err == nil { //success
		obj = response{Res: CommonSuccess, Msg: `success`, Data: data}
	} else {
		obj = response{Res: CommonError, Msg: err.Error(), Data: data}
	}
	jsonStr, err := tool.JsonEncode(obj)
	if err != nil {
		logs.Error(`data:%v,err:%#v`, data, err)
		return `{"res":1,"msg":"system exception, please try again later","data":""}`
	}
	return jsonStr
}

//
//func FmtError(c *gin.Context, msg string, params ...string) {
//	data := struct{}{}
//	err := errors.New(i18n.Show(common.GetLang(c), msg, params))
//	c.String(http.StatusOK, FmtJson(data, err))
//	c.Abort()
//}
//
//func FmtOk(c *gin.Context, data interface{}, params ...string) {
//	msg := "success"
//	err := errors.New(i18n.Show(common.GetLang(c), msg, params))
//	c.String(http.StatusOK, FmtJson(data, err))
//	c.Abort()
//}
