// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package define

import (
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
)

type MenuJsonStruct struct {
	Content  string   `json:"content"`
	Question []string `json:"question"`
}

type UploadInfo struct {
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	Ext    string `json:"ext"`
	Link   string `json:"link"`
	Online bool   `json:"-"`
	DocUrl string `json:"-"`
	Custom bool   `json:"-"`
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
	Error      error
	Lang       string
	Question   string
	DialogueId int
	Prompt     string
	LibraryIds string
	IsClose    *bool
}

type DocSplitItem struct {
	Number    int    `json:"number"`
	PageNum   int    `json:"page_num"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	WordTotal int    `json:"word_total"`
}

type SplitParams struct {
	IsTableFile    int      `json:"is_table_file"`
	IsDiySplit     int      `json:"is_diy_split"`
	SeparatorsNo   string   `json:"separators_no"`
	Separators     []string `json:"-"`
	ChunkSize      int      `json:"chunk_size"`
	ChunkOverlap   int      `json:"chunk_overlap"`
	IsQaDoc        int      `json:"is_qa_doc"`
	QuestionLable  string   `json:"question_lable"`
	AnswerLable    string   `json:"answer_lable"`
	QuestionColumn string   `json:"question_column"`
	AnswerColumn   string   `json:"answer_column"`
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
