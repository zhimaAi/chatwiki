// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"baliance.com/gooxml/document"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/xuri/excelize/v2"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/pdf"
)

func MultDocSplit(split textsplitter.TextSplitter, items []define.DocSplitItem) []define.DocSplitItem {
	list := make([]define.DocSplitItem, 0)
	for _, item := range items {
		contents, _ := split.SplitText(item.Content)
		for _, content := range contents {
			if len(content) > 0 {
				list = append(list, define.DocSplitItem{PageNum: item.PageNum, Content: content})
			}
		}
	}
	return list
}

func ReadDocx(fileUrl string) ([]define.DocSplitItem, int, error) {
	if len(fileUrl) == 0 {
		return nil, 0, errors.New(`file_url cannot be empty`)
	}
	doc, err := document.Open(GetFileByLink(fileUrl))
	if err != nil {
		return nil, 0, err
	}
	var content string
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			content += run.Text()
		}
		content += "\r\n"
	}
	list := []define.DocSplitItem{{Content: content}}
	return list, utf8.RuneCountInString(content), nil
}

func ReadTxt(fileUrl string, stripTags bool) ([]define.DocSplitItem, int, error) {
	if len(fileUrl) == 0 {
		return nil, 0, errors.New(`file_url cannot be empty`)
	}
	content, err := tool.ReadFile(GetFileByLink(fileUrl))
	if err != nil {
		return nil, 0, err
	}
	if !utf8.ValidString(content) {
		content = tool.Convert(content, `gbk`, `utf-8`)
	}
	if stripTags { //clean up html tags
		content = strip.StripTags(content)
	}
	list := []define.DocSplitItem{{Content: content}}
	return list, utf8.RuneCountInString(content), nil
}

func ReadPdf(pdfUrl string) ([]define.DocSplitItem, int, error) {
	if len(pdfUrl) == 0 {
		return nil, 0, errors.New(`file link cannot be empty`)
	}
	//read pdf
	file, reader, err := pdf.Open(GetFileByLink(pdfUrl))
	if err != nil {
		return nil, 0, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	//paging collection
	list := make([]define.DocSplitItem, 0)
	wordTotal := 0
	for num := 1; num <= reader.NumPage(); num++ {
		p := reader.Page(num)
		if p.V.IsNull() {
			continue
		}
		content, err := p.GetPlainText(nil)
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		if len(content) > 0 {
			wordTotal += utf8.RuneCountInString(content)
			list = append(list, define.DocSplitItem{PageNum: num, Content: content})
		}
	}
	return list, wordTotal, nil
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
		var q, a string
		if len(row) > answerIndex {
			a = row[answerIndex]
		}
		if len(row) > questionIndex {
			q = row[questionIndex]
		}
		if len(a) == 0 || len(q) == 0 {
			continue
		}

		wordTotal += utf8.RuneCountInString(q + a)
		//question := fmt.Sprintf("%s:%s", rows[0][questionIndex], rows[i+1][questionIndex])
		//answer := fmt.Sprintf("%s:%s", rows[0][answerIndex], rows[i+1][answerIndex])
		list = append(list, define.DocSplitItem{PageNum: i + 1, Question: q, Answer: a})
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
			var q, a string
			if len(qa) == 0 {
				continue
			}
			if len(qa) == 1 {
				q = qa[0]
			}
			if len(qa) == 2 {
				q = qa[0]
				a = qa[1]
			}
			list = append(list, define.DocSplitItem{PageNum: i + 1, Question: q, Answer: a})
		}
	}
	return list
}
