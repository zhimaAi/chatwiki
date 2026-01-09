// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

// GetRobotTemplateCategoryList 获取模板分类列表
func GetRobotTemplateCategoryList(c *gin.Context) {
	resp, err := requestXiaokefu(`kf/ChatWiki/CommonGetRobotTemplateCategoryList`, map[string]any{`switch`: 1})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	respData, ok := resp.Data.([]interface{})
	if !ok {
		err = errors.New(`invalid data format`)
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(respData, nil))
}

// GetRobotTemplateList 获取模板列表
func GetRobotTemplateList(c *gin.Context) {
	body := make(map[string]any)

	keyword := strings.TrimSpace(c.Query(`keyword`))
	categoryId := cast.ToInt(c.Query(`category_id`))
	page := cast.ToInt(c.Query(`page`))
	if len(keyword) > 0 {
		body[`keyword`] = keyword
	}
	if categoryId > 0 {
		body[`category_id`] = categoryId
	}
	if page > 0 {
		body[`page`] = page
	}
	body[`switch`] = 1
	resp, err := requestXiaokefu(`kf/ChatWiki/CommonGetRobotTemplateList`, body)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(resp.Data, nil))
}

// UseRobotTemplate 使用模板
func UseRobotTemplate(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	templateId := cast.ToInt(c.PostForm(`template_id`))
	cslUrl := strings.TrimSpace(c.PostForm(`csl_url`))
	if templateId == 0 || len(cslUrl) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	// 下载csl文件
	resp, err := http.Get(cslUrl)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logs.Error(err.Error())
		}
	}(resp.Body)

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	if len(bs) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(`文件内容不能为空`)))
		return
	}
	robotCsl, err := common.ParseRobotCsl(string(bs))
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	//获取模版最新信息
	body := make(map[string]any)
	body[`template_id`] = templateId
	respTemplateInfo, infoErr := requestXiaokefu(`kf/ChatWiki/CommonGetRobotTemplateDetail`, body)
	if infoErr != nil {
		logs.Error(infoErr.Error())
		common.FmtError(c, infoErr.Error())
		return
	}
	templateInfo, ok := respTemplateInfo.Data.(map[string]any)
	if ok {
		//更新Csl中的机器人信息
		robotCsl.Robot[`robot_name`] = cast.ToString(templateInfo[`name`])
		robotCsl.Robot[`robot_intro`] = cast.ToString(templateInfo[`description`])
		robotCsl.Robot[`robot_avatar`] = cast.ToString(templateInfo[`avatar`])
	}
	// 增加使用次数
	_, err = requestXiaokefu("kf/ChatWiki/CommonUseTemplate", map[string]any{`template_id`: templateId})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	// 导入模板
	token := c.GetHeader(`token`)
	if len(token) == 0 {
		token = c.Query(`token`)
	}
	c.String(http.StatusOK, lib_web.FmtJson(ApplyRobotCsl(adminUserId, getLoginUserId(c), token, robotCsl)))
}

// requestXiaokefu 封装小客服请求接口
func requestXiaokefu(api string, data map[string]any) (lib_web.Response, error) {
	domain := define.Config.Xiaokefu[`domain`]
	body, err := tool.JsonEncode(data)
	if err != nil {
		return lib_web.Response{}, err
	}
	if len(body) == 0 {
		body = `{}`
	}
	var (
		link    string
		request *curl.Request
	)
	link = fmt.Sprintf("%s/%s", domain, api)
	request = curl.Post(link)
	for key, item := range data {
		request.Param(key, cast.ToString(item))
	}

	resp, err := request.Response()
	if err != nil {
		return lib_web.Response{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return lib_web.Response{}, fmt.Errorf(`SYSTEM ERROR:%d`, resp.StatusCode)
	}
	code := lib_web.Response{}
	// 获取响应体
	//returnBody, returnErr := ioutil.ReadAll(resp.Body)
	//logs.Error(fmt.Sprintf(`结果：%v ,错误:%v`, string(returnBody), returnErr))
	if err = request.ToJSON(&code); err != nil {
		return lib_web.Response{}, err
	}
	if code.Res != lib_web.CommonSuccess {
		return code, errors.New(code.Msg)
	}
	return code, nil
}

// CommonGetRobotTemplateDetail 获取机器人模板详情
func CommonGetRobotTemplateDetail(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}

	robotId := cast.ToInt(c.Query(`robot_id`))
	if robotId <= 0 {
		common.FmtError(c, `param_err`, `robot_id`)
		return
	}

	body := make(map[string]any)
	body[`robot_id`] = robotId
	resp, err := requestXiaokefu(`kf/ChatWiki/CommonGetRobotTemplateDetail`, body)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, err.Error())
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(resp.Data, nil))
}
