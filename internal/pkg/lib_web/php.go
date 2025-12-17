// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package lib_web

import (
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

type Response struct {
	Res  int         `json:"res"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func FmtJson(data interface{}, err error) string {
	if data == nil {
		data = ``
	}
	var obj Response
	if err == nil { //success
		obj = Response{Res: CommonSuccess, Msg: `success`, Data: data}
	} else {
		obj = Response{Res: CommonError, Msg: err.Error(), Data: data}
	}
	jsonStr, err := tool.JsonEncode(obj)
	if err != nil {
		logs.Error(`data:%v,err:%#v`, data, err)
		return `{"res":1,"msg":"system exception, please try again later","data":""}`
	}
	return jsonStr
}

func FmtJsonWithCode(code int, data interface{}, err error) string {
	if data == nil {
		data = ``
	}
	var obj Response
	if err == nil { //success
		obj = Response{Res: CommonSuccess, Msg: `success`, Data: data}
	} else {
		if code == 0 {
			code = CommonError
		}
		obj = Response{Res: code, Msg: err.Error(), Data: data}
	}
	jsonStr, err := tool.JsonEncode(obj)
	if err != nil {
		logs.Error(`data:%v,err:%#v`, data, err)
		return `{"res":1,"msg":"system exception, please try again later","data":""}`
	}
	return jsonStr
}
