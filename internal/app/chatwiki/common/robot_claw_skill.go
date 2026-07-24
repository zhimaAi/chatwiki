// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"archive/zip"
	"chatwiki/internal/app/chatwiki/define"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"gopkg.in/yaml.v3"
)

type ClawbotUserSkillItem struct {
	SkillId        int64  `json:"skill_id"`
	RobotSkillId   int64  `json:"robot_skill_id"`
	SkillName      string `json:"skill_name"`
	RemarkName     string `json:"remark_name"`
	Intro          string `json:"intro"`
	Description    string `json:"description"`
	FileSize       int    `json:"file_size"`
	OriginFileName string `json:"origin_file_name"`
	IsSelected     int    `json:"is_selected"`
	IsMine         int    `json:"is_mine"`
	CreateTime     int    `json:"create_time"`
	UpdateTime     int    `json:"update_time"`
}

type ClawbotSkillUploadMeta struct {
	AdminUserId    int    `json:"admin_user_id"`
	OriginFileName string `json:"origin_file_name"`
	FileSize       int64  `json:"file_size"`
	SkillName      string `json:"skill_name"`
	Description    string `json:"description"`
	UploadTime     int64  `json:"upload_time"`
	ExpireTime     int64  `json:"expire_time"`
}

func privateSkillsDir(robotKey string) string {
	return strings.ReplaceAll(define.PrivateSkillsDir, `<robot_key>`, robotKey)
}

func userSkillsDir(adminUserId int) string {
	return strings.ReplaceAll(define.UserSkillsDir, `<admin_user_id>`, cast.ToString(adminUserId))
}

func userSkillTmpRoot(adminUserId int) string {
	return filepath.Join(userSkillsDir(adminUserId), define.SkillTmpDir)
}

func skillTmpRoot(robotKey string) string {
	return filepath.Join(privateSkillsDir(robotKey), define.SkillTmpDir)
}

type frontMatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

func parseSkillFrontmatter(content string) (name, description, body string, err error) {
	const delimiter = `---`
	data := strings.TrimSpace(content)
	if !strings.HasPrefix(data, delimiter) {
		return ``, ``, ``, errors.New(`file does not start with frontmatter delimiter`)
	}
	rest := data[len(delimiter):]
	endIdx := strings.Index(rest, "\n"+delimiter)
	if endIdx == -1 {
		return ``, ``, ``, errors.New(`frontmatter closing delimiter not found`)
	}
	fmText := strings.TrimSpace(rest[:endIdx])
	body = strings.TrimPrefix(rest[endIdx+len("\n"+delimiter):], "\n")
	var fm frontMatter
	if err = yaml.Unmarshal([]byte(fmText), &fm); err != nil {
		return ``, ``, ``, err
	}
	return strings.TrimSpace(fm.Name), strings.TrimSpace(fm.Description), body, nil
}

// rewriteSkillMdFrontmatter updates the name and/or description fields in a
// SKILL.md frontmatter. An empty argument leaves that field untouched, so
// callers can rewrite only what changed (e.g. rename keeps the description,
// editing the intro only rewrites the description).
func rewriteSkillMdFrontmatter(skillMdPath, newName, newDescription string) error {
	content, err := tool.ReadFile(skillMdPath)
	if err != nil {
		return err
	}
	_, _, body, err := parseSkillFrontmatter(content)
	if err != nil {
		return err
	}
	const delimiter = `---`
	data := strings.TrimSpace(content)
	rest := data[len(delimiter):]
	endIdx := strings.Index(rest, "\n"+delimiter)
	fmText := strings.TrimSpace(rest[:endIdx])
	var node yaml.Node
	if err = yaml.Unmarshal([]byte(fmText), &node); err != nil {
		return err
	}
	if newName != `` {
		setYamlMapValue(&node, `name`, newName)
	}
	if newDescription != `` {
		setYamlMapValue(&node, `description`, newDescription)
	}
	newFm, err := yaml.Marshal(&node)
	if err != nil {
		return err
	}
	out := delimiter + "\n" + strings.TrimRight(string(newFm), "\n") + "\n" + delimiter + "\n" + body
	return tool.WriteFile(skillMdPath, out)
}

func setYamlMapValue(doc *yaml.Node, key, value string) {
	if doc.Kind == yaml.DocumentNode && len(doc.Content) > 0 {
		doc = doc.Content[0]
	}
	if doc.Kind != yaml.MappingNode {
		doc.Kind = yaml.MappingNode
	}
	for i := 0; i+1 < len(doc.Content); i += 2 {
		if doc.Content[i].Value == key {
			doc.Content[i+1].Kind = yaml.ScalarNode
			doc.Content[i+1].Tag = `!!str`
			doc.Content[i+1].Value = value
			return
		}
	}
	doc.Content = append(doc.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Tag: `!!str`, Value: key},
		&yaml.Node{Kind: yaml.ScalarNode, Tag: `!!str`, Value: value},
	)
}

type SkillUploadResult struct {
	SkillId        int64  `json:"skill_id"`
	UploadKey      string `json:"upload_key"`
	SkillName      string `json:"skill_name"`
	RemarkName     string `json:"remark_name"`
	Intro          string `json:"intro"`
	Description    string `json:"description"`
	FileSize       int64  `json:"file_size"`
	OriginFileName string `json:"origin_file_name"`
	ExpireTime     int64  `json:"expire_time"`
}

func UploadAndStageUserClawbotSkillZip(adminUserId int, originFileName, srcZipPath string, fileSize int64) (_ *SkillUploadResult, errKey string, err error) {
	reader, err := zip.OpenReader(srcZipPath)
	if err != nil {
		return nil, `clawbot_skill_zip_invalid`, nil
	}
	defer func() { _ = reader.Close() }()

	InitClawbotUserDirs(adminUserId)
	cleanExpiredSkillTmpRoot(userSkillTmpRoot(adminUserId))

	now := time.Now()
	uploadKey := fmt.Sprintf(`skup_%s_%s`, now.Format(`20060102150405`), tool.Random(16))
	tmpDir := filepath.Join(userSkillTmpRoot(adminUserId), uploadKey)
	extractDir := filepath.Join(tmpDir, `extract`)
	if err = tool.MkDirAll(extractDir); err != nil {
		return nil, ``, err
	}
	committed := false
	defer func() {
		if !committed {
			_ = os.RemoveAll(tmpDir)
		}
	}()

	if errKey, err = unzipSkillTo(reader, extractDir); err != nil || errKey != `` {
		return nil, errKey, err
	}
	skillRoot, err := resolveSkillRoot(extractDir)
	if err != nil {
		return nil, `clawbot_skill_md_missing`, nil
	}
	skillMdPath, errKey := normalizeSkillMdFile(skillRoot)
	if errKey != `` {
		return nil, errKey, nil
	}
	content, err := tool.ReadFile(skillMdPath)
	if err != nil {
		return nil, ``, err
	}
	skillName, description, _, perr := parseSkillFrontmatter(content)
	if perr != nil || skillName == `` || description == `` {
		return nil, `clawbot_skill_frontmatter_missing`, nil
	}
	if !define.SkillNameRegexp.MatchString(skillName) {
		return nil, `clawbot_skill_name_invalid`, nil
	}
	if skillName == define.SkillReservedName {
		return nil, `clawbot_skill_name_reserved`, nil
	}

	stageSkillDir := filepath.Join(tmpDir, `skill`)
	if err = os.Rename(skillRoot, stageSkillDir); err != nil {
		if err = copyDir(skillRoot, stageSkillDir); err != nil {
			return nil, ``, err
		}
	}
	meta := ClawbotSkillUploadMeta{
		AdminUserId:    adminUserId,
		OriginFileName: originFileName,
		FileSize:       fileSize,
		SkillName:      skillName,
		Description:    description,
		UploadTime:     now.Unix(),
		ExpireTime:     now.Add(define.SkillUploadKeyExpire).Unix(),
	}
	if err = writeUserSkillUploadMeta(tmpDir, meta); err != nil {
		return nil, ``, err
	}
	committed = true
	return &SkillUploadResult{
		UploadKey:      uploadKey,
		SkillName:      skillName,
		RemarkName:     skillName,
		Intro:          description,
		Description:    description,
		FileSize:       fileSize,
		OriginFileName: originFileName,
		ExpireTime:     meta.ExpireTime,
	}, ``, nil
}

type SaveUserClawbotSkillParam struct {
	AdminUserId int
	SkillId     int64
	SkillName   string
	RemarkName  string
	Description string
	UploadKey   string
}

func SaveUserClawbotSkill(p SaveUserClawbotSkillParam) (_ *ClawbotUserSkillItem, errKey string, err error) {
	InitClawbotUserDirs(p.AdminUserId)
	var meta *ClawbotSkillUploadMeta
	var tmpDir string
	if p.UploadKey != `` {
		meta, tmpDir, errKey = ReadUserSkillUploadMeta(p.AdminUserId, p.UploadKey)
		if errKey != `` {
			return nil, errKey, nil
		}
	} else if p.SkillId <= 0 {
		return nil, `clawbot_skill_upload_key_invalid`, nil
	}

	baseDir := userSkillsDir(p.AdminUserId)
	destDir := filepath.Join(baseDir, p.SkillName)
	tmpRoot := userSkillTmpRoot(p.AdminUserId)

	now := tool.Time2Int()
	// description is the single source of truth for the skill description:
	// it is stored in the DB and rewritten into SKILL.md frontmatter on
	// publish, so the value the agent reads always matches what the user
	// edited. The legacy `intro` column is no longer written; the json
	// `intro` field is still returned for backward compatibility.
	data := msql.Datas{
		`admin_user_id`: p.AdminUserId,
		`skill_name`:    p.SkillName,
		`remark_name`:   p.RemarkName,
		`description`:   p.Description,
		`update_time`:   now,
	}
	if meta != nil {
		data[`file_size`] = meta.FileSize
		data[`origin_file_name`] = meta.OriginFileName
	}

	old, err := getUserClawbotSkillTarget(p.AdminUserId, p.SkillId, p.SkillName)
	if err != nil {
		return nil, ``, err
	}
	if p.SkillId > 0 && len(old) == 0 {
		return nil, `no_data`, nil
	}
	if p.SkillId > 0 {
		dup, derr := msql.Model(define.TableChatAiClawbotUserSkill, define.Postgres).
			Where(`admin_user_id`, cast.ToString(p.AdminUserId)).
			Where(`skill_name`, p.SkillName).
			Where(`id`, `!=`, cast.ToString(p.SkillId)).
			Find()
		if derr != nil {
			return nil, ``, derr
		}
		if len(dup) > 0 {
			return nil, `clawbot_skill_system_dir_exists`, nil
		}
	}

	var cleanup func()
	rollback := func() {}
	if len(old) > 0 {
		oldName := old[`skill_name`]
		oldDir := filepath.Join(baseDir, oldName)
		if meta != nil {
			cleanup, rollback, err = replaceSkillDirWithBackup(tmpRoot, oldDir, destDir, filepath.Join(tmpDir, `skill`), p.SkillName, p.Description)
			if err != nil {
				return nil, ``, err
			}
		} else if oldName != p.SkillName {
			cleanup, rollback, err = renameUserSkillDirWithBackup(tmpRoot, oldDir, destDir, p.SkillName, p.Description)
			if err != nil {
				return nil, ``, err
			}
		} else if old[`description`] != p.Description {
			// only the description changed: rewrite SKILL.md in place
			if mdPath := findSkillMdInDir(oldDir); mdPath != `` {
				if err = rewriteSkillMdFrontmatter(mdPath, p.SkillName, p.Description); err != nil {
					return nil, ``, err
				}
			}
		}
		_, err = msql.Model(define.TableChatAiClawbotUserSkill, define.Postgres).
			Where(`id`, old[`id`]).
			Where(`admin_user_id`, cast.ToString(p.AdminUserId)).
			Update(data)
	} else {
		if meta == nil {
			return nil, `clawbot_skill_upload_key_invalid`, nil
		}
		cleanup, rollback, err = replaceSkillDirWithBackup(tmpRoot, ``, destDir, filepath.Join(tmpDir, `skill`), p.SkillName, p.Description)
		if err != nil {
			return nil, ``, err
		}
		data[`create_time`] = now
		_, err = msql.Model(define.TableChatAiClawbotUserSkill, define.Postgres).Insert(data, `id`)
	}
	if err != nil {
		rollback()
		return nil, ``, err
	}
	if cleanup != nil {
		cleanup()
	}
	if tmpDir != `` {
		_ = os.RemoveAll(tmpDir)
	}
	item, ierr := getUserClawbotSkillTarget(p.AdminUserId, p.SkillId, p.SkillName)
	if ierr != nil {
		return nil, ``, ierr
	}
	if len(item) == 0 {
		return nil, `no_data`, nil
	}
	out := BuildClawbotUserSkillItem(item)
	// propagate the updated user skill to every robot that has it mounted;
	// agent conversations only scan skills_robot/ and would otherwise keep
	// showing the stale name until the robot skill config is re-saved
	robotRows, rerr := msql.Model(define.TableChatAiClawbotSkill, define.Postgres).
		Where(`user_skill_id`, item[`id`]).
		Where(`source_type`, cast.ToString(define.SkillSourceTypeUpload)).
		Select()
	if rerr != nil {
		logs.Error(fmt.Sprintf(`sync user skill to robots failed (query), admin_user_id=%d skill_id=%s: %v`, p.AdminUserId, item[`id`], rerr))
	}
	for _, robotRow := range robotRows {
		if sErrKey, sErr := syncOneRobotSkill(p.AdminUserId, robotRow[`robot_key`], item, robotRow); sErrKey != `` || sErr != nil {
			logs.Error(fmt.Sprintf(`sync user skill to robot failed, admin_user_id=%d skill_id=%s robot_key=%s: errKey=%s err=%v`, p.AdminUserId, item[`id`], robotRow[`robot_key`], sErrKey, sErr))
		}
	}
	return &out, ``, nil
}

// ReadClawbotSkillZipMeta reads only the root SKILL.md. Generated-task completion
// must not apply the package size limits that belong to skill installation.
func ReadClawbotSkillZipMeta(srcZipPath string) (skillName, description string, err error) {
	reader, err := zip.OpenReader(srcZipPath)
	if err != nil {
		return ``, ``, errors.New(`invalid skill zip file`)
	}
	defer func() { _ = reader.Close() }()
	content, err := readSkillMarkdownFromZip(reader)
	if err != nil {
		return ``, ``, err
	}
	skillName, description, _, err = parseSkillFrontmatter(content)
	if err != nil || skillName == `` || description == `` {
		return ``, ``, errors.New(`skill name or description is missing from SKILL.md frontmatter`)
	}
	if len([]rune(skillName)) > 500 {
		return ``, ``, errors.New(`skill name exceeds 500 characters`)
	}
	return skillName, description, nil
}

func InstallUserClawbotSkillZip(adminUserId int, originFileName, srcZipPath, expectedSkillName, expectedDescription string, overwrite bool) (_ *ClawbotUserSkillItem, errKey string, err error) {
	if strings.ToLower(filepath.Ext(srcZipPath)) != `.zip` {
		return nil, `clawbot_skill_zip_only`, nil
	}
	fileInfo, err := os.Stat(srcZipPath)
	if err != nil {
		return nil, ``, err
	}
	if fileInfo.IsDir() || fileInfo.Size() <= 0 {
		return nil, `clawbot_skill_zip_invalid`, nil
	}
	if fileInfo.Size() > int64(define.MaxSkillZipSize) {
		return nil, `clawbot_skill_zip_too_big`, nil
	}
	if !define.SkillNameRegexp.MatchString(expectedSkillName) || expectedSkillName == define.SkillReservedName {
		return nil, `clawbot_skill_name_invalid`, nil
	}

	m := msql.Model(define.TableChatAiClawbotUserSkill, define.Postgres)
	dup, err := m.Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`skill_name`, expectedSkillName).
		Find()
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return nil, ``, err
	}
	if len(dup) > 0 && !overwrite {
		return nil, `clawbot_skill_name_exists`, nil
	}
	if len(dup) == 0 && tool.IsDir(filepath.Join(userSkillsDir(adminUserId), expectedSkillName)) {
		return nil, `clawbot_skill_system_dir_exists`, nil
	}

	result, errKey, err := UploadAndStageUserClawbotSkillZip(adminUserId, originFileName, srcZipPath, fileInfo.Size())
	if err != nil || errKey != `` {
		return nil, errKey, err
	}
	cleanupUpload := func() { _ = os.RemoveAll(filepath.Join(userSkillTmpRoot(adminUserId), result.UploadKey)) }
	if result.SkillName != expectedSkillName || result.Description != expectedDescription {
		cleanupUpload()
		return nil, ``, errors.New(`generated skill metadata does not match task record`)
	}
	skillId := int64(0)
	remarkName := result.RemarkName
	if len(dup) > 0 {
		skillId = cast.ToInt64(dup[`id`])
		if dup[`remark_name`] != `` {
			remarkName = dup[`remark_name`]
		}
	}
	item, errKey, err := SaveUserClawbotSkill(SaveUserClawbotSkillParam{
		AdminUserId: adminUserId,
		SkillId:     skillId,
		SkillName:   result.SkillName,
		RemarkName:  remarkName,
		Description: result.Description,
		UploadKey:   result.UploadKey,
	})
	if err != nil || errKey != `` {
		cleanupUpload()
		return nil, errKey, err
	}
	return item, ``, nil
}

func getUserClawbotSkillTarget(adminUserId int, skillId int64, skillName string) (msql.Params, error) {
	m := msql.Model(define.TableChatAiClawbotUserSkill, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId))
	if skillId > 0 {
		return m.Where(`id`, cast.ToString(skillId)).Find()
	}
	return m.Where(`skill_name`, skillName).Find()
}

func DeleteUserClawbotSkill(adminUserId int, skillId int64) (errKey string, err error) {
	item, ierr := msql.Model(define.TableChatAiClawbotUserSkill, define.Postgres).
		Where(`id`, cast.ToString(skillId)).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Find()
	if ierr != nil {
		return ``, ierr
	}
	if len(item) == 0 {
		return `no_data`, nil
	}
	skillName := item[`skill_name`]
	if !define.SkillNameRegexp.MatchString(skillName) {
		return `clawbot_skill_name_invalid`, nil
	}
	skillDir := filepath.Join(userSkillsDir(adminUserId), skillName)
	deleteDir, moved, merr := moveSkillDirToHiddenTmpRoot(userSkillTmpRoot(adminUserId), skillDir, `.delete_`)
	if merr != nil {
		return ``, merr
	}
	restoreDeleteDir := func() {
		if moved && !tool.IsDir(skillDir) {
			if rerr := os.Rename(deleteDir, skillDir); rerr != nil {
				logs.Error(fmt.Sprintf(`restore deleted user skill dir failed, skill_id=%d: %v`, skillId, rerr))
			}
		}
	}
	robotRows, rerr := msql.Model(define.TableChatAiClawbotSkill, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`source_type`, cast.ToString(define.SkillSourceTypeUpload)).
		Where(`user_skill_id`, cast.ToString(skillId)).
		Select()
	if rerr != nil {
		restoreDeleteDir()
		return ``, rerr
	}
	for _, row := range robotRows {
		if errKey, err = deleteRobotSkillRowAndDir(row[`robot_key`], row); err != nil || errKey != `` {
			restoreDeleteDir()
			return errKey, err
		}
	}
	rows, derr := msql.Model(define.TableChatAiClawbotUserSkill, define.Postgres).
		Where(`id`, cast.ToString(skillId)).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Delete()
	if derr != nil {
		restoreDeleteDir()
		return ``, derr
	}
	if rows <= 0 {
		restoreDeleteDir()
		return `no_data`, nil
	}
	if moved {
		_ = os.RemoveAll(deleteDir)
	}
	return ``, nil
}

func writeUserSkillUploadMeta(tmpDir string, meta ClawbotSkillUploadMeta) error {
	bs, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	return tool.WriteFile(filepath.Join(tmpDir, `metadata.json`), string(bs))
}

func ReadUserSkillUploadMeta(adminUserId int, uploadKey string) (_ *ClawbotSkillUploadMeta, tmpDir string, errKey string) {
	if !define.SkillUploadKeyRegexp.MatchString(uploadKey) {
		return nil, ``, `clawbot_skill_upload_key_invalid`
	}
	tmpDir = filepath.Join(userSkillTmpRoot(adminUserId), uploadKey)
	metaPath := filepath.Join(tmpDir, `metadata.json`)
	if !tool.IsFile(metaPath) {
		return nil, ``, `clawbot_skill_upload_key_invalid`
	}
	content, err := tool.ReadFile(metaPath)
	if err != nil {
		logs.Error(err.Error())
		return nil, ``, `clawbot_skill_upload_key_invalid`
	}
	var meta ClawbotSkillUploadMeta
	if err = json.Unmarshal([]byte(content), &meta); err != nil {
		logs.Error(err.Error())
		return nil, ``, `clawbot_skill_upload_key_invalid`
	}
	if meta.AdminUserId != adminUserId {
		return nil, ``, `clawbot_skill_upload_key_invalid`
	}
	if meta.ExpireTime > 0 && time.Now().Unix() > meta.ExpireTime {
		return nil, ``, `clawbot_skill_upload_key_expired`
	}
	return &meta, tmpDir, ``
}

func readSkillMarkdownFromZip(reader *zip.ReadCloser) (string, error) {
	topDirs := make(map[string]struct{})
	topFiles := make([]string, 0)
	skillFiles := make(map[string]*zip.File)
	for _, file := range reader.File {
		cleanName := filepath.ToSlash(filepath.Clean(file.Name))
		if cleanName == `.` || cleanName == `` || filepath.IsAbs(file.Name) || strings.HasPrefix(cleanName, `/`) || strings.HasPrefix(cleanName, `../`) {
			continue
		}
		parts := strings.Split(cleanName, `/`)
		if len(parts) == 0 || parts[0] == `__MACOSX` || strings.HasPrefix(parts[0], `.`) {
			continue
		}
		if len(parts) == 1 {
			if file.FileInfo().IsDir() {
				topDirs[parts[0]] = struct{}{}
			} else {
				topFiles = append(topFiles, parts[0])
			}
		} else {
			topDirs[parts[0]] = struct{}{}
		}
		if !file.FileInfo().IsDir() && strings.EqualFold(parts[len(parts)-1], define.SkillMdFileName) {
			skillFiles[cleanName] = file
		}
	}

	expectedPath := define.SkillMdFileName
	if len(topDirs) == 1 && len(topFiles) == 0 {
		for root := range topDirs {
			expectedPath = root + `/` + define.SkillMdFileName
		}
	}
	var skillFile *zip.File
	for name, file := range skillFiles {
		if !strings.EqualFold(name, expectedPath) {
			continue
		}
		if skillFile != nil {
			return ``, errors.New(`multiple SKILL.md files found at generated skill root`)
		}
		skillFile = file
	}
	if skillFile == nil {
		return ``, errors.New(`SKILL.md is missing from generated skill zip`)
	}
	fileReader, err := skillFile.Open()
	if err != nil {
		return ``, err
	}
	defer func() { _ = fileReader.Close() }()
	content, err := io.ReadAll(fileReader)
	if err != nil {
		return ``, err
	}
	return string(content), nil
}

// installation-time zip extraction guards
func unzipSkillTo(reader *zip.ReadCloser, dest string) (errKey string, err error) {
	cleanDest := filepath.Clean(dest)
	for _, f := range reader.File {
		cleanName := filepath.Clean(f.Name)
		if cleanName == `__MACOSX` || strings.HasPrefix(cleanName, `__MACOSX`+string(os.PathSeparator)) {
			continue
		}
		if cleanName == `.` || cleanName == `` || strings.HasPrefix(cleanName, `/`) || filepath.IsAbs(cleanName) {
			continue
		}
		fpath := filepath.Join(dest, cleanName)
		// path traversal guard
		if fpath != cleanDest && !strings.HasPrefix(fpath, cleanDest+string(os.PathSeparator)) {
			return `clawbot_skill_zip_invalid`, nil
		}
		// reject symlinks
		if !f.FileInfo().Mode().IsRegular() && !f.FileInfo().IsDir() {
			continue
		}
		if f.FileInfo().IsDir() {
			if err = os.MkdirAll(fpath, 0755); err != nil {
				return ``, err
			}
			continue
		}
		if err = os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return ``, err
		}
		rc, oerr := f.Open()
		if oerr != nil {
			return ``, oerr
		}
		outFile, oerr := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if oerr != nil {
			_ = rc.Close()
			return ``, oerr
		}
		_, cerr := io.Copy(outFile, rc)
		_ = outFile.Close()
		_ = rc.Close()
		if cerr != nil {
			return ``, cerr
		}
	}
	return ``, nil
}

func resolveSkillRoot(extractDir string) (string, error) {
	entries, err := os.ReadDir(extractDir)
	if err != nil {
		return ``, err
	}
	var validDirs, validFiles []string
	for _, e := range entries {
		name := e.Name()
		if name == `__MACOSX` || strings.HasPrefix(name, `.`) {
			continue
		}
		if e.IsDir() {
			validDirs = append(validDirs, name)
		} else {
			validFiles = append(validFiles, name)
		}
	}
	root := extractDir
	if len(validDirs) == 1 && len(validFiles) == 0 {
		root = filepath.Join(extractDir, validDirs[0])
	}
	if findSkillMdInDir(root) == `` {
		return ``, errors.New(`SKILL.md not found at skill root`)
	}
	return root, nil
}

func findSkillMdInDir(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ``
	}
	for _, e := range entries {
		if !e.IsDir() && strings.EqualFold(e.Name(), define.SkillMdFileName) {
			return filepath.Join(dir, e.Name())
		}
	}
	return ``
}

func normalizeSkillMdFile(skillRoot string) (string, string) {
	current := findSkillMdInDir(skillRoot)
	if current == `` {
		return ``, `clawbot_skill_md_missing`
	}
	target := filepath.Join(skillRoot, define.SkillMdFileName)
	if current != target {
		if err := os.Rename(current, target); err != nil {
			logs.Error(err.Error())
			return ``, `clawbot_skill_md_missing`
		}
	}
	return target, ``
}

func cleanExpiredSkillTmpRoot(root string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return
	}
	now := time.Now().Unix()
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		tmpDir := filepath.Join(root, e.Name())
		expired := false
		if strings.HasPrefix(name, `.upload_`) || strings.HasPrefix(name, `.publish_`) || strings.HasPrefix(name, `.bak_`) || strings.HasPrefix(name, `.delete_`) {
			if info, ierr := e.Info(); ierr == nil && now-info.ModTime().Unix() > int64((2*define.SkillUploadKeyExpire).Seconds()) {
				expired = true
			}
		} else {
			continue
		}
		if expired {
			_ = os.RemoveAll(tmpDir)
		}
	}
}

func GetClawbotUserSkillList(adminUserId int, robotKey string) ([]ClawbotUserSkillItem, error) {
	rows, err := msql.Model(define.TableChatAiClawbotUserSkill, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Order(`id desc`).
		Select()
	if err != nil {
		return nil, err
	}
	selected := make(map[int64]int64)
	if robotKey != `` {
		robotRows, rerr := msql.Model(define.TableChatAiClawbotSkill, define.Postgres).
			Where(`robot_key`, robotKey).
			Where(`source_type`, cast.ToString(define.SkillSourceTypeUpload)).
			Where(`user_skill_id`, `!=`, `0`).
			Select()
		if rerr != nil {
			return nil, rerr
		}
		for _, row := range robotRows {
			selected[cast.ToInt64(row[`user_skill_id`])] = cast.ToInt64(row[`id`])
		}
	}
	list := make([]ClawbotUserSkillItem, 0, len(rows))
	for _, row := range rows {
		item := BuildClawbotUserSkillItem(row)
		if robotSkillId := selected[item.SkillId]; robotSkillId > 0 {
			item.IsSelected = 1
			item.RobotSkillId = robotSkillId
		}
		list = append(list, item)
	}
	return list, nil
}

func GetClawbotUserSkillInfo(adminUserId int, skillId int64) (ClawbotUserSkillItem, error) {
	row, err := msql.Model(define.TableChatAiClawbotUserSkill, define.Postgres).
		Where(`id`, cast.ToString(skillId)).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Find()
	if err != nil || len(row) == 0 {
		return ClawbotUserSkillItem{}, err
	}
	return BuildClawbotUserSkillItem(row), nil
}

func BuildClawbotUserSkillItem(row msql.Params) ClawbotUserSkillItem {
	return ClawbotUserSkillItem{
		SkillId:        cast.ToInt64(row[`id`]),
		SkillName:      row[`skill_name`],
		RemarkName:     row[`remark_name`],
		Intro:          row[`description`],
		Description:    row[`description`],
		FileSize:       cast.ToInt(row[`file_size`]),
		OriginFileName: row[`origin_file_name`],
		IsMine:         1,
		CreateTime:     cast.ToInt(row[`create_time`]),
		UpdateTime:     cast.ToInt(row[`update_time`]),
	}
}

func publishSkillDirWithTmpRoot(stageRoot, srcSkillDir, destDir, skillName, description string) error {
	if !tool.IsDir(srcSkillDir) {
		return errors.New(`staged skill dir not found`)
	}
	if err := tool.MkDirAll(stageRoot); err != nil {
		return err
	}
	staging := filepath.Join(stageRoot, `.publish_`+tool.Random(8))
	if err := copyDir(srcSkillDir, staging); err != nil {
		_ = os.RemoveAll(staging)
		return err
	}
	if mdPath := findSkillMdInDir(staging); mdPath != `` {
		curName, curDesc, _, perr := readSkillMdName(mdPath)
		if perr == nil && (curName != skillName || curDesc != description) {
			if err := rewriteSkillMdFrontmatter(mdPath, skillName, description); err != nil {
				_ = os.RemoveAll(staging)
				return err
			}
		}
	}
	if err := os.Rename(staging, destDir); err != nil {
		_ = os.RemoveAll(staging)
		return err
	}
	return nil
}

func SyncClawbotRobotSkills(adminUserId int, robotKey string, selectedSkillIds []int64) (errKey string, err error) {
	InitClawbotDirs(robotKey)
	InitClawbotUserDirs(adminUserId)

	selected := make(map[int64]struct{})
	selectedIds := make([]string, 0, len(selectedSkillIds))
	for _, id := range selectedSkillIds {
		if id <= 0 {
			continue
		}
		if _, ok := selected[id]; ok {
			continue
		}
		selected[id] = struct{}{}
		selectedIds = append(selectedIds, cast.ToString(id))
	}

	userRows := make([]msql.Params, 0)
	if len(selectedIds) > 0 {
		userRows, err = msql.Model(define.TableChatAiClawbotUserSkill, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`id`, `in`, strings.Join(selectedIds, `,`)).
			Select()
		if err != nil {
			return ``, err
		}
		if len(userRows) != len(selectedIds) {
			return `no_data`, nil
		}
	}
	existingRows, err := msql.Model(define.TableChatAiClawbotSkill, define.Postgres).
		Where(`robot_key`, robotKey).
		Where(`source_type`, cast.ToString(define.SkillSourceTypeUpload)).
		Where(`user_skill_id`, `!=`, `0`).
		Select()
	if err != nil {
		return ``, err
	}
	existingByUserId := make(map[int64]msql.Params, len(existingRows))
	for _, row := range existingRows {
		existingByUserId[cast.ToInt64(row[`user_skill_id`])] = row
	}

	for _, row := range existingRows {
		userSkillId := cast.ToInt64(row[`user_skill_id`])
		if _, ok := selected[userSkillId]; ok {
			continue
		}
		if errKey, err = deleteRobotSkillRowAndDir(robotKey, row); err != nil || errKey != `` {
			return errKey, err
		}
	}

	for _, userRow := range userRows {
		userSkillId := cast.ToInt64(userRow[`id`])
		if existing, ok := existingByUserId[userSkillId]; ok {
			if errKey, err = syncOneRobotSkill(adminUserId, robotKey, userRow, existing); err != nil || errKey != `` {
				return errKey, err
			}
			continue
		}
		if errKey, err = syncOneRobotSkill(adminUserId, robotKey, userRow, nil); err != nil || errKey != `` {
			return errKey, err
		}
	}
	return ``, nil
}

func syncOneRobotSkill(adminUserId int, robotKey string, userRow msql.Params, oldRow msql.Params) (errKey string, err error) {
	userSkillId := cast.ToInt64(userRow[`id`])
	skillName := userRow[`skill_name`]
	if !define.SkillNameRegexp.MatchString(skillName) || skillName == define.SkillReservedName {
		return `clawbot_skill_name_invalid`, nil
	}
	srcDir := filepath.Join(userSkillsDir(adminUserId), skillName)
	if !tool.IsDir(srcDir) {
		return `clawbot_skill_dir_missing`, nil
	}
	var oldId int64
	oldName := ``
	if len(oldRow) > 0 {
		oldId = cast.ToInt64(oldRow[`id`])
		oldName = oldRow[`skill_name`]
	}
	dup, derr := msql.Model(define.TableChatAiClawbotSkill, define.Postgres).
		Where(`robot_key`, robotKey).
		Where(`source_type`, cast.ToString(define.SkillSourceTypeUpload)).
		Where(`skill_name`, skillName).
		Find()
	if derr != nil {
		return ``, derr
	}
	if len(dup) > 0 && cast.ToInt64(dup[`id`]) != oldId {
		return `clawbot_skill_system_dir_exists`, nil
	}

	baseDir := privateSkillsDir(robotKey)
	oldDir := ``
	if oldName != `` {
		oldDir = filepath.Join(baseDir, oldName)
	}
	destDir := filepath.Join(baseDir, skillName)
	cleanup, rollback, err := replaceSkillDirWithBackup(skillTmpRoot(robotKey), oldDir, destDir, srcDir, skillName, userRow[`description`])
	if err != nil {
		return ``, err
	}
	now := tool.Time2Int()
	data := msql.Datas{
		`admin_user_id`: adminUserId,
		`robot_key`:     robotKey,
		`source_type`:   define.SkillSourceTypeUpload,
		`user_skill_id`: userSkillId,
		`skill_name`:    skillName,
		`remark_name`:   userRow[`remark_name`],
		`description`:   userRow[`description`],
		`file_size`:     cast.ToInt64(userRow[`file_size`]),
		`update_time`:   now,
	}
	if oldId > 0 {
		_, err = msql.Model(define.TableChatAiClawbotSkill, define.Postgres).
			Where(`id`, cast.ToString(oldId)).
			Where(`robot_key`, robotKey).
			Where(`source_type`, cast.ToString(define.SkillSourceTypeUpload)).
			Update(data)
	} else {
		data[`create_time`] = now
		_, err = msql.Model(define.TableChatAiClawbotSkill, define.Postgres).Insert(data, `id`)
	}
	if err != nil {
		rollback()
		return ``, err
	}
	cleanup()
	return ``, nil
}

func deleteRobotSkillRowAndDir(robotKey string, row msql.Params) (errKey string, err error) {
	skillId := cast.ToInt64(row[`id`])
	skillName := row[`skill_name`]
	if !define.SkillNameRegexp.MatchString(skillName) {
		return `clawbot_skill_name_invalid`, nil
	}
	skillDir := filepath.Join(privateSkillsDir(robotKey), skillName)
	deleteDir, moved, merr := moveSkillDirToHiddenTmp(robotKey, skillDir, `.delete_`)
	if merr != nil {
		return ``, merr
	}
	restoreDeleteDir := func() {
		if moved && !tool.IsDir(skillDir) {
			if rerr := os.Rename(deleteDir, skillDir); rerr != nil {
				logs.Error(fmt.Sprintf(`restore deleted skill dir failed, skill_id=%d: %v`, skillId, rerr))
			}
		}
	}
	rows, derr := msql.Model(define.TableChatAiClawbotSkill, define.Postgres).
		Where(`id`, cast.ToString(skillId)).
		Where(`robot_key`, robotKey).
		Where(`source_type`, cast.ToString(define.SkillSourceTypeUpload)).
		Where(`user_skill_id`, `!=`, `0`).
		Delete()
	if derr != nil {
		restoreDeleteDir()
		return ``, derr
	}
	if rows <= 0 {
		restoreDeleteDir()
		return `no_data`, nil
	}
	if moved {
		_ = os.RemoveAll(deleteDir)
	}
	return ``, nil
}

func replaceSkillDirWithBackup(tmpRoot, oldDir, destDir, srcDir, skillName, description string) (cleanup func(), rollback func(), err error) {
	backupDirs := make([]string, 0, 2)
	restorePairs := make([][2]string, 0, 2)
	moveForReplace := func(dir string) error {
		if dir == `` || !tool.IsDir(dir) {
			return nil
		}
		backupDir, moved, merr := moveSkillDirToHiddenTmpRoot(tmpRoot, dir, `.bak_`)
		if merr != nil {
			return merr
		}
		if moved {
			backupDirs = append(backupDirs, backupDir)
			restorePairs = append(restorePairs, [2]string{backupDir, dir})
		}
		return nil
	}
	if err = moveForReplace(oldDir); err != nil {
		return nil, nil, err
	}
	if destDir != oldDir {
		if err = moveForReplace(destDir); err != nil {
			for i := len(restorePairs) - 1; i >= 0; i-- {
				backupDir := restorePairs[i][0]
				targetDir := restorePairs[i][1]
				if !tool.IsDir(targetDir) {
					if rerr := os.Rename(backupDir, targetDir); rerr != nil {
						logs.Error(fmt.Sprintf(`restore skill backup failed, target=%s: %v`, targetDir, rerr))
					}
				}
			}
			return nil, nil, err
		}
	}
	cleanup = func() {
		for _, dir := range backupDirs {
			_ = os.RemoveAll(dir)
		}
	}
	rollback = func() {
		_ = os.RemoveAll(destDir)
		for i := len(restorePairs) - 1; i >= 0; i-- {
			backupDir := restorePairs[i][0]
			targetDir := restorePairs[i][1]
			if !tool.IsDir(targetDir) {
				if rerr := os.Rename(backupDir, targetDir); rerr != nil {
					logs.Error(fmt.Sprintf(`restore skill backup failed, target=%s: %v`, targetDir, rerr))
				}
			}
		}
	}
	if err = publishSkillDirWithTmpRoot(tmpRoot, srcDir, destDir, skillName, description); err != nil {
		rollback()
		return nil, nil, err
	}
	return cleanup, rollback, nil
}

func renameUserSkillDirWithBackup(tmpRoot, oldDir, destDir, skillName, description string) (cleanup func(), rollback func(), err error) {
	if !tool.IsDir(oldDir) {
		return nil, nil, errors.New(`staged skill dir not found`)
	}
	if oldDir == destDir {
		return func() {}, func() {}, nil
	}
	if tool.IsDir(destDir) {
		return nil, nil, errors.New(`target skill dir exists`)
	}
	backupDir, moved, err := moveSkillDirToHiddenTmpRoot(tmpRoot, oldDir, `.bak_`)
	if err != nil {
		return nil, nil, err
	}
	rollback = func() {
		_ = os.RemoveAll(destDir)
		if moved && !tool.IsDir(oldDir) {
			if rerr := os.Rename(backupDir, oldDir); rerr != nil {
				logs.Error(fmt.Sprintf(`restore user skill backup failed, target=%s: %v`, oldDir, rerr))
			}
		}
	}
	if err = publishSkillDirWithTmpRoot(tmpRoot, backupDir, destDir, skillName, description); err != nil {
		rollback()
		return nil, nil, err
	}
	cleanup = func() {
		if moved {
			_ = os.RemoveAll(backupDir)
		}
	}
	return cleanup, rollback, nil
}

func readSkillMdName(skillMdPath string) (name, description, body string, err error) {
	content, rerr := tool.ReadFile(skillMdPath)
	if rerr != nil {
		return ``, ``, ``, rerr
	}
	return parseSkillFrontmatter(content)
}

func moveSkillDirToHiddenTmp(robotKey, srcDir, prefix string) (hiddenDir string, moved bool, err error) {
	return moveSkillDirToHiddenTmpRoot(skillTmpRoot(robotKey), srcDir, prefix)
}

func moveSkillDirToHiddenTmpRoot(root, srcDir, prefix string) (hiddenDir string, moved bool, err error) {
	if !tool.IsDir(srcDir) {
		return ``, false, nil
	}
	if err = tool.MkDirAll(root); err != nil {
		return ``, false, err
	}
	for i := 0; i < 5; i++ {
		candidate := filepath.Join(root, prefix+tool.Random(8))
		if _, serr := os.Lstat(candidate); serr == nil {
			continue
		} else if !os.IsNotExist(serr) {
			return ``, false, serr
		}
		if err = os.Rename(srcDir, candidate); err != nil {
			return ``, false, err
		}
		return candidate, true, nil
	}
	return ``, false, errors.New(`failed to allocate hidden skill tmp dir`)
}

func copyDir(src, dest string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, rerr := filepath.Rel(src, path)
		if rerr != nil {
			return rerr
		}
		target := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		return copyFile(path, target)
	})
}
