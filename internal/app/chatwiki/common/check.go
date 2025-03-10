// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"encoding/json"
	"errors"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

var InvalidLibraryImageError = errors.New("invalid library image")

func CheckChatRequest(c *gin.Context) (*define.ChatBaseParam, error) {
	//source check
	appType := strings.TrimSpace(c.GetHeader(`App-Type`))
	if len(appType) == 0 {
		appType = lib_define.AppYunH5 //default value
	}
	if !tool.InArrayString(appType, []string{lib_define.AppYunH5, lib_define.AppYunPc}) {
		return nil, errors.New(i18n.Show(GetLang(c), `param_invalid`, `app_type`))
	}
	//format check
	robotKey := strings.TrimSpace(c.PostForm(`robot_key`))
	if len(robotKey) == 0 {
		robotKey = strings.TrimSpace(c.Query(`robot_key`))
	}
	if !CheckRobotKey(robotKey) {
		return nil, errors.New(i18n.Show(GetLang(c), `param_invalid`, `robot_key`))
	}
	openid := strings.TrimSpace(c.PostForm(`openid`))
	if len(openid) == 0 {
		openid = strings.TrimSpace(c.Query(`openid`))
	}
	if !IsChatOpenid(openid) {
		return nil, errors.New(i18n.Show(GetLang(c), `param_invalid`, `openid`))
	}
	//data check
	robot, err := GetRobotInfo(robotKey)
	if err != nil {
		logs.Error(err.Error())
		return nil, errors.New(i18n.Show(GetLang(c), `sys_err`))
	}
	if len(robot) == 0 {
		return nil, errors.New(i18n.Show(GetLang(c), `no_data`))
	}
	adminUserId := cast.ToInt(robot[`admin_user_id`])
	customer, err := GetCustomerInfo(openid, adminUserId)
	if err != nil {
		logs.Error(err.Error())
		return nil, errors.New(i18n.Show(GetLang(c), `sys_err`))
	}
	return &define.ChatBaseParam{AppType: appType, Openid: openid, AdminUserId: adminUserId, Robot: robot, Customer: customer}, nil

}

func CheckMenuJson(menuJson string) (string, error) {
	info := define.MenuJsonStruct{}
	_ = tool.JsonDecodeUseNumber(menuJson, &info)
	questions := make([]string, 0)
	for _, question := range info.Question {
		question := strings.TrimSpace(question)
		if len(question) > 0 {
			questions = append(questions, question)
		}
	}
	info.Question = questions
	return tool.JsonEncode(info)
}

func CheckCommonQuestionJson(c *gin.Context, commonQuestionList string) (string, error) {
	var commonQuestionListArray []string
	err := tool.JsonDecode(commonQuestionList, &commonQuestionListArray)
	if err != nil {
		return "", errors.New(i18n.Show(GetLang(c), `param_invalid`, `common_question_list`))
	}
	for _, commonQuestion := range commonQuestionListArray {
		if len(commonQuestion) == 0 {
			return "", errors.New(i18n.Show(GetLang(c), `param_invalid`, `common_question_list`))
		}
	}
	if len(commonQuestionListArray) > 10 {
		return "", errors.New(i18n.Show(GetLang(c), `param_invalid`, `common_question_list`))
	}
	return tool.JsonEncode(commonQuestionListArray)
}

func CheckIds(ids string) bool {
	ok, err := regexp.MatchString(`^(\d+)(,\d+)*$`, ids)
	if err == nil && ok {
		return true
	}
	return false
}

func CheckRobotKey(robotKey string) bool {
	ok, err := regexp.MatchString(`^[a-zA-Z0-9]{10}$`, robotKey)
	if err == nil && ok {
		return true
	}
	return false
}

func IsChatOpenid(openid string) bool {
	ok, err := regexp.MatchString(`^[a-zA-Z0-9_\-]{1,78}$`, openid)
	if err == nil && ok {
		return true
	}
	return false
}

func IsVariableName(key string) bool {
	ok, err := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_\-]{0,99}$`, key)
	if err == nil && ok {
		return true
	}
	return false
}

func IsVariableNames(variable string) bool {
	ok, err := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_\-.]*$`, variable)
	if err == nil && ok {
		return true
	}
	return false
}

func GetImgInMessage(message string) (string, []string) {
	imgRE := regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']>`)
	imgs := imgRE.FindAllStringSubmatch(message, -1)
	out := make([]string, len(imgs))
	for i := range out {
		out[i] = imgs[i][1]
		if !strings.Contains(out[i], "http://") && !strings.Contains(out[i], "https://") {
			out[i] = define.Config.WebService["api_domain"] + out[i]
		}
	}
	message = imgRE.ReplaceAllString(message, "")
	return message, out
}

func CheckSplitParams(splitParams define.SplitParams, lang string) (define.SplitParams, error) {
	//diy split
	if len(splitParams.SeparatorsNo) == 0 {
		return splitParams, errors.New(i18n.Show(lang, `param_empty`, `separators_no`))
	}
	if splitParams.ChunkSize < 200 || splitParams.ChunkSize > 2000 {
		return splitParams, errors.New(i18n.Show(lang, `chunk_size_err`, 200, 2000))
	}
	maxChunkOverlap := splitParams.ChunkSize / 2
	if splitParams.ChunkOverlap < 0 || splitParams.ChunkOverlap > maxChunkOverlap {
		return splitParams, errors.New(i18n.Show(lang, `chunk_overlap_err`, 0, maxChunkOverlap))
	}
	for i, noStr := range strings.Split(splitParams.SeparatorsNo, `,`) {
		no := cast.ToInt(noStr)
		if no < 1 || no > len(define.SeparatorsList) {
			return splitParams, errors.New(i18n.Show(lang, `param_invalid`, `separators_no.`+cast.ToString(i)))
		}
		code := define.SeparatorsList[no-1][`code`]
		if realCode, ok := code.([]string); ok {
			splitParams.Separators = append(splitParams.Separators, realCode...)
		} else {
			splitParams.Separators = append(splitParams.Separators, cast.ToString(code))
		}
	}
	//qa_doc
	if splitParams.IsQaDoc == define.DocTypeQa {
		if splitParams.IsTableFile == define.FileIsTable {
			if len(splitParams.QuestionColumn) == 0 {
				return splitParams, errors.New(i18n.Show(lang, `param_empty`, `question_column`))
			}
			if len(splitParams.AnswerColumn) == 0 {
				return splitParams, errors.New(i18n.Show(lang, `param_empty`, `answer_column`))
			}
		} else {
			if len(splitParams.QuestionLable) == 0 {
				return splitParams, errors.New(i18n.Show(lang, `param_empty`, `question_lable`))
			}
			if len(splitParams.AnswerLable) == 0 {
				return splitParams, errors.New(i18n.Show(lang, `param_empty`, `answer_lable`))
			}
		}
	} else {
		splitParams.QuestionLable = ``
		splitParams.AnswerLable = ``
	}
	return splitParams, nil
}

func CheckLibraryImage(images []string) (string, error) {
	extensions := strings.Join(define.ImageAllowExt, "|")
	pattern := `^\/upload\/chat_ai\/\d+\/library_image\/\d+\/[a-f0-9]{32}\.(` + extensions + `)$`
	re := regexp.MustCompile(pattern)
	for _, image := range images {
		if IsUrl(image) { //oss file
			ext := strings.ToLower(strings.TrimLeft(filepath.Ext(image), `.`))
			if !tool.InArrayString(ext, define.ImageAllowExt) {
				return ``, InvalidLibraryImageError
			}
		} else { //local file
			if !re.MatchString(image) {
				return "", InvalidLibraryImageError
			}
		}
	}
	jsonImages, err := json.Marshal(images)
	if err != nil {
		return "[]", err
	}
	if string(jsonImages) == "null" {
		return "[]", nil
	}
	return string(jsonImages), nil
}
