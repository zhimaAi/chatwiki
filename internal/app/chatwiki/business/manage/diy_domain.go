// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/lib_web"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"net/url"
	"path/filepath"
	"strings"
)

type SaveDiyDomainReq struct {
	Id  int64  `form:"id" json:"id"`
	Url string `form:"url" json:"url" binding:"required"`
}

func SaveDiyDomain(c *gin.Context) {
	var (
		req         = SaveDiyDomainReq{}
		adminUserId int
		err         error
	)
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	if err = c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	if strings.Contains(req.Url, define.DefaultCustomDomain) {
		common.FmtError(c, `param_invalid`, `url`)
		return
	}
	id := req.Id
	m := msql.Model(`chat_ai_user_domain`, define.Postgres)
	if req.Id > 0 {
		_, err = m.Where(`id`, cast.ToString(req.Id)).Update(msql.Datas{
			`url`:         req.Url,
			`update_time`: tool.Time2Int(),
		})
	} else {
		id, err = m.Insert(msql.Datas{
			`url`:           req.Url,
			`admin_user_id`: adminUserId,
			`create_time`:   tool.Time2Int(),
			`update_time`:   tool.Time2Int(),
		}, `id`)
	}
	common.FmtOk(c, id)
}

type DeleteDiyDomainReq struct {
	Id int64 `form:"id" json:"id" binding:"required"`
}

func DeleteDiyDomain(c *gin.Context) {
	var (
		req         = DeleteDiyDomainReq{}
		adminUserId int
		err         error
	)
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	if err = c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	m := msql.Model(`chat_ai_user_domain`, define.Postgres)
	nums, err := m.Where(`id`, cast.ToString(req.Id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Delete()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if nums <= 0 {
		common.FmtError(c, `operate_err`)
		return
	}
	common.FmtOk(c, req.Id)
}

func DiyDomainList(c *gin.Context) {
	var (
		adminUserId int
		err         error
	)
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	m := msql.Model(`chat_ai_user_domain`, define.Postgres)
	data, err := m.Where(`admin_user_id`, cast.ToString(adminUserId)).Select()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, data)
}

type UploadCertificateReq struct {
	Id                int64  `form:"id" json:"id" binding:"required"`
	SslCertificate    string `form:"ssl_certificate" json:"ssl_certificate" binding:"required"`
	SslCertificateKey string `form:"ssl_certificate_key" json:"ssl_certificate_key" binding:"required"`
}

func UploadCertificate(c *gin.Context) {
	var (
		req         = UploadCertificateReq{}
		adminUserId int
		err         error
	)
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	if err = c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	data, err := msql.Model(`chat_ai_user_domain`, define.Postgres).Where(`id`, cast.ToString(req.Id)).Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(data) <= 0 {
		common.FmtError(c, `operate_err`)
		return
	}
	if adminUserId != cast.ToInt(data[`admin_user_id`]) {
		common.FmtError(c, `auth_no_permission`)
		return
	}
	// http request to write files
	uri, err := ParseUrl(data[`url`])
	if err != nil {
		common.FmtError(c, `param_invalid`, `url`)
		return
	}
	body := map[string]string{
		`url`:                 uri,
		`ssl_certificate`:     req.SslCertificate,
		`ssl_certificate_key`: req.SslCertificateKey,
	}
	srvAddr := define.Config.UserDomainService[`domain`] + `/manage/save_cert`
	resp := &lib_web.Response{}
	request, _ := curl.Post(srvAddr).JSONBody(body)
	err = request.ToJSON(resp)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `operate_err`)
		return
	}
	if resp.Res != define.StatusOK {
		logs.Error("resp:%#v", resp)
		common.FmtError(c, `operate_err`)
		return
	}
	msql.Model(`chat_ai_user_domain`, define.Postgres).Where(`id`, cast.ToString(req.Id)).Update(msql.Datas{`is_upload`: 1})
	common.FmtOk(c, nil)
}

type UploadCheckFileReq struct {
	Id          int64  `form:"id" json:"id" binding:"required"`
	FileName    string `form:"file_name" json:"file_name" binding:"required"`
	FileContent string `form:"file_content" json:"file_content" binding:"required"`
}

func UploadCheckFile(c *gin.Context) {
	var (
		req         = UploadCheckFileReq{}
		adminUserId int
		err         error
	)
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	if err = c.ShouldBind(&req); err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	data, err := msql.Model(`chat_ai_user_domain`, define.Postgres).Where(`id`, cast.ToString(req.Id)).Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(data) <= 0 {
		common.FmtError(c, `operate_err`)
		return
	}
	if adminUserId != cast.ToInt(data[`admin_user_id`]) {
		common.FmtError(c, `auth_no_permission`)
		return
	}
	ext := strings.ToLower(strings.TrimLeft(filepath.Ext(req.FileName), `.`))
	if !tool.InArrayString(ext, []string{`txt`, `html`}) {
		common.FmtError(c, `cert_file_ext_err`)
		return
	}
	// save
	m := msql.Model(`chat_ai_file_info`, define.Postgres)
	info, err := m.Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`file_name`, req.FileName).Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if cast.ToInt(info[`id`]) > 0 {
		_, err = m.Where(`id`, cast.ToString(info[`id`])).Update(msql.Datas{
			`file_name`:    req.FileName,
			`file_content`: req.FileContent,
			`update_time`:  tool.Time2Int(),
		})
	} else {
		_, err = m.Insert(msql.Datas{
			`admin_user_id`: adminUserId,
			`file_name`:     req.FileName,
			`file_content`:  req.FileContent,
			`create_time`:   tool.Time2Int(),
			`update_time`:   tool.Time2Int(),
		})
	}
	common.FmtOk(c, nil)
}

func ParseUrl(addr string) (string, error) {
	// parse URL
	parsedURL, err := url.Parse(addr)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}
	// get host
	host := parsedURL.Hostname()
	return host, nil
}
