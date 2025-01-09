// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"io/fs"
	"os"
	"path/filepath"
	"time"

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
		Where(`doc_last_renew_time`, "<=", cast.ToString(time.Now().Add(-24*time.Hour).Unix())).
		Field(`id,admin_user_id,doc_auto_renew_frequency,doc_last_renew_time,doc_url`).
		Select()
	if err != nil {
		logs.Error(err.Error())
		return
	}

	for _, doc := range docs {
		shouldRenew := false
		if cast.ToInt(doc[`doc_auto_renew_frequency`]) == 2 && cast.ToInt64(doc[`doc_last_renew_time`]) <= time.Now().Add(-24*time.Hour).Unix() { // everyday
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
