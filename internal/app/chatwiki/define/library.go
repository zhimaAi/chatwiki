// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package define

type ChunkParam struct {
	SetChunk                         string `form:"set_chunk"`
	ChunkType                        string `form:"chunk_type"`
	NormalChunkDefaultSeparatorsNo   string `form:"normal_chunk_default_separators_no"`
	NormalChunkDefaultChunkSize      string `form:"normal_chunk_default_chunk_size"`
	NormalChunkDefaultChunkOverlap   string `form:"normal_chunk_default_chunk_overlap"`
	NormalChunkDefaultNotMergedText  string `form:"normal_chunk_default_not_merged_text"`
	SemanticChunkDefaultChunkSize    string `form:"semantic_chunk_default_chunk_size"`
	SemanticChunkDefaultChunkOverlap string `form:"semantic_chunk_default_chunk_overlap"`
	SemanticChunkDefaultThreshold    string `form:"semantic_chunk_default_threshold"`
	AiChunkPrumpt                    string `form:"ai_chunk_prumpt"`
	AiChunkModel                     string `form:"ai_chunk_model"`
	AiChunkModelConfigId             string `form:"ai_chunk_model_config_id"`
	AiChunkSize                      string `form:"ai_chunk_size"`
	QaIndexType                      string `form:"qa_index_type"`
	FatherChunkParagraphType         string `form:"father_chunk_paragraph_type"`
	FatherChunkSeparatorsNo          string `form:"father_chunk_separators_no"`
	FatherChunkChunkSize             string `form:"father_chunk_chunk_size"`
	SonChunkSeparatorsNo             string `form:"son_chunk_separators_no"`
	SonChunkChunkSize                string `form:"son_chunk_chunk_size"`
}

type AddFileParam struct {
	LibraryId             string `form:"library_id"`
	LibraryKey            string `form:"library_key"`
	DocType               string `form:"doc_type"`
	Urls                  string `form:"urls"`
	FileName              string `form:"file_name"`
	Content               string `form:"content"`
	Title                 string `form:"title"`
	IsQaDoc               string `form:"is_qa_doc"`
	QaIndexType           string `form:"qa_index_type"`
	DocAutoRenewFrequency string `form:"doc_auto_renew_frequency"`
	DocAutoRenewMinute    string `form:"doc_auto_renew_minute"`
	AnswerLable           string `form:"answer_lable"`
	AnswerColumn          string `form:"answer_column"`
	QuestionLable         string `form:"question_lable"`
	QuestionColumn        string `form:"question_column"`
	SimilarColumn         string `form:"similar_column"`
	SimilarLabel          string `form:"similar_label"`
	PdfParseType          string `form:"pdf_parse_type"`
	GroupId               string `form:"group_id"`
}
