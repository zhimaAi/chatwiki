// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/textsplitter"
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func PassiveReplyLogNotify(logid int64, question, answer string) {
	//答案抽取图片,内容按500分段
	content, images := common.GetImgInMessage(answer, false)
	split := textsplitter.NewRecursiveCharacter()
	split.ChunkSize, split.ChunkOverlap = 500, 0
	replys, _ := split.SplitText(content)
	if define.IsDev && question == `debug` { //测试环境调试数据
		replys = []string{`这是AI回复的第1段`, `这是AI回复的第2段`, `这是AI回复的第3段`, `这是AI回复的第4段`}
		images = []string{`/upload/model_icon/yiyan.png`, `/upload/model_icon/tongyi.png`, `/upload/model_icon/doubao.png`}
	}
	if len(images) > 1 {
		replys = append(replys, ` ----- 查看图片 ----- `)
	}
	//图片兼容非oss链接处理
	for i, image := range images {
		if !common.IsUrl(image) {
			images[i] = define.Config.WebService[`h5_domain`] + image
		}
	}
	//信息入库处理
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
	//完成消息通知
	err = define.Redis.Publish(context.Background(), lib_define.RedisPrefixPassiveSubscribe, logid).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error(`passive_id:%s,err:%s`, logid, err.Error())
	}
}
