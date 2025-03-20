// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_redis"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetCatalog(c *gin.Context) {
	var adminUserId int
	//get params
	//libraryId := cast.ToInt(c.Query(`library_id`))
	libraryKey := cast.ToString(c.Query(`library_key`))
	libraryId := cast.ToInt(common.ParseLibraryKey(libraryKey))
	pid := cast.ToInt(c.Query(`pid`))
	info, err := common.GetLibraryData(libraryId)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	if cast.ToInt(info[`type`]) != define.OpenLibraryType {
		common.FmtError(c, `library_ids_err`)
		return
	}
	adminUserId = getAdminUserId(c)
	fetchAll := true
	// private access rights
	if cast.ToInt(info[`access_rights`]) != define.OpenLibraryAccessRights {
		if adminUserId <= 0 {
			common.FmtError(c, `query_no_permission`)
			return
		}
		if cast.ToInt(info[`admin_user_id`]) != adminUserId {
			common.FmtError(c, `query_no_permission`)
			return
		}
	} else {
		if cast.ToInt(info[`admin_user_id`]) != adminUserId {
			fetchAll = false
		}
	}

	docInfo := common.GetLibDocCateLog(pid, libraryId, fetchAll)
	common.FmtOk(c, docInfo)
}

func GetLibDocInfo(c *gin.Context) {
	//get params
	var adminUserId int
	//get params
	libraryKey := cast.ToString(c.Query(`library_key`))
	libraryId := cast.ToInt(common.ParseLibraryKey(libraryKey))
	docId := cast.ToInt(c.Query(`doc_id`))
	info, err := common.GetLibraryData(libraryId)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	if cast.ToInt(info[`type`]) != define.OpenLibraryType {
		common.FmtError(c, `library_ids_err`)
		return
	}
	// private access rights
	if cast.ToInt(info[`access_rights`]) != define.OpenLibraryAccessRights {
		if adminUserId = GetAdminUserId(c); adminUserId <= 0 {
			return
		}
		if cast.ToInt(info[`admin_user_id`]) != adminUserId {
			common.FmtError(c, `query_no_permission`)
			return
		}
	}
	docInfo := common.GetLibDocInfo(docId)
	if cast.ToInt(docInfo[`is_draft`]) == define.IsDraft {
		docInfo[`content`] = docInfo[`draft_content`]
	}
	common.FmtOk(c, docInfo)
}

func LibDocSearch(c *gin.Context) {
	var adminUserId int
	//get params
	libraryKey := cast.ToString(c.Query(`library_key`))
	libraryId := cast.ToInt(common.ParseLibraryKey(libraryKey))
	//search := strings.TrimSpace(c.Query(`search`))
	info, err := common.GetLibraryData(libraryId)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	if cast.ToInt(info[`type`]) != define.OpenLibraryType {
		common.FmtError(c, `library_ids_err`)
		return
	}
	// private access rights
	if cast.ToInt(info[`access_rights`]) != define.OpenLibraryAccessRights {
		if adminUserId = GetAdminUserId(c); adminUserId <= 0 {
			return
		}
		if cast.ToInt(info[`admin_user_id`]) != adminUserId {
			common.FmtError(c, `query_no_permission`)
			return
		}
	}
	//docInfo, err := common.LibDocSearch(libraryId, search, info)
	//if err != nil {
	//	common.FmtError(c, `sys_err`)
	//	return
	//}
	//common.FmtOk(c, docInfo)
}

func DraftLibDoc(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	//get params
	libraryKey := cast.ToString(c.PostForm(`library_key`))
	libraryId := cast.ToInt(common.ParseLibraryKey(libraryKey))
	pid := cast.ToInt(c.PostForm(`pid`))
	docId := cast.ToInt(c.PostForm(`doc_id`))
	isIndex := cast.ToInt(c.PostForm(`is_index`))
	title := strings.TrimSpace(c.PostForm(`title`))
	content := strings.TrimSpace(c.PostForm(`content`))
	if len(title) <= 0 {
		column := `title`
		common.FmtError(c, `param_empty`, column)
		return
	}
	if !checkPartnerRights(c, define.PartnerRightsEdit) {
		common.FmtError(c, `auth_no_permission`)
		return
	}
	info, err := common.GetLibraryInfo(libraryId, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	if cast.ToInt(info[`type`]) != define.OpenLibraryType {
		common.FmtError(c, `library_ids_err`)
		return
	}
	if pid > 0 {
		docInfo := common.GetLibDocFields(pid, `id`)
		if len(docInfo) <= 0 {
			common.FmtError(c, `param_err`, `pid`)
			return
		}
	}
	if docId > 0 {
		docInfo := common.GetLibDocInfo(docId)
		editUser := cast.ToInt(docInfo[`edit_user`])
		if editUser > 0 && editUser != getLoginUserId(c) {
			data, err := msql.Model(define.TableUser, define.Postgres).Where("id", cast.ToString(editUser)).Field(`user_name`).Find()
			if err != nil {
				common.FmtError(c, `sys_err`)
				return
			}
			if data == nil {
				common.FmtError(c, `user_not_exist`)
				return
			}
			common.FmtError(c, `edit_lock`, data[`user_name`])
			return
		}
	}
	id, err := common.SaveLibDoc(adminUserId, getLoginUserId(c), libraryId, docId, 0, pid, isIndex, define.IsDraft, title, info, content)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `doc_save_err`)
		return
	}
	common.FmtOk(c, map[string]any{"doc_id": id, `sort`: time.Now().UnixMilli()})
}

func SaveLibDoc(c *gin.Context) {
	// permission check
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	//get params
	libraryKey := cast.ToString(c.PostForm(`library_key`))
	libraryId := cast.ToInt(common.ParseLibraryKey(libraryKey))
	isIndex := cast.ToInt(c.PostForm(`is_index`))
	pid := cast.ToInt(c.PostForm(`pid`))
	docId := cast.ToInt(c.PostForm(`doc_id`))
	title := strings.TrimSpace(c.PostForm(`title`))
	content := strings.TrimSpace(c.PostForm(`content`))
	docKey := ""
	if libraryId <= 0 {
		common.FmtError(c, `library_ids_err`)
		return
	}
	if len(title) <= 0 {
		column := `title`
		common.FmtError(c, `param_empty`, column)
		return
	}
	if !checkPartnerRights(c, define.PartnerRightsEdit) {
		common.FmtError(c, `auth_no_permission`)
		return
	}
	info, err := common.GetLibraryInfo(libraryId, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	if cast.ToInt(info[`type`]) != define.OpenLibraryType {
		common.FmtError(c, `library_ids_err`)
		return
	}
	if pid > 0 {
		docPidInfo := common.GetLibDocFields(pid, `id,file_id`)
		if len(docPidInfo) <= 0 {
			common.FmtError(c, `param_err`, `pid`)
			return
		}
	}
	fileId := 0
	if docId > 0 {
		docInfo := common.GetLibDocFields(docId, `id,file_id,edit_user,doc_key,admin_user_id`)
		editUser := cast.ToInt(docInfo[`edit_user`])
		if editUser > 0 && editUser != getLoginUserId(c) {
			data, err := msql.Model(define.TableUser, define.Postgres).Where("id", cast.ToString(editUser)).Field(`user_name`).Find()
			if err != nil {
				common.FmtError(c, `sys_err`)
				return
			}
			if data == nil {
				common.FmtError(c, `user_not_exist`)
				return
			}
			common.FmtError(c, `edit_lock`, data[`user_name`])
			return
		}
		docKey = docInfo[`doc_key`]
		fileId = cast.ToInt(docInfo[`file_id`])
	}
	if isIndex != define.LibDocIndex {
		if fileId > 0 {
			err = common.DeleteLibDocRelationData(fileId)
			if err != nil {
				logs.Error(err.Error())
				common.FmtError(c, `sys_err`)
				return
			}
		}
		//common save
		fileIds, err := addLibFile(c, adminUserId, libraryId, cast.ToInt(info[`type`]))
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
		if len(fileIds) <= 0 {
			common.FmtError(c, `doc_save_err`)
			return
		}
		fileId = int(fileIds[0])
	}

	id, err := common.SaveLibDoc(adminUserId, getLoginUserId(c), libraryId, docId, fileId, pid, isIndex, 0, title, info, content)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `doc_save_err`)
		return
	}
	lib_redis.DelCacheData(define.Redis, &common.LibDocCacheBuildHandler{DocKey: docKey})
	lib_redis.DelCacheData(define.Redis, &common.LibraryCatalogCacheBuildHandle{LibraryId: libraryId})
	common.FmtOk(c, map[string]any{"doc_id": id})
}

func ChangeLibDoc(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	//get params
	libraryKey := cast.ToString(c.PostForm(`library_key`))
	libraryId := cast.ToInt(common.ParseLibraryKey(libraryKey))
	sort := cast.ToInt(c.PostForm(`sort`))
	top := cast.ToInt(c.PostForm(`top`)) // 1
	pid := cast.ToInt(c.PostForm(`pid`))
	docId := cast.ToInt(c.PostForm(`doc_id`))
	info, err := common.GetLibraryInfo(libraryId, adminUserId)
	if !checkPartnerRights(c, define.PartnerRightsEdit) {
		common.FmtError(c, `auth_no_permission`)
		return
	}
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	if cast.ToInt(info[`type`]) != define.OpenLibraryType {
		common.FmtError(c, `library_ids_err`)
		return
	}
	docInfo := common.GetLibDocFields(docId, `pid,sort`)
	if len(docInfo) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	if top == 1 {
		sort += 1
	} else if sort == 0 {
		sort = int(time.Now().UnixMilli())
	} else {
		sort -= 1
	}
	updateData := msql.Datas{
		`pid`:  pid,
		`sort`: sort,
	}
	id, err := common.ChangeLibDoc(docId, updateData)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `doc_save_err`)
		return
	}
	lib_redis.DelCacheData(define.Redis, &common.LibraryCatalogCacheBuildHandle{LibraryId: libraryId})
	common.FmtOk(c, map[string]any{"doc_id": id, `sort`: sort})
}

func SaveLibDocSeo(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	if !checkPartnerRights(c, define.PartnerRightsEdit) {
		common.FmtError(c, `auth_no_permission`)
		return
	}
	//get params
	libraryKey := cast.ToString(c.PostForm(`library_key`))
	libraryId := cast.ToInt(common.ParseLibraryKey(libraryKey))
	docId := cast.ToInt(c.PostForm(`doc_id`))
	info, err := common.GetLibraryInfo(libraryId, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	if cast.ToInt(info[`type`]) != define.OpenLibraryType {
		common.FmtError(c, `library_ids_err`)
		return
	}
	docInfo := common.GetLibDocFields(docId, `id`)
	if len(docInfo) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	updateData := msql.Datas{
		`seo_title`:    strings.TrimSpace(c.PostForm(`seo_title`)),
		`seo_desc`:     strings.TrimSpace(c.PostForm(`seo_desc`)),
		`seo_keywords`: strings.TrimSpace(c.PostForm(`seo_keywords`)),
	}
	id, err := common.ChangeLibDoc(docId, updateData)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `doc_save_err`)
		return
	}
	common.FmtOk(c, map[string]any{"doc_id": id})
}

func UploadLibDoc(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	if !checkPartnerRights(c, define.PartnerRightsEdit) {
		common.FmtError(c, `auth_no_permission`)
		return
	}
	//get params
	libraryKey := cast.ToString(c.PostForm(`library_key`))
	libraryId := cast.ToInt(common.ParseLibraryKey(libraryKey))
	pid := cast.ToInt(c.PostForm(`pid`))
	docId := cast.ToInt(c.PostForm(`doc_id`))
	if libraryId <= 0 {
		common.FmtError(c, `library_ids_err`)
		return
	}
	info, err := common.GetLibraryInfo(libraryId, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	if cast.ToInt(info[`type`]) != define.OpenLibraryType {
		common.FmtError(c, `library_ids_err`)
		return
	}
	if pid > 0 {
		docInfo := common.GetLibDocFields(pid, `id`)
		if len(docInfo) <= 0 {
			common.FmtError(c, `param_err`, `pid`)
			return
		}
	}
	//common save
	fileIds, err := addLibFile(c, adminUserId, libraryId, cast.ToInt(info[`type`]))
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(fileIds) <= 0 {
		common.FmtError(c, `upload_empty`)
		return
	}
	// save upload file
	if c.Request.MultipartForm == nil || len(c.Request.MultipartForm.File) == 0 {
		common.FmtError(c, `upload_empty`)
		return
	}
	for _, fileHeader := range c.Request.MultipartForm.File[`library_files`] {
		if fileHeader == nil {
			common.FmtError(c, `upload_empty`)
			return
		}
		uploadInfo, err := common.ReadUploadedFile(fileHeader, define.LibFileLimitSize, define.LibDocFileAllowExt)
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `file_err`)
			return
		}
		if uploadInfo == nil || uploadInfo.Columns == "" {
			common.FmtError(c, `file_data_err`)
			return
		}
		title := ""
		fileName := strings.Split(uploadInfo.Name, `.`)
		fileName = fileName[:len(fileName)-1]
		for _, item := range fileName {
			title += item
		}
		docId, err := common.SaveLibDoc(adminUserId, getLoginUserId(c), libraryId, docId, int(fileIds[0]), pid, 0, 0, title, info, uploadInfo.Columns)
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `file_data_err`)
			return
		}
		common.FmtOk(c, map[string]any{"doc_id": docId})
		return
	}
}

func ExportLibDoc(c *gin.Context) {
	//get params
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	if !checkPartnerRights(c, define.PartnerRightsEdit) {
		common.FmtError(c, `auth_no_permission`)
		return
	}
	//get params
	libraryKey := cast.ToString(c.Query(`library_key`))
	libraryId := cast.ToInt(common.ParseLibraryKey(libraryKey))
	docId := cast.ToInt(c.Query(`doc_id`))
	info, err := common.GetLibraryInfo(libraryId, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	docInfo := common.GetLibDocInfo(docId)
	if len(docInfo) <= 0 {
		common.FmtError(c, `param_err`, `doc_id`)
		return
	}
	filePath := fmt.Sprintf(`/tmp/%d-%d.md`, docId, time.Now().UnixNano())
	if err = tool.WriteFile(filePath, docInfo[`content`]); err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	//download
	c.FileAttachment(filePath, docInfo[`title`]+".md")

	//delay delete file
	go func(filePath string) {
		time.Sleep(30 * time.Minute)
		_ = os.Remove(filePath)
	}(filePath)
}

func DeleteLibDoc(c *gin.Context) {
	//get params
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	if !checkPartnerRights(c, define.PartnerRightsEdit) {
		common.FmtError(c, `auth_no_permission`)
		return
	}
	//get params
	libraryKey := cast.ToString(c.PostForm(`library_key`))
	libraryId := cast.ToInt(common.ParseLibraryKey(libraryKey))
	docId := cast.ToInt(c.PostForm(`doc_id`))
	info, err := common.GetLibraryInfo(libraryId, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	// check user operate rights creator admin rights=1

	children := common.GetLibDocAllChildren([]string{cast.ToString(docId)}, true)
	docInfo := common.GetLibDocInfo(docId)
	children = append(children, docInfo)
	for _, item := range children {
		if err := common.DeleteLibDocInfo(cast.ToInt(item[`id`]), item); err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
		if err := common.DeleteLibDocRelationData(cast.ToInt(item[`file_id`])); err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `sys_err`)
			return
		}
	}

	common.FmtOk(c, docId)
}

func SaveQuestionGuide(c *gin.Context) {
	//get params
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	if !checkPartnerRights(c, define.PartnerRightsEdit) {
		common.FmtError(c, `auth_no_permission`)
		return
	}
	//get params
	qId := cast.ToInt64(c.PostForm(`id`))
	libraryKey := cast.ToString(c.PostForm(`library_key`))
	question := cast.ToString(c.PostForm(`question`))
	libraryId := cast.ToInt(common.ParseLibraryKey(libraryKey))
	count, _ := msql.Model(`chat_ai_library_question_guide`, define.Postgres).Where(`library_id`, cast.ToString(libraryId)).Count()
	if count > 10 {
		common.FmtError(c, `over_limits`, `10`)
		return
	}
	id, err := common.SaveQuestionGuide(int64(adminUserId), int64(libraryId), qId, question)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `doc_save_err`)
		return
	}
	common.FmtOk(c, map[string]any{"qid": id})
}

func DeleteQuestionGuide(c *gin.Context) {
	//get params
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	if !checkPartnerRights(c, define.PartnerRightsEdit) {
		common.FmtError(c, `auth_no_permission`)
		return
	}
	//get params
	qId := cast.ToInt64(c.PostForm(`id`))
	libraryKey := cast.ToString(c.PostForm(`library_key`))
	libraryId := cast.ToInt(common.ParseLibraryKey(libraryKey))
	info, err := common.GetLibraryInfo(libraryId, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	_, err = msql.Model(`chat_ai_library_question_guide`, define.Postgres).Where(`id`, cast.ToString(qId)).Delete()
	if err != nil {
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, map[string]any{"qid": qId})
}

func QuestionGuideList(c *gin.Context) {
	libraryKey := cast.ToString(c.Query(`library_key`))
	libraryId := cast.ToInt(common.ParseLibraryKey(libraryKey))
	common.FmtOk(c, common.GetQuestionGuideList(libraryId))
}

func LibDocPartnerList(c *gin.Context) {
	//get params
	var (
		userId      int
		adminUserId int
	)
	if userId = getLoginUserId(c); userId == 0 {
		common.FmtError(c, `user_no_login`)
		return
	}
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	page := cast.ToInt(c.Query(`page`))
	size := cast.ToInt(c.Query(`size`))
	libraryKey := cast.ToString(c.Query(`library_key`))
	libraryId := cast.ToString(common.ParseLibraryKey(libraryKey))
	list, total, err := common.LibDocPartnerList([]string{libraryId}, page, size)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, map[string]any{
		`list`:     list,
		`total`:    total,
		`has_more`: page*size < int(total),
	})
}

func SaveLibDocPartner(c *gin.Context) {
	//get params
	var (
		userId      int
		adminUserId int
	)
	if userId = getLoginUserId(c); userId == 0 {
		common.FmtError(c, `user_no_login`)
		return
	}
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	libraryKey := cast.ToString(c.PostForm(`library_key`))
	libraryId := cast.ToInt(common.ParseLibraryKey(libraryKey))
	operateRights := cast.ToInt(c.PostForm(`operate_rights`))
	typ := cast.ToInt(c.PostForm(`type`))
	userIds := strings.Split(c.PostForm(`user_ids`), `,`)
	if !tool.InArray(operateRights, []int{define.PartnerRightsManage, define.PartnerRightsEdit}) {
		common.FmtError(c, `param_invalid`, `operate_rights`)
		return
	}
	info, err := common.GetLibraryInfo(libraryId, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	// check user operate rights creator admin rights=1
	if !checkPartnerRights(c, define.PartnerRightsManage) {
		common.FmtError(c, `auth_no_permission`)
		return
	}
	err = common.SaveLibDocPartner(userId, libraryId, operateRights, typ, userIds)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, nil)
}

func GetLibDocPartner(c *gin.Context) {
	//get params
	var (
		userId      int
		adminUserId int
	)
	if userId = getLoginUserId(c); userId == 0 {
		common.FmtError(c, `user_no_login`)
		return
	}
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	libraryId := cast.ToInt(c.Query(`library_id`))
	partner := common.GetPartnerInfo(userId, libraryId)
	common.FmtOk(c, partner)
}

func SyncPartnerHistory(c *gin.Context) {
	list, err := msql.Model(`chat_ai_library`, define.Postgres).Where(`type`, cast.ToString(define.OpenLibraryType)).Select()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	for _, item := range list {
		for _, user := range []string{item[`admin_user_id`], item[`creator`]} {
			data := common.GetPartnerInfo(cast.ToInt(user), cast.ToInt(item[`id`]))
			if len(data) > 0 {
				continue
			}
			common.SaveLibDocPartner(cast.ToInt(user), cast.ToInt(item[`id`]), define.PartnerRightsManage, 1, []string{user})
		}
	}
	common.FmtOk(c, nil)
}

func DeleteLibDocPartner(c *gin.Context) {
	//get params
	var (
		userId      int
		adminUserId int
	)
	if userId = getLoginUserId(c); userId == 0 {
		common.FmtError(c, `user_no_login`)
		return
	}
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	libraryKey := cast.ToString(c.PostForm(`library_key`))
	libraryId := cast.ToInt(common.ParseLibraryKey(libraryKey))
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		common.FmtError(c, `param_invalid`, `id`)
		return
	}
	info, err := common.GetLibraryInfo(libraryId, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(info) == 0 {
		common.FmtError(c, `no_data`)
		return
	}
	// check user operate rights creator admin rights=1
	if !checkPartnerRights(c, define.PartnerRightsManage) {
		common.FmtError(c, `auth_no_permission`)
		return
	}
	err = common.DeleteLibDocPartner(id)
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	common.FmtOk(c, nil)
}

func checkPartnerRights(c *gin.Context, rights int) bool {
	libraryKey := cast.ToString(c.PostForm(`library_key`))
	if c.Request.Method == http.MethodGet {
		libraryKey = cast.ToString(c.Query(`library_key`))
	}
	libraryId := cast.ToInt(common.ParseLibraryKey(libraryKey))
	return checkIsPartner(c, libraryId, rights)
}

func checkIsPartner(c *gin.Context, libraryId, rights int) bool {
	userId := getLoginUserId(c)
	//adminUserId := GetAdminUserId(c)
	//info, _ := common.GetLibraryInfo(libraryId, adminUserId)
	partnerInfo := common.GetPartnerInfo(userId, libraryId)
	// check user operate rights creator admin rights
	if len(partnerInfo) > 0 && cast.ToInt(partnerInfo[`operate_rights`]) >= rights {
		return true
	}
	return false
}
