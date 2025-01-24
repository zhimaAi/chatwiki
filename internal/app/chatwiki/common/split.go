// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"bytes"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/spf13/cast"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/xuri/excelize/v2"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"golang.org/x/image/webp"
)

func GetLibFileSplit(userId, fileId int, splitParams define.SplitParams, lang string) (list []define.DocSplitItem, wordTotal int, err error) {
	info, err := GetLibFileInfo(fileId, userId)
	if err != nil {
		err = errors.New(i18n.Show(lang, `sys_err`))
		return
	}
	if len(info) == 0 {
		err = errors.New(i18n.Show(lang, `no_data`))
		return
	}
	if !tool.InArrayInt(cast.ToInt(info[`status`]), []int{define.FileStatusWaitSplit, define.FileStatusLearned}) {
		err = errors.New(i18n.Show(lang, `status_exception`))
		return
	}
	splitParams.IsTableFile = cast.ToInt(info[`is_table_file`])
	splitParams, err = CheckSplitParams(splitParams, lang)
	if err != nil {
		return
	}

	if cast.ToInt(info[`is_table_file`]) == define.FileIsTable && splitParams.IsQaDoc == define.DocTypeQa {
		list, wordTotal, err = ReadQaTab(info[`file_url`], info[`file_ext`], splitParams)
	} else if cast.ToInt(info[`is_table_file`]) == define.FileIsTable && splitParams.IsQaDoc != define.DocTypeQa {
		list, wordTotal, err = ReadTab(info[`file_url`], info[`file_ext`])
	} else {
		if len(info[`html_url`]) == 0 || !LinkExists(info[`html_url`]) { //compatible with old data
			list, wordTotal, err = ConvertAndReadHtmlContent(cast.ToInt(info[`id`]), info[`file_url`], userId)
		} else {
			list, wordTotal, err = ReadHtmlContent(info[`html_url`], userId)
		}
	}

	if err != nil {
		return
	}
	if len(list) == 0 || wordTotal == 0 {
		err = errors.New(i18n.Show(lang, `doc_empty`))
		return
	}

	//initialize RecursiveCharacter
	split := textsplitter.NewRecursiveCharacter()
	split.Separators = append(splitParams.Separators)
	split.ChunkSize = splitParams.ChunkSize
	split.ChunkOverlap = splitParams.ChunkOverlap
	// split by document type
	if splitParams.IsQaDoc == define.DocTypeQa {
		if cast.ToInt(info[`is_table_file`]) != define.FileIsTable {
			list = QaDocSplit(splitParams, list)
		}
	} else {
		list = MultDocSplit(split, list)
	}

	for i := range list {
		list[i].Number = i + 1 //serial number
		if splitParams.IsQaDoc == define.DocTypeQa {
			list[i].WordTotal = utf8.RuneCountInString(list[i].Question) + utf8.RuneCountInString(list[i].Answer)
		} else {
			list[i].WordTotal = utf8.RuneCountInString(list[i].Content)
		}
	}

	return
}

func SaveLibFileSplit(userId, fileId, wordTotal, qaIndexType int, splitParams define.SplitParams, list []define.DocSplitItem, lang string) error {
	info, err := GetLibFileInfo(fileId, userId)
	if err != nil {
		logs.Error(err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(info) == 0 {
		return errors.New(i18n.Show(lang, `no_data`))
	}
	if !tool.InArrayInt(cast.ToInt(info[`status`]), []int{define.FileStatusWaitSplit, define.FileStatusLearned}) {
		return errors.New(i18n.Show(lang, `status_exception`))
	}

	//check params
	if splitParams.IsQaDoc == define.DocTypeQa { // qa
		for i := range list {
			list[i].Number = i + 1 //serial number
			list[i].WordTotal = utf8.RuneCountInString(list[i].Question + list[i].Answer)
			if utf8.RuneCountInString(list[i].Question) < 1 || utf8.RuneCountInString(list[i].Question) > MaxContent {
				return errors.New(i18n.Show(lang, `length_err`, i+1))
			}
			if utf8.RuneCountInString(list[i].Answer) < 1 || utf8.RuneCountInString(list[i].Answer) > MaxContent {
				return errors.New(i18n.Show(lang, `length_err`, i+1))
			}
		}
	} else {
		for i := range list {
			list[i].Number = i + 1 //serial number
			list[i].WordTotal = utf8.RuneCountInString(list[i].Content)
			if list[i].WordTotal < 1 || list[i].WordTotal > MaxContent {
				return errors.New(i18n.Show(lang, `length_err`, i+1))
			}
		}
	}

	if splitParams.IsQaDoc == define.DocTypeQa {
		if qaIndexType != define.QAIndexTypeQuestionAndAnswer && qaIndexType != define.QAIndexTypeQuestion {
			return errors.New(i18n.Show(lang, `param_invalid`, `qa_index_type`))
		}
	}

	//add lock dispose
	if !lib_redis.AddLock(define.Redis, define.LockPreKey+`SaveLibFileSplit`+cast.ToString(fileId), time.Minute*5) {
		err = errors.New(i18n.Show(lang, `op_lock`))
	}
	defer func(fileId int) {
		lib_redis.UnLock(define.Redis, define.LockPreKey+`SaveLibFileSplit`+cast.ToString(fileId))
	}(fileId)

	//database dispose
	m := msql.Model(`chat_ai_library_file`, define.Postgres)
	err = m.Begin()
	defer func() {
		_ = m.Rollback()
	}()

	if err != nil {
		logs.Error(err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}

	data := msql.Datas{
		`status`:               define.FileStatusLearning,
		`errmsg`:               `success`,
		`word_total`:           wordTotal,
		`split_total`:          len(list),
		`is_qa_doc`:            splitParams.IsQaDoc,
		`is_diy_split`:         splitParams.IsDiySplit,
		`separators_no`:        splitParams.SeparatorsNo,
		`chunk_size`:           splitParams.ChunkSize,
		`chunk_overlap`:        splitParams.ChunkOverlap,
		`question_lable`:       splitParams.QuestionLable,
		`answer_lable`:         splitParams.AnswerLable,
		`question_column`:      splitParams.QuestionColumn,
		`answer_column`:        splitParams.AnswerColumn,
		`enable_extract_image`: splitParams.EnableExtractImage,
		`update_time`:          tool.Time2Int(),
	}
	if qaIndexType != 0 {
		data[`qa_index_type`] = qaIndexType
	}

	_, err = m.Where(`id`, cast.ToString(fileId)).Update(data)
	if err != nil {
		logs.Error(err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &LibFileCacheBuildHandler{FileId: fileId})

	//database dispose
	vm := msql.Model("chat_ai_library_file_data", define.Postgres)
	_, err = vm.Where(`admin_user_id`, cast.ToString(userId)).Where(`file_id`, cast.ToString(fileId)).Delete()
	if err != nil {
		logs.Error(err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	_, err = msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Where(`admin_user_id`, cast.ToString(userId)).Where(`file_id`, cast.ToString(fileId)).Delete()
	if err != nil {
		logs.Error(err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}

	var indexIds []int64
	for i, item := range list {
		if utf8.RuneCountInString(item.Content) > MaxContent || utf8.RuneCountInString(item.Question) > MaxContent || utf8.RuneCountInString(item.Answer) > MaxContent {
			return errors.New(i18n.Show(lang, `length_err`, i+1))
		}

		data := msql.Datas{
			`admin_user_id`: info[`admin_user_id`],
			`library_id`:    info[`library_id`],
			`file_id`:       fileId,
			`number`:        item.Number,
			`page_num`:      item.PageNum,
			`title`:         item.Title,
			`word_total`:    item.WordTotal,
			`create_time`:   tool.Time2Int(),
			`update_time`:   tool.Time2Int(),
		}
		if splitParams.IsQaDoc == define.DocTypeQa {
			if splitParams.IsTableFile == define.FileIsTable {
				data[`type`] = define.ParagraphTypeExcelQA
			} else {
				data[`type`] = define.ParagraphTypeDocQA
			}
			data[`question`] = strings.TrimSpace(item.Question)
			data[`answer`] = strings.TrimSpace(item.Answer)
			if len(item.Images) > 0 {
				jsonImages, err := CheckLibraryImage(item.Images)
				if err != nil {
					_ = m.Rollback()
					return errors.New(i18n.Show(lang, `param_invalid`, `images`))
				}
				data[`images`] = jsonImages
			} else {
				data[`images`] = `[]`
			}
			id, err := vm.Insert(data, `id`)
			if err != nil {
				logs.Error(err.Error())
				return errors.New(i18n.Show(lang, `sys_err`))
			}
			vectorID, err := SaveVector(
				cast.ToInt64(info[`admin_user_id`]),
				cast.ToInt64(info[`library_id`]),
				cast.ToInt64(fileId),
				id,
				cast.ToString(define.VectorTypeQuestion),
				strings.TrimSpace(item.Question),
			)
			if err != nil {
				logs.Error(err.Error())
				return errors.New(i18n.Show(lang, `sys_err`))
			}
			indexIds = append(indexIds, vectorID)
			if qaIndexType == define.QAIndexTypeQuestionAndAnswer {
				vectorID, err = SaveVector(
					cast.ToInt64(info[`admin_user_id`]),
					cast.ToInt64(info[`library_id`]),
					cast.ToInt64(fileId),
					id,
					cast.ToString(define.VectorTypeAnswer),
					strings.TrimSpace(item.Answer),
				)
				if err != nil {
					logs.Error(err.Error())
					return errors.New(i18n.Show(lang, `sys_err`))
				}
				indexIds = append(indexIds, vectorID)
			}
		} else {
			data[`type`] = define.ParagraphTypeNormal
			data[`content`] = strings.TrimSpace(item.Content)
			if len(item.Images) > 0 {
				jsonImages, err := CheckLibraryImage(item.Images)
				if err != nil {
					return errors.New(i18n.Show(lang, `param_invalid`, `images`))
				}
				data[`images`] = jsonImages
			} else {
				data[`images`] = `[]`
			}
			id, err := vm.Insert(data, `id`)
			if err != nil {
				logs.Error(err.Error())
				return errors.New(i18n.Show(lang, `sys_err`))
			}
			vectorID, err := SaveVector(
				cast.ToInt64(info[`admin_user_id`]),
				cast.ToInt64(info[`library_id`]),
				cast.ToInt64(fileId),
				id,
				cast.ToString(define.VectorTypeParagraph),
				strings.TrimSpace(item.Content),
			)
			if err != nil {
				logs.Error(err.Error())
				return errors.New(i18n.Show(lang, `sys_err`))
			}
			indexIds = append(indexIds, vectorID)
		}
	}
	err = m.Commit()
	if err != nil {
		logs.Error(err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}

	//async task:convert vector
	for _, id := range indexIds {
		if message, err := tool.JsonEncode(map[string]any{`id`: id, `file_id`: fileId}); err != nil {
			logs.Error(err.Error())
		} else if err := AddJobs(define.ConvertVectorTopic, message); err != nil {
			logs.Error(err.Error())
		}
	}

	return nil
}

func MultDocSplit(split textsplitter.TextSplitter, items []define.DocSplitItem) []define.DocSplitItem {
	list := make([]define.DocSplitItem, 0)
	for _, item := range items {
		contents, _ := split.SplitText(item.Content)
		for _, content := range contents {
			if len(content) == 0 {
				continue
			}
			content, images := ExtractTextImages(content)
			if len(content) == 0 {
				continue
			}
			list = append(list, define.DocSplitItem{PageNum: item.PageNum, Content: content, Images: images})
		}
	}
	return list
}

func ConvertAndReadHtmlContent(fileId int, fileUrl string, userId int) ([]define.DocSplitItem, int, error) {
	htmlUrl, err := ConvertHtml(fileUrl, userId)
	if err != nil {
		return nil, 0, err
	}

	_, err = msql.Model(`chat_ai_library_file`, define.Postgres).Where(`id`, cast.ToString(fileId)).Update(msql.Datas{
		`html_url`:    htmlUrl,
		`status`:      define.FileStatusWaitSplit,
		`update_time`: tool.Time2Int(),
	})
	if err != nil {
		return nil, 0, err
	}
	lib_redis.DelCacheData(define.Redis, &LibFileCacheBuildHandler{FileId: fileId})

	return ReadHtmlContent(htmlUrl, userId)
}

func ConvertHtml(link string, userId int) (string, error) {
	if !LinkExists(link) {
		return ``, errors.New(`file not exist:` + link)
	}
	ext := strings.ToLower(strings.TrimLeft(filepath.Ext(link), `.`))
	request := curl.Post(define.Config.WebService[`converter`]+`/convert`).
		PostFile(`file`, GetFileByLink(link)).
		Param(`from_format`, ext).
		Param(`to_format`, `html`)
	content, err := request.String()
	if err != nil {
		return ``, err
	}
	resp, err := request.Response()
	if err != nil {
		return ``, err
	}
	if resp.StatusCode != http.StatusOK {
		return ``, errors.New(content)
	}
	objectKey := fmt.Sprintf(`chat_ai/%d/%s/%s/%s.html`, userId,
		`convert`, tool.Date(`Ym`), tool.MD5(content))
	url, err := WriteFileByString(objectKey, content)
	if err != nil {
		return ``, err
	}
	return url, nil
}

func ReadHtmlContent(htmlUrl string, userId int) ([]define.DocSplitItem, int, error) {
	if !LinkExists(htmlUrl) {
		return nil, 0, errors.New(`file not exist:` + htmlUrl)
	}

	content, err := tool.ReadFile(GetFileByLink(htmlUrl))
	if err != nil {
		return nil, 0, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return nil, 0, err
	}

	doc.Find("img").Each(func(index int, item *goquery.Selection) {
		src, exists := item.Attr("src")
		if exists && strings.HasPrefix(src, "data:image") {
			parts := strings.Split(src, ";")
			if len(parts) < 2 {
				logs.Debug(fmt.Sprintf("could not find base64 data"))
				return
			}
			format := strings.TrimPrefix(parts[0], "data:image/")
			base64Data := strings.TrimPrefix(parts[1], "base64,")
			imgData, err := base64.StdEncoding.DecodeString(base64Data)
			if err != nil {
				logs.Error(err.Error())
				return
			}
			if format == `webp` {
				imgData, err = ConvertWebPToPNG(imgData)
				if err != nil {
					logs.Error(err.Error())
					return
				}
				format = `png`
			}

			// save to file
			objectKey := fmt.Sprintf(`chat_ai/%d/%s/%s/%s.%s`, userId, `library_image`, tool.Date(`Ym`), tool.MD5(string(imgData)), format)
			imgUrl, err := WriteFileByString(objectKey, string(imgData))
			if err != nil {
				logs.Error(err.Error())
				return
			}

			// Replace img tag with a span tag
			newTag := fmt.Sprintf("<b>{{!!%s!!}}</b>", imgUrl)
			item.ReplaceWithHtml(newTag)
		}
	})

	content, err = doc.Html()
	if err != nil {
		logs.Error(err.Error())
		return nil, 0, err
	}

	if !utf8.ValidString(content) {
		content = tool.Convert(content, `gbk`, `utf-8`)
	}
	content = strip.StripTags(content)
	list := []define.DocSplitItem{{Content: content}}
	return list, utf8.RuneCountInString(content), nil
}

func ParseTabFile(fileUrl, fileExt string) ([][]string, error) {
	if len(fileUrl) == 0 {
		return nil, errors.New(`file link cannot be empty`)
	}
	if !LinkExists(fileUrl) {
		return nil, errors.New(`file not exist:` + fileUrl)
	}
	//read excel
	var rows = make([][]string, 0)
	if fileExt == `csv` {
		content, err := tool.ReadFile(GetFileByLink(fileUrl))
		if err != nil {
			return nil, err
		}
		if !utf8.ValidString(content) {
			content = tool.Convert(content, `gbk`, `utf-8`)
		}
		for _, val := range strings.Fields(content) {
			rows = append(rows, strings.Split(val, `,`))
		}
	} else {
		f, err := excelize.OpenFile(GetFileByLink(fileUrl))
		if err != nil {
			return nil, err
		}
		rows, err = f.GetRows(f.GetSheetName(f.GetActiveSheetIndex()))
		if err != nil {
			return nil, err
		}
	}
	return rows, nil
}

func ReadTab(fileUrl, fileExt string) ([]define.DocSplitItem, int, error) {
	rows, err := ParseTabFile(fileUrl, fileExt)
	if err != nil {
		return nil, 0, err
	}
	if len(rows) < 2 {
		return nil, 0, errors.New(`excel no less than 2 lines`)
	}
	//line collection
	list := make([]define.DocSplitItem, 0)
	wordTotal := 0
	for i := range rows {
		pairs := make([]string, 0)
		for j := range rows[i] {
			wordTotal += utf8.RuneCountInString(rows[i][j])
			if i == 0 { //excel head
				continue
			}
			if len(rows[i][j]) == 0 || len(rows[0]) <= j {
				continue
			}
			pairs = append(pairs, fmt.Sprintf(`%s:%s`, rows[0][j], rows[i][j]))
		}
		if len(pairs) == 0 {
			continue
		}
		list = append(list, define.DocSplitItem{PageNum: i, Content: strings.Join(pairs, `;`)})
	}
	return list, wordTotal, nil
}

func ColumnIndexFromIdentifier(identifier string) (int, error) {
	if identifier == "" {
		return -1, errors.New("identifier cannot be empty")
	}
	colIndex := 0
	for _, char := range identifier {
		if !unicode.IsUpper(char) || char < 'A' || char > 'Z' {
			return -1, errors.New("invalid identifier: identifiers must be uppercase letters (A-Z)")
		}
		colIndex = colIndex*26 + int(char-'A'+1)
	}
	return colIndex - 1, nil
}

func IdentifierFromColumnIndex(index int) (string, error) {
	if index < 0 {
		return "", errors.New("index cannot be negative")
	}
	var identifier strings.Builder
	for index >= 0 {
		charIndex := (index % 26) + 'A'
		identifier.WriteByte(byte(charIndex))
		index = (index / 26) - 1
	}
	result := identifier.String()
	runes := []rune(result)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes), nil
}

func ReadQaTab(fileUrl, fileExt string, splitParams define.SplitParams) ([]define.DocSplitItem, int, error) {
	rows, err := ParseTabFile(fileUrl, fileExt)
	if err != nil {
		return nil, 0, err
	}
	if len(rows) < 2 {
		return nil, 0, errors.New(`excel no less than 2 lines`)
	}
	questionIndex, err := ColumnIndexFromIdentifier(splitParams.QuestionColumn)
	if err != nil {
		return nil, 0, err
	}
	answerIndex, err := ColumnIndexFromIdentifier(splitParams.AnswerColumn)
	if err != nil {
		return nil, 0, err
	}
	if questionIndex == answerIndex {
		return nil, 0, errors.New(`excel question index cannot be equal to answer`)
	}

	//line collection
	list := make([]define.DocSplitItem, 0)
	wordTotal := 0
	for i, row := range rows[1:] {
		var question, answer string
		if len(row) > answerIndex {
			answer = row[answerIndex]
		}
		if len(row) > questionIndex {
			question = row[questionIndex]
		}
		answer, images := ExtractTextImages(answer)
		if len(answer) == 0 || len(question) == 0 {
			continue
		}

		wordTotal += utf8.RuneCountInString(question + answer)
		list = append(list, define.DocSplitItem{PageNum: i + 1, Question: question, Answer: answer, Images: images})
	}

	return list, wordTotal, nil
}

func QaDocSplit(splitParams define.SplitParams, items []define.DocSplitItem) []define.DocSplitItem {
	list := make([]define.DocSplitItem, 0)
	for i, item := range items {
		for _, section := range strings.Split(item.Content, splitParams.QuestionLable) {
			if len(strings.TrimSpace(section)) == 0 {
				continue
			}
			qa := strings.SplitN(section, splitParams.AnswerLable, 2)
			var question, answer string
			if len(qa) == 0 {
				continue
			} else if len(qa) == 1 {
				question = qa[0]
				continue
			} else if len(qa) == 2 {
				question = qa[0]
				answer = qa[1]
			} else {
				continue
			}
			answer, images := ExtractTextImages(answer)
			if len(question) == 0 || len(answer) == 0 {
				continue
			}
			list = append(list, define.DocSplitItem{PageNum: i + 1, Question: question, Answer: answer, Images: images})
		}
	}
	return list
}

func ExtractTextImages(content string) (string, []string) {
	re := regexp.MustCompile(`\{\{\!\!(.+?)\!\!\}\}`)
	matches := re.FindAllStringSubmatch(content, -1)
	images := make([]string, 0)
	for _, match := range matches {
		if len(match) > 1 {
			images = append(images, match[1])
		}
	}
	content = re.ReplaceAllString(content, "")
	return content, images
}

func EmbTextImages(content string, images []string) string {
	var imgTags []string
	for _, image := range images {
		imgTags = append(imgTags, fmt.Sprintf(`<img src="%s">`, image))
	}

	return content + "\n" + strings.Join(imgTags, " ")
}

func MbSubstr(s string, start, length int) string {
	runes := []rune(s)
	if start >= len(runes) {
		return ""
	}
	if start+length > len(runes) {
		length = len(runes) - start
	}
	return string(runes[start : start+length])
}

func ConvertWebPToPNG(webpData []byte) ([]byte, error) {
	img, err := webp.Decode(bytes.NewReader(webpData))
	if err != nil {
		return nil, fmt.Errorf("error decoding WebP image: %w", err)
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("error encoding PNG image: %w", err)
	}

	return buf.Bytes(), nil
}
