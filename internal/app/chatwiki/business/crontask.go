// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
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
