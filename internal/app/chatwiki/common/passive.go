// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/textsplitter"
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func PassiveReplyLogNotify(lang string, logid int64, question, answer string) {
	// Extract images from answers, content segmented by 500
	content, images := GetImgInMessage(answer, false)
	split := textsplitter.NewRecursiveCharacter()
	split.ChunkSize, split.ChunkOverlap = 500, 0
	replys, _ := split.SplitText(content)
	if define.IsDev && question == `debug` { // Test environment debug data
		replys = []string{`ai_reply_1`, `ai_reply_2`, `ai_reply_3`, `ai_reply_4`}
		images = []string{`/upload/model_icon/yiyan.png`, `/upload/model_icon/doubao.png`}
	}
	if len(images) > 1 {
		replys = append(replys, fmt.Sprintf(` ----- %s ----- `, i18n.Show(lang, `view_images`)))
	}
	// Process image compatibility for non-OSS links
	for i, image := range images {
		if !IsUrl(image) {
			images[i] = define.Config.WebService[`image_domain`] + image
		}
	}
	// Process data storage
	passiveModel := msql.Model(`chat_ai_passive_reply_log`, define.Postgres)
	_, err := passiveModel.Where(`id`, cast.ToString(logid)).Update(msql.Datas{
		`status`:      define.SwitchOn,
		`replys`:      tool.JsonEncodeNoError(replys),
		`images`:      tool.JsonEncodeNoError(images),
		`update_time`: tool.Time2Int(),
	})
	if err != nil {
		logs.Error(`sql:%s,err:%s`, passiveModel.GetLastSql(), err.Error())
	}
	// Complete message notification
	err = define.Redis.Publish(context.Background(), lib_define.RedisPrefixPassiveSubscribe, logid).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error(`passive_id:%s,err:%s`, logid, err.Error())
	}
}
