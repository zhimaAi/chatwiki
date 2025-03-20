// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package business

import (
	"bytes"
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/app/user_domain_service/define"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type SaveCertReq struct {
	Url               string `form:"url" json:"url" binding:"required"`
	SslCertificate    string `form:"ssl_certificate" json:"ssl_certificate" binding:"required"`
	SslCertificateKey string `form:"ssl_certificate_key" json:"ssl_certificate_key" binding:"required"`
	Upstream          string `form:"upstream" json:"upstream"`
}

func SaveCert(c *gin.Context) {
	var (
		req = SaveCertReq{}
		err error
	)
	if err = c.ShouldBind(&req); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_err`, err.Error()))))
		return
	}
	if req.Upstream == "" {
		req.Upstream = define.Config.ChatWiki[`host`] + `:` + define.Config.ChatWiki[`port`]
	}
	if err = writeConf(req.Url, req.Upstream); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `写入配置文件失败`))))
		return
	}

	if err = writeCert(req.Url, req.SslCertificate, req.SslCertificateKey); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `写入配置文件失败`))))
		return
	}
	if output, code, err := reloadNginx(); err != nil {
		logs.Error("output:%s,code:%d,err:%s", output, code, err.Error())
		_ = deleteFile(req.Url)
		reloadNginx()
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `Nginx配置文件启动失败,请检查后再试`+err.Error()))))
		return
	}

	common.FmtOk(c, nil)
}

func writeConf(domain, upstream string) error {
	path := "/etc/nginx/conf.d"
	// check file exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// tpl
	conf := `
server {
	listen 443 ssl;
	server_name %s;

	ssl_certificate cert/%s.crt;
	ssl_certificate_key cert/%s.key;
	ssl_session_timeout 5m;
	ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE:ECDH:AES:HIGH:!NULL:!aNULL:!MD5:!ADH:!RC4;
	ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
	ssl_prefer_server_ciphers on;

	#error_log /var/log/nginx/%s.error.log;
	#access_log /var/log/nginx/%s.access.log;

	location / {
		proxy_pass https://%s; 
	  %s  proxy_set_header Host $http_host;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	}

}
`
	// replace
	note := "#"
	conf = fmt.Sprintf(conf, upstream, domain, domain, domain, domain, domain, note)

	// conf path
	confFilePath := path + `/` + fmt.Sprintf("%s.conf", domain)
	confFilePathBackup := path + `/` + fmt.Sprintf("%s.conf.bak%d", domain, time.Now().Unix())

	// backup
	if _, err := os.Stat(confFilePath); err == nil {
		if err := tool.WriteFile(confFilePathBackup, conf); err != nil {
			return fmt.Errorf("failed to create backup file: %w", err)
		}
	}

	// write conf
	if err := tool.WriteFile(confFilePath, conf); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	return nil
}

// writeCert write SSL
func writeCert(domain, crt, key string) error {
	path := "/etc/nginx/cert/"
	// check file exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	crtFilePath := path + `/` + fmt.Sprintf("%s.crt", domain)
	crtFilePathBackup := path + `/` + fmt.Sprintf("%s.crt.bak%d", domain, time.Now().Unix())

	// backup
	if _, err := os.Stat(crtFilePath); err == nil {
		if err := tool.WriteFile(crtFilePathBackup, crt); err != nil {
			return fmt.Errorf("failed to create backup file for certificate: %w", err)
		}
	}

	// write cert
	if err := tool.WriteFile(crtFilePath, crt); err != nil {
		return fmt.Errorf("failed to write certificate file: %w", err)
	}

	// key path
	keyFilePath := path + `/` + fmt.Sprintf("%s.key", domain)
	keyFilePathBackup := path + `/` + fmt.Sprintf("%s.key.bak%d", domain, time.Now().Unix())

	// backup
	if _, err := os.Stat(keyFilePath); err == nil {
		if err := tool.WriteFile(keyFilePathBackup, key); err != nil {
			return fmt.Errorf("failed to create backup file for key: %w", err)
		}
	}

	// write key
	if err := tool.WriteFile(keyFilePath, key); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}

	return nil
}

func deleteFile(domain string) error {
	path := "/etc/nginx/conf.d/"
	pathCert := "/etc/nginx/cert/"
	confFilePath := filepath.Join(path, fmt.Sprintf("%s.conf", domain))
	if _, err := os.Stat(confFilePath); err == nil {
		if err := os.Remove(confFilePath); err != nil {
			return fmt.Errorf("failed to delete %s: %w", confFilePath, err)
		}
	}

	crtFilePath := filepath.Join(pathCert, fmt.Sprintf("%s.crt", domain))
	if _, err := os.Stat(crtFilePath); err == nil {
		if err := os.Remove(crtFilePath); err != nil {
			return fmt.Errorf("failed to delete %s: %w", crtFilePath, err)
		}
	}

	keyFilePath := filepath.Join(pathCert, fmt.Sprintf("%s.key", domain))
	if _, err := os.Stat(keyFilePath); err == nil {
		if err := os.Remove(keyFilePath); err != nil {
			return fmt.Errorf("failed to delete %s: %w", keyFilePath, err)
		}
	}

	return nil
}

func reloadNginx() (string, int, error) {
	// exec ...
	cmd := exec.Command("nginx", "-s", "reload")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return out.String(), -1, fmt.Errorf("failed to reload nginx: %w", err)
	}

	code := cmd.ProcessState.ExitCode()

	return out.String(), code, nil
}
