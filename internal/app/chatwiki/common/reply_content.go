// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import "chatwiki/internal/pkg/lib_define"

// ReplyContent response content
type ReplyContent struct {
	ReplyType       string               `json:"reply_type" form:"reply_type"`
	ThumbURL        string               `json:"thumb_url" form:"thumb_url"`
	Title           string               `json:"title" form:"title"`
	Description     string               `json:"description" form:"description"`
	URL             string               `json:"url" form:"url"`
	PagePath        string               `json:"page_path" form:"page_path"`
	Appid           string               `json:"appid" form:"appid"`
	Status          string               `json:"status" form:"status"`
	MoreImgTextJSON string               `json:"more_img_text_json" form:"more_img_text_json"`
	MediaID         string               `json:"media_id" form:"media_id"`
	Type            string               `json:"type" form:"type"`
	Pic             string               `json:"pic,omitempty" form:"pic,omitempty"`
	SmartMenuId     string               `json:"smart_menu_id" form:"smart_menu_id"`
	SmartMenu       lib_define.SmartMenu `json:"smart_menu,omitempty" form:"smart_menu"` // Smart menu, passed when outputting
	SendSource      string               `json:"send_source" form:"send_source"`
}

const (
	ReplyTypeImageText = `imageText` // Image-text
	ReplyTypeText      = `text`      // Text
	ReplyTypeUrl       = `url`       // URL
	ReplyTypeImg       = `image`     // Image
	ReplyTypeCard      = `card`      // Mini program card
	ReplyTypeSmartMenu = `smartMenu` // Smart menu
)

const (
	TypeText  = `text`
	TypeImage = `image`
	TypeFile  = `file`
	TypeAudio = `audio`
	TypeVideo = `video`
)
