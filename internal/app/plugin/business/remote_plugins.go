// Copyright © 2016- 2025 Sesame Network Technology all right reserved

package business

import (
	"archive/zip"
	"chatwiki/internal/app/plugin/define"
	"chatwiki/internal/app/plugin/php"
	"chatwiki/internal/pkg/lib_web"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func GetRemotePluginList(c *gin.Context) {
	resp, err := requestXiaokefu(`kf/ChatWiki/CommonGetPluginList`, nil)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	respData, ok := resp.Data.([]interface{})
	if !ok {
		err = errors.New(`invalid data format`)
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(respData, nil))
}

func GetRemotePluginDetail(c *gin.Context) {
	resp, err := requestXiaokefu(`kf/ChatWiki/CommonGetPluginDetail`, nil)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	respData, ok := resp.Data.([]interface{})
	if !ok {
		err = errors.New(`invalid data format`)
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(respData, nil))
}

func DownloadRemotePlugin(c *gin.Context) {
	adminUserId := c.GetHeader(`admin_user_id`)
	if adminUserId == "" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(`缺少admin_user_id`)))
		return
	}
	var req struct {
		URL       string `form:"url" binding:"required"`
		VersionId int    `form:"version_id" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "plugin_download_*")
	if err != nil {
		logs.Error(fmt.Sprintf("创建临时目录失败: %v", err))
		c.String(http.StatusOK, lib_web.FmtJson(nil, fmt.Errorf("创建临时目录失败: %w", err)))
		return
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			logs.Error(err.Error())
		}
	}(tmpDir) // 清理临时目录

	// 下载压缩包
	zipPath := filepath.Join(tmpDir, "plugin.zip")
	if err := downloadFile(req.URL, zipPath); err != nil {
		logs.Error(fmt.Sprintf("下载插件失败: %v", err))
		c.String(http.StatusOK, lib_web.FmtJson(nil, fmt.Errorf("下载插件失败: %w", err)))
		return
	}

	// 解压到临时目录
	extractDir := filepath.Join(tmpDir, "extracted")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		logs.Error(fmt.Sprintf("创建解压目录失败: %v", err))
		c.String(http.StatusOK, lib_web.FmtJson(nil, fmt.Errorf("创建解压目录失败: %w", err)))
		return
	}

	if err := unzipFile(zipPath, extractDir); err != nil {
		logs.Error(fmt.Sprintf("解压插件失败: %v", err))
		c.String(http.StatusOK, lib_web.FmtJson(nil, fmt.Errorf("解压插件失败: %w", err)))
		return
	}

	// 读取 manifest.json
	manifestPath := filepath.Join(extractDir, "manifest.json")
	manifest, err := readManifest(manifestPath)
	if err != nil {
		logs.Error(fmt.Sprintf("读取 manifest.json 失败: %v", err))
		c.String(http.StatusOK, lib_web.FmtJson(nil, fmt.Errorf("读取 manifest.json 失败: %w", err)))
		return
	}

	// 获取工作目录
	workDir, err := os.Getwd()
	if err != nil {
		logs.Error(fmt.Sprintf("获取工作目录失败: %v", err))
		c.String(http.StatusOK, lib_web.FmtJson(nil, fmt.Errorf("获取工作目录失败: %w", err)))
		return
	}

	// 目标目录
	targetDir := filepath.Join(workDir, "php", "plugins", manifest.Name)

	// 如果目标目录已存在，先删除
	if tool.IsDir(targetDir) {
		if err := os.RemoveAll(targetDir); err != nil {
			logs.Error(fmt.Sprintf("删除已存在的插件目录失败: %v", err))
			c.String(http.StatusOK, lib_web.FmtJson(nil, fmt.Errorf("删除已存在的插件目录失败: %w", err)))
			return
		}
	}

	// 拷贝解压后的目录到目标位置
	if err := copyDirectory(extractDir, targetDir); err != nil {
		logs.Error(fmt.Sprintf("拷贝插件目录失败: %v", err))
		c.String(http.StatusOK, lib_web.FmtJson(nil, fmt.Errorf("拷贝插件目录失败: %w", err)))
		return
	}

	// 增加安装次数
	if err := increaseInstallCount(req.VersionId); err != nil {
		logs.Error(fmt.Sprintf("增加安装次数失败: %v", err))
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	// 增加数据库记录
	info, err := msql.Model(`plugin_config`, define.Postgres).Where(`admin_user_id`, adminUserId).Where(`name`, manifest.Name).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	if len(info) == 0 {
		_, err = msql.Model(`plugin_config`, define.Postgres).Insert(msql.Datas{
			`create_time`:   tool.Time2Int(),
			`update_time`:   tool.Time2Int(),
			`admin_user_id`: adminUserId,
			`name`:          manifest.Name,
			`type`:          manifest.Type,
			`has_loaded`:    false,
		})
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

// downloadFile 下载文件
func downloadFile(url, filePath string) error {
	request := curl.Get(url)
	response, err := request.Response()
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，状态码: %d", response.StatusCode)
	}
	err = request.ToFile(filePath)
	if err != nil && tool.IsFile(filePath) {
		_ = os.Remove(filePath) // 删除错误文件
	}
	return err
}

// unzipFile 解压 zip 文件
func unzipFile(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		cleanName := filepath.Clean(f.Name)
		if cleanName == "__MACOSX" || strings.HasPrefix(cleanName, "__MACOSX"+string(os.PathSeparator)) {
			continue
		}

		// 防止路径遍历攻击
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("非法的文件路径: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			// 创建目录
			if err := os.MkdirAll(fpath, f.FileInfo().Mode()); err != nil {
				return err
			}
			continue
		}

		// 创建父目录
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}

		// 打开 zip 文件中的文件
		rc, err := f.Open()
		if err != nil {
			return err
		}

		// 创建目标文件
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode())
		if err != nil {
			rc.Close()
			return err
		}

		// 复制内容
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// readManifest 读取 manifest.json
func readManifest(manifestPath string) (*php.PluginManifest, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("读取 manifest.json 失败: %w", err)
	}

	var manifest php.PluginManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("解析 manifest.json 失败: %w", err)
	}

	if manifest.Name == "" {
		return nil, errors.New("manifest.json 中缺少 name 字段")
	}

	return &manifest, nil
}

// copyDirectory 拷贝目录
func copyDirectory(src, dest string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dest, relPath)

		info, err := d.Info()
		if err != nil {
			return err
		}

		if d.IsDir() {
			// 创建目录，保证至少有 0755 权限
			mode := info.Mode().Perm()
			if mode == 0 {
				mode = 0755
			}
			return os.MkdirAll(destPath, mode)
		}

		// 创建父目录
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		// 拷贝文件
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}

		fileMode := info.Mode().Perm()
		if fileMode == 0 {
			fileMode = 0644
		}

		destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileMode)
		if err != nil {
			_ = srcFile.Close()
			return err
		}

		_, err = io.Copy(destFile, srcFile)
		closeErr := srcFile.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
		closeErr = destFile.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
		return err
	})
}

// increaseInstallCount 增加安装次数
func increaseInstallCount(versionId int) error {
	_, err := requestXiaokefu(`kf/ChatWiki/CommonAddPluginInstalledCount`, map[string]any{
		"version_id": versionId,
	})
	if err != nil {
		return err
	}
	return nil
}

func requestXiaokefu(api string, data map[string]any) (lib_web.Response, error) {
	domain := define.Config.Xiaokefu[`domain`]
	body, err := tool.JsonEncode(data)
	if err != nil {
		return lib_web.Response{}, err
	}
	if len(body) == 0 {
		body = `{}`
	}
	var (
		link    string
		request *curl.Request
	)
	link = fmt.Sprintf("%s/%s", domain, api)
	request = curl.Post(link)
	for key, item := range data {
		request.Param(key, cast.ToString(item))
	}

	resp, err := request.Response()
	if err != nil {
		return lib_web.Response{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return lib_web.Response{}, fmt.Errorf(`SYSTEM ERROR:%d`, resp.StatusCode)
	}
	code := lib_web.Response{}
	if err = request.ToJSON(&code); err != nil {
		return lib_web.Response{}, err
	}
	if code.Res != lib_web.CommonSuccess {
		return code, errors.New(code.Msg)
	}
	return code, nil
}
