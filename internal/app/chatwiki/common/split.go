// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cast"
	"github.com/syyongx/php2go"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/msql"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/xuri/excelize/v2"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

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

	return ReadHtmlContent(htmlUrl, userId)
}

func ConvertHtml(link string, userId int) (string, error) {
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
			// parse base64 image
			dataPos := php2go.Strpos(src, `base64,`, 0)
			if dataPos < 0 {
				logs.Debug(fmt.Sprintf("could not find base64 data"))
				return
			}
			base64Data := php2go.Substr(src, uint(dataPos)+7, -1)
			imgData, err := base64.StdEncoding.DecodeString(base64Data)
			if err != nil {
				logs.Error(err.Error())
				return
			}

			// save to png file
			objectKey := fmt.Sprintf(`chat_ai/%d/%s/%s/%s.png`, userId, `library_image`, tool.Date(`Ym`), tool.MD5(string(imgData)))
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
