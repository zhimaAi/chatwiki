// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package define

import "github.com/zhimaAi/go_tools/mq"

const ConvertPdfTopic = `chatwiki_convert_pdf_topic`
const ConvertPdfChannel = `convert_pdf_channel`

const ConvertVectorTopic = `chatwiki_convert_vector_topic`
const ConvertVectorChannel = `convert_vector_channel`

const CrawlArticleTopic = `chatwiki_crawl_article_topic`
const CrawlArticleChannel = `chatwiki_crawl_article_channel`

var ConsumerHandle *mq.ConsumerHandle
var ProducerHandle *mq.ProducerHandle
