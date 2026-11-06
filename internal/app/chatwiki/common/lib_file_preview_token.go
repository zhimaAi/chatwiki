// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"chatwiki/internal/app/chatwiki/define"
)

const LibFilePagePreviewTTL = 10 * time.Minute

// LibFilePagePreview identifies one PDF page already disclosed by an answer.
// AdminAccess is set only after checking a management JWT, never from request data.
type LibFilePagePreview struct {
	AdminUserID int    `json:"admin_user_id"`
	RobotID     int    `json:"robot_id"`
	OpenID      string `json:"openid"`
	MessageID   int    `json:"message_id"`
	FileID      int    `json:"file_id"`
	ParagraphID int    `json:"paragraph_id"`
	Page        int    `json:"page"`
	AdminAccess bool   `json:"admin_access,omitempty"`
	ExpiresAt   int64  `json:"exp"`
}

var errInvalidLibFilePreview = errors.New("invalid library page preview token")

func (p LibFilePagePreview) valid() bool {
	return p.AdminUserID > 0 && p.RobotID > 0 && IsXkfOpenid(p.OpenID) &&
		p.MessageID > 0 && p.FileID > 0 && p.ParagraphID > 0 && p.Page > 0
}

func libFilePreviewMAC(payload string) []byte {
	mac := hmac.New(sha256.New, []byte(define.JwtKey))
	// A distinct purpose and non-JWT format prevent use as a login credential.
	_, _ = mac.Write([]byte("chatwiki/lib-file-page-preview/v1\x00" + payload))
	return mac.Sum(nil)
}

func SignLibFilePagePreview(preview LibFilePagePreview, now time.Time) (string, error) {
	if define.JwtKey == "" || !preview.valid() {
		return "", errInvalidLibFilePreview
	}
	preview.ExpiresAt = now.Add(LibFilePagePreviewTTL).Unix()
	data, err := json.Marshal(preview)
	if err != nil {
		return "", err
	}
	payload := "lp1." + base64.RawURLEncoding.EncodeToString(data)
	return payload + "." + base64.RawURLEncoding.EncodeToString(libFilePreviewMAC(payload)), nil
}

func ParseLibFilePagePreview(token string, now time.Time) (LibFilePagePreview, error) {
	var preview LibFilePagePreview
	if define.JwtKey == "" || len(token) > 4096 {
		return preview, errInvalidLibFilePreview
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 || parts[0] != "lp1" {
		return preview, errInvalidLibFilePreview
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || !hmac.Equal(signature, libFilePreviewMAC(parts[0]+"."+parts[1])) {
		return preview, errInvalidLibFilePreview
	}
	data, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil || json.Unmarshal(data, &preview) != nil || !preview.valid() ||
		preview.ExpiresAt <= now.Unix() || preview.ExpiresAt > now.Add(LibFilePagePreviewTTL).Unix() {
		return LibFilePagePreview{}, errInvalidLibFilePreview
	}
	return preview, nil
}
