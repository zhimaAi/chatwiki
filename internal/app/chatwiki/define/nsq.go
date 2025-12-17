// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

import "github.com/zhimaAi/go_tools/mq"

const ConvertHtmlTopic = `chatwiki_convert_html_topic`
const ConvertHtmlChannel = `convert_html_channel`

const ConvertVectorTopic = `chatwiki_convert_vector_topic`
const ConvertVectorChannel = `convert_vector_channel`

const ConvertGraphTopic = `chatwiki_convert_graph_topic`
const ConvertGraphChannel = `convert_graph_channel`

const CrawlArticleTopic = `chatwiki_crawl_article_topic`
const CrawlArticleChannel = `chatwiki_crawl_article_channel`

const ExportTaskTopic = `chatwiki_export_task_topic`
const ExportTaskChannel = `chatwiki_export_task_channel`

const ExtractFaqFilesTopic = `chatwiki_extract_faq_files_topic`
const ExtractFaqFilesChannel = `chatwiki_extract_faq_files_channel`

const ImportFAQFileTopic = `chatwiki_import_faq_file_topic`
const ImportFAQFileChannel = `chatwiki_import_faq_file_channel`

const OfficialAccountDraftSyncTopic = `chatwiki_official_account_draft_sync_topic`
const OfficialAccountDraftSyncChannel = `chatwiki_official_account_draft_sync_channel`

const OfficialAccountBatchSendTopic = `chatwiki_official_account_batch_send_topic`
const OfficialAccountBatchSendChannel = `chatwiki_official_account_batch_send_channel`

const OfficialAccountCommentSyncTopic = `chatwiki_official_account_comment_sync_topic`
const OfficialAccountCommentSyncChannel = `chatwiki_official_account_comment_sync_channel`

const OfficialAccountCommentAiCheckTopic = `chatwiki_official_account_comment_ai_check_topic`
const OfficialAccountCommentAiCheckChannel = `chatwiki_official_account_comment_ai_check_channel`

var ConsumerHandle *mq.ConsumerHandle
var ProducerHandle *mq.ProducerHandle
