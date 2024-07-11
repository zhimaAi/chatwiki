// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/tool"
)

var InvalidLibraryImageError = errors.New("invalid library image")

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

func CheckSplitParams(c *gin.Context, isTableFile int) (define.SplitParams, error) {
	splitParams := define.SplitParams{
		IsTableFile:        isTableFile,
		IsDiySplit:         cast.ToInt(c.Query(`is_diy_split`)),
		SeparatorsNo:       strings.TrimSpace(c.Query(`separators_no`)),
		Separators:         make([]string, 0),
		ChunkSize:          cast.ToInt(c.Query(`chunk_size`)),
		ChunkOverlap:       cast.ToInt(c.Query(`chunk_overlap`)),
		IsQaDoc:            cast.ToInt(c.Query(`is_qa_doc`)),
		QuestionLable:      strings.TrimSpace(c.Query(`question_lable`)),
		AnswerLable:        strings.TrimSpace(c.Query(`answer_lable`)),
		QuestionColumn:     strings.TrimSpace(c.Query(`question_column`)),
		AnswerColumn:       strings.TrimSpace(c.Query(`answer_column`)),
		EnableExtractImage: cast.ToBool(c.Query(`enable_extract_image`)),
	}
	if splitParams.IsTableFile == define.FileIsTable {
		if splitParams.IsDiySplit == define.SplitTypeDiy {
			return splitParams, errors.New(i18n.Show(GetLang(c), `not_support`))
		} else {
			if splitParams.IsQaDoc == define.DocTypeQa {
				//return splitParams, errors.New(i18n.Show(GetLang(c), `not_support`))
			} else {
				//category1:excel+smart+general
			}
		}
	} else {
		if splitParams.IsDiySplit == define.SplitTypeDiy {
			if splitParams.IsQaDoc == define.DocTypeQa {
				return splitParams, errors.New(i18n.Show(GetLang(c), `not_support`))
			} else {
				//category2:general+diy+general
			}
		} else {
			if splitParams.IsQaDoc == define.DocTypeQa {
				//category3:general+smart+qa_doc
			} else {
				//category4:general+smart+general
			}
		}
	}
	//diy split
	if splitParams.IsDiySplit == define.SplitTypeDiy {
		if len(splitParams.SeparatorsNo) == 0 {
			return splitParams, errors.New(i18n.Show(GetLang(c), `param_empty`, `separators_no`))
		}
		if splitParams.ChunkSize < 200 || splitParams.ChunkSize > 2000 {
			return splitParams, errors.New(i18n.Show(GetLang(c), `chunk_size_err`, 200, 2000))
		}
		maxChunkOverlap := splitParams.ChunkSize / 2
		if splitParams.ChunkOverlap < 2 || splitParams.ChunkOverlap > maxChunkOverlap {
			return splitParams, errors.New(i18n.Show(GetLang(c), `chunk_overlap_err`, 2, maxChunkOverlap))
		}
		for i, noStr := range strings.Split(splitParams.SeparatorsNo, `,`) {
			no := cast.ToInt(noStr)
			if no < 1 || no > len(define.SeparatorsList) {
				return splitParams, errors.New(i18n.Show(GetLang(c), `param_invalid`, `separators_no.`+cast.ToString(i)))
			}
			code := define.SeparatorsList[no-1][`code`]
			if realCode, ok := code.([]string); ok {
				splitParams.Separators = append(splitParams.Separators, realCode...)
			} else {
				splitParams.Separators = append(splitParams.Separators, cast.ToString(code))
			}
		}
	} else {
		splitParams.SeparatorsNo = ``
		splitParams.ChunkSize = 0
		splitParams.ChunkOverlap = 0
	}
	//qa_doc
	if splitParams.IsQaDoc == define.DocTypeQa {
		if splitParams.IsTableFile == define.FileIsTable {
			if len(splitParams.QuestionColumn) == 0 {
				return splitParams, errors.New(i18n.Show(GetLang(c), `param_empty`, `question_column`))
			}
			if len(splitParams.AnswerColumn) == 0 {
				return splitParams, errors.New(i18n.Show(GetLang(c), `param_empty`, `answer_column`))
			}
		} else {
			if len(splitParams.QuestionLable) == 0 {
				return splitParams, errors.New(i18n.Show(GetLang(c), `param_empty`, `question_lable`))
			}
			if len(splitParams.AnswerLable) == 0 {
				return splitParams, errors.New(i18n.Show(GetLang(c), `param_empty`, `answer_lable`))
			}
		}
	} else {
		splitParams.QuestionLable = ``
		splitParams.AnswerLable = ``
	}
	return splitParams, nil
}

func CheckLibraryImage(images []string) (string, error) {
	pattern := `^\/upload\/chat_ai\/\d+\/library_image\/\d+\/[a-f0-9]{32}\.png$`
	re := regexp.MustCompile(pattern)
	for _, image := range images {
		if !re.MatchString(image) {
			return "", InvalidLibraryImageError
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
