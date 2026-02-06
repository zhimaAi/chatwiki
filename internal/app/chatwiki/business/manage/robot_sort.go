// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
)

// GetMaxRobotNum gets max sort number
func GetMaxRobotNum(adminUserId int) (int, error) {
	maxSortNum, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Max(`sort_num`)
	return cast.ToInt(maxSortNum), err
}

// MoveRobotSort moves robot sort order
func MoveRobotSort(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	moveRobotId := cast.ToInt(c.DefaultPostForm(`move_robot_id`, `0`))
	if moveRobotId <= 0 { // Filter by app type
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `move_robot_id`, define.MaxRobotNum))))
		return
	}

	toRobotId := cast.ToInt(c.DefaultPostForm(`to_robot_id`, `0`))

	if toRobotId <= 0 { // Filter by app type
		// Move to front, get max sort_num
		maxSortNum, _ := GetMaxRobotNum(adminUserId)
		_, err := msql.Model(`chat_ai_robot`, define.Postgres).
			Where(`id`, cast.ToString(moveRobotId)).
			Update(msql.Datas{`sort_num`: maxSortNum + 1})
		if err != nil {
			common.FmtError(c, `sys_err`, err.Error())
			return
		}
		common.FmtOk(c, nil)
		return
	} else {
		// Move after specified position, sync replace position sort_num and is_top config, increment sort_num for items greater than specified position
		toRobotInfo, err := msql.Model(`chat_ai_robot`, define.Postgres).
			Where(`id`, cast.ToString(toRobotId)).
			Field(`id,is_top,sort_num`).Find()
		if err != nil {
			common.FmtError(c, `sys_err`, err.Error())
			return
		}
		// Increment sort_num for items greater than specified position
		_, err = msql.Model(`chat_ai_robot`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`sort_num`, `>=`, toRobotInfo[`sort_num`]).
			Update2(`sort_num=sort_num+1`)
		if err != nil {
			common.FmtError(c, `sys_err`, err.Error())
			return
		}
		// Sync replace position sort_num and is_top config
		_, err = msql.Model(`chat_ai_robot`, define.Postgres).
			Where(`id`, cast.ToString(moveRobotId)).
			Update(msql.Datas{`sort_num`: toRobotInfo[`sort_num`], `is_top`: toRobotInfo[`is_top`]})
		if err != nil {
			common.FmtError(c, `sys_err`, err.Error())
			return
		}
		common.FmtOk(c, nil)
		return
	}

}
