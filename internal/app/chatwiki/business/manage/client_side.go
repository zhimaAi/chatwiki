// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetClientSideLoginSwitch(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	info := map[string]int{`client_side_login_switch`: cast.ToInt(common.ClientSideNeedLogin(userId))}
	c.String(http.StatusOK, lib_web.FmtJson(info, nil))
}

func SetClientSideLoginSwitch(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	clientSideLoginSwitch := cast.ToInt(c.PostForm(`client_side_login_switch`))
	if !tool.InArrayInt(clientSideLoginSwitch, []int{define.SwitchOff, define.SwitchOn}) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `client_side_login_switch`))))
		return
	}
	//database dispose
	m := msql.Model(define.TableUser, define.Postgres)
	up := msql.Datas{
		`client_side_login_switch`: clientSideLoginSwitch,
		`update_time`:              tool.Time2Int(),
	}
	if _, err := m.Where(`id`, cast.ToString(userId)).Where(`is_deleted`, define.Normal).Update(up); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func getDomain(c *gin.Context) (string, error) {
	domain := strings.TrimSpace(c.PostForm(`domain`))
	referer := c.Request.Referer()
	if len(domain) == 0 && len(referer) > 0 {
		urlobj, err := url.Parse(referer)
		if err == nil && urlobj.Host == c.Request.Host {
			domain = urlobj.Scheme + `://` + urlobj.Host
		}
	}
	if len(domain) == 0 {
		domain = `http` + `://` + c.Request.Host
	}
	urlobj, err := url.Parse(domain)
	if err != nil {
		return ``, err
	}
	return urlobj.Scheme + `://` + urlobj.Host, nil
}

func ClientSideDownload(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	domain, err := getDomain(c)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `domain`))))
		return
	}
	version := lib_define.GetElectronVersion()
	md5 := tool.MD5(cast.ToString(userId) + domain)
	exeUrl := fmt.Sprintf(`/public/client_side/%s/%s/chatwiki.exe`, version, md5)
	zipUrl := fmt.Sprintf(`/public/client_side/%s/%s/chatwiki.zip`, version, md5)
	if !tool.IsFile(`static`+exeUrl) || !tool.IsFile(`static`+zipUrl) {
		if lib_redis.AddLock(define.Redis, define.LockPreKey+`client_side_build.`+version+md5, time.Minute) {
			//async task:build
			if message, err := tool.JsonEncode(map[string]any{`version`: version, `domain`: domain, `admin_user_id`: userId, `exe_url`: exeUrl, `zip_url`: zipUrl}); err != nil {
				logs.Error(err.Error())
			} else if err := common.AddJobs(lib_define.ClientSideBuildTopic, message); err != nil {
				logs.Error(err.Error())
			}
		}
		// unpacked front end waiting
		exeUrl = ``
		zipUrl = ``
	}
	data := map[string]any{`domain`: domain, `admin_user_id`: userId, `exe_url`: exeUrl, `zip_url`: zipUrl, `file_url`: zipUrl}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}
