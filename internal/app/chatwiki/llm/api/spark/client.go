// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package spark

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	APIKey    string
	APPID     string
	APISecret string
	Model     string
}
type endpoint struct {
	URL    string
	Domain string
}

var modelToEndpoint = map[string]endpoint{
	"Spark4.0 Ultra": {URL: "wss://spark-api.xf-yun.com/v4.0/chat", Domain: "4.0Ultra"},
	"Spark Max":      {URL: "wss://spark-api.xf-yun.com/v3.5/chat", Domain: "generalv3.5"},
	"Spark Pro":      {URL: "wss://spark-api.xf-yun.com/v3.1/chat", Domain: "generalv3"},
	"Spark V2.0":     {URL: "wss://spark-api.xf-yun.com/v2.1/chat", Domain: "generalv2"},
	"Spark Lite":     {URL: "wss://spark-api.xf-yun.com/v1.1/chat", Domain: "general"},
}

func NewClient(apiKey, appID, appSecret, model string) *Client {
	return &Client{
		APIKey:    apiKey,
		APPID:     appID,
		APISecret: appSecret,
		Model:     model,
	}
}

func (c *Client) CreateChatCompletion(req ChatCompletionRequest) (ChatCompletionResponse, error) {
	// connect websocket
	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	endpoint, ok := modelToEndpoint[c.Model]
	if !ok {
		return ChatCompletionResponse{}, errors.New("unsupported model")
	}

	conn, resp, err := d.Dial(assembleAuthUrl1(endpoint.URL, c.APIKey, c.APISecret), nil)
	if err != nil {
		return ChatCompletionResponse{}, errors.New(readResp(resp) + err.Error())
	} else if resp.StatusCode != 101 {
		return ChatCompletionResponse{}, errors.New(readResp(resp))
	}

	// send message
	req.Header.APPID = c.APPID
	req.Parameter.Chat.Domain = endpoint.Domain
	req.Parameter.Chat.Auditing = "default"
	req.Parameter.Chat.Stream = false
	req.Parameter.Chat.TopK = 5
	err = conn.WriteJSON(req)
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	_, msg, err := conn.ReadMessage()
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)
	if err != nil {
		return ChatCompletionResponse{}, err
	}

	var response ChatCompletionResponse
	err = json.Unmarshal(msg, &response)
	if err != nil {
		return ChatCompletionResponse{}, err
	}
	if len(response.Payload.Choices.Text) <= 0 {
		return ChatCompletionResponse{}, errors.New("spark response no response text")
	}
	return response, nil
}

func (c *Client) CreateChatCompletionStream(req ChatCompletionRequest) (*ChatCompletionStream, error) {
	// connect websocket
	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	endpoint, ok := modelToEndpoint[c.Model]
	if !ok {
		return nil, errors.New("unsupported model")
	}
	conn, resp, err := d.Dial(assembleAuthUrl1(endpoint.URL, c.APIKey, c.APISecret), nil)
	if err != nil {
		return nil, errors.New(readResp(resp) + err.Error())
	} else if resp.StatusCode != 101 {
		return nil, errors.New(readResp(resp))
	}

	// send message
	req.Header.APPID = c.APPID
	req.Parameter.Chat.Domain = endpoint.Domain
	req.Parameter.Chat.Auditing = "default"
	req.Parameter.Chat.Stream = true
	req.Parameter.Chat.TopK = 5
	conn.WriteJSON(req)

	return &ChatCompletionStream{conn: conn}, nil
}

// create authorization url  apikey is hmac username
func assembleAuthUrl1(hosturl string, apiKey, apiSecret string) string {
	ul, err := url.Parse(hosturl)
	if err != nil {
		fmt.Println(err)
	}
	date := time.Now().UTC().Format(time.RFC1123)
	//date = "Tue, 28 May 2019 09:10:42 MST"
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	sgin := strings.Join(signString, "\n")
	// fmt.Println(sgin)
	sha := HmacWithShaTobase64("hmac-sha256", sgin, apiSecret)
	// fmt.Println(sha)
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", sha)
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	callurl := hosturl + "?" + v.Encode()
	return callurl
}

func HmacWithShaTobase64(_, data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func readResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}
