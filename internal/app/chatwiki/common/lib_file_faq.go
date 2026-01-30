// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
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
	"github.com/zhimaAi/go_tools/curl"
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
		return fmt.Errorf(`extract_type error : %d`, splitParams.ExtractType)
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
					Content: prompt + define.ExtractLibFaqFilesPrompt,
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
				define.LangEnUs,
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

func ExportFAQFileAllQA(lang string, list []msql.Params, ext, source string) (string, error) {
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
			para.AddFormattedText(i18n.Show(lang, `faq_question`)+`：`, &document.TextFormat{Bold: true, FontSize: 14})
			para.AddFormattedText(lineBreak+params[`question`], &document.TextFormat{FontSize: 12})
			//相似问法
			para.AddFormattedText(lineBreak+i18n.Show(lang, `faq_similar_questions`)+`：`, &document.TextFormat{Bold: true, FontSize: 14})
			for _, item := range DisposeStringList(params[`similar_questions`]) {
				para.AddFormattedText(lineBreak+item, &document.TextFormat{FontSize: 12})
			}
			//答案
			para.AddFormattedText(lineBreak+i18n.Show(lang, `faq_answer`)+`：`, &document.TextFormat{Bold: true, FontSize: 14})
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
		fields := tool.Fields{
			{Field: `group_name`, Header: i18n.Show(lang, `faq_group_name`)},
			{Field: `question`, Header: i18n.Show(lang, `faq_question`)},
			{Field: `similar_questions`, Header: i18n.Show(lang, `faq_similar_questions`)},
			{Field: `answer`, Header: i18n.Show(lang, `faq_answer`)},
		}
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

func ImportFAQFile(adminUserId, libraryId, fileId int, ids, token string, isSync bool) {
	if isSync {
		if message, err := tool.JsonEncode(map[string]any{
			`admin_user_id`: adminUserId,
			`token`:         token,
			`library_id`:    libraryId,
			`file_id`:       fileId,
			"ids":           ids,
		}); err != nil {
			logs.Error(err.Error())
		} else if err := AddJobs(define.ImportFAQFileTopic, message); err != nil {
			logs.Error(err.Error())
		}
		return
	}

	var qaList []msql.Params
	var err error

	if len(ids) > 0 {
		// 查询单条QA数据
		qaList, err = msql.Model("chat_ai_faq_files_data_qa", define.Postgres).
			Where("id", `in`, cast.ToString(ids)).
			Where("admin_user_id", cast.ToString(adminUserId)).
			Select()
	} else {
		// 查询文件下所有QA数据
		qaList, err = msql.Model("chat_ai_faq_files_data_qa", define.Postgres).
			Where("file_id", cast.ToString(fileId)).
			Where("admin_user_id", cast.ToString(adminUserId)).
			Select()
	}
	if err != nil {
		logs.Error(err.Error())
		return
	}
	if len(qaList) == 0 {
		return
	}
	var qaIds []string
	var res lib_web.Response
	for _, item := range qaList {
		var imgs []string
		req := curl.Post(fmt.Sprintf(`http://127.0.0.1:%s/manage/addParagraph`, define.Config.WebService[`port`])).Header(`token`, token)
		req.Param(`library_id`, cast.ToString(libraryId))
		req.Param(`similar_questions`, `[]`)
		req.Param(`question`, item[`question`])
		req.Param(`answer`, item[`answer`])
		if len(item[`images`]) > 0 {
			if err = tool.JsonDecodeUseNumber(cast.ToString(item[`images`]), &imgs); err != nil {
				logs.Error(err.Error())
				continue
			}
			for _, img := range imgs {
				req.Param(`images`, img)
			}
		}
		if err = req.ToJSON(&res); err != nil {
			logs.Error(err.Error())
			continue
		}
		if res.Res != define.StatusOK {
			logs.Error(fmt.Sprintf(`%s`, cast.ToString(res.Msg)))
			continue
		}
		qaIds = append(qaIds, cast.ToString(item[`id`]))
	}
	// 更新数据
	_, err = msql.Model("chat_ai_faq_files_data_qa", define.Postgres).
		Where("id", `in`, strings.Join(qaIds, `,`)).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Update(msql.Datas{
			"is_import":   define.SwitchOn,
			`library_id`:  libraryId,
			"update_time": tool.Time2Int(),
		})
	if err != nil {
		logs.Error(err.Error())
	}
}
