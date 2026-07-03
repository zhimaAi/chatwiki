// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func showSkillErr(c *gin.Context, errKey string, err error) {
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), errKey))))
}

func GetClawbotSkillList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	id := cast.ToInt64(c.Query(`id`))
	robotKey := ``
	if id > 0 {
		var ok bool
		robotKey, ok = common.GetClawbotRobotKey(c, adminUserId, id)
		if !ok {
			return
		}
	}
	list, err := common.GetClawbotUserSkillList(adminUserId, robotKey)
	if err != nil {
		showSkillErr(c, ``, err)
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func GetClawbotSkillInfo(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	skillId := cast.ToInt64(c.Query(`skill_id`))
	if skillId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	item, err := common.GetClawbotUserSkillInfo(adminUserId, skillId)
	if err != nil {
		showSkillErr(c, ``, err)
		return
	}
	if item.SkillId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(item, nil))
}

func UploadClawbotSkillZip(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	fileHeader, err := c.FormFile(`file`)
	if err != nil || fileHeader == nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if fileHeader.Size > define.MaxSkillZipSize {
		showSkillErr(c, `clawbot_skill_zip_too_big`, nil)
		return
	}
	if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), `.zip`) {
		showSkillErr(c, `clawbot_skill_zip_only`, nil)
		return
	}
	// save file
	if err = tool.MkDirAll(define.UploadDir); err != nil {
		showSkillErr(c, ``, err)
		return
	}
	tmpZip := filepath.Join(define.UploadDir, fmt.Sprintf(`skill_upload_%d_%s.zip`, time.Now().Unix(), tool.Random(12)))
	if err = c.SaveUploadedFile(fileHeader, tmpZip); err != nil {
		showSkillErr(c, ``, err)
		return
	}
	defer func() { _ = os.Remove(tmpZip) }()

	result, errKey, err := common.UploadAndStageUserClawbotSkillZip(adminUserId, fileHeader.Filename, tmpZip, fileHeader.Size)
	if err != nil || errKey != `` {
		showSkillErr(c, errKey, err)
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(result, nil))
}

func SaveClawbotSkill(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	skillId := cast.ToInt64(c.PostForm(`skill_id`))
	skillName := strings.TrimSpace(c.PostForm(`skill_name`))
	remarkName := strings.TrimSpace(c.PostForm(`remark_name`))
	intro := strings.TrimSpace(c.PostForm(`intro`))
	uploadKey := strings.TrimSpace(c.PostForm(`upload_key`))
	if skillId <= 0 && uploadKey == `` {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if !define.SkillNameRegexp.MatchString(skillName) {
		showSkillErr(c, `clawbot_skill_name_invalid`, nil)
		return
	}
	if skillName == define.SkillReservedName {
		showSkillErr(c, `clawbot_skill_name_reserved`, nil)
		return
	}
	if remarkName == `` || len([]rune(remarkName)) > 50 {
		showSkillErr(c, `clawbot_skill_name_invalid`, nil)
		return
	}
	if len([]rune(intro)) > 500 {
		showSkillErr(c, `clawbot_skill_name_invalid`, nil)
		return
	}
	lockKey := define.LockPreKey + `ClawbotUserSkill.` + cast.ToString(adminUserId)
	if !lib_redis.AddLock(define.Redis, lockKey, time.Minute*5) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `op_lock`))))
		return
	}
	defer lib_redis.UnLock(define.Redis, lockKey)

	item, errKey, err := common.SaveUserClawbotSkill(common.SaveUserClawbotSkillParam{
		AdminUserId: adminUserId,
		SkillId:     skillId,
		SkillName:   skillName,
		RemarkName:  remarkName,
		Intro:       intro,
		UploadKey:   uploadKey,
	})
	if err != nil || errKey != `` {
		showSkillErr(c, errKey, err)
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(item, nil))
}

func DeleteClawbotSkill(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	skillId := cast.ToInt64(c.PostForm(`skill_id`))
	if skillId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	lockKey := define.LockPreKey + `ClawbotUserSkill.` + cast.ToString(adminUserId)
	if !lib_redis.AddLock(define.Redis, lockKey, time.Minute*5) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `op_lock`))))
		return
	}
	defer lib_redis.UnLock(define.Redis, lockKey)

	errKey, err := common.DeleteUserClawbotSkill(adminUserId, skillId)
	if err != nil || errKey != `` {
		showSkillErr(c, errKey, err)
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func SaveClawbotRobotSkills(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	id := cast.ToInt64(c.PostForm(`id`))
	robotKey, ok := common.GetClawbotRobotKey(c, adminUserId, id)
	if !ok {
		return
	}
	skillIds, valid := parseClawbotSkillIds(c.PostForm(`skill_ids`))
	if !valid {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `skill_ids`))))
		return
	}
	lockKey := define.LockPreKey + `ClawbotSkill.` + robotKey
	if !lib_redis.AddLock(define.Redis, lockKey, time.Minute*5) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `op_lock`))))
		return
	}
	defer lib_redis.UnLock(define.Redis, lockKey)

	errKey, err := common.SyncClawbotRobotSkills(adminUserId, robotKey, skillIds)
	if err != nil || errKey != `` {
		showSkillErr(c, errKey, err)
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func parseClawbotSkillIds(value string) ([]int64, bool) {
	value = strings.TrimSpace(value)
	if value == `` {
		return []int64{}, true
	}
	parts := strings.Split(value, `,`)
	ids := make([]int64, 0, len(parts))
	seen := make(map[int64]struct{})
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == `` {
			continue
		}
		id := cast.ToInt64(part)
		if id <= 0 || cast.ToString(id) != part {
			return nil, false
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids, true
}
