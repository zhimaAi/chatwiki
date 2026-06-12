// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/ArtisanCloud/PowerLibs/v3/os"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"gopkg.in/yaml.v3"
)

func InitClawbotDirs(robotKey string) {
	_ = tool.MkDirAll(define.PublicSkillsDir)
	_ = tool.MkDirAll(strings.ReplaceAll(define.PrivateSkillsDir, `<robot_key>`, robotKey))
	privateFileDir := strings.ReplaceAll(define.PrivateFileDir, `<robot_key>`, robotKey)
	_ = tool.MkDirAll(privateFileDir)
	_ = tool.MkDirAll(strings.ReplaceAll(define.PrivateWorkDir, `<robot_key>`, robotKey))
	// copy query-local-docs skill template
	if !tool.IsFile(privateFileDir + `/../SKILL.md`) {
		err := os.CopyFile(define.AppRoot+`data/template/query-local-docs.md`, privateFileDir+`/../SKILL.md`)
		if err != nil {
			logs.Error(err.Error())
		}
	}
	// generate query-local-docs skill index.yaml
	if !tool.IsFile(privateFileDir + `/index.yaml`) {
		err := tool.WriteFile(privateFileDir+`/index.yaml`, define.QueryLocalDocsIndexDesc)
		if err != nil {
			logs.Error(err.Error())
		}
	}
}

func GetClawbotRobotKey(c *gin.Context, adminUserId int, id int64) (_ string, _ bool) {
	// check required
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(GetLang(c), `param_lack`))))
		return
	}
	// data check
	m := msql.Model(`chat_ai_robot`, define.Postgres)
	robotKey, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`application_type`, cast.ToString(define.ApplicationTypeClaw)).Value(`robot_key`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(GetLang(c), `sys_err`))))
		return
	}
	if len(robotKey) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(GetLang(c), `no_data`))))
		return
	}
	return robotKey, true
}

func SaveClawbotLocalDoc(fileHeader *multipart.FileHeader, robotKey string) (*define.ClawbotLocalDocInfo, error) {
	if fileHeader == nil {
		return nil, errors.New(`file header is nil`)
	}
	name, ok := NormalizeClawbotLocalDocName(fileHeader.Filename)
	if !ok {
		return nil, errors.New(`file name is invalid`)
	}
	ext := strings.ToLower(strings.TrimLeft(path.Ext(name), `.`))
	if !tool.InArrayString(ext, define.ClawbotLocalDocAllowExt) {
		return nil, errors.New(ext + ` not allow`)
	}
	if fileHeader.Size > int64(define.LibFileLimitSize) {
		return nil, errors.New(`file size too big`)
	}
	reader, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer func(reader multipart.File) {
		_ = reader.Close()
	}(reader)
	bs, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if len(bs) == 0 {
		return nil, errors.New(`file content is empty`)
	}
	filePath := strings.ReplaceAll(define.PrivateFileDir, `<robot_key>`, robotKey) + `/` + name
	if err = tool.WriteFile(filePath, string(bs)); err != nil {
		return nil, err
	}
	return &define.ClawbotLocalDocInfo{Name: name, Size: fileHeader.Size, Ext: ext}, nil
}

func SaveClawbotLocalDocIndex(robotKey string, docInfo define.ClawbotLocalDocInfo, description string, keywords []string) error {
	return updateClawbotLocalDocIndex(robotKey, func(list []define.ClawbotLocalDocIndexItem) []define.ClawbotLocalDocIndexItem {
		list = removeClawbotLocalDocIndexItem(list, docInfo.Name)
		return append(list, define.ClawbotLocalDocIndexItem{
			File:        docInfo.Name,
			Type:        clawbotLocalDocIndexType(docInfo.Ext),
			Description: strings.TrimSpace(description),
			Keywords:    normalizeClawbotLocalDocKeywords(keywords...),
			Updated:     time.Now().Format(`2006-01-02`),
		})
	})
}

func DeleteClawbotLocalDocIndex(robotKey string, name string) error {
	return updateClawbotLocalDocIndex(robotKey, func(list []define.ClawbotLocalDocIndexItem) []define.ClawbotLocalDocIndexItem {
		return removeClawbotLocalDocIndexItem(list, name)
	})
}

func normalizeClawbotLocalDocKeywords(values ...string) []string {
	keywords := make([]string, 0)
	seen := make(map[string]struct{})
	for _, value := range values {
		for _, keyword := range splitClawbotLocalDocKeywordValue(value) {
			if _, ok := seen[keyword]; ok {
				continue
			}
			seen[keyword] = struct{}{}
			keywords = append(keywords, keyword)
		}
	}
	return keywords
}

func NormalizeClawbotLocalDocName(name string) (string, bool) {
	if name == `index.yaml` {
		return ``, false // uploading and deleting index files are prohibited
	}
	name = strings.ReplaceAll(strings.TrimSpace(name), `\`, `/`)
	if name == `` || name == `/` || name == `.` || name == `..` || name != path.Base(name) || strings.ContainsRune(name, 0) {
		return ``, false
	}
	return name, true
}

func updateClawbotLocalDocIndex(robotKey string, modify func([]define.ClawbotLocalDocIndexItem) []define.ClawbotLocalDocIndexItem) error {
	indexPath := strings.ReplaceAll(define.PrivateFileDir, `<robot_key>`, robotKey) + `/index.yaml`
	list := make([]define.ClawbotLocalDocIndexItem, 0)
	if tool.IsFile(indexPath) {
		content, err := tool.ReadFile(indexPath)
		if err != nil {
			return err
		}
		if strings.TrimSpace(content) != `` {
			if err = yaml.Unmarshal([]byte(content), &list); err != nil {
				return err
			}
		}
	}

	list = modify(list)
	content := define.QueryLocalDocsIndexDesc
	if len(list) > 0 {
		bs, err := yaml.Marshal(list)
		if err != nil {
			return err
		}
		content += string(bs)
	}
	return tool.WriteFile(indexPath, content)
}

func removeClawbotLocalDocIndexItem(list []define.ClawbotLocalDocIndexItem, name string) []define.ClawbotLocalDocIndexItem {
	out := make([]define.ClawbotLocalDocIndexItem, 0, len(list))
	for _, item := range list {
		if item.File == name {
			continue
		}
		out = append(out, item)
	}
	return out
}

func clawbotLocalDocIndexType(ext string) string {
	switch strings.ToLower(strings.TrimLeft(ext, `.`)) {
	case `md`:
		return `Markdown`
	case `txt`:
		return `Text`
	case `pdf`:
		return `PDF`
	case `docx`:
		return `DOCX`
	case `xlsx`:
		return `XLSX`
	default:
		return strings.ToUpper(ext)
	}
}

func splitClawbotLocalDocKeywordValue(value string) []string {
	value = strings.TrimSpace(value)
	if value == `` {
		return nil
	}
	if strings.HasPrefix(value, `[`) {
		var arr []string
		if err := json.Unmarshal([]byte(value), &arr); err == nil {
			return normalizeClawbotLocalDocKeywords(arr...)
		}
	}
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == '，' || r == ';' || r == '；' || r == '\n' || r == '\r'
	})
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != `` {
			out = append(out, part)
		}
	}
	return out
}
