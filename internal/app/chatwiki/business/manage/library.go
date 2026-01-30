// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"bytes"
	"chatwiki/internal/app/chatwiki/middlewares"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/wechat"
	"errors"
	"fmt"
	"net/http"
	netURL "net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/go-shiori/go-readability"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"

	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
)

var officialAccountLibraryCreateLock sync.Map

func GetLibraryList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	req := BridgeLibraryListReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	list, httpStatus, err := BridgeGetLibraryList(adminUserId, userId, common.GetLang(c), &req)
	if httpStatus != 0 {
		common.FmtBridgeResponse(c, list, httpStatus, err)
		return
	}

	libIds := make([]string, 0, len(list))
	for _, one := range list {
		libIds = append(libIds, one[`id`])
	}
	schemaByLib := make(map[string][]msql.Params)
	if len(libIds) > 0 {
		schemaList, e := msql.Model(`library_meta_schema`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`library_id`, `in`, strings.Join(libIds, `,`)).
			Order(`id asc`).
			Field(`id,library_id,name,key,type,is_show`).
			Select()
		if e == nil {
			for _, s := range schemaList {
				schemaByLib[s[`library_id`]] = append(schemaByLib[s[`library_id`]], s)
			}
		}
	}
	out := make([]map[string]any, 0, len(list))
	builtinMetaSchemaList := common.GetBuiltinMetaSchemaList(common.GetLang(c))
	for _, params := range list {
		obj := make(map[string]any, len(params)+1)
		for k, v := range params {
			obj[k] = v
		}
		metaList := make([]map[string]any, 0, len(builtinMetaSchemaList)+len(schemaByLib[params[`id`]]))
		for _, b := range builtinMetaSchemaList {
			isShow := 0
			switch b.Key {
			case define.BuiltinMetaKeySource:
				isShow = cast.ToInt(params[`show_meta_source`])
			case define.BuiltinMetaKeyUpdateTime:
				isShow = cast.ToInt(params[`show_meta_update_time`])
			case define.BuiltinMetaKeyCreateTime:
				isShow = cast.ToInt(params[`show_meta_create_time`])
			case define.BuiltinMetaKeyGroup:
				isShow = cast.ToInt(params[`show_meta_group`])
			}
			metaList = append(metaList, map[string]any{
				`name`:       b.Name,
				`key`:        b.Key,
				`type`:       b.Type,
				`value`:      ``,
				`is_show`:    isShow,
				`is_builtin`: 1,
			})
		}
		for _, s := range schemaByLib[params[`id`]] {
			metaList = append(metaList, map[string]any{
				`id`:         cast.ToInt(s[`id`]),
				`library_id`: cast.ToInt(s[`library_id`]),
				`name`:       s[`name`],
				`key`:        s[`key`],
				`type`:       cast.ToInt(s[`type`]),
				`value`:      ``,
				`is_show`:    cast.ToInt(s[`is_show`]),
				`is_builtin`: 0,
			})
		}
		obj[`meta_list`] = metaList
		out = append(out, obj)
	}
	common.FmtBridgeResponse(c, out, httpStatus, err)
}

func GetLibraryInfo(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.Query(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	info, err := msql.Model(`chat_ai_library`, define.Postgres).
		Alias("a").
		Join(`chat_ai_model_config b`, "a.model_config_id=b.id", "left").
		Where(`a.id`, cast.ToString(id)).
		Where(`a.admin_user_id`, cast.ToString(userId)).
		Field(`a.*`).
		Field(`b.model_define`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	data := msql.Datas{}
	if cast.ToInt(info[`ai_chunk_size`]) == 0 {
		info[`ai_chunk_size`] = cast.ToString(define.SplitAiChunkMaxSize)
	}
	for k, v := range info {
		data[k] = v
	}
	data[`is_offline`] = false
	data[`library_key`] = common.BuildLibraryKey(cast.ToInt(data[`id`]), cast.ToInt(data[`create_time`]))
	robotInfo, err := common.GetLibraryRobotInfo(userId, id)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	data[`robot_nums`] = len(robotInfo)
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func GetLibraryRobotInfo(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.Query(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	data, err := common.GetLibraryRobotInfo(userId, id)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func CreateLibrary(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	req := BridgeCreateLibraryReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		logs.Error(`bind error %s`, err.Error())
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	logs.Debug(`CreateLibrary req %#v`, req)
	req.FileAvatar, _ = c.FormFile(`avatar`)
	list, httpStatus, err := BridgeCreateLibrary(adminUserId, userId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, list, httpStatus, err)
}

func CreateOfficialLibrary(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	req := BridgeCreateLibraryReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		logs.Error(`bind error %s`, err.Error())
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	appIdList := strings.TrimSpace(c.PostForm("app_id_list"))
	syncOfficialHistoryType := c.PostForm("sync_official_history_type")
	enableCronSyncOfficialContent := c.PostForm("enable_cron_sync_official_content")
	logs.Debug(`CreateOfficialLibrary req %#v`, req)
	req.FileAvatar, _ = c.FormFile(`avatar`)

	// 参数检查
	if !tool.InArrayInt(cast.ToInt(syncOfficialHistoryType), define.SyncOfficialHistoryTypeList[:]) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_err`, `sync_official_history_type`))))
		return
	}
	if enableCronSyncOfficialContent != `0` && enableCronSyncOfficialContent != `1` {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_err`, `enable_cron_fetch_official_content`))))
		return
	}

	appIdListArr := strings.Split(appIdList, ",")
	if len(appIdListArr) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`, `app_id_list`))))
		return
	}

	// 避免重复创建相同的公众号知识库
	key := fmt.Sprintf("official_account_library_create_lock_%d", adminUserId)
	value, _ := officialAccountLibraryCreateLock.LoadOrStore(key, &sync.Mutex{})
	lock := value.(*sync.Mutex)
	lock.Lock()
	defer lock.Unlock()

	// 重复性检查
	for _, appId := range appIdListArr {
		appInfo, err := common.GetWechatAppInfo(`app_id`, appId)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}

		old, err := msql.Model(`chat_ai_library`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`official_app_id`, appId).
			Find()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(old) > 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `official_account_library_exist`, appInfo[`app_name`]))))
			return
		}
	}

	// 一次性选择多个公众号需要创建多个公众号知识库，应该放到事务里
	m := msql.Model(`chat_ai_library`, define.Postgres)
	err = m.Begin()
	defer func() {
		err = m.Rollback()
		if err != nil {
			logs.Error(err.Error())
		}
	}()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	var libIdList []int
	for _, appId := range appIdListArr {
		appInfo, err := common.GetWechatAppInfo(`app_id`, appId)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}

		req.LibraryName = appInfo[`app_name`]
		req.LibraryIntro = appInfo[`app_name`]
		lib, _, err := BridgeCreateLibrary(adminUserId, userId, common.GetLang(c), &req)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		libIdList = append(libIdList, cast.ToInt(lib[`id`]))
		_, err = m.Where(`id`, cast.ToString(lib[`id`])).Update(msql.Datas{
			`official_app_id`:                   appId,
			`avatar`:                            appInfo[`app_avatar`],
			`sync_official_history_type`:        cast.ToInt(syncOfficialHistoryType),
			`enable_cron_sync_official_content`: cast.ToBool(enableCronSyncOfficialContent),
		})
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	}
	if err = m.Commit(); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	// 同步公众号知识库
	for _, libId := range libIdList {
		go SyncOfficialLibrary(common.GetLang(c), libId)
	}

	c.String(http.StatusOK, lib_web.FmtJson(libIdList, nil))
}

func DeleteLibrary(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	req := BridgeDeleteLibraryReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	list, httpStatus, err := BridgeDeleteLibrary(adminUserId, userId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, list, httpStatus, err)
}

func EditLibrary(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	req := BridgeEditLibraryReq{}
	err := common.RequestParamsBind(&req, c)
	if err != nil {
		common.FmtError(c, `param_err`, middlewares.GetValidateErr(req, err, common.GetLang(c)).Error())
		return
	}
	req.FileAvatar, _ = c.FormFile(`avatar`)
	list, httpStatus, err := BridgeEditLibrary(c, adminUserId, userId, common.GetLang(c), &req)
	common.FmtBridgeResponse(c, list, httpStatus, err)
}

func LibraryRecallTest(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}

	modelConfigId := cast.ToInt(c.PostForm(`model_config_id`))
	useModel := strings.TrimSpace(c.PostForm(`use_model`))
	libraryIds := cast.ToString(c.PostForm(`id`))
	question := strings.TrimSpace(c.PostForm(`question`))
	size := cast.ToInt(c.PostForm(`size`))
	similarity := cast.ToFloat64(c.PostForm(`similarity`))
	searchType := cast.ToInt(c.DefaultPostForm(`search_type`, `1`))
	rrfWeight := strings.TrimSpace(c.PostForm(`rrf_weight`))
	rerankModelConfigID := cast.ToInt(c.PostForm(`rerank_model_config_id`))
	rerankUseModel := strings.TrimSpace(c.PostForm(`rerank_use_model`))
	rerankStatus := strings.TrimSpace(c.DefaultPostForm(`rerank_status`, `1`))
	recallType := cast.ToString(c.PostForm(`recall_type`))
	if len(libraryIds) <= 0 || len(question) == 0 || size <= 0 || similarity <= 0 || similarity > 1 || searchType == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if !tool.InArrayInt(searchType, []int{define.SearchTypeMixed, define.SearchTypeVector, define.SearchTypeFullText, define.SearchTypeGraph}) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `search_type`))))
		return
	}
	if err := common.CheckRrfWeight(rrfWeight, common.GetLang(c)); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	if modelConfigId > 0 || useModel != "" {
		//check model_config_id and use_model
		if ok := common.CheckModelIsValid(userId, modelConfigId, useModel, common.Llm); !ok {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `use_model`))))
			return
		}
	}
	robot := msql.Params{
		`recall_type`:   recallType,
		`rrf_weight`:    rrfWeight,
		`admin_user_id`: cast.ToString(userId),
	}
	for _, libraryId := range strings.Split(libraryIds, `,`) {
		info, err := common.GetLibraryInfo(cast.ToInt(libraryId), userId)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(info) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
		robotName, err := msql.Model(`chat_ai_robot`, define.Postgres).Where(`rerank_status`, `1`).Where(`rerank_model_config_id`, cast.ToString(rerankModelConfigID)).Value(`robot_name`)
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		}
		if rerankModelConfigID > 0 && cast.ToInt(rerankStatus) == define.SwitchOn {
			robot[`rerank_status`] = cast.ToString(rerankStatus)
			robot[`rerank_model_config_id`] = cast.ToString(rerankModelConfigID)
			robot[`rerank_use_model`] = cast.ToString(rerankUseModel)
			robot[`robot_name`] = robotName
		}
		if searchType == define.SearchTypeGraph {
			if !cast.ToBool(info[`graph_switch`]) {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `graph is not enabled`))))
				return
			}
			robot[`admin_user_id`] = info[`admin_user_id`]
			robot[`model_config_id`] = info[`graph_model_config_id`]
			robot[`use_model`] = info[`graph_use_model`]
			robot[`id`] = strconv.Itoa(0)
		}
		if modelConfigId > 0 && useModel != "" {
			robot[`model_config_id`] = cast.ToString(modelConfigId)
			robot[`use_model`] = useModel
		}
	}

	list, _, err := common.GetMatchLibraryParagraphList(common.GetLang(c), cast.ToString(userId), lib_define.AppYunH5, question, []string{}, libraryIds, size, similarity, searchType, robot)
	for _, item := range list {
		library, err := common.GetLibraryInfo(cast.ToInt(item[`library_id`]), userId)
		if err != nil {
			logs.Error(err.Error())
		}
		item[`library_name`] = library[`library_name`]
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, err))
}

func RelationRobot(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	robotIds := strings.TrimSpace(c.PostForm(`robot_ids`))
	libraryId := cast.ToInt(c.PostForm(`library_id`))
	if len(robotIds) == 0 || libraryId <= 0 {
		common.FmtError(c, `param_lack`)
		return
	}
	data, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Field(`id,robot_key,library_ids`).
		Where(`admin_user_id`, cast.ToString(userId)).Where(`id`, `in`, robotIds).Select()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	for _, item := range data {
		ids := strings.Split(item[`library_ids`], ",")
		if !tool.InArrayString(cast.ToString(libraryId), ids) {
			ids = append(ids, cast.ToString(libraryId))
		}
		_, err = msql.Model(`chat_ai_robot`, define.Postgres).Where(`id`, cast.ToString(item[`id`])).Update(msql.Datas{`library_ids`: strings.Join(ids, ",")})
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
		lib_redis.DelCacheData(define.Redis, &common.RobotCacheBuildHandler{RobotKey: item[`robot_key`]})
	}
	common.FmtOk(c, nil)
}

func CronSyncOfficialContent() {
	libraries, err := msql.Model(`chat_ai_library`, define.Postgres).Where(`enable_cron_sync_official_content`, `1`).Select()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	for _, library := range libraries {
		SyncOfficialLibrary(define.LangEnUs, cast.ToInt(library[`id`]))
	}
}

func SyncOfficialLibrary(lang string, libraryId int) {
	libraryInfo, err := msql.Model(`chat_ai_library`, define.Postgres).Where(`id`, cast.ToString(libraryId)).Find()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	if len(libraryInfo) == 0 {
		logs.Error(`library not exist`)
		return
	}
	if cast.ToInt(libraryInfo[`type`]) != define.OfficialLibraryType {
		logs.Error(`only support official account library`)
		return
	}

	appInfo, err := common.GetWechatAppInfo(`app_id`, libraryInfo[`official_app_id`])
	if err != nil {
		updateSyncOfficialLibraryStatus(libraryId, define.SyncOfficialContentStatusFailed, fmt.Sprintf(`find official account info failed: %v`, err.Error()))
		return
	}
	app, err := wechat.GetApplication(msql.Params{`app_type`: lib_define.AppOfficeAccount, `app_id`: libraryInfo[`official_app_id`], `app_secret`: appInfo[`app_secret`]})
	if err != nil {
		updateSyncOfficialLibraryStatus(libraryId, define.SyncOfficialContentStatusFailed, fmt.Sprintf(`corresponding official account info not found: %v`, err.Error()))
		return
	}

	officialApp, ok := app.(wechat.OfficialAccountInterface)
	if !ok {
		updateSyncOfficialLibraryStatus(libraryId, define.SyncOfficialContentStatusFailed, `official account instance init failed`)
		return
	}

	// 根据 sync_official_history_type 计算时间范围
	syncType := cast.ToInt(libraryInfo[`sync_official_history_type`])
	var startTime int64
	now := time.Now()

	switch syncType {
	case define.SyncOfficialHistoryTypeHalfYear:
		startTime = now.AddDate(0, -6, 0).Unix()
	case define.SyncOfficialHistoryTypeOneYear:
		startTime = now.AddDate(-1, 0, 0).Unix()
	case define.SyncOfficialHistoryTypeThreeYear:
		startTime = now.AddDate(-3, 0, 0).Unix()
	case define.SyncOfficialHistoryTypeAll:
		startTime = 0 // 0 表示不限制开始时间
	default:
		startTime = now.AddDate(-1, 0, 0).Unix() // 默认一年
	}

	updateSyncOfficialLibraryStatus(libraryId, define.SyncOfficialContentStatusWorking, ``)

	offset := 0
	count := 20

	returnDir := filepath.Join(fmt.Sprintf(`%s/upload/chat_ai/%s/%s/%s`, filepath.Dir(define.AppRoot), libraryInfo[`admin_user_id`], `library_file`, tool.Date(`Ym`)))
	thumbFolderPath := filepath.Join(fmt.Sprintf(`/upload/chat_ai/%s/%s/%s/`, libraryInfo[`admin_user_id`], `library_file`, tool.Date(`Ym`)))

	for {
		resp, err := officialApp.GetPublishedMessageList(offset, count, 0)
		if err != nil {
			updateSyncOfficialLibraryStatus(libraryId, define.SyncOfficialContentStatusFailed, fmt.Sprintf(`request wechat official platform to get published message list failed: %v`, err.Error()))
			return
		}
		if len(resp.Item) == 0 {
			logs.Info(`no official account message content obtained, skip`)
			break
		}
		// 处理当前批次的数据
		for _, item := range resp.Item {
			// 只处理在指定时间范围内的消息
			if startTime > 0 && item.UpdateTime < startTime {
				continue // 跳过早于开始时间的消息
			}

			for _, newsItem := range item.Content.NewsItem {
				ThumbPath, err := officialApp.GetMaterial(newsItem.ThumbMediaId, returnDir, thumbFolderPath) //获取一下封面
				if err != nil {
					logs.Error("error：" + err.Error())
				}
				// 检查URL地址
				parsedURL, err := netURL.Parse(newsItem.Url)
				if err != nil || parsedURL == nil {
					logs.Error(`invalid url: %v`, err.Error())
					continue
				}

				// 去除html中无用标签
				blockTags := "</(div|p|h[1-6]|article|section|header|footer|blockquote|ul|ol|li|nav|aside)>"
				brTag := "<br[^>]*>"
				reBlock := regexp.MustCompile(blockTags)
				reBr := regexp.MustCompile(brTag)
				html := reBlock.ReplaceAllString(newsItem.Content, "$0\n")
				html = reBr.ReplaceAllString(html, "$0\n")
				article, err := readability.FromReader(bytes.NewReader([]byte(html)), parsedURL)
				if err != nil {
					logs.Error(fmt.Sprintf("failed to parse readability article: %v\n", err.Error()))
					continue
				}

				content, err := common.ProcessHTMLImages(strings.ReplaceAll(article.Content, "data-src", "src"), cast.ToInt(libraryInfo[`admin_user_id`]))
				if err != nil {
					logs.Error(fmt.Sprintf(`extract html image failed:%v`, err.Error()))
				}

				req := BridgeAddLibraryFileReq{}
				req.LibraryId = strconv.Itoa(libraryId)
				req.DocType = cast.ToString(define.DocTypeOfficial)
				req.OfficialArticleId = item.ArticleId
				req.OfficialArticleUpdateTime = item.UpdateTime
				req.LibraryId = cast.ToString(libraryId)
				req.Content = content
				req.ThumbPath = ThumbPath

				req.Title = newsItem.Title
				_, err = addLibFile(nil, lang, cast.ToInt(libraryInfo[`admin_user_id`]), libraryId, cast.ToInt(libraryInfo[`type`]), &define.ChunkParam{}, &req)
				if err != nil {
					updateSyncOfficialLibraryStatus(libraryId, define.SyncOfficialContentStatusFailed, fmt.Sprintf(`add knowledge base file failed: %v`, err.Error()))
					return
				}
			}
		}
		if len(resp.Item) < count {
			break
		}
		offset += count
	}

	updateSyncOfficialLibraryStatus(libraryId, define.SyncOfficialContentStatusNotStart, ``)
}

func updateSyncOfficialLibraryStatus(libraryId int, status int, errMsg string) {
	var data = make(msql.Datas)
	data[`sync_official_content_status`] = status
	if status == define.SyncOfficialContentStatusFailed {
		data[`sync_official_content_last_err_msg`] = errMsg
	}
	_, err := msql.Model(`chat_ai_library`, define.Postgres).Where(`id`, cast.ToString(libraryId)).Update(data)
	if err != nil {
		logs.Error(err.Error())
	}
}

func GetLibraryMetaSchemaList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	libraryId := cast.ToInt(c.Query(`library_id`))
	if libraryId <= 0 {
		common.FmtError(c, `param_lack`)
		return
	}
	libraryInfo, err := msql.Model(`chat_ai_library`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(libraryId)).
		Field(`id,show_meta_source,show_meta_update_time,show_meta_create_time,show_meta_group`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(libraryInfo) == 0 {
		common.FmtError(c, `no_data`)
		return
	}

	// 内置元数据：不落库，只由 chat_ai_library 的开关字段控制展示
	list := make([]map[string]any, 0, 8)
	builtinMetaSchemaList := common.GetBuiltinMetaSchemaList(common.GetLang(c))
	for _, b := range builtinMetaSchemaList {
		isShow := 0
		switch b.Key {
		case define.BuiltinMetaKeySource:
			isShow = cast.ToInt(libraryInfo[`show_meta_source`])
		case define.BuiltinMetaKeyUpdateTime:
			isShow = cast.ToInt(libraryInfo[`show_meta_update_time`])
		case define.BuiltinMetaKeyCreateTime:
			isShow = cast.ToInt(libraryInfo[`show_meta_create_time`])
		case define.BuiltinMetaKeyGroup:
			isShow = cast.ToInt(libraryInfo[`show_meta_group`])
		}
		list = append(list, map[string]any{
			`id`:         0,
			`name`:       b.Name,
			`key`:        b.Key,
			`type`:       b.Type,
			`is_show`:    isShow,
			`is_builtin`: 1,
		})
	}

	// 自定义元数据：来自 library_meta_schema
	customList, err := msql.Model(`library_meta_schema`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`library_id`, cast.ToString(libraryId)).
		Order(`id asc`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	for _, item := range customList {
		obj := make(map[string]any, len(item)+2)
		for k, v := range item {
			obj[k] = v
		}
		// 数字化字段：前后端统一用 int
		obj[`id`] = cast.ToInt(item[`id`])
		obj[`library_id`] = cast.ToInt(item[`library_id`])
		obj[`admin_user_id`] = cast.ToInt(item[`admin_user_id`])
		obj[`type`] = cast.ToInt(item[`type`])
		obj[`is_show`] = cast.ToInt(item[`is_show`])
		obj[`create_time`] = cast.ToInt(item[`create_time`])
		obj[`update_time`] = cast.ToInt(item[`update_time`])
		obj[`is_builtin`] = 0
		list = append(list, obj)
	}

	common.FmtOk(c, list)
}

// GetLibraryMultiMetaSchemaList 获取一个或多个知识库“共有”的元数据 schema（交集）
// 入参（GET）：
// - library_ids: 知识库ID集合（英文逗号分隔），兼容 library_id 单个 id
// 返回结构与 GetRobotMetaSchemaList 对齐：
// - 内置元数据：只返回一份（is_show=1）
// - 自定义元数据：仅返回所有指定知识库都存在的 name（按 name 求交集）
func GetLibraryMultiMetaSchemaList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	libraryIdsStr := strings.TrimSpace(c.Query(`library_ids`))
	if libraryIdsStr == `` {
		// 兼容单个
		libraryIdsStr = strings.TrimSpace(c.Query(`library_id`))
	}
	if libraryIdsStr == `` {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if !common.CheckIds(libraryIdsStr) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `library_ids`))))
		return
	}

	idArr := strings.Split(libraryIdsStr, ",")
	// 去重 & 保序不重要：这里只需要数量
	libIdSet := make(map[string]struct{}, len(idArr))
	libIds := make([]string, 0, len(idArr))
	for _, s := range idArr {
		s = strings.TrimSpace(s)
		if s == `` {
			continue
		}
		if _, ok := libIdSet[s]; ok {
			continue
		}
		libIdSet[s] = struct{}{}
		libIds = append(libIds, s)
	}
	if len(libIds) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson([]map[string]any{}, nil))
		return
	}

	existIds, err := msql.Model(`chat_ai_library`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, `in`, strings.Join(libIds, `,`)).
		ColumnArr(`id`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(existIds) != len(libIds) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}

	// 结果（按 key 去重）
	seen := make(map[string]bool)
	result := make([]map[string]any, 0, 32)

	// 内置 meta：只保留一份（与 GetRobotMetaSchemaList 一致）
	builtinMetaSchemaList := common.GetBuiltinMetaSchemaList(common.GetLang(c))
	for _, b := range builtinMetaSchemaList {
		k := b.Key
		if seen[k] {
			continue
		}
		seen[k] = true
		result = append(result, map[string]any{
			`id`:         0,
			`name`:       b.Name,
			`key`:        b.Key,
			`type`:       b.Type,
			`is_show`:    1,
			`is_builtin`: 1,
		})
	}

	// 自定义 meta：按 name 求交集（必须每个库都存在；type 必须一致）
	customList, err := msql.Model(`library_meta_schema`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`library_id`, `in`, strings.Join(libIds, `,`)).
		Order(`id asc`).
		Field(`id,library_id,name,key,type,is_show`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	type agg struct {
		count    int
		typ      int
		name     string
		key      string
		conflict bool
		seenLib  map[string]struct{}
	}
	aggMap := make(map[string]*agg, 64) // name => agg
	for _, item := range customList {
		name := strings.TrimSpace(item[`name`])
		if name == `` {
			continue
		}
		k := strings.TrimSpace(item[`key`])
		if k == `` || define.IsBuiltinMetaKey(k) {
			continue
		}
		libId := strings.TrimSpace(item[`library_id`])
		if libId == `` {
			libId = cast.ToString(item[`library_id`])
		}
		typ := cast.ToInt(item[`type`])
		a, ok := aggMap[name]
		if !ok {
			a = &agg{
				typ:     typ,
				name:    name,
				key:     k,
				seenLib: make(map[string]struct{}, len(libIds)),
			}
			aggMap[name] = a
		}
		// 同一个库重复 name 不重复计数
		if _, ok := a.seenLib[libId]; ok {
			continue
		}
		a.seenLib[libId] = struct{}{}
		a.count++
		if a.typ != typ {
			a.conflict = true
		}
	}

	for _, a := range aggMap {
		if a == nil {
			continue
		}
		if a.conflict {
			continue
		}
		if a.count != len(libIds) {
			continue
		}
		result = append(result, map[string]any{
			`id`:         0,
			`library_id`: 0,
			`name`:       a.name,
			`key`:        a.key,
			`type`:       a.typ,
			`is_show`:    1,
			`is_builtin`: 0,
		})
	}

	c.String(http.StatusOK, lib_web.FmtJson(result, nil))
}

func SaveLibraryMetaSchema(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	libraryId := cast.ToInt(c.PostForm(`library_id`))
	if libraryId <= 0 {
		common.FmtError(c, `param_lack`)
		return
	}

	// 校验知识库归属
	libraryInfo, err := msql.Model(`chat_ai_library`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(libraryId)).
		Field(`id`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(libraryInfo) == 0 {
		common.FmtError(c, `no_data`)
		return
	}

	// 内置元数据 is_show：更新到 chat_ai_library.show_meta_*
	// 约定：前端提交内置项时传 key=source/update_time/create_time/group 与 is_show
	key := strings.TrimSpace(c.PostForm(`key`))
	isShow := cast.ToInt(c.PostForm(`is_show`))
	if isShow != 0 && isShow != 1 {
		common.FmtError(c, `param_err`)
		return
	}
	if define.IsBuiltinMetaKey(key) {
		update := msql.Datas{`update_time`: tool.Time2Int()}
		switch key {
		case define.BuiltinMetaKeySource:
			update[`show_meta_source`] = isShow
		case define.BuiltinMetaKeyUpdateTime:
			update[`show_meta_update_time`] = isShow
		case define.BuiltinMetaKeyCreateTime:
			update[`show_meta_create_time`] = isShow
		case define.BuiltinMetaKeyGroup:
			update[`show_meta_group`] = isShow
		}
		_, err := msql.Model(`chat_ai_library`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`id`, cast.ToString(libraryId)).
			Update(update)
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
		common.FmtOk(c, nil)
		return
	}

	id := cast.ToInt(c.PostForm(`id`))
	name := strings.TrimSpace(c.PostForm(`name`))
	if len(name) == 0 {
		common.FmtError(c, `param_lack`)
		return
	}
	if utf8.RuneCountInString(name) > 20 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `meta_name_too_long`, 20))))
		return
	}
	typeCode := cast.ToInt(c.PostForm(`type`))
	if c.PostForm(`type`) == `` {
		typeCode = define.LibraryMetaTypeString
	}
	if !define.IsLibraryMetaTypeValid(typeCode) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `meta_type_invalid`))))
		return
	}
	// isShow 已在上方统一校验过

	now := tool.Time2Int()
	m := msql.Model(`library_meta_schema`, define.Postgres)
	data := msql.Datas{
		`name`:        name,
		`type`:        typeCode,
		`is_show`:     isShow,
		`update_time`: now,
	}

	if id > 0 {
		old, err := m.Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`library_id`, cast.ToString(libraryId)).
			Where(`id`, cast.ToString(id)).
			Field(`id,key`).
			Find()
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
		if len(old) == 0 {
			common.FmtError(c, `no_data`)
			return
		}
		// key 不可变：更新时不写 key
		_, err = m.Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`library_id`, cast.ToString(libraryId)).
			Where(`id`, cast.ToString(id)).
			Update(data)
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
		common.FmtOk(c, id)
		return
	}

	// 新增：先插入拿到 id，再把 key 更新成 key_{id}（后续不可变）
	tempKey := fmt.Sprintf("tmp_%d", time.Now().UnixNano())
	insertData := msql.Datas{
		`create_time`:   now,
		`update_time`:   now,
		`admin_user_id`: adminUserId,
		`library_id`:    libraryId,
		`name`:          name,
		`key`:           tempKey,
		`type`:          typeCode,
		`is_show`:       isShow,
	}
	newId, err := m.Insert(insertData, `id`)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	stableKey := fmt.Sprintf("key_%d", newId)
	_, err = m.Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`library_id`, cast.ToString(libraryId)).
		Where(`id`, cast.ToString(newId)).
		Update(msql.Datas{`key`: stableKey})
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, newId)
}

func DeleteLibraryMetaSchema(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `builtin_meta_cannot_delete`))))
		return
	}
	m := msql.Model(`library_meta_schema`, define.Postgres)
	item, err := m.Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(id)).
		Field(`id`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(item) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	_, err = m.Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(id)).
		Delete()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, nil)
}
