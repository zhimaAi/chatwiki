// Copyright © 2016- 2024 Sesame Network Technology all right reserved

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

var ConsumerHandle *mq.ConsumerHandle
var ProducerHandle *mq.ProducerHandle
