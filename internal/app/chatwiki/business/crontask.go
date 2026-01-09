// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/chatwiki/business/manage"
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func RenewCrawl() {
	docs, err := msql.Model(`chat_ai_library_file`, `postgres`).
		Where(`doc_type`, cast.ToString(define.DocTypeOnline)).
		Where(`status`, cast.ToString(define.FileStatusLearned)).
		Where(`doc_auto_renew_frequency`, ">", "1").
		Field(`id,admin_user_id,doc_auto_renew_frequency,doc_auto_renew_minute,doc_last_renew_time,doc_url`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		return
	}

	timestamp := time.Now().Unix()
	t := time.Unix(timestamp, 0)
	currentMinute := t.Hour()*60 + t.Minute()

	for _, doc := range docs {
		shouldRenew := false
		if cast.ToInt(doc[`doc_auto_renew_minute`]) == currentMinute {
			if cast.ToInt(doc[`doc_auto_renew_frequency`]) == 2 { // everyday
				shouldRenew = true
			} else if cast.ToInt(doc[`doc_auto_renew_frequency`]) == 3 && cast.ToInt64(doc[`doc_last_renew_time`]) <= time.Now().Add(-3*24*time.Hour).Unix() { //every 3 days
				shouldRenew = true
			} else if cast.ToInt(doc[`doc_auto_renew_frequency`]) == 4 && cast.ToInt64(doc[`doc_last_renew_time`]) <= time.Now().Add(-24*7*time.Hour).Unix() { //every 7 days
				shouldRenew = true
			} else if cast.ToInt(doc[`doc_auto_renew_frequency`]) == 5 && cast.ToInt64(doc[`doc_last_renew_time`]) <= time.Now().Add(-24*30*time.Hour).Unix() { //every 30 days
				shouldRenew = true
			} else {
				shouldRenew = false
			}
		}

		if shouldRenew {
			if message, err := tool.JsonEncode(map[string]any{`file_id`: doc[`id`], `admin_user_id`: doc[`admin_user_id`]}); err != nil {
				logs.Error(err.Error())
			} else if err := common.AddJobs(define.CrawlArticleTopic, message); err != nil {
				logs.Error(err.Error())
			}
		}
	}
}

func DeleteFormEntry() {
	_, err := msql.Model(`form_entry`, define.Postgres).
		Where(`delete_time`, `>`, `0`).
		Where(`delete_time`, `<`, cast.ToString(time.Now().Add(-time.Hour*24*7).Unix())).
		Delete()
	if err != nil {
		logs.Error(err.Error())
		return
	}
}

func DeleteExportFile() {
	err := filepath.WalkDir(`static/public/export`, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil {
			return err //path not exist
		}
		if info, err := d.Info(); err == nil && !info.IsDir() && info.ModTime().Unix() < time.Now().Unix()-86400*7 {
			if err = os.Remove(path); err != nil {
				logs.Error(err.Error())
			}
		}
		return nil
	})
	if err != nil {
		logs.Error(err.Error())
	}
}

func DeleteConvertHtml() {
	err := filepath.WalkDir(define.UploadDir+`chat_ai`, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil {
			return err //path not exist
		}
		if !strings.Contains(path, `convert`) || strings.ToLower(filepath.Ext(path)) != `.html` {
			return nil //not convert create html file
		}
		if info, err := d.Info(); err == nil && !info.IsDir() && info.ModTime().Unix() < time.Now().Unix()-86400 {
			if err = os.Remove(path); err != nil {
				logs.Error(err.Error())
			}
		}
		return nil
	})
	if err != nil {
		logs.Error(err.Error())
	}
}

func DeleteClientSide() {
	err := filepath.WalkDir(`static/public/client_side`, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil {
			return err //path not exist
		}
		if info, err := d.Info(); err == nil && !info.IsDir() && info.ModTime().Unix() < time.Now().Unix()-86400 {
			if err = os.Remove(path); err != nil {
				logs.Error(err.Error())
			}
		}
		return nil
	})
	if err != nil {
		logs.Error(err.Error())
	}
}

func DeleteDownloadFile() {
	err := filepath.WalkDir(`static/public/download`, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil || d.Name() == `.gitignore` {
			return err //path not exist
		}
		if info, err := d.Info(); err == nil && !info.IsDir() && info.ModTime().Unix() < time.Now().Unix()-3600 {
			if err = os.Remove(path); err != nil {
				logs.Error(err.Error())
			}
		}
		return nil
	})
	if err != nil {
		logs.Error(err.Error())
	}
}

func CheckAliOcrRequest() {
	files, err := msql.Model(`chat_ai_library_file`, define.Postgres).
		Where(`status`, cast.ToString(define.FileStatusInitial)).
		Where(`pdf_parse_type`, cast.ToString(define.PdfParseTypeOcrAli)).
		Where(`create_time`, `>`, cast.ToString(time.Now().Add(-24*time.Hour).Unix())).
		Where(`ali_ocr_job_id`, `<>`, ``).
		Select()
	if err != nil {
		logs.Error(err.Error())
		return
	}

	for _, file := range files {
		company, err := msql.Model(`company`, define.Postgres).Where(`parent_id`, file[`admin_user_id`]).Find()
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		if len(company) == 0 || cast.ToInt(company[`ali_ocr_switch`]) != 1 {
			continue
		}
		if err := common.QueryAndParseAliOcrRequest(file, company[`ali_ocr_key`], company[`ali_ocr_secret`]); err != nil {
			logs.Error(err.Error())
			continue
		}
	}
}

func UpdateLibraryFileData() {
	// 每次处理的批次大小
	const batchSize = 1000
	logs.Debug("开始更新文档分段数据")
	var (
		total = 0
		page  = 0
	)
	for {
		// 获取一批数据
		rows, err := msql.Model("chat_ai_library_file_data", define.Postgres).
			Limit((page * batchSize), batchSize).
			Order(`id asc`).
			ColumnArr(`id`)

		if err != nil {
			logs.Error(err.Error())
			return
		}

		// 如果没有数据了就退出循环
		if len(rows) == 0 {
			break
		}
		minId := cast.ToInt(rows[0])
		maxId := cast.ToInt(rows[len(rows)-1])
		// 更新数据
		_, err = msql.Model("chat_ai_library_file_data", define.Postgres).
			Where("id", ">=", cast.ToString(minId)).
			Where("id", "<=", cast.ToString(maxId)).
			Update2(fmt.Sprintf(`yesterday_hits = today_hits, today_hits = 0,update_time=%v`, tool.Time2Int()))
		if err != nil {
			logs.Error(err.Error())
			return
		}
		page++
		total += len(rows)
	}
	logs.Debug("结束更新文档分段数据,共:%v", total)
}

func DeleteLlmRequestLogs() {
	endTime := tool.GetTimestamp(tool.GetYmdBeforeDay(7)) - 1
	m := msql.Model(common.GetLlmRequestLogsTableName(endTime), define.Postgres)
	info, err := m.Where(`create_time`, `<=`, cast.ToString(endTime)).
		Field(`min(id) minid,max(id) maxid`).Find()
	if err != nil {
		var sqlerr *pq.Error
		if errors.As(err, &sqlerr) && sqlerr.Code == `42P01` {
			return //表不存在,不报错了
		}
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return
	}
	minId, maxId := cast.ToInt(info[`minid`]), cast.ToInt(info[`maxid`])
	if minId <= 0 || maxId <= 0 {
		return //没有可以清理的数据
	}
	var size = 1000 //每一批次数
	for i := 0; ; i++ {
		start, end := minId+i*size, min(maxId, minId+(i+1)*size)
		affect, err := m.Where(`id`, `>=`, cast.ToString(start)).
			Where(`id`, `<=`, cast.ToString(end)).Delete()
		if err != nil {
			logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
			return
		}
		logs.Debug(`清理结果:第%d轮(%d~%d),affect(%d)`, i+1, start, end, affect)
		if end >= maxId {
			break //处理完毕,结束循环
		}
	}
}

func UpdateBatchSendData() {

	list, err := msql.Model("wechat_official_account_batch_send_task", define.Postgres).
		Where(`send_status`, cast.ToString(define.BatchSendStatusExec)).Select()
	if err != nil {
		logs.Error("任务列表查询错误:" + err.Error())
		return
	}

	for _, params := range list {
		res, err := manage.BridgeGetSendStatus(params["access_key"], params["send_msg_id"])
		if err != nil {
			logs.Error("查询错误：" + err.Error())
		}

		if res != nil && res.MsgStatus == "SEND_SUCCESS" {
			msql.Model("wechat_official_account_batch_send_task", define.Postgres).Where(`id`, params["id"]).Update(msql.Datas{
				`update_time`: time.Now().Unix(),
				`send_status`: define.BatchSendStatusSucc,
			})
		}
	}

	//	获取发送成功的任务，继续刷新评论
	succList, err := msql.Model("wechat_official_account_batch_send_task", define.Postgres).
		Where(`send_status`, cast.ToString(define.BatchSendStatusSucc)).Where(`comment_status`, define.BaseOpen).Select()
	if err != nil {
		logs.Error("任务列表查询错误:" + err.Error())
		return
	}

	for _, task := range succList {
		adminConfig := common.GetAdminConfig(cast.ToInt(task["admin_user_id"]))

		sendTime := cast.ToInt(task["send_time"])
		if sendTime == 0 {
			sendTime = cast.ToInt(task["create_time"])
		}
		//如果超过拉取评论时间，
		if int(time.Now().Unix())-sendTime > cast.ToInt(adminConfig["comment_pull_days"])*86400 {
			continue
		}

		timeLimit := 60
		if define.IsDev {
			timeLimit = 10
		}

		//未到同步评论时间间隔
		if (cast.ToInt(task["last_comment_sync_time"]) + (cast.ToInt(adminConfig["comment_pull_limit"])-1)*timeLimit) > int(time.Now().Unix()) {
			logs.Debug("未到同步时间，下次同步时间为：" + cast.ToString(cast.ToInt(task["last_comment_sync_time"])+cast.ToInt(adminConfig["comment_pull_limit"])*timeLimit))
			continue
		}

		delayTime := int64(5)
		common.AddDelayTask(define.DelayTaskEvent{
			BaseDelayTask: define.BaseDelayTask{Type: define.OfficialAccountBatchSendSyncCommentTask},
			AdminUserId:   cast.ToInt(task["admin_user_id"]),
			TaskId:        cast.ToInt(task["id"]),
		}, delayTime)

	}

}

// CheckRobotPaymentDurationAuthCode 检查应用收费时长套餐
func CheckRobotPaymentDurationAuthCode() {
	logs.Debug(`开始检查应用收费时长套餐`)
	robotIdList, err := msql.Model(`robot_payment_setting`, define.Postgres).
		Group(`robot_id`).
		ColumnArr(`robot_id`)
	if err != nil {
		logs.Error(err.Error())
		return
	}

	for _, robotId := range robotIdList {
		exchangeOpenIdList, err := msql.Model(`robot_payment_auth_code`, define.Postgres).
			Where(`robot_id`, robotId).
			Group(`exchanger_openid`).
			ColumnArr(`exchanger_openid`)
		if err != nil {
			logs.Error(err.Error())
			return
		}
		for _, exchangeOpenId := range exchangeOpenIdList {
			authCodeItem, err := msql.Model(`robot_payment_auth_code`, define.Postgres).
				Where(`robot_id`, cast.ToString(robotId)).
				Where(`exchanger_openid`, exchangeOpenId).
				Where(`package_type`, cast.ToString(define.RobotPaymentPackageTypeDuration)).
				Where(`usage_status`, cast.ToString(define.RobotPaymentAuthCodeUsageStatusExchanged)).
				Order(`id asc`).
				Find()
			if err != nil {
				logs.Error(err.Error())
				return
			}
			if len(authCodeItem) == 0 {
				continue
			}
			packageDuration := cast.ToInt(authCodeItem[`package_duration`])
			usedDuration := cast.ToInt(authCodeItem[`used_duration`])

			data := msql.Datas{
				`update_time`: tool.Time2Int(),
			}
			usedDuration += 1
			if usedDuration < packageDuration {
				data[`used_duration`] = usedDuration
			} else if usedDuration == packageDuration {
				data[`used_duration`] = usedDuration
				data[`use_time`] = tool.Time2Int()
				data[`use_date`] = time.Now().Format(`2006-01-02`)
				data[`usage_status`] = define.RobotPaymentAuthCodeUsageStatusUsed
			} else {
				data[`use_time`] = tool.Time2Int()
				data[`use_date`] = time.Now().Format(`2006-01-02`)
				data[`usage_status`] = define.RobotPaymentAuthCodeUsageStatusUsed
			}

			_, err = msql.Model(`robot_payment_auth_code`, define.Postgres).
				Where(`id`, authCodeItem[`id`]).
				Update(data)
			if err != nil {
				logs.Error(err.Error())
			}
		}
	}
}
