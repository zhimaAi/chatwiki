// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package adaptor

type Meta struct {
	Corp       string
	EndPoint   string
	Model      string
	APIVersion string
	APIKey     string
	SecretKey  string

	Region string //for volcengine and tencent

	APPID string //for spark
}

type Adaptor struct {
	meta Meta
}

func (a *Adaptor) Init(meta Meta) {
	a.meta.Corp = meta.Corp
	a.meta.EndPoint = meta.EndPoint
	a.meta.APIVersion = meta.APIVersion
	a.meta.APIKey = meta.APIKey
	a.meta.SecretKey = meta.SecretKey
	a.meta.Model = meta.Model
	a.meta.Region = meta.Region
	a.meta.APPID = meta.APPID
}
