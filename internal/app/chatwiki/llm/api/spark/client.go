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
	EndPoint  string
	APIKey    string
	APPID     string
	APISecret string
}

func NewClient(EndPoint, apiKey, appID, appSecret string) *Client {
	return &Client{
		EndPoint:  EndPoint,
		APIKey:    apiKey,
		APPID:     appID,
		APISecret: appSecret,
	}
}

func (c *Client) CreateChatCompletion(req ChatCompletionRequest) (ChatCompletionResponse, error) {
	// connect websocket
	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	conn, resp, err := d.Dial(assembleAuthUrl1(c.EndPoint, c.APIKey, c.APISecret), nil)
	if err != nil {
		return ChatCompletionResponse{}, errors.New(readResp(resp) + err.Error())
	} else if resp.StatusCode != 101 {
		return ChatCompletionResponse{}, errors.New(readResp(resp))
	}

	// send message
	req.Header.APPID = c.APPID
	req.Parameter.Chat.Domain = "general"
	req.Parameter.Chat.Auditing = "default"
	req.Parameter.Chat.Stream = false
	req.Parameter.Chat.TopK = 5
	conn.WriteJSON(req)

	// read message
	_, msg, err := conn.ReadMessage()
	defer conn.Close()
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
	conn, resp, err := d.Dial(assembleAuthUrl1(c.EndPoint, c.APIKey, c.APISecret), nil)
	if err != nil {
		return nil, errors.New(readResp(resp) + err.Error())
	} else if resp.StatusCode != 101 {
		return nil, errors.New(readResp(resp))
	}

	// send message
	req.Header.APPID = c.APPID
	req.Parameter.Chat.Domain = "general"
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
