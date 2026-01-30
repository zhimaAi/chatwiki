// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_web"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

var categoryList = []string{"a", "b", "c", "d", "e"}

func initCategory(adminUserId int) {
	for _, t := range categoryList {
		res, err := msql.Model(`category`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`type`, t).Find()
		if err != nil {
			logs.Error(err.Error())
		}
		if len(res) == 0 {
			data := msql.Datas{
				`admin_user_id`: adminUserId,
				`type`:          t,
				`create_time`:   tool.Time2Int(),
				`update_time`:   tool.Time2Int(),
			}
			if t == `a` {
				data[`name`] = lib_define.Choiceness
			}
			if _, err := msql.Model(`category`, define.Postgres).Insert(data); err != nil {
				logs.Error(err.Error())
			}
		} else {
			if t == `a` && len(res[`name`]) == 0 {
				if _, err := msql.Model(`category`, define.Postgres).Where(`id`, res[`id`]).Update(msql.Datas{`name`: lib_define.Choiceness}); err != nil {
					logs.Error(err.Error())
				}
			}
		}
	}
}

func GetCategoryList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	initCategory(userId)
	fileId := cast.ToInt(c.Query(`file_id`))
	libraryId := cast.ToInt(c.Query(`library_id`))
	m := msql.Model(`category`, define.Postgres).
		Alias(`c`).
		Where(`c.admin_user_id`, cast.ToString(userId)).
		Field(`c.id,c.name,c.type,count(d.*) as data_count`).
		Order(`c.type asc`).
		Group(`c.id`)
	if fileId > 0 {
		m.Join(`chat_ai_library_file_data d`, `c.id = d.category_id and d.file_id = `+cast.ToString(fileId)+` and d.isolated=false and d.delete_time =0`, `left`)
	} else if libraryId > 0 {
		m.Join(`chat_ai_library_file_data d`, `c.id = d.category_id and d.library_id = `+cast.ToString(libraryId)+` and d.delete_time =0`, `left`)
	} else {
		m.Join(`chat_ai_library_file_data d`, `c.id = d.category_id and d.isolated=false and d.delete_time =0`, `left`)
	}
	res, err := m.Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(res, nil))
}

func SaveCategory(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	type CategoryItem struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	}

	var categoryItems []CategoryItem
	jsonData := c.PostForm("data")
	if jsonData == "" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	if err := json.Unmarshal([]byte(jsonData), &categoryItems); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `data`))))
		return
	}
	if len(categoryItems) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	m := msql.Model(`category`, define.Postgres)
	err := m.Begin()
	defer func() {
		_ = m.Rollback()
	}()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	for _, item := range categoryItems {
		name := strings.TrimSpace(item.Name)
		_, err := msql.Model(`category`, define.Postgres).
			Where(`id`, cast.ToString(item.ID)).
			Where(`admin_user_id`, cast.ToString(userId)).
			Select()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		_, err = msql.Model(`category`, define.Postgres).
			Where(`id`, cast.ToString(item.ID)).
			Where(`admin_user_id`, cast.ToString(userId)).
			Update(msql.Datas{
				`name`:        name,
				`update_time`: tool.Time2Int(),
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

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}
