// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package volcenginev2

import (
	OfficialApi "github.com/volcengine/volc-sdk-golang/service/maas/models/api/v2"
	OfficialClient "github.com/volcengine/volc-sdk-golang/service/maas/v2"
)

type Client struct {
	Host       string
	EndPointID string
	AK         string
	SK         string
	Region     string
}

func NewClient(Host, EndPointID, AK, SK, Region string) *Client {
	return &Client{
		Host:       Host,
		EndPointID: EndPointID,
		AK:         AK,
		SK:         SK,
		Region:     Region,
	}
}

func (c *Client) CreateEmbeddings(req EmbeddingRequest) (OfficialApi.EmbeddingsResp, error) {
	r := OfficialClient.NewInstance(c.Host, c.Region)
	r.SetAccessKey(c.AK)
	r.SetSecretKey(c.SK)

	req2 := &OfficialApi.EmbeddingsReq{Input: req.Input}
	response, _, err := r.Embeddings(c.EndPointID, req2)
	if err != nil {
		return OfficialApi.EmbeddingsResp{}, err
	}

	return *response, err
}
