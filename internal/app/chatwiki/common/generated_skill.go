// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/zhimaAi/go_tools/tool"
)

func generatedSkillFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return ``, err
	}
	defer func() { _ = file.Close() }()
	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		return ``, err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func saveGeneratedSkillZipFile(adminUserId int, source string, zipPath string, allowExt []string) (*define.UploadInfo, error) {
	info, err := os.Stat(zipPath)
	if err != nil {
		return nil, err
	}
	if info.IsDir() || info.Size() <= 0 {
		return nil, errors.New(`invalid zip file size`)
	}
	ext := strings.ToLower(strings.TrimLeft(filepath.Ext(zipPath), `.`))
	if !tool.InArrayString(ext, allowExt) {
		return nil, errors.New(ext + ` not allow`)
	}
	md5Hash, err := generatedSkillFileMD5(zipPath)
	if err != nil {
		return nil, err
	}
	objectKey := fmt.Sprintf(`chat_ai/%v/%s/%s/%s.%s`, adminUserId, source, tool.Date(`Ym`), md5Hash, ext)
	link, err := WriteFileByFile(objectKey, zipPath)
	if err != nil {
		return nil, err
	}
	return &define.UploadInfo{Name: filepath.Base(zipPath), Size: info.Size(), Ext: ext, Link: link}, nil
}
