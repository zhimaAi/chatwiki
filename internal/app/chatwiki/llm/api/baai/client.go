// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package baai

type Client struct {
	EndPoint string
	Model    string
	APIKey   string
}

func NewClient(endPoint, Model, APIKey string) *Client {
	return &Client{
		EndPoint: endPoint,
		Model:    Model,
		APIKey:   APIKey,
	}
}
