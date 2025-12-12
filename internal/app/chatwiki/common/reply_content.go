package common

import "chatwiki/internal/pkg/lib_define"

// ReplyContent 回复内容
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
	SmartMenu       lib_define.SmartMenu `json:"smart_menu,omitempty" form:"smart_menu"` //智能菜单 输出时候传递
	SendSource      string               `json:"send_source" form:"send_source"`
}

const (
	ReplyTypeImageText = `imageText` //图文
	ReplyTypeText      = `text`      //文本
	ReplyTypeUrl       = `url`       //网址
	ReplyTypeImg       = `image`     //图片
	ReplyTypeCard      = `card`      //小程序卡片
	ReplyTypeSmartMenu = `smartMenu` //智能菜单
)

const (
	TypeText  = `text`
	TypeImage = `image`
	TypeFile  = `file`
	TypeAudio = `audio`
	TypeVideo = `video`
)
