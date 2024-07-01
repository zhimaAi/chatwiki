// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"chatwiki/internal/pkg/lib_define"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

func ClientsidebuildWin(msg string, _ ...string) error {
	logs.Debug(`nsq:%s`, msg)
	data := make(map[string]any)
	if err := tool.JsonDecode(msg, &data); err != nil {
		logs.Error(`parsing failure:%s/%s`, msg, err.Error())
		return nil
	}
	version := lib_define.GetElectronVersion()
	if version != cast.ToString(data[`version`]) {
		logs.Error(`version inconformity:%s/%s`, msg, version)
		return nil
	}
	domain := cast.ToString(data[`domain`])
	adminUserId := cast.ToInt(data[`admin_user_id`])
	exeUrl := cast.ToString(data[`exe_url`])
	zipUrl := cast.ToString(data[`zip_url`])
	if len(domain) == 0 || adminUserId <= 0 || len(exeUrl) == 0 || len(zipUrl) == 0 {
		logs.Error(`params error:%s/%d/%s/%s`, domain, adminUserId, exeUrl, zipUrl)
		return nil
	}
	if tool.IsFile(`static`+exeUrl) && tool.IsFile(`static`+zipUrl) {
		logs.Debug(`already packed:%s/%s`, exeUrl, zipUrl)
		return nil
	}
	envFile := lib_define.ElectronPath + `.env.production`
	content := fmt.Sprintf("VITE_BASE_API_URL = '%s'\r\nVITE_ADMIN_USER_ID = %d", domain, adminUserId)
	if err := tool.WriteFile(envFile, content); err != nil {
		logs.Error(`write envFile failure:%s`, err.Error())
		return nil
	}
	cmd := exec.Command(`sh`, `build_shell.sh`)
	cmd.Dir = lib_define.ElectronPath
	if err := cmd.Run(); err != nil {
		logs.Error(`build shell failure:%s`, err.Error())
		return nil
	}
	if !tool.IsDir(tool.DirName(`static` + exeUrl)) {
		if err := tool.MkDirAll(tool.DirName(`static` + exeUrl)); err != nil {
			logs.Error(`make dir failure:%s`, err.Error())
			return nil
		}
	}
	exeFile := fmt.Sprintf(`%sdist/ChatWiki-%s-setup.exe`, lib_define.ElectronPath, version)
	if err := os.Rename(exeFile, `static`+exeUrl); err != nil {
		logs.Error(`move file failure:%s`, err.Error())
		return nil
	}
	if !tool.IsDir(tool.DirName(`static` + zipUrl)) {
		if err := tool.MkDirAll(tool.DirName(`static` + zipUrl)); err != nil {
			logs.Error(`make dir failure:%s`, err.Error())
			return nil
		}
	}
	if err := os.Rename(lib_define.ElectronPath+`dist/chatwiki.zip`, `static`+zipUrl); err != nil {
		logs.Error(`move file failure:%s`, err.Error())
		return nil
	}
	logs.Debug(`build finish:%s`, msg)
	return nil
}
