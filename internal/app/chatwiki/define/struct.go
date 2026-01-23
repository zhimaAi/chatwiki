// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

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
	AppInfo     msql.Params
	Openid      string
	AdminUserId int
	Robot       msql.Params
	Customer    msql.Params
	RelUserId   int
}

type ChatRequestParam struct {
	*ChatBaseParam
	Error                error
	Lang                 string
	Question             string
	MsgId                string
	PassiveId            int64
	ReceivedMessageType  string
	ReceivedMessage      map[string]any
	MediaIdToOssUrl      string
	ThumbMediaIdToOssUrl string
	OpenApiContent       string
	WechatappAppid       string
	DialogueId           int
	Prompt               string
	LibraryIds           string
	IsClose              *bool
	WorkFlowGlobal       map[string]any
	QuoteLib             bool
	LoopTestParams       []any
	BatchTestParams      []any
	HeaderToken          string
	ChatPromptVariables  string
	TestParams           []any
}

type DocSplitItem struct {
	FileDataId          int      `json:"file_data_id"`
	Number              int      `json:"number"`
	PageNum             int      `json:"page_num"`
	Title               string   `json:"title"`
	Content             string   `json:"content"`
	Question            string   `json:"question"`
	SimilarQuestionList []string `json:"similar_question_list"`
	Answer              string   `json:"answer"`
	WordTotal           int      `json:"word_total"`
	Images              []string `json:"images"`
	AiChunkErrMsg       string   `json:"ai_chunk_err_msg"`
	//父子分段
	FatherChunkParagraphNumber int `json:"father_chunk_paragraph_number"`
}

type DocSplitItems []DocSplitItem

func (list *DocSplitItems) UnifySetNumber() {
	if len(*list) == 0 {
		return
	}
	var curIdx int
	var lastFn = (*list)[0].FatherChunkParagraphNumber
	for i := range *list {
		if i > 0 && lastFn != (*list)[i].FatherChunkParagraphNumber {
			curIdx = 0
			lastFn = (*list)[i].FatherChunkParagraphNumber
		}
		curIdx++
		(*list)[i].Number = curIdx //serial number
	}
}

type SplitParams struct {
	SplitPreviewParams
	IsTableFile                int      `json:"is_table_file"`
	IsDiySplit                 int      `json:"is_diy_split"`
	SeparatorsNo               string   `json:"separators_no"`
	Separators                 []string `json:"-"`
	ChunkSize                  int      `json:"chunk_size"`
	ChunkOverlap               int      `json:"chunk_overlap"`
	IsQaDoc                    int      `json:"is_qa_doc"`
	QaIndexType                int      `json:"qa_index_type"`
	QuestionLable              string   `json:"question_lable"`
	SimilarLabel               string   `json:"similar_label"`
	AnswerLable                string   `json:"answer_lable"`
	QuestionColumn             string   `json:"question_column"`
	SimilarColumn              string   `json:"similar_column"`
	AnswerColumn               string   `json:"answer_column"`
	EnableExtractImage         bool     `json:"enable_extract_image"`
	ChunkType                  int      `json:"chunk_type"`
	SemanticChunkUseModel      string   `json:"semantic_chunk_use_model"`
	SemanticChunkModelConfigId int      `json:"semantic_chunk_model_config_id"`
	SemanticChunkSize          int      `json:"semantic_chunk_size"`
	SemanticChunkOverlap       int      `json:"semantic_chunk_overlap"`
	SemanticChunkThreshold     int      `json:"semantic_chunk_threshold"`
	PdfParseType               int      `json:"pdf_parse_type"`
	AiChunkPrumpt              string   `json:"ai_chunk_prumpt"`
	AiChunkModel               string   `json:"ai_chunk_model"`
	AiChunkModelConfigId       int      `json:"ai_chunk_model_config_id"`
	AiChunkTaskId              string   `json:"ai_chunk_task_id"`
	AiChunkSize                int      `json:"ai_chunk_size"`
	AiChunkNew                 bool     `json:"ai_chunk_new"`
	ParagraphChunk             bool     `json:"paragraph_chunk"`
	ChunkAsync                 bool     `json:"chunk_async"` // 异步
	FileExt                    string   `json:"file_ext"`
	NotMergedText              bool     `json:"not_merged_text"`
	//父子分段
	FatherChunkParagraphType int    `json:"father_chunk_paragraph_type"`
	FatherChunkSeparatorsNo  string `json:"father_chunk_separators_no"`
	FatherChunkChunkSize     int    `json:"father_chunk_chunk_size"`
	SonChunkSeparatorsNo     string `json:"son_chunk_separators_no"`
	SonChunkChunkSize        int    `json:"son_chunk_chunk_size"`
}

type SplitPreviewParams struct {
	ChunkPreview     bool `json:"chunk_preview"`
	ChunkPreviewSize int  `json:"chunk_preview_size"`
}

type FormFilterCondition struct {
	FormFieldId int    `json:"form_field_id"`
	Rule        string `json:"rule"`
	RuleValue1  string `json:"rule_value1"`
	RuleValue2  string `json:"rule_value2"`
}

func (f FormFilterCondition) Check(_type string, workFlow ...bool) error {
	isWorkFlow := len(workFlow) > 0 && workFlow[0]
	if _type == `string` {
		if !tool.InArrayString(f.Rule, []string{`string_eq`, `string_neq`, `string_contain`, `string_not_contain`, `string_empty`, `string_not_empty`}) {
			return errors.New(`rule value must be one of string_eq, string_neq, string_contain, string_not_contain, string_empty, string_not_empty when type is String`)
		}
		if len(f.RuleValue2) > 0 {
			return errors.New(`rule_value2 not empty when type is String`)
		}
	} else if _type == `integer` {
		if !tool.InArrayString(f.Rule, []string{`integer_gt`, `integer_gte`, `integer_lt`, `integer_lte`, `integer_eq`, `integer_between`}) {
			return errors.New(`rule value must be on of integer_gt, integer_gte, integer_lt, integer_lte, integer_eq, integer_between when type is integer or number`)
		}
		if len(f.RuleValue1) == 0 {
			return errors.New(`rule_value1 is empty when type is Integer or Number`)
		}
		if _, err := strconv.Atoi(f.RuleValue1); !isWorkFlow && err != nil {
			return errors.New(`rule_value1 is not integer when type is Integer`)
		}
		if f.Rule == `number_between` {
			if len(f.RuleValue2) == 0 {
				return errors.New(`rule_value2 is empty when rule is between`)
			}
			if _, err := strconv.Atoi(f.RuleValue2); !isWorkFlow && err != nil {
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
		if _, err := strconv.ParseFloat(f.RuleValue1, 64); !isWorkFlow && err != nil {
			return errors.New(`rule_value1 is invalid number when rule is number`)
		}
		if f.Rule == `number_between` {
			if len(f.RuleValue2) == 0 {
				return errors.New(`rule_value2 is empty when rule is between and type is number`)
			}
			if _, err := strconv.ParseFloat(f.RuleValue2, 64); !isWorkFlow && err != nil {
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

type RobotDefaultAppConfig struct {
	Name      string `json:"name"`
	Id        int    `json:"id"`
	AutoReply int    `json:"auto_reply"`
}

type UploadFormFile struct {
	Total     int              `json:"total"`
	Processed int              `json:"processed"`
	Finish    bool             `json:"finish"`
	Success   int              `json:"success"`
	ErrData   []map[string]any `json:"err_data"`
}

type SplitFaqParams struct {
	SeparatorsNo       string `json:"separators_no"`
	ChunkSize          int    `json:"chunk_size"`
	ChunkType          int    `json:"chunk_type"`
	ChunkPrompt        string `json:"chunk_prompt"`
	ChunkModel         string `json:"chunk_model"`
	ChunkModelConfigId int    `json:"chunk_model_config_id"`
	IsQaDoc            int    `json:"is_qa_doc"`
	FileExt            string `json:"file_ext"`
	ExtractType        int    `json:"extract_type"`
}

type OfficialAccountDraftListStruct struct {
	Item []struct {
		Content struct {
			CreateTime int `json:"create_time"`
			NewsItem   []struct {
				ArticleType              string `json:"article_type"`
				NeedOpenCofansCanComment int    `json:"need_open_cofans_can_comment"`
				ThumbUrl                 string `json:"thumb_url"`
				Title                    string `json:"title"`
				Digest                   string `json:"digest"`
			} `json:"news_item"`
			UpdateTime int `json:"update_time"`
		} `json:"content"`
		MediaId    string `json:"media_id"`
		UpdateTime int    `json:"update_time"`
	} `json:"item"`
	ItemCount  int `json:"item_count"`
	TotalCount int `json:"total_count"`
}

type OfficialAccountComment struct {
	CommentType int    `json:"comment_type"`
	Content     string `json:"content"`
	CreateTime  int    `json:"create_time"`
	Openid      string `json:"openid"`
	Reply       struct {
		Content    string `json:"content"`
		CreateTime int    `json:"create_time"`
	} `json:"reply"`
	UserCommentId int `json:"user_comment_id"`
}

type OfficialAccountCommentResp struct {
	Comment    []OfficialAccountComment `json:"comment"`
	CommentCnt int                      `json:"comment_cnt"`
	Errcode    int                      `json:"errcode"`
	Errmsg     string                   `json:"errmsg"`
	Total      int                      `json:"total"`
}

type RobotPaymentCountPackage struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Count int     `json:"count"`
	Price float32 `json:"price"`
}
type RobotPaymentDurationPackage struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Duration int     `json:"duration"`
	Count    int     `json:"count"`
	Price    float32 `json:"price"`
}

// 远程插件
type RemotePluginStruct struct {
	Author              string `json:"author"`
	CreateTime          string `json:"create_time"`
	Description         string `json:"description"`
	Enabled             string `json:"enabled"`
	FilterType          string `json:"filter_type"`
	FilterTypeTitle     string `json:"filter_type_title"`
	HelpUrl             string `json:"help_url"`
	Icon                string `json:"icon"`
	Id                  string `json:"id"`
	Images              string `json:"images"`
	InstalledCount      string `json:"installed_count"`
	Introduction        string `json:"introduction"`
	LatestCompatible    string `json:"latest_compatible"`
	LatestVersion       string `json:"latest_version"`
	LatestVersionDetail struct {
		Compatible         string `json:"compatible"`
		CreateTime         string `json:"create_time"`
		DownloadUrl        string `json:"download_url"`
		Enabled            string `json:"enabled"`
		Id                 string `json:"id"`
		InstalledCount     string `json:"installed_count"`
		Memory             string `json:"memory"`
		Name               string `json:"name"`
		Permission         string `json:"permission"`
		Remark             string `json:"remark"`
		UpdateTime         string `json:"update_time"`
		UpgradeDescription string `json:"upgrade_description"`
		Version            string `json:"version"`
	} `json:"latest_version_detail"`
	LatestVersionHistory []struct {
		Compatible         string `json:"compatible"`
		CreateTime         string `json:"create_time"`
		DownloadUrl        string `json:"download_url"`
		Enabled            string `json:"enabled"`
		Id                 string `json:"id"`
		InstalledCount     string `json:"installed_count"`
		Memory             string `json:"memory"`
		Name               string `json:"name"`
		Permission         string `json:"permission"`
		Remark             string `json:"remark"`
		UpdateTime         string `json:"update_time"`
		UpgradeDescription string `json:"upgrade_description"`
		Version            string `json:"version"`
	} `json:"latest_version_history"`
	Name       string `json:"name"`
	Overview   string `json:"overview"`
	Remark     string `json:"remark"`
	Sort       string `json:"sort"`
	Title      string `json:"title"`
	Type       string `json:"type"`
	UpdateTime string `json:"update_time"`
}
