// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func SaveClawbotConf(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	// get params
	id := cast.ToInt64(c.PostForm(`id`))
	searchKnowledgeClose := cast.ToInt(cast.ToBool(c.PostForm(`search_knowledge_close`)))
	queryLocalDocsClose := cast.ToInt(cast.ToBool(c.PostForm(`query_local_docs_close`)))
	openAgentWriteFileTool := cast.ToInt(cast.ToBool(c.PostForm(`open_agent_write_file_tool`)))
	openAgentExecuteTool := cast.ToInt(cast.ToBool(c.PostForm(`open_agent_execute_tool`)))
	openAgentEditFileTool := cast.ToInt(cast.ToBool(c.PostForm(`open_agent_edit_file_tool`)))
	goodsLibRecommendSwitch := cast.ToInt(c.PostForm(`goods_lib_recommend_switch`))
	goodsLibRecommendGroupIds := strings.TrimSpace(c.PostForm(`goods_lib_recommend_group_ids`))
	if len(goodsLibRecommendGroupIds) > 0 && !common.CheckIds(goodsLibRecommendGroupIds) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `goods_lib_recommend_group_ids`))))
		return
	}
	// check required
	robotKey, ok := common.GetClawbotRobotKey(c, adminUserId, id)
	if !ok {
		return
	}
	// database dispose
	data := msql.Datas{
		`search_knowledge_close`:        searchKnowledgeClose,
		`query_local_docs_close`:        queryLocalDocsClose,
		`open_agent_write_file_tool`:    openAgentWriteFileTool,
		`open_agent_execute_tool`:       openAgentExecuteTool,
		`open_agent_edit_file_tool`:     openAgentEditFileTool,
		`goods_lib_recommend_switch`:    goodsLibRecommendSwitch,
		`goods_lib_recommend_group_ids`: goodsLibRecommendGroupIds,
		`update_time`:                   tool.Time2Int(),
	}
	m := msql.Model(`chat_ai_robot`, define.Postgres)
	if _, err := m.Where(`id`, cast.ToString(id)).Update(data); err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	// clear cached data
	lib_redis.DelCacheData(define.Redis, &common.RobotCacheBuildHandler{RobotKey: robotKey})
	c.String(http.StatusOK, lib_web.FmtJson(common.GetRobotInfo(robotKey)))
}

func UploadClawbotLocalDoc(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	// get params
	id := cast.ToInt64(c.PostForm(`id`))
	description := c.PostForm(`description`)
	keywords := append(c.PostFormArray(`keywords`), c.PostFormArray(`keywords[]`)...)
	// check required
	robotKey, ok := common.GetClawbotRobotKey(c, adminUserId, id)
	if !ok {
		return
	}
	fileHeader, err := c.FormFile(`file`)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	// init dirs
	common.InitClawbotDirs(robotKey)
	// add lock
	lockKey := define.LockPreKey + `ClawbotLocalDocIndex.` + robotKey
	if !lib_redis.AddLock(define.Redis, lockKey, time.Minute*5) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `op_lock`))))
		return
	}
	defer lib_redis.UnLock(define.Redis, lockKey)
	// save file
	docInfo, err := common.SaveClawbotLocalDoc(fileHeader, robotKey)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	// save index
	err = common.SaveClawbotLocalDocIndex(robotKey, *docInfo, description, keywords)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(*docInfo, nil))
}

func GetClawbotLocalDocList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	// get params
	id := cast.ToInt64(c.Query(`id`))
	// check required
	robotKey, ok := common.GetClawbotRobotKey(c, adminUserId, id)
	if !ok {
		return
	}
	// get file list
	list := make([]define.ClawbotLocalDocInfo, 0)
	entries, _ := os.ReadDir(strings.ReplaceAll(define.PrivateFileDir, `<robot_key>`, robotKey))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if entry.Name() == `index.yaml` {
			continue // do not list index files
		}
		var size int64
		var modTime time.Time
		if info, err := entry.Info(); err == nil {
			size = info.Size()
			modTime = info.ModTime()
		}
		list = append(list, define.ClawbotLocalDocInfo{
			Name: entry.Name(), Size: size, Time: modTime,
			Ext: strings.TrimLeft(path.Ext(entry.Name()), `.`),
		})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Time.After(list[j].Time)
	})
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func DeleteClawbotLocalDoc(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	// get params
	id := cast.ToInt64(c.PostForm(`id`))
	name, nameOk := common.NormalizeClawbotLocalDocName(c.PostForm(`name`))
	if !nameOk {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `name`))))
		return
	}
	ext := strings.ToLower(strings.TrimLeft(path.Ext(name), `.`))
	if !tool.InArrayString(ext, define.ClawbotLocalDocAllowExt) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(ext+` not allow`)))
		return
	}
	// check required
	robotKey, ok := common.GetClawbotRobotKey(c, adminUserId, id)
	if !ok {
		return
	}
	// add lock
	lockKey := define.LockPreKey + `ClawbotLocalDocIndex.` + robotKey
	if !lib_redis.AddLock(define.Redis, lockKey, time.Minute*5) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `op_lock`))))
		return
	}
	defer lib_redis.UnLock(define.Redis, lockKey)
	// delete index
	if err := common.DeleteClawbotLocalDocIndex(robotKey, name); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	// delete file
	filePath := strings.ReplaceAll(define.PrivateFileDir, `<robot_key>`, robotKey) + `/` + name
	if err := os.Remove(filePath); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}
