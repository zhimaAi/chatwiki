// Copyright © 2016- 2024 Sesame Network Technology all right reserved

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
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

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
	common.FmtBridgeResponse(c, list, httpStatus, err)
	return
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
	searchType := cast.ToInt(c.PostForm(`search_type`))
	rerankModelConfigID := cast.ToInt(c.PostForm(`rerank_model_config_id`))
	rerankUseModel := strings.TrimSpace(c.PostForm(`rerank_use_model`))
	rerankStatus := strings.TrimSpace(c.DefaultPostForm(`rerank_status`, `1`))
	recallType := cast.ToString(c.PostForm(`recall_type`))
	if len(libraryIds) <= 0 || len(question) == 0 || size <= 0 || similarity <= 0 || similarity > 1 || searchType == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if searchType != define.SearchTypeMixed && searchType != define.SearchTypeVector && searchType != define.SearchTypeFullText && searchType != define.SearchTypeGraph {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `search_type`))))
		return
	}
	if modelConfigId > 0 || useModel != "" {
		//check model_config_id and use_model
		config, err := common.GetModelConfigInfo(modelConfigId, userId)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		modelInfo, _ := common.GetModelInfoByDefine(config[`model_define`])
		if !tool.InArrayString(useModel, modelInfo.LlmModelList) && !common.IsMultiConfModel(config["model_define"]) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `use_model`))))
			return
		}
		if len(config) == 0 || !tool.InArrayString(common.Llm, strings.Split(config[`model_types`], `,`)) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `model_config_id`))))
			return
		}
	}
	robot := msql.Params{
		`recall_type`:   recallType,
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

	list, _, err := common.GetMatchLibraryParagraphList("", "", question, []string{}, libraryIds, size, similarity, searchType, robot)
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
		SyncOfficialLibrary(define.LangZhCn, cast.ToInt(library[`id`]))
	}
}

func SyncOfficialLibrary(lang string, libraryId int) {
	libraryInfo, err := msql.Model(`chat_ai_library`, define.Postgres).Where(`id`, cast.ToString(libraryId)).Find()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	if len(libraryInfo) == 0 {
		logs.Error("知识库不存在")
		return
	}
	if cast.ToInt(libraryInfo[`type`]) != define.OfficialLibraryType {
		logs.Error(`仅支持公众号知识库`)
		return
	}

	appInfo, err := common.GetWechatAppInfo(`app_id`, libraryInfo[`official_app_id`])
	if err != nil {
		updateSyncOfficialLibraryStatus(libraryId, define.SyncOfficialContentStatusFailed, fmt.Sprintf(`查找公众号信息失败: %v`, err.Error()))
		return
	}
	app, err := wechat.GetApplication(msql.Params{`app_type`: lib_define.AppOfficeAccount, `app_id`: libraryInfo[`official_app_id`], `app_secret`: appInfo[`app_secret`]})
	if err != nil {
		updateSyncOfficialLibraryStatus(libraryId, define.SyncOfficialContentStatusFailed, fmt.Sprintf(`没找到对应的公众号信息: %v`, err.Error()))
		return
	}

	officialApp, ok := app.(wechat.OfficialAccountInterface)
	if !ok {
		updateSyncOfficialLibraryStatus(libraryId, define.SyncOfficialContentStatusFailed, `公众号示例初始化失败`)
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
	for {
		resp, err := officialApp.GetPublishedMessageList(offset, count, 0)
		if err != nil {
			updateSyncOfficialLibraryStatus(libraryId, define.SyncOfficialContentStatusFailed, fmt.Sprintf(`请求微信公众平台获取已发布的消息列表失败: %v`, err.Error()))
			return
		}
		if len(resp.Item) == 0 {
			logs.Info("没有获取到公众号消息内容,跳过")
			break
		}
		// 处理当前批次的数据
		for _, item := range resp.Item {
			// 只处理在指定时间范围内的消息
			if startTime > 0 && item.UpdateTime < startTime {
				continue // 跳过早于开始时间的消息
			}

			for _, newsItem := range item.Content.NewsItem {
				// 检查URL地址
				parsedURL, err := netURL.Parse(newsItem.Url)
				if err != nil || parsedURL == nil {
					logs.Error("URL地址不合法: %v", err.Error())
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
					logs.Error(fmt.Sprintf("提取html图片失败:%v", err.Error()))
				}

				req := BridgeAddLibraryFileReq{}
				req.LibraryId = strconv.Itoa(libraryId)
				req.DocType = cast.ToString(define.DocTypeOfficial)
				req.OfficialArticleId = item.ArticleId
				req.OfficialArticleUpdateTime = item.UpdateTime
				req.LibraryId = cast.ToString(libraryId)
				req.Content = content

				req.Title = newsItem.Title
				_, err = addLibFile(nil, lang, cast.ToInt(libraryInfo[`admin_user_id`]), libraryId, cast.ToInt(libraryInfo[`type`]), &define.ChunkParam{}, &req)
				if err != nil {
					updateSyncOfficialLibraryStatus(libraryId, define.SyncOfficialContentStatusFailed, fmt.Sprintf(`添加知识库文件失败: %v`, err.Error()))
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
