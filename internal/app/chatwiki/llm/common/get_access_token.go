// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Token struct {
	APIKey      string
	SecretKey   string
	accessToken string
	expireTime  time.Time
}

var (
	tokenManagerInstance *TokenManager
	once                 sync.Once
)

type TokenManager struct {
	mutex  sync.Mutex
	tokens map[string]*Token
}

func GetTokenManagerInstance() *TokenManager {
	once.Do(func() {
		tokenManagerInstance = &TokenManager{
			tokens: make(map[string]*Token),
		}
	})
	return tokenManagerInstance
}

func (manager *TokenManager) GetBaiduAccessToken(EndPoint, APIKey, SecretKey string) (string, error) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if token, exists := manager.tokens[APIKey]; exists && time.Now().Before(token.expireTime) {
		return token.accessToken, nil
	}

	// Create new token if not exists or expired
	token := &Token{
		APIKey:     APIKey,
		SecretKey:  SecretKey,
		expireTime: time.Now(), // Initialize with current time to force refresh
	}
	accessToken, err := token.refreshBaiduToken(EndPoint)
	if err != nil {
		return "", err
	}
	manager.tokens[APIKey] = token
	return accessToken, nil
}

func (t *Token) refreshBaiduToken(EndPoint string) (string, error) {
	link := EndPoint + "/oauth/2.0/token"
	params := []Param{
		{Key: "grant_type", Value: "client_credentials"},
		{Key: "client_id", Value: t.APIKey},
		{Key: "client_secret", Value: t.SecretKey},
	}
	resp, err := HttpPost(link, nil, params, nil)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var respData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", err
	}

	if errorDesc, exists := respData["error_description"]; exists {
		return "", fmt.Errorf("baidu OAuth API error: %s", errorDesc)
	}

	if accessToken, exists := respData["access_token"].(string); exists {
		expiresIn, ok := respData["expires_in"].(float64)
		if !ok {
			return "", fmt.Errorf("invalid expiration time format")
		}
		t.accessToken = accessToken
		refreshBuffer := time.Hour
		t.expireTime = time.Now().Add(time.Duration(expiresIn) * time.Second).Add(-refreshBuffer)
		return accessToken, nil
	}

	return "", fmt.Errorf("access token not found in the response")
}

func (manager *TokenManager) GetVolcengineAccessToken(EndPoint, Region, Model, AccessKeyID, SecretAccessKey string) (string, error) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if token, exists := manager.tokens[AccessKeyID]; exists && time.Now().Before(token.expireTime) {
		return token.accessToken, nil
	}

	// Create new token if not exists or expired
	token := &Token{
		APIKey:     AccessKeyID,
		SecretKey:  SecretAccessKey,
		expireTime: time.Now(), // Initialize with current time to force refresh
	}
	accessToken, err := token.refreshVolcengineToken(EndPoint, Region, Model)
	if err != nil {
		return "", err
	}
	manager.tokens[AccessKeyID] = token
	return accessToken, nil
}

func (t *Token) refreshVolcengineToken(EndPoint, Region, Model string) (string, error) {
	const (
		Service         = "ark"
		Action          = "GetApiKey"
		Version         = "2024-01-01"
		Path            = "/"
		DurationSeconds = 2592000
	)

	body := []byte(fmt.Sprintf("DurationSeconds=%d&ResourceType=endpoint&ResourceIds=%s", DurationSeconds, Model))

	queries := make(url.Values)
	queries.Set("Action", Action)
	queries.Set("Version", Version)
	requestAddr := fmt.Sprintf("%s%s?%s", EndPoint, Path, queries.Encode())
	request, err := http.NewRequest(http.MethodPost, requestAddr, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	// 2. build signature material
	now := time.Now()
	date := now.UTC().Format("20060102T150405Z")
	authDate := date[:8]
	request.Header.Set("X-Date", date)

	enc, err := hashSHA256(body)
	if err != nil {
		return "", err
	}
	payload := hex.EncodeToString(enc)
	request.Header.Set("X-Content-Sha256", payload)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	queryString := strings.Replace(queries.Encode(), "+", "%20", -1)
	signedHeaders := []string{"host", "x-date", "x-content-sha256", "content-type"}
	var headerList []string
	for _, header := range signedHeaders {
		if header == "host" {
			headerList = append(headerList, header+":"+request.Host)
		} else {
			v := request.Header.Get(header)
			headerList = append(headerList, header+":"+strings.TrimSpace(v))
		}
	}
	headerString := strings.Join(headerList, "\n")

	canonicalString := strings.Join([]string{
		http.MethodPost,
		Path,
		queryString,
		headerString + "\n",
		strings.Join(signedHeaders, ";"),
		payload,
	}, "\n")

	enc, err = hashSHA256([]byte(canonicalString))
	if err != nil {
		return "", nil
	}
	hashedCanonicalString := hex.EncodeToString(enc)

	credentialScope := authDate + "/" + Region + "/" + Service + "/request"
	signString := strings.Join([]string{
		"HMAC-SHA256",
		date,
		credentialScope,
		hashedCanonicalString,
	}, "\n")

	// 3. build the authentication request header
	signedKey := getSignedKey(t.SecretKey, authDate, Region, Service)
	signature := hex.EncodeToString(hmacSHA256(signedKey, signString))

	authorization := "HMAC-SHA256" +
		" Credential=" + t.APIKey + "/" + credentialScope +
		", SignedHeaders=" + strings.Join(signedHeaders, ";") +
		", Signature=" + signature
	request.Header.Set("Authorization", authorization)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}

	respBody, err := io.ReadAll(resp.Body)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return "", err
	}

	// check error
	if metadata, ok := data["ResponseMetadata"].(map[string]interface{}); ok {
		if errorInfo, ok := metadata["Error"].(map[string]interface{}); ok {
			return "", errors.New(errorInfo["Message"].(string))
		}
	}

	// fetch ApiKey
	if result, ok := data["Result"].(map[string]interface{}); ok {
		if apiKey, ok := result["ApiKey"].(string); ok {
			return apiKey, nil
		}
	}

	return "", errors.New("no volcengine apikey or error found in the response")
}

func hmacSHA256(key []byte, content string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(content))
	return mac.Sum(nil)
}

func getSignedKey(secretKey, date, region, service string) []byte {
	kDate := hmacSHA256([]byte(secretKey), date)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	kSigning := hmacSHA256(kService, "request")

	return kSigning
}

func hashSHA256(data []byte) ([]byte, error) {
	hash := sha256.New()
	if _, err := hash.Write(data); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}
