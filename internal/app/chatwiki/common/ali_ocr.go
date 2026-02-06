// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"bytes"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"errors"
	"fmt"
	"os"
	"sort"

	openClient "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/darabonba-openapi/v2/models"
	"github.com/alibabacloud-go/docmind-api-20220711/client"
	"github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/spf13/cast"
	"github.com/yuin/goldmark"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetOcrConfig(key, secret string) (models.Config, error) {
	akCredential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId(key).
		SetAccessKeySecret(secret))
	if err != nil {
		return models.Config{}, err
	}
	accessKeyId, err := akCredential.GetAccessKeyId()
	accessSecret, err := akCredential.GetAccessKeySecret()
	var endpoint = "docmind-api.cn-hangzhou.aliyuncs.com"
	return openClient.Config{AccessKeyId: accessKeyId, AccessKeySecret: accessSecret, Endpoint: &endpoint}, nil
}

func CheckAliOcr(key, secret string) error {
	config, err := GetOcrConfig(key, secret)
	if err != nil {
		return err
	}
	cli, err := client.NewClient(&config)
	if err != nil {
		return err
	}
	id := "test"
	request := client.GetDocStructureResultRequest{Id: &id}
	_, err = cli.GetDocStructureResult(&request)
	if err != nil {
		return err
	}
	return nil
}

func SubmitOdcParserJob(lang string, userId int, fileUrl string) (string, error) {
	company, err := msql.Model(`company`, define.Postgres).Where(`parent_id`, cast.ToString(userId)).Find()
	if err != nil {
		return "", err
	}

	if len(company) == 0 || cast.ToInt(company[`ali_ocr_switch`]) != 1 {
		return "", errors.New(i18n.Show(lang, `ali_ocr_not_open`))
	}

	config, err := GetOcrConfig(company[`ali_ocr_key`], company[`ali_ocr_secret`])
	if err != nil {
		return "", err
	}

	cli, err := client.NewClient(&config)
	if err != nil {
		return "", err
	}

	url := GetFileByLink(fileUrl)
	f, err := os.Open(url)
	if err != nil {
		return "", err
	}

	request := client.SubmitDocParserJobAdvanceRequest{
		FileName:      &url,
		FileUrlObject: f,
	}

	response, err := cli.SubmitDocParserJobAdvance(&request, &service.RuntimeOptions{})
	if err != nil {
		return "", err
	}

	return *response.Body.Data.Id, nil
}

func QueryAliOcrResult(aliOcrKey, aliOcrSecret, aliOcrJobId string) (*client.GetDocParserResultResponse, error) {
	config, err := GetOcrConfig(aliOcrKey, aliOcrSecret)
	if err != nil {
		return nil, err
	}

	cli, err := client.NewClient(&config)
	if err != nil {
		return nil, err
	}

	var layoutStepSize int32 = 3000
	var layoutNum int32 = 0
	request := client.GetDocParserResultRequest{Id: &aliOcrJobId, LayoutNum: &layoutNum, LayoutStepSize: &layoutStepSize}
	response, err := cli.GetDocParserResult(&request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func QueryAliOcrProgress(aliOcrKey, aliOcrSecret, aliOcrJobId string) (*client.QueryDocParserStatusResponse, error) {
	config, err := GetOcrConfig(aliOcrKey, aliOcrSecret)
	if err != nil {
		return nil, err
	}

	cli, err := client.NewClient(&config)
	if err != nil {
		return nil, err
	}

	request := client.QueryDocParserStatusRequest{Id: &aliOcrJobId}
	response, err := cli.QueryDocParserStatus(&request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func QueryAndParseAliOcrRequest(file msql.Params, aliOcrKey, aliOcrSecret string) error {
	progress, err := QueryAliOcrProgress(aliOcrKey, aliOcrSecret, file[`ali_ocr_job_id`])
	if err != nil {
		return err
	}

	if *progress.Body.Data.Status != "success" {
		_, err := msql.Model(`chat_ai_library_file`, define.Postgres).
			Where(`id`, file[`id`]).
			Update(msql.Datas{`ocr_pdf_index`: *progress.Body.Data.PageCountEstimate + 1})
		if err != nil {
			logs.Error(err.Error())
			return err
		}
		return nil
	}

	response, err := QueryAliOcrResult(aliOcrKey, aliOcrSecret, file[`ali_ocr_job_id`])
	if err != nil {
		return err
	}
	htmlContent := generateOcrHtmlContent(response.Body.Data)
	htmlContent, err = ReplaceRemoteImg(htmlContent, cast.ToInt(file[`admin_user_id`]))
	if err != nil {
		return err
	}

	objectKey := fmt.Sprintf(`chat_ai/%d/%s/%s/%s.html`, cast.ToInt(file[`admin_user_id`]),
		`convert`, tool.Date(`Ym`), tool.MD5(htmlContent))
	url, err := WriteFileByString(objectKey, htmlContent)

	// Update file status to pending split
	_, err = msql.Model(`chat_ai_library_file`, define.Postgres).
		Where(`id`, file[`id`]).
		Update(msql.Datas{
			`status`:   define.FileStatusWaitSplit,
			`html_url`: url,
		})
	if err != nil {
		return err
	}

	//create default lib file split
	splitParams := DefaultSplitParams()
	if len(file[`async_split_params`]) > 0 {
		if err = tool.JsonDecodeUseNumber(file[`async_split_params`], &splitParams); err != nil {
			logs.Error(err.Error())
		}
	}
	AutoSplitLibFile(cast.ToInt(file[`admin_user_id`]), cast.ToInt(file[`id`]), splitParams)

	return nil
}

// generateOcrHtmlContent generates HTML content from OCR recognition result
func generateOcrHtmlContent(data map[string]interface{}) string {
	// Build HTML header
	htmlContent := `<html><head><meta charset="utf-8"></head><body>`

	// Process layouts data
	layouts := extractLayouts(data)

	// Group layouts by page number
	pageLayouts := groupLayoutsByPage(layouts)

	// Get all page numbers and sort them
	pageNums := make([]int, 0, len(pageLayouts))
	for pageNum := range pageLayouts {
		pageNums = append(pageNums, pageNum)
	}
	sort.Ints(pageNums)

	// Process each page's layouts in order to ensure all page numbers are included
	for pageNum := 0; pageNum <= pageNums[len(pageNums)-1]; pageNum++ {
		htmlContent += "<meta charset=\"UTF-8\"/>\n"

		layouts, exists := pageLayouts[pageNum]
		if exists {
			htmlContent += "<div>" + processPageLayouts(layouts) + "</div>"
		} else {
			// Add empty content for missing page numbers
			htmlContent += "<div></div>"
		}
	}

	htmlContent += "</body></html>"
	return htmlContent
}

// extractLayouts extracts layouts data
func extractLayouts(data map[string]interface{}) []map[string]interface{} {
	var layouts []map[string]interface{}

	// Get layouts array
	layoutsInterface, hasLayouts := data["layouts"]
	if !hasLayouts {
		return layouts
	}

	// Convert to array type
	layoutsArray, ok := layoutsInterface.([]interface{})
	if !ok {
		return layouts
	}

	// Convert each layout to map
	for _, layoutInterface := range layoutsArray {
		if layout, ok := layoutInterface.(map[string]interface{}); ok {
			layouts = append(layouts, layout)
		}
	}

	return layouts
}

// groupLayoutsByPage groups layouts by page number
func groupLayoutsByPage(layouts []map[string]interface{}) map[int][]map[string]interface{} {
	pageLayouts := make(map[int][]map[string]interface{})

	for _, layout := range layouts {
		pageNum := cast.ToInt(layout["pageNum"])
		pageLayouts[pageNum] = append(pageLayouts[pageNum], layout)
	}

	return pageLayouts
}

// processPageLayouts processes page layouts
func processPageLayouts(layouts []map[string]interface{}) string {
	md := ""
	for _, layout := range layouts {
		md += cast.ToString(layout["markdownContent"])
	}

	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		logs.Error(err.Error())
		return ""
	}

	return buf.String()
}
