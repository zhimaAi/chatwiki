// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

// GetPaymentSetting gets app payment settings
func GetPaymentSetting(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	robotId := cast.ToInt(c.Query(`robot_id`))
	if robotId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	res, err := msql.Model(`robot_payment_setting`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(res, nil))
}

// SavePaymentSetting saves app payment settings
func SavePaymentSetting(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	robotId := cast.ToInt(c.PostForm(`robot_id`))
	tryCount := cast.ToInt(c.PostForm(`try_count`))
	packageType := cast.ToInt(c.PostForm(`package_type`))
	countPackage := strings.TrimSpace(c.PostForm(`count_package`))
	durationPackage := strings.TrimSpace(c.PostForm(`duration_package`))
	contactQrcode := strings.TrimSpace(c.PostForm(`contact_qrcode`))
	packagePoster := strings.TrimSpace(c.PostForm(`package_poster`))
	old, err := msql.Model(`robot_payment_setting`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if robotId == 0 || tryCount < 0 || (packageType != define.RobotPaymentPackageTypeCount && packageType != define.RobotPaymentPackageTypeDuration) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if tryCount < 1 || tryCount > 50 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `try_count`))))
	}
	if packageType == define.RobotPaymentPackageTypeCount {
		if len(countPackage) <= 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `count_package`))))
			return
		}
		var countPackageInfoList []*define.RobotPaymentCountPackage
		err = json.Unmarshal([]byte(countPackage), &countPackageInfoList)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `count_package`))))
			return
		}
		for index, countPackageInfo := range countPackageInfoList {
			if countPackageInfo.Id == 0 {
				countPackageInfo.Id = 20000 + index
			}
			if len(countPackageInfo.Name) == 0 {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `count_package`))))
				return
			}
			if countPackageInfo.Count <= 0 || countPackageInfo.Price <= 0 {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `count_package`))))
				return
			}
			t, err := json.Marshal(countPackageInfo)
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				return
			}
			countPackage = string(t)
		}
		t, err := json.Marshal(countPackageInfoList)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		countPackage = string(t)
	} else {
		if len(durationPackage) <= 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `duration_package`))))
			return
		}
		var durationPackageInfoList []*define.RobotPaymentDurationPackage
		err = json.Unmarshal([]byte(durationPackage), &durationPackageInfoList)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `duration_package`))))
			return
		}
		for index, durationPackageInfo := range durationPackageInfoList {
			if cast.ToInt(durationPackageInfo.Id) == 0 {
				durationPackageInfo.Id = 10000 + index
			}
			if len(durationPackageInfo.Name) == 0 {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `duration_package`))))
				return
			}
			if durationPackageInfo.Count <= 0 || durationPackageInfo.Price <= 0 {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `duration_package`))))
				return
			}
			if durationPackageInfo.Duration <= 0 {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `duration_package`))))
				return
			}
		}
		t, err := json.Marshal(durationPackageInfoList)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		durationPackage = string(t)
	}

	if len(old) > 0 {
		_, err = msql.Model(`robot_payment_setting`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(userId)).
			Where(`robot_id`, cast.ToString(robotId)).
			Update(msql.Datas{
				`package_type`:     packageType,
				`try_count`:        tryCount,
				`count_package`:    countPackage,
				`duration_package`: durationPackage,
				`contact_qrcode`:   contactQrcode,
				`package_poster`:   packagePoster,
				`update_time`:      tool.Time2Int(),
			})
	} else {
		_, err = msql.Model(`robot_payment_setting`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(userId)).
			Where(`robot_id`, cast.ToString(robotId)).
			Insert(msql.Datas{
				`admin_user_id`:    userId,
				`robot_id`:         robotId,
				`package_type`:     packageType,
				`try_count`:        tryCount,
				`count_package`:    countPackage,
				`duration_package`: durationPackage,
				`package_poster`:   packagePoster,
				`create_time`:      tool.Time2Int(),
				`update_time`:      tool.Time2Int(),
			})
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, err))
}

// CopyPaymentSetting copies settings to other apps
func CopyPaymentSetting(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	fromRobotId := cast.ToInt(c.PostForm(`from_robot_id`))
	toRobotId := cast.ToInt(c.PostForm(`to_robot_id`))
	setting, err := msql.Model(`robot_payment_setting`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`robot_id`, cast.ToString(fromRobotId)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
	}
	if len(setting) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `config_current_robot`))))
		return
	}
	toRobotList, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`id`, `in`, "'"+cast.ToString(toRobotId)+"'").
		Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
	}
	if len(toRobotList) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `to_robot_id`))))
		return
	}

	m := msql.Model(`robot_payment_setting`, define.Postgres)
	err = m.Begin()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	for _, toRobot := range toRobotList {
		old, err := msql.Model(`robot_payment_setting`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(userId)).
			Where(`robot_id`, toRobot[`id`]).
			Find()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			_ = m.Rollback()
			return
		}
		if len(old) == 0 {
			_, err = msql.Model(`robot_payment_setting`, define.Postgres).Insert(msql.Datas{
				`admin_user_id`:    userId,
				`robot_id`:         toRobotId,
				`try_count`:        setting[`try_count`],
				`package_type`:     setting[`package_type`],
				`count_package`:    setting[`count_package`],
				`duration_package`: setting[`duration_package`],
				`contact_qrcode`:   setting[`contact_qrcode`],
				`create_time`:      tool.Time2Int(),
				`update_time`:      tool.Time2Int(),
			})
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				_ = m.Rollback()
				return
			}
		} else {
			_, err = msql.Model(`robot_payment_setting`, define.Postgres).Where(`id`, old[`id`]).Update(msql.Datas{
				`try_count`:        setting[`try_count`],
				`package_type`:     setting[`package_type`],
				`count_package`:    setting[`count_package`],
				`duration_package`: setting[`duration_package`],
				`contact_qrcode`:   setting[`contact_qrcode`],
				`update_time`:      tool.Time2Int(),
			})
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				_ = m.Rollback()
				return
			}
		}
	}
	err = m.Commit()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

// GetAuthCodeList gets authorization code list
func GetAuthCodeList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	page := cast.ToInt(c.Query(`page`))
	size := cast.ToInt(c.Query(`size`))
	robotId := cast.ToInt(c.Query(`robot_id`))
	content := strings.TrimSpace(c.Query(`content`))
	openid := strings.TrimSpace(c.Query(`openid`))
	usageStatus := cast.ToInt(c.Query(`usage_status`))
	if robotId == 0 || page == 0 || size == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	m := msql.Model(`robot_payment_auth_code`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`robot_id`, cast.ToString(robotId))
	if len(content) > 0 {
		m = m.Where(`content`, content)
	}
	if len(openid) > 0 {
		m = m.Where(`exchanger_openid`, openid)
	}
	if usageStatus > 0 {
		m = m.Where(`usage_status`, cast.ToString(usageStatus))
	}

	list, total, err := m.Order(`id desc`).Paginate(page, size)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	data := map[string]any{`list`: list, `total`: total, `page`: page, `size`: size}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

// AddAuthCode adds authorization code
func AddAuthCode(c *gin.Context) {
	var userId, adminUserId int
	if userId = getLoginUserId(c); userId == 0 {
		return
	}
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	robotId := cast.ToInt(c.PostForm(`robot_id`))
	packageId := cast.ToInt(c.PostForm(`package_id`))
	count := cast.ToInt(c.PostForm(`count`))
	remark := strings.TrimSpace(c.PostForm(`remark`))
	if robotId == 0 || packageId == 0 || count <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	// Get configuration
	setting, err := msql.Model(`robot_payment_setting`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(setting) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `robot_id`))))
		return
	}

	var packageName string
	var packageDuration int
	var packageCount int
	var packagePrice float32
	if cast.ToInt(setting[`package_type`]) == define.RobotPaymentPackageTypeCount {
		var countPackageInfoList []*define.RobotPaymentCountPackage
		err = json.Unmarshal([]byte(setting[`count_package`]), &countPackageInfoList)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		for _, countPackageInfo := range countPackageInfoList {
			if countPackageInfo.Id == packageId {
				packageName = countPackageInfo.Name
				packageCount = countPackageInfo.Count
				packagePrice = countPackageInfo.Price
				break
			}
		}
	} else {
		var durationPackageInfoList []define.RobotPaymentDurationPackage
		err = json.Unmarshal([]byte(setting[`duration_package`]), &durationPackageInfoList)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		for _, durationPackageInfo := range durationPackageInfoList {
			if durationPackageInfo.Id == packageId {
				packageName = durationPackageInfo.Name
				packageDuration = durationPackageInfo.Duration
				packageCount = durationPackageInfo.Count
				packagePrice = durationPackageInfo.Price
				break
			}
		}
	}
	if len(packageName) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `package_not_found`))))
		return
	}

	m := msql.Model(`robot_payment_auth_code`, define.Postgres)
	err = m.Begin()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	defer func() {
		err = m.Rollback()
		if err != nil {
			logs.Error(err.Error())
		}
	}()

	creatorName := ``
	creatorInfo := common.GetUserInfo(userId)
	if len(creatorInfo) > 0 && len(creatorInfo[`user_name`]) > 0 {
		creatorName = creatorInfo[`user_name`]
	}

	var result []*msql.Datas
	for range count {
		content := `###`
		if packageDuration > 0 {
			content += fmt.Sprintf(`%dD`, packageDuration)
		}
		if packageCount > 0 {
			content += fmt.Sprintf(`%dC`, packageCount)
		}
		content += tool.Random(15) + "###"
		item := msql.Datas{
			`admin_user_id`:    adminUserId,
			`creator_id`:       userId,
			`creator_name`:     creatorName,
			`robot_id`:         robotId,
			`content`:          content,
			`package_type`:     setting[`package_type`],
			`package_id`:       packageId,
			`package_name`:     packageName,
			`package_duration`: packageDuration,
			`package_count`:    packageCount,
			`package_price`:    packagePrice,
			`usage_status`:     define.RobotPaymentAuthCodeUsageStatusPending,
			`remark`:           remark,
			`create_date`:      time.Now().Format(`2006-01-02`),
			`exchange_date`:    ``,
			`use_time`:         0,
			`use_date`:         ``,
			`create_time`:      tool.Time2Int(),
			`update_time`:      tool.Time2Int(),
		}
		id, err := msql.Model(`robot_payment_auth_code`, define.Postgres).Insert(item)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			return
		}
		item[`id`] = id
		result = append(result, &item)
	}
	if err = m.Commit(); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(result, nil))
}

// DeleteAuthCode deletes authorization code
func DeleteAuthCode(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	_, err := msql.Model(`robot_payment_auth_code`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(id)).
		Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

// GetAuthCodeStats gets authorization code statistics
func GetAuthCodeStats(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	robotId := cast.ToInt(c.Query(`robot_id`))
	if robotId == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_error`))))
		return
	}

	// Today's generated count
	todayGenerateStats, err := msql.Model(`robot_payment_auth_code`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`create_date`, time.Now().Format(`2006-01-02`)).
		Field(`count(*) as count, sum(package_price) as price`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	// Yesterday's generated count
	yesterdayGenerateStats, err := msql.Model(`robot_payment_auth_code`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`create_date`, time.Now().AddDate(0, 0, -1).Format(`2006-01-02`)).
		Field(`count(*) as count, sum(package_price) as price`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	// Cumulative generated count
	totalStats, err := msql.Model(`robot_payment_auth_code`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Field(`count(*) as count, sum(package_price) as price`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	// Today's redeemed count
	todayExchangeStats, err := msql.Model(`robot_payment_auth_code`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`exchange_date`, time.Now().Format(`2006-01-02`)).
		Field(`count(*) as count, sum(package_price) as price`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	// Yesterday's redeemed count
	yesterdayExchangeStats, err := msql.Model(`robot_payment_auth_code`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`exchange_date`, time.Now().AddDate(0, 0, -1).Format(`2006-01-02`)).
		Field(`count(*) as count, sum(package_price) as price`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	// Cumulative redeemed count
	totdalExchangeStats, err := msql.Model(`robot_payment_auth_code`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`usage_status`, `>`, cast.ToString(define.RobotPaymentAuthCodeUsageStatusPending)).
		Field(`count(*) as count, sum(package_price) as price`).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{
		"today_generate_count":     todayGenerateStats[`count`],
		"today_generate_price":     todayGenerateStats[`price`],
		"today_exchange_count":     todayExchangeStats[`count`],
		"today_exchange_price":     todayExchangeStats[`price`],
		`yesterday_generate_count`: yesterdayGenerateStats[`count`],
		`yesterday_generate_price`: yesterdayGenerateStats[`price`],
		`yesterday_exchange_count`: yesterdayExchangeStats[`count`],
		`yesterday_exchange_price`: yesterdayExchangeStats[`price`],
		`total_generate_count`:     totalStats[`count`],
		`total_generate_price`:     totalStats[`price`],
		`total_exchange_count`:     totdalExchangeStats[`count`],
		`total_exchange_price`:     totdalExchangeStats[`price`],
	}, nil))
}

// GetAuthCodeManager gets authorization code manager list
func GetAuthCodeManager(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	robotId := cast.ToInt(c.Query(`robot_id`))
	list, err := msql.Model(`robot_payment_auth_code_manager`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

// AddAuthCodeManager adds authorization code manager
func AddAuthCodeManager(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	robotId := cast.ToInt(c.PostForm(`robot_id`))
	if robotId == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_err`))))
		return
	}
	val := define.RobotPaymentAuthCodePrefix + tool.Random(10) + define.RobotPaymentAuthCodeSuffix
	key := define.RobotPaymentAuthCodeManagerCachePrefix + cast.ToString(robotId)
	err := define.Redis.Set(c.Request.Context(), key, val, 3*time.Minute).Err()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(map[string]any{"code": val}, nil))
	return
}

// DeleteAuthCodeManager deletes authorization code manager
func DeleteAuthCodeManager(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	robotId := cast.ToInt(c.PostForm(`robot_id`))
	id := cast.ToInt(c.PostForm(`id`))
	if robotId == 0 || id == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_err`))))
		return
	}
	_, err := msql.Model(`robot_payment_auth_code_manager`, define.Postgres).
		Where(`id`, cast.ToString(id)).
		Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}
