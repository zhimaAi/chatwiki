// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/textsplitter"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/ZeroHawkeye/wordZero/pkg/document"
	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func UpdateLibFileFaq(id, adminUserId int, data msql.Datas) error {
	if len(data) == 0 {
		return nil
	}
	data[`update_time`] = tool.Time2Int()
	_, err := msql.Model(`chat_ai_faq_fils`, define.Postgres).Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Update(data)
	return err
}

func UpdateLibFileFaqStatus(id, adminUserId, status int, errMsg string) error {
	_, err := msql.Model(`chat_ai_faq_files`, define.Postgres).Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Update(msql.Datas{
		`status`:      status,
		`errmsg`:      errMsg,
		`update_time`: tool.Time2Int()})
	return err
}

func GetLibFileFaqSplit(fileId, userId int, splitParams define.SplitFaqParams) (list define.DocSplitItems, err error) {
	info, err := GetFaqFilesInfo(fileId, userId)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	if len(info) == 0 {
		return nil, err
	}
	splitParams.FileExt = info[`file_ext`]
	if define.IsDocxFile(info[`file_ext`]) {
		list, _, err = ReadDocx(info[`file_url`], userId)
	} else if define.IsTxtFile(info[`file_ext`]) || define.IsMdFile(info[`file_ext`]) {
		list, _, err = ReadTxt(info[`file_url`])
	} else {
		list, _, err = ReadHtmlContent(info[`html_url`], userId)
	}
	if err != nil {
		return
	}
	if len(list) == 0 {
		return
	}
	return
}

func MultSplitFaqFiles(list define.DocSplitItems, splitParams define.SplitFaqParams) (listSplit define.DocSplitItems, err error) {
	if len(list) == 0 {
		return
	}
	// split by document type
	if splitParams.ChunkType == define.FAQChunkTypeLength {
		spliter := NewAiSpliterClient(splitParams.ChunkSize)
		for _, item := range list {
			chunks, err := spliter.SplitText(item.Content)
			if err != nil {
				logs.Error(err.Error())
				continue
			}
			for _, chunk := range chunks {
				content, images := ExtractTextImagesPlaceholders(chunk)
				if len(content) == 0 {
					continue
				}
				if len(images) == 0 && len(item.Images) > 0 {
					images = append(images, item.Images...)
				}
				spliteItem := define.DocSplitItem{
					PageNum: item.PageNum,
					Content: strings.TrimSpace(content),
					Images:  images,
				}
				listSplit = append(listSplit, spliteItem)
			}
		}
	} else if splitParams.ChunkType == define.FAQChunkTypeSeparatorsNo {
		spliter := textsplitter.NewRecursiveCharacter()
		separators, _ := GetSeparatorsByNo(splitParams.SeparatorsNo, ``)
		spliter.Separators = append(separators, spliter.Separators...)
		spliter.ChunkSize = splitParams.ChunkSize
		for _, item := range list {
			chunks, err := spliter.SplitText(item.Content)
			if err != nil {
				logs.Error(err.Error())
				continue
			}
			for _, chunk := range chunks {
				content, images := ExtractTextImagesPlaceholders(chunk)
				if len(content) == 0 {
					continue
				}
				if len(images) == 0 && len(item.Images) > 0 {
					images = append(images, item.Images...)
				}
				spliteItem := define.DocSplitItem{
					PageNum: item.PageNum,
					Content: strings.TrimSpace(content),
					Images:  images,
				}
				listSplit = append(listSplit, spliteItem)
			}
		}
	}
	listSplit.UnifySetNumber() //对number统一编号
	for i := range listSplit {
		if splitParams.IsQaDoc == define.DocTypeQa {
			listSplit[i].WordTotal = utf8.RuneCountInString(listSplit[i].Question) + utf8.RuneCountInString(listSplit[i].Answer)
		} else {
			listSplit[i].WordTotal = utf8.RuneCountInString(listSplit[i].Content)
		}
	}
	return listSplit, nil
}

func ExtractLibFaqFiles(adminUserId int, splitParams define.SplitFaqParams, submitList define.DocSplitItems, results chan define.DocSplitItem) error {
	defer close(results)
	if splitParams.ExtractType != define.FAQExtractTypeAI {
		return errors.New(`提取方式不对`)
	}
	wg := &sync.WaitGroup{}
	lock := &sync.Mutex{}
	contents := ""
	var errMsg = ""
	contentMap := make(map[int]define.DocSplitItem)
	index := 1
	maxToken := 0
	currChan := make(chan struct{}, 8)
	if !CheckModelIsValid(adminUserId, cast.ToInt(splitParams.ChunkModelConfigId), splitParams.ChunkModel, Llm) {
		errMsg = `model not valid`
		return errors.New(errMsg)
	}
	prompt := define.PromptFaqFileAiChunk
	if splitParams.ChunkPrompt != "" {
		prompt = splitParams.ChunkPrompt
	}
	for _, item := range submitList {
		contents = item.Content
		// ai split
		if len(contents) == 0 {
			continue
		}
		if len(contents) <= 2 {
			lock.Lock()
			contentMap[index] = define.DocSplitItem{
				PageNum: item.PageNum,
				Content: contents,
				Images:  item.Images,
			}
			lock.Unlock()
			continue
		}
		currChan <- struct{}{}
		wg.Add(1)
		go func(wg *sync.WaitGroup, contentMap map[int]define.DocSplitItem, index int, contents string, errMsg *string) {
			defer wg.Done()
			defer func() { <-currChan }()
			messages := []adaptor.ZhimaChatCompletionMessage{
				{
					Role:    `system`,
					Content: prompt + `\n 如果文本中有[图片占位符]请原样返回,严格按照输出示例返回,输出示例:[{"question":"问题1","answer":"回答1"},{"question":"问题2","answer":"回答2"}]`,
				},
				{
					Role:    `user`,
					Content: contents,
				},
			}
			ctx, cancel := context.WithCancel(context.Background())
			retryTimes := 0
			chanStream := make(chan sse.Event, 100)
			go func() {
				for {
					select {
					case <-chanStream:
						continue
					case <-ctx.Done():
						return
					}
				}
			}()
		Retry:
			chatResp, _, err := RequestChatStream(
				adminUserId,
				"",
				msql.Params{},
				"",
				cast.ToInt(splitParams.ChunkModelConfigId),
				splitParams.ChunkModel,
				messages,
				nil,
				chanStream,
				0.1,
				maxToken,
			)
			docSplitItem := define.DocSplitItem{
				PageNum: item.PageNum,
				Content: contents,
				Images:  item.Images,
			}
			if err != nil && len(chatResp.Result) <= 0 {
				logs.Error(err.Error())
				if retryTimes <= 3 {
					retryTimes++
					goto Retry
				}
				*errMsg = err.Error()
				docSplitItem.AiChunkErrMsg = err.Error()
			} else {
				chatResp.Result, _ = strings.CutPrefix(chatResp.Result, "```json")
				chatResp.Result, _ = strings.CutSuffix(chatResp.Result, "```")
				docSplitItem.Answer = chatResp.Result
			}
			results <- docSplitItem
			cancel()
		}(wg, contentMap, index, contents, &errMsg)
		index++
	}
	wg.Wait()
	return nil
}

func ExportFAQFileAllQA(list []msql.Params, ext, source string) (string, error) {
	fileSavePath := `static/public/export/` + source
	if define.IsDocxFile(ext) {
		var lineBreak = "\r" //docx的换行符,比较特殊
		doc := document.New()
		var imageConfig = &document.ImageConfig{Size: &document.ImageSize{Width: 145, KeepAspectRatio: true}}
		for idx, params := range list {
			if idx > 0 { //添加两个换行
				doc.AddParagraph(lineBreak + lineBreak)
			}
			//问题
			para := doc.AddParagraph(``)
			para.AddFormattedText(`问题：`, &document.TextFormat{Bold: true, FontSize: 14})
			para.AddFormattedText(lineBreak+params[`question`], &document.TextFormat{FontSize: 12})
			//相似问法
			para.AddFormattedText(lineBreak+`相似问法：`, &document.TextFormat{Bold: true, FontSize: 14})
			for _, item := range DisposeStringList(params[`similar_questions`]) {
				para.AddFormattedText(lineBreak+item, &document.TextFormat{FontSize: 12})
			}
			//答案
			para.AddFormattedText(lineBreak+`答案：`, &document.TextFormat{Bold: true, FontSize: 14})
			para.AddFormattedText(lineBreak+params[`answer`], &document.TextFormat{FontSize: 12})
			//图片
			for _, imgUrl := range DisposeStringList(params[`images`]) {
				if !LinkExists(imgUrl) {
					continue
				}
				if _, err := doc.AddImageFromFile(GetFileByLink(imgUrl), imageConfig); err != nil {
					logs.Error(err.Error())
				}
			}
		}
		md := tool.MD5(tool.JsonEncodeNoError(list) + time.Now().String() + tool.Random(10))
		filepath := fileSavePath + `/` + md[:2] + `/` + md[2:] + `.docx`
		return filepath, doc.Save(filepath)
	} else {
		fields := tool.Fields{{Field: "group_name", Header: "问答分组"}, {Field: "question", Header: "问题"}, {Field: "similar_questions", Header: "相似问法"}, {Field: "answer", Header: "答案"}}
		data := make([]map[string]any, len(list))
		for idx, params := range list {
			for _, imgUrl := range DisposeStringList(params[`images`]) {
				params[`answer`] += fmt.Sprintf("\r\n{{!!%s!!}}", imgUrl)
			}
			data[idx] = map[string]any{
				`group_name`: params[`group_name`],
				`question`:   params[`question`], `answer`: params[`answer`],
				`similar_questions`: strings.Join(DisposeStringList(params[`similar_questions`]), "\r\n"),
			}
		}
		filepath, _, err := tool.ExcelExportPro(data, fields, `FAQ`, fileSavePath)
		return filepath, err
	}
}
