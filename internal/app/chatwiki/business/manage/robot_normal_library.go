// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"database/sql"
	"errors"

	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

func EnsureRobotNormalLibrary(c *gin.Context) {
	adminID, userID := GetAdminUserId(c), getLoginUserId(c)
	if adminID == 0 || userID == 0 {
		return
	}
	robotID := cast.ToInt(c.PostForm("robot_id"))
	lang := common.GetLang(c)
	if robotID <= 0 {
		common.FmtError(c, "param_lack", "robot_id")
		return
	}
	if !common.CheckOperateRights(adminID, userID, define.IdentityTypeUser, robotID, define.ObjectTypeRobot, define.PermissionEditRights, lang) {
		common.FmtError(c, "auth_no_permission")
		return
	}
	db, err := msql.GetDB(define.Postgres)
	if err != nil {
		common.FmtError(c, "sys_err")
		return
	}
	var validationErr error
	result, err := ensureNormalLibraryTx(c.Request.Context(), db, adminID, robotID, func(name string) (msql.Datas, error) {
		modelID, model, exists := common.GetDefaultEmbeddingConfig(lang, adminID)
		if !exists {
			validationErr = errors.New(i18n.Show(lang, "param_invalid", "model_config_id"))
			return nil, validationErr
		}
		req := BridgeCreateLibraryReq{
			LibraryName: name, Type: "0", ModelConfigId: cast.ToString(modelID), UseModel: model,
			ChunkType: cast.ToString(define.ChunkTypeNormal), NormalChunkDefaultSeparatorsNo: "12,11",
			NormalChunkDefaultChunkSize: "512", NormalChunkDefaultChunkOverlap: "50",
			IsDefault: cast.ToString(define.NotDefault),
		}
		data, err := prepareLibraryCreation(adminID, userID, lang, &req)
		validationErr = err
		return data, err
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			common.FmtError(c, "no_data")
		} else if validationErr != nil {
			common.FmtBridgeResponse(c, nil, -1, validationErr)
		} else {
			logs.Error(err.Error())
			common.FmtError(c, "sys_err")
		}
		return
	}
	if result.Created {
		finishLibraryCreation(adminID, userID, result.LibraryID, result.data, true)
	}
	lib_redis.DelCacheData(define.Redis, &common.RobotCacheBuildHandler{RobotKey: result.robotKey})
	common.FmtOk(c, result)
}
