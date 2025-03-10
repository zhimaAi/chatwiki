// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package define

import (
	"errors"
	"strconv"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

type RequestContext struct {
	Lang        string
	AdminUserID string
}

type MenuJsonStruct struct {
	Content  string   `json:"content"`
	Question []string `json:"question"`
}

type UploadInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Ext     string `json:"ext"`
	Link    string `json:"link"`
	Online  bool   `json:"-"`
	DocUrl  string `json:"-"`
	Custom  bool   `json:"-"`
	Columns string `json:"columns"`
}

func (u *UploadInfo) GetDocType() int {
	if u.Custom {
		return DocTypeCustom
	}
	if u.Online {
		return DocTypeOnline
	}
	return DocTypeLocal
}

type ChatBaseParam struct {
	AppType     string
	Openid      string
	AdminUserId int
	Robot       msql.Params
	Customer    msql.Params
}

type ChatRequestParam struct {
	*ChatBaseParam
	Error          error
	Lang           string
	Question       string
	OpenApiContent string
	DialogueId     int
	Prompt         string
	LibraryIds     string
	IsClose        *bool
}

type DocSplitItem struct {
	Number    int      `json:"number"`
	PageNum   int      `json:"page_num"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Question  string   `json:"question"`
	Answer    string   `json:"answer"`
	WordTotal int      `json:"word_total"`
	Images    []string `json:"images"`
}

type SplitParams struct {
	IsTableFile        int      `json:"is_table_file"`
	IsDiySplit         int      `json:"is_diy_split"`
	SeparatorsNo       string   `json:"separators_no"`
	Separators         []string `json:"-"`
	ChunkSize          int      `json:"chunk_size"`
	ChunkOverlap       int      `json:"chunk_overlap"`
	IsQaDoc            int      `json:"is_qa_doc"`
	QuestionLable      string   `json:"question_lable"`
	AnswerLable        string   `json:"answer_lable"`
	QuestionColumn     string   `json:"question_column"`
	AnswerColumn       string   `json:"answer_column"`
	EnableExtractImage bool     `json:"enable_extract_image"`
}

type FormFilterCondition struct {
	FormFieldId int    `json:"form_field_id"`
	Rule        string `json:"rule"`
	RuleValue1  string `json:"rule_value1"`
	RuleValue2  string `json:"rule_value2"`
}

func (f FormFilterCondition) Check(_type string) error {
	if _type == `string` {
		if !tool.InArrayString(f.Rule, []string{`string_eq`, `string_neq`, `string_contain`, `string_not_contain`, `string_empty`, `string_not_empty`}) {
			return errors.New(`rule value must be one of string_eq, string_neq, string_contain, string_not_contain, string_empty, string_not_empty when type is String`)
		}
		if len(f.RuleValue2) > 0 {
			return errors.New(`rule_value1 is empty or rule_value2 not empty when type is String`)
		}
	} else if _type == `integer` {
		if !tool.InArrayString(f.Rule, []string{`integer_gt`, `integer_gte`, `integer_lt`, `integer_lte`, `integer_eq`, `integer_between`}) {
			return errors.New(`rule value must be on of integer_gt, integer_gte, integer_lt, integer_lte, integer_eq, integer_between when type is integer or number`)
		}
		if len(f.RuleValue1) == 0 {
			return errors.New(`rule_value1 is empty when type is Integer or Number`)
		}
		if _, err := strconv.Atoi(f.RuleValue1); err != nil {
			return errors.New(`rule_value1 is not integer when type is Integer`)
		}
		if f.Rule == `number_between` {
			if len(f.RuleValue2) == 0 {
				return errors.New(`rule_value2 is empty when rule is between`)
			}
			if _, err := strconv.Atoi(f.RuleValue2); err != nil {
				return errors.New(`rule_value2 is invalid integer when rule is between and type is integer`)
			}
		}
	} else if _type == `number` {
		if !tool.InArrayString(f.Rule, []string{`number_gt`, `number_gte`, `number_lt`, `number_lte`, `number_eq`, `number_between`}) {
			return errors.New(`rule value must be on of number_gt, number_gte, number_lt, number_lte, number_eq, number_between when type is integer or number`)
		}
		if len(f.RuleValue1) == 0 {
			return errors.New(`rule_value1 is empty when type is integer or number`)
		}
		if _, err := strconv.ParseFloat(f.RuleValue1, 64); err != nil {
			return errors.New(`rule_value1 is invalid number when rule is number`)
		}
		if f.Rule == `number_between` {
			if len(f.RuleValue2) == 0 {
				return errors.New(`rule_value2 is empty when rule is between and type is number`)
			}
			if _, err := strconv.ParseFloat(f.RuleValue2, 64); err != nil {
				return errors.New(`rule_value2 is invalid number when rule is between and type is number`)
			}
		}
	} else if _type == `boolean` {
		if !tool.InArrayString(f.Rule, []string{`boolean_true`, `boolean_false`}) {
			return errors.New(`rule value must be on of boolean_true, boolean_false when type is boolean`)
		}
	}
	return nil
}

type SimilarityResult []msql.Params

func (m SimilarityResult) Len() int {
	return len(m)
}

func (m SimilarityResult) Less(i, j int) bool {
	return cast.ToFloat64(m[i][`similarity`]) > cast.ToFloat64(m[j][`similarity`])
}

func (m SimilarityResult) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

type CommonQuestion struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type UploadFormFile struct {
	Total     int              `json:"total"`
	Processed int              `json:"processed"`
	Finish    bool             `json:"finish"`
	Success   int              `json:"success"`
	ErrData   []map[string]any `json:"err_data"`
}
