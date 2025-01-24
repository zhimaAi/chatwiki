// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elliotchance/orderedmap/v2"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

func GetFormList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	m := msql.Model(`form`, define.Postgres).
		Alias(`f`).
		Join(`form_entry e`, `f.id = e.form_id and e.delete_time = 0`, `left`).
		Where(`f.admin_user_id`, cast.ToString(adminUserId)).
		Where(`f.delete_time`, `0`).
		Order(`id desc`).
		Field(`f.*,count(e) as entry_count`).
		Group(`f.id`)

	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	userInfo, err := msql.Model(define.TableUser, define.Postgres).
		Alias(`u`).
		Join(`role r`, `u.user_roles::integer=r.id`, `left`).
		Where(`u.id`, cast.ToString(userId)).
		Field(`u.*,r.role_type`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(userInfo) == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	if !tool.InArrayInt(cast.ToInt(userInfo[`role_type`]), []int{define.RoleTypeRoot, define.RoleTypeAdmin}) {
		managedFormIdList := GetUserManagedData(userId, `managed_form_list`)
		m.Where(`f.id`, `in`, strings.Join(managedFormIdList, `,`))
	}

	list, err := m.Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func GetFormInfo(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	id := cast.ToInt(c.Query(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	userId := getLoginUserId(c)
	if userId <= 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	userInfo, err := msql.Model(define.TableUser, define.Postgres).
		Alias(`u`).
		Join(`role r`, `u.user_roles::integer=r.id`, `left`).
		Where(`u.id`, cast.ToString(userId)).
		Field(`u.*,r.role_type`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(userInfo) == 0 {
		common.FmtErrorWithCode(c, http.StatusUnauthorized, `user_no_login`)
		return
	}
	if !tool.InArrayInt(cast.ToInt(userInfo[`role_type`]), []int{define.RoleTypeRoot, define.RoleTypeAdmin}) {
		managedFormIdList := GetUserManagedData(userId, `managed_form_list`)
		if !tool.InArrayString(cast.ToString(id), managedFormIdList) {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
	}

	form, err := common.GetFormInfo(id, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(form, nil))
}

func SaveForm(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	id := cast.ToInt64(c.PostForm(`id`))
	name := strings.TrimSpace(c.PostForm(`name`))
	description := strings.TrimSpace(c.PostForm(`description`))
	if id < 0 || len(name) == 0 || len(description) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	matched, err := regexp.MatchString(`^[a-z][a-z0-9_]*$`, name)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if !matched || utf8.RuneCountInString(name) > 64 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `name`))))
		return
	}
	if utf8.RuneCountInString(description) > 500 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `description`))))
		return
	}
	m := msql.Model(`form`, define.Postgres)
	if id > 0 {
		form, err := m.Where(`id`, cast.ToString(id)).Find()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(form) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
		data := msql.Datas{
			`admin_user_id`: userId,
			`name`:          name,
			`description`:   description,
			`update_time`:   tool.Time2Int(),
		}
		_, err = m.Where(`id`, cast.ToString(id)).Update(data)
	} else {
		data := msql.Datas{
			`admin_user_id`: userId,
			`name`:          name,
			`description`:   description,
			`update_time`:   tool.Time2Int(),
		}
		data[`create_time`] = tool.Time2Int()
		id, err = m.Insert(data, `id`)
		_ = AddUserMangedData(getLoginUserId(c), `managed_form_list`, id)
	}
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{`id`: id}, nil))
}

func DelForm(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	id := cast.ToInt64(c.PostForm(`id`))
	if id < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	_, err := msql.Model(`form`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`id`, cast.ToString(id)).
		Update(msql.Datas{`delete_time`: tool.Time2Int()})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func GetFormFieldList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	formId := cast.ToInt(c.Query(`form_id`))
	if formId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	list, err := msql.Model(`form_field`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`form_id`, cast.ToString(formId)).
		Order(`id asc`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
	return
}

func SaveFormField(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	id := cast.ToInt(c.PostForm(`id`))
	formId := cast.ToInt(c.PostForm(`form_id`))
	name := cast.ToString(c.PostForm(`name`))
	description := cast.ToString(c.PostForm(`description`))
	_type := cast.ToString(c.PostForm(`type`))
	required := cast.ToBool(c.PostForm(`required`))
	if id < 0 || formId < 0 || len(name) == 0 || len(description) == 0 || len(_type) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	matched, err := regexp.MatchString(`^[a-z][a-z0-9_]*$`, name)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if !matched || len(name) > 64 || name == `id` {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `name`))))
		return
	}
	if utf8.RuneCountInString(description) > 64 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `description`))))
		return
	}

	if !tool.InArrayString(_type, []string{`string`, `integer`, `number`, `boolean`}) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `type`))))
		return
	}

	form, err := msql.Model(`form`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`id`, cast.ToString(formId)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(form) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	if id > 0 {
		formField, err := msql.Model(`form_field`, define.Postgres).Where(`id`, cast.ToString(id)).Find()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(formField) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
		duplicated, err := msql.Model(`form_field`, define.Postgres).
			Where(`form_id`, cast.ToString(formId)).
			Where(`name`, name).
			Where(`id`, `!=`, cast.ToString(id)).
			Count()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if duplicated > 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `duplicated_field`))))
			return
		}
		_, err = msql.Model(`form_field`, define.Postgres).Where(`id`, cast.ToString(id)).Update(msql.Datas{
			`name`:        name,
			`description`: description,
			`type`:        _type,
			`required`:    required,
			`update_time`: tool.Time2Int(),
		})
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	} else {
		duplicated, err := msql.Model(`form_field`, define.Postgres).Where(`form_id`, cast.ToString(formId)).Where(`name`, name).Count()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if duplicated > 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `duplicated_field`))))
			return
		}

		data := msql.Datas{
			`admin_user_id`: userId,
			`form_id`:       formId,
			`name`:          name,
			`description`:   description,
			`type`:          _type,
			`required`:      required,
			`create_time`:   tool.Time2Int(),
			`update_time`:   tool.Time2Int(),
		}
		_, err = msql.Model(`form_field`, define.Postgres).Insert(data)
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func UpdateFormRequired(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	id := cast.ToInt(c.PostForm(`id`))
	formId := cast.ToInt(c.PostForm(`form_id`))
	required := cast.ToBool(c.PostForm(`required`))
	if id < 0 || formId < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	form, err := msql.Model(`form`, define.Postgres).Where(`admin_user_id`, cast.ToString(userId)).Where(`id`, cast.ToString(formId)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(form) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	if _, err = msql.Model(`form_field`, define.Postgres).Where(`id`, cast.ToString(id)).Update(msql.Datas{
		`required`:    required,
		`update_time`: tool.Time2Int(),
	}); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func DelFormField(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	if id < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	formField, err := msql.Model(`form_field`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`id`, cast.ToString(id)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(formField) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}

	m := msql.Model(`form_field`, define.Postgres)
	err = m.Begin()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	_, err = msql.Model(`form_field`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`id`, cast.ToString(id)).
		Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		_ = m.Rollback()
		return
	}
	count, err := msql.Model(`form_filter_condition`, define.Postgres).
		Where(`form_field_id`, cast.ToString(id)).
		Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		_ = m.Rollback()
		return
	}
	if count > 0 {
		subSql := `not exists (
			select 1 
			from form_filter_condition 
			where form_filter_id = form_filter.id
		)`
		_, err = msql.Model(`form_filter`, define.Postgres).
			Where(`form_id`, formField[`form_id`]).
			Where(subSql).
			Delete()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			_ = m.Rollback()
			return
		}
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func GetFormEntryList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	formId := cast.ToInt(c.Query(`form_id`))
	filterId := cast.ToInt(c.Query(`filter_id`))
	page := max(1, cast.ToInt(c.Query(`page`)))
	size := max(1, cast.ToInt(c.Query(`size`)))
	if formId < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	list, total, err := common.GetFormEntryList(userId, formId, filterId, page, size, 0, 0)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	data := map[string]any{`list`: list, `total`: total, `page`: page, `size`: size}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}
func UploadFormFile(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	//get params
	formId := cast.ToInt(c.PostForm(`form_id`))
	if formId <= 0 {
		common.FmtError(c, `param_invalid`, `form_id`)
		return
	}

	fieldsList, err := msql.Model(`form_field`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`form_id`, cast.ToString(formId)).
		Order(`id asc`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		common.FmtError(c, `sys_err`)
		return
	}
	if len(fieldsList) <= 0 {
		common.FmtError(c, `no_fields`)
		return
	}
	fieldsMap := make(map[string]msql.Params, 0)
	for _, item := range fieldsList {
		fieldsMap[item[`name`]] = item
	}
	// save upload file
	uploadInfoMap := make([]map[string]any, 0)
	if c.Request.MultipartForm == nil || len(c.Request.MultipartForm.File) == 0 {
		common.FmtError(c, `upload_empty`)
		return
	}
	for _, fileHeader := range c.Request.MultipartForm.File[`form_files`] {
		if fileHeader == nil {
			common.FmtError(c, `upload_empty`)
			return
		}
		uploadInfo, err := common.ReadUploadedFile(fileHeader, define.LibFileLimitSize, define.FormFileAllowExt)
		if err != nil {
			logs.Error(err.Error())
			common.FmtError(c, `file_err`)
			return
		}
		if uploadInfo == nil || uploadInfo.Columns == "" {
			common.FmtError(c, `upload_empty`)
			return
		}
		if uploadInfo.Ext == `json` {
			if err = tool.JsonDecode(uploadInfo.Columns, &uploadInfoMap); err != nil {
				logs.Error(err.Error())
				common.FmtError(c, `file_data_err`)
				return
			}
		} else {
			splitData := strings.Split(uploadInfo.Columns, "\r\n")
			title := make([]string, 0)
			for key, item := range splitData {
				upData := strings.Split(item, ",")
				if len(upData) < len(title) {
					continue
				}
				if key == 0 {
					title = upData
					continue
				}
				var data = make(map[string]any)
				for k, v := range title {
					data[cast.ToString(v)] = upData[k]
				}
				uploadInfoMap = append(uploadInfoMap, data)
			}
		}

		if len(uploadInfoMap) <= 0 {
			common.FmtError(c, `file_data_err`)
			return
		}
		if len(uploadInfoMap) > 10000 {
			common.FmtError(c, `file_data_limits`)
			return
		}
		guuid := tool.MD5(cast.ToString(adminUserId) + cast.ToString(formId) + cast.ToString(time.Now().UnixMicro()))
		go func() {
			err = common.UploadFormFile(adminUserId, formId, guuid, uploadInfo.Ext, fieldsMap, uploadInfoMap)
		}()
		common.FmtOk(c, map[string]string{"task_id": guuid})
		return
	}
}

func GetUploadFormFileProc(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	//get params
	taskId := cast.ToString(c.PostForm(`task_id`))
	if len(taskId) <= 0 {
		common.FmtError(c, `param_invalid`, `task_id`)
		return
	}
	procData, err := common.GetUploadFormFileProc(taskId)
	if err != nil || procData == nil {
		common.FmtError(c, `param_invalid`, `task_id`)
		return
	}
	// finish
	if procData.Finish {
		procData.Success = procData.Total - len(procData.ErrData)
		common.SetUploadFormFileProc(taskId, nil, time.Duration(3))
	}
	common.FmtOk(c, procData)
}

func SaveFormEntry(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	formId := cast.ToInt(c.PostForm(`form_id`))
	formEntryId := cast.ToInt(c.PostForm(`id`))
	if formId < 0 || formEntryId < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	form, err := msql.Model(`form`, define.Postgres).Where(`id`, cast.ToString(formId)).Where(`admin_user_id`, cast.ToString(userId)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(form) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	entryValues := make(map[string]any)
	for key, values := range c.Request.Form {
		if key != `form_id` && key != `form_entry_id` && len(values) > 0 {
			entryValues[key] = values[0]
		}
	}

	err = common.SaveFormEntry(userId, formId, formEntryId, entryValues)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func DelFormEntry(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	id := cast.ToInt64(c.PostForm(`id`))
	if id < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	_, err := msql.Model(`form_entry`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`id`, cast.ToString(id)).
		Update(msql.Datas{`delete_time`: tool.Time2Int()})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func EmptyFormEntry(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	formId := cast.ToInt64(c.PostForm(`form_id`))
	if formId < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	form, err := msql.Model(`form`, define.Postgres).Where(`admin_user_id`, cast.ToString(userId)).Where(`id`, cast.ToString(formId)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(form) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	_, err = msql.Model(`form_entry`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`form_id`, cast.ToString(formId)).
		Update(msql.Datas{`delete_time`: tool.Time2Int()})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func ExportFormEntry(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	formId := cast.ToInt(c.Query(`form_id`))
	startDate := cast.ToString(c.Query(`start_date`))
	endDate := cast.ToString(c.Query(`end_date`))
	if formId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	form, err := msql.Model(`form`, define.Postgres).Where(`admin_user_id`, cast.ToString(userId)).Where(`id`, cast.ToString(formId)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, i18n.Show(common.GetLang(c), `sys_err`))
		return
	}
	if len(form) == 0 {
		c.String(http.StatusOK, i18n.Show(common.GetLang(c), `no_data`))
		return
	}
	formFields, err := msql.Model(`form_field`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`form_id`, cast.ToString(formId)).
		Order(`id asc`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, i18n.Show(common.GetLang(c), `sys_err`))
		return
	}
	if len(formFields) == 0 {
		c.String(http.StatusOK, i18n.Show(common.GetLang(c), `no_data`))
		return
	}

	var list []msql.Params
	if len(startDate) > 0 && len(endDate) > 0 {
		startTimeStamp, err := now.Parse(startDate)
		if err != nil {
			c.String(http.StatusOK, i18n.Show(common.GetLang(c), `param_invalid`, `start_date`))
			return
		}
		endTimeStamp, err := now.Parse(endDate + ` 23:59:59`)
		if err != nil {
			c.String(http.StatusOK, i18n.Show(common.GetLang(c), `param_invalid`, `end_date`))
			return
		}
		list, _, err = common.GetFormEntryList(userId, formId, 0, 1, 10000, startTimeStamp.Unix(), endTimeStamp.Unix())
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, i18n.Show(common.GetLang(c), `sys_err`))
			return
		}
	} else {
		list, _, err = common.GetFormEntryList(userId, formId, 0, 1, 10000, 0, 0)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, i18n.Show(common.GetLang(c), `sys_err`))
			return
		}
	}
	for _, item := range list {
		if _, ok := item[`id`]; ok {
			tmp := []msql.Params{{
				`name`: `id`,
				`type`: `integer`,
			}}
			tmp = append(tmp, formFields...)
			formFields = tmp
			break
		}
	}
	//if len(list) == 0 {
	//	c.String(http.StatusOK, i18n.Show(common.GetLang(c), `no_entry_data`))
	//	return
	//}

	//save file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			logs.Error(err.Error())
		}
	}()
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, i18n.Show(common.GetLang(c), `sys_err`))
		return
	}

	//set header
	columnIndex := 0
	nameToColumnMap := orderedmap.NewOrderedMap[string, any]()
	for _, field := range formFields {
		column, _ := common.IdentifierFromColumnIndex(columnIndex)
		nameToColumnMap.Set(field[`name`], column)
		cell := fmt.Sprintf(`%s1`, column)
		_ = f.SetCellValue("Sheet1", cell, field[`name`])
		columnIndex = columnIndex + 1
	}
	//set content
	for row, item := range list {
		for k, v := range item {
			column, ok := nameToColumnMap.Get(k)
			if ok {
				cell := fmt.Sprintf(`%s%d`, column, row+2)
				_ = f.SetCellValue("Sheet1", cell, v)
			}
		}
	}
	f.SetActiveSheet(index)
	filePath := fmt.Sprintf(`/tmp/%d-%d.xlsx`, formId, time.Now().UnixNano())
	if err := f.SaveAs(filePath); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, i18n.Show(common.GetLang(c), `sys_err`))
		return
	}

	//download
	c.FileAttachment(filePath, form[`name`]+".xlsx")

	//delay delete file
	go func(filePath string) {
		time.Sleep(1 * time.Minute)
		_ = os.Remove(filePath)
	}(filePath)
}

func GetFormFilterList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	formId := cast.ToInt(c.Query(`form_id`))
	if formId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	form, err := msql.Model(`form`, define.Postgres).Where(`id`, cast.ToString(formId)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(form) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	formFilters, err := msql.Model(`form_filter`, define.Postgres).Where(`form_id`, cast.ToString(formId)).Order(`sort desc,id desc`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	var formFilterIds []string
	for _, item := range formFilters {
		formFilterIds = append(formFilterIds, item[`id`])
	}
	formFilterConditions, err := msql.Model(`form_filter_condition`, define.Postgres).
		Where(`form_filter_id`, `in`, strings.Join(formFilterIds, `,`)).
		Select()
	for _, item := range formFilters {
		var conditions []define.FormFilterCondition
		for _, condition := range formFilterConditions {
			if cast.ToInt(condition[`form_filter_id`]) == cast.ToInt(item[`id`]) {
				conditions = append(conditions, define.FormFilterCondition{
					FormFieldId: cast.ToInt(condition[`form_field_id`]),
					Rule:        condition[`rule`],
					RuleValue1:  condition[`rule_value1`],
					RuleValue2:  condition[`rule_value2`],
				})
			}
		}
		if len(conditions) == 0 {
			item[`conditions`] = `[]`
			item[`entry_count`] = `0`
		} else {
			item[`conditions`], _ = tool.JsonEncode(conditions)
			entryCount, err := common.GetFormEntryCountByFilter(userId, formId, cast.ToInt(item[`type`]), conditions)
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				return
			}
			item[`entry_count`] = cast.ToString(entryCount)
		}
	}

	c.String(http.StatusOK, lib_web.FmtJson(formFilters, nil))
}

func GetFormFilterInfo(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	formId := cast.ToInt64(c.Query(`form_id`))
	id := cast.ToInt(c.Query(`id`))
	if formId <= 0 || id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	form, err := msql.Model(`form`, define.Postgres).Where(`id`, cast.ToString(formId)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(form) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	formFilter, err := msql.Model(`form_filter`, define.Postgres).Where(`form_id`, cast.ToString(formId)).Where(`id`, cast.ToString(id)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	formFilterConditions, err := msql.Model(`form_filter_condition`, define.Postgres).Where(`form_filter_id`, cast.ToString(id)).Select()
	var conditions []define.FormFilterCondition
	for _, condition := range formFilterConditions {
		if cast.ToInt(condition[`form_filter_id`]) == cast.ToInt(formFilter[`id`]) {
			conditions = append(conditions, define.FormFilterCondition{
				FormFieldId: cast.ToInt(condition[`form_field_id`]),
				Rule:        condition[`rule`],
				RuleValue1:  condition[`rule_value1`],
				RuleValue2:  condition[`rule_value2`],
			})
		}
	}
	if len(conditions) == 0 {
		formFilter[`conditions`] = `[]`
	} else {
		formFilter[`conditions`], err = tool.JsonEncode(conditions)
	}

	c.String(http.StatusOK, lib_web.FmtJson(formFilter, nil))
}

func SaveFormFilter(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	id := cast.ToInt64(c.PostForm(`id`))
	formId := cast.ToInt64(c.PostForm(`form_id`))
	name := cast.ToString(c.PostForm(`name`))
	_type := cast.ToInt64(c.PostForm(`type`))
	condition := strings.TrimSpace(c.PostForm(`condition`))
	if id < 0 || formId <= 0 || len(name) == 0 || (_type != 1 && _type != 2) || len(condition) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	var formFilterConditions []define.FormFilterCondition
	if err := json.Unmarshal([]byte(condition), &formFilterConditions); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `condition`))))
		return
	}
	duplicateMap := make(map[int]bool)
	for _, formFilterCondition := range formFilterConditions {
		if _, exists := duplicateMap[formFilterCondition.FormFieldId]; exists {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `condition`))))
			return
		}
		duplicateMap[formFilterCondition.FormFieldId] = true
	}

	form, err := msql.Model(`form`, define.Postgres).Where(`admin_user_id`, cast.ToString(userId)).Where(`id`, cast.ToString(formId)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(form) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	formFields, err := msql.Model(`form_field`, define.Postgres).Where(`form_id`, cast.ToString(formId)).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(formFields) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}

	var formFieldIdMap = make(map[int]msql.Params)
	for _, formField := range formFields {
		formFieldIdMap[cast.ToInt(formField[`id`])] = formField
	}

	for _, formFilterCondition := range formFilterConditions {
		formField, ok := formFieldIdMap[formFilterCondition.FormFieldId]
		if !ok {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `condition`))))
			return
		}
		if err := formFilterCondition.Check(formField[`type`]); err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `condition`))))
			return
		}
	}
	m := msql.Model(`form_filter`, define.Postgres)
	err = m.Begin()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if id > 0 { // update
		formFilter, err := msql.Model(`form_filter`, define.Postgres).Where(`form_id`, cast.ToString(formId)).Find()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(formFilter) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
		_, err = msql.Model(`form_filter`, define.Postgres).
			Where(`id`, cast.ToString(id)).
			Update(msql.Datas{`name`: name, `type`: _type})
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			_ = m.Rollback()
			return
		}

		_, err = msql.Model(`form_filter_condition`, define.Postgres).Where(`form_filter_id`, cast.ToString(id)).Delete()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			_ = m.Rollback()
			return
		}
		for _, formFilterCondition := range formFilterConditions {
			_, err = msql.Model(`form_filter_condition`, define.Postgres).
				Insert(msql.Datas{
					`form_filter_id`: cast.ToString(id),
					`form_field_id`:  cast.ToString(formFilterCondition.FormFieldId),
					`rule`:           cast.ToString(formFilterCondition.Rule),
					`rule_value1`:    cast.ToString(formFilterCondition.RuleValue1),
					`rule_value2`:    cast.ToString(formFilterCondition.RuleValue2),
					`create_time`:    tool.Time2Int(),
					`update_time`:    tool.Time2Int(),
				})
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				_ = m.Rollback()
				return
			}
		}
	} else { // insert
		formFilterId, err := msql.Model(`form_filter`, define.Postgres).
			Insert(msql.Datas{
				`form_id`:     formId,
				`name`:        name,
				`type`:        _type,
				`create_time`: tool.Time2Int(),
				`update_time`: tool.Time2Int(),
			}, `id`)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			_ = m.Rollback()
			return
		}

		for _, formFilterCondition := range formFilterConditions {
			_, err = msql.Model(`form_filter_condition`, define.Postgres).
				Insert(msql.Datas{
					`form_filter_id`: formFilterId,
					`form_field_id`:  formFilterCondition.FormFieldId,
					`rule`:           formFilterCondition.Rule,
					`rule_value1`:    formFilterCondition.RuleValue1,
					`rule_value2`:    formFilterCondition.RuleValue2,
					`create_time`:    tool.Time2Int(),
					`update_time`:    tool.Time2Int(),
				})
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				_ = m.Rollback()
				return
			}
		}
	}
	if err = m.Commit(); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func UpdateFormFilterEnabled(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	formId := cast.ToInt64(c.PostForm(`form_id`))
	id := cast.ToInt64(c.PostForm(`id`))
	enabled := cast.ToBool(c.PostForm(`enabled`))
	if formId < 0 || id < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	form, err := msql.Model(`form`, define.Postgres).Where(`admin_user_id`, cast.ToString(userId)).Where(`id`, cast.ToString(formId)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(form) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	_, err = msql.Model(`form_filter`, define.Postgres).Where(`form_id`, cast.ToString(formId)).Where(`id`, cast.ToString(id)).Update(msql.Datas{
		`enabled`:     enabled,
		`update_time`: tool.Time2Int(),
	})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func UpdateFormFilterSort(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	formId := cast.ToInt64(c.PostForm(`form_id`))
	filterSorts := strings.TrimSpace(c.PostForm(`filter_sort`))
	if formId < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	type filterSort struct {
		Id   int `json:"id"`
		Sort int `json:"sort"`
	}
	var filterSortsArr []filterSort
	err := json.Unmarshal([]byte(filterSorts), &filterSortsArr)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `filter_sort`))))
		return
	}

	m := msql.Model(`form_filter`, define.Postgres)
	err = m.Begin()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	for _, filterSort := range filterSortsArr {
		_, err = msql.Model(`form_filter`, define.Postgres).Where(`form_id`, cast.ToString(formId)).Where(`id`, cast.ToString(filterSort.Id)).Update(msql.Datas{
			`sort`:        filterSort.Sort,
			`update_time`: tool.Time2Int(),
		})
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			_ = m.Rollback()
			return
		}
	}
	if err = m.Commit(); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func DelFormFilter(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	formId := cast.ToInt64(c.PostForm(`form_id`))
	id := cast.ToInt64(c.PostForm(`id`))
	if formId < 0 || id < 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	form, err := msql.Model(`form`, define.Postgres).Where(`admin_user_id`, cast.ToString(userId)).Where(`id`, cast.ToString(formId)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(form) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}

	m := msql.Model(`form_filter`, define.Postgres)

	formFilter, err := m.Where(`form_id`, cast.ToString(formId)).Where(`id`, cast.ToString(id)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(formFilter) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}

	err = m.Begin()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	_, err = m.Where(`form_id`, cast.ToString(formId)).Where(`id`, cast.ToString(id)).Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		_ = m.Rollback()
		return
	}
	_, err = msql.Model(`form_filter_condition`, define.Postgres).Where(`form_filter_id`, cast.ToString(id)).Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		_ = m.Rollback()
		return
	}
	err = m.Commit()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}
