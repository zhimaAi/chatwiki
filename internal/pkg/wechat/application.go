// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package wechat

import (
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/wechat/dingtalk_robot"
	"chatwiki/internal/pkg/wechat/feishu_robot"
	"chatwiki/internal/pkg/wechat/mini_program"
	"chatwiki/internal/pkg/wechat/official_account"
	"chatwiki/internal/pkg/wechat/wechat_kefu"
	"errors"
	"net/http"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	publishresponse "github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/publish/response"
	openresponse "github.com/ArtisanCloud/PowerWeChat/v3/src/openPlatform/authorizer/miniProgram/account/response"
	"github.com/zhimaAi/go_tools/msql"
)

type ApplicationInterface interface {
	SetTyping(customer, command string) (int, error)
	SendText(customer, content string, push *lib_define.PushMessage) (int, error)
	GetToken(refresh bool) (*response.ResponseGetToken, int, error)
	SendMsgOnEvent(code, content string) (int, error)
	GetCustomerInfo(customer string) (map[string]any, int, error)
	UploadTempImage(filePath string) (string, int, error)
	SendImage(customer, filePath string, push *lib_define.PushMessage) (int, error)
	GetFileByMedia(mediaId string, push *lib_define.PushMessage) ([]byte, http.Header, int, error)
	SendUrl(customer, url, title string, push *lib_define.PushMessage) (int, error)                                               // 发送链接
	SendMiniProgramPage(customer, appid, title, pagePath, localThumbURL string, push *lib_define.PushMessage) (int, error)        // 发送小程序卡片
	SendImageTextLink(customer, url, title, description, localThumbURL, picurl string, push *lib_define.PushMessage) (int, error) // 发送图文链接
	SendSmartMenu(customer string, smartMenu lib_define.SmartMenu, push *lib_define.PushMessage) (int, error)                     // 发送智能菜单
	GetAccountBasicInfo() (*openresponse.ResponseGetBasicInfo, int, error)
	SendVoice(customer, filePath string, push *lib_define.PushMessage) (int, error)
}

// OfficialAccountInterface 专门用于公众号的方法
type OfficialAccountInterface interface {
	GetSubscribeScene(openid string) (string, error)
	GetPublishedMessageList(offset, count, noContent int) (*publishresponse.ResponseBatchGet, error)
	GetPublishedArticle(articleId string) (*publishresponse.ResponsePublishGetArticle, error)
	GetAccountClient() (*officialAccount.OfficialAccount, error)
}

type FeishuInterface interface {
	GetDocFileList(folderToken string) ([]feishu_robot.FeishuDocFile, error)
	GetDocFileTree(userAccessToken string, folderToken string) ([]*feishu_robot.FeishuDocFileTree, error)
	GetDocFileDetail(documentId string) (string, string, error)
	BuildUserAuthLoginUrl(redirectUri, frontRedirectUrl string) (string, error)
	GetUserAccessToken(code, redirectUri string) (*feishu_robot.FeishuUserAccessToken, error)
}

func GetApplication(appInfo msql.Params) (ApplicationInterface, error) {
	if len(appInfo) == 0 {
		return nil, errors.New(`app info is empty`)
	}
	switch appInfo[`app_type`] {
	case lib_define.AppOfficeAccount:
		return &official_account.Application{AppID: appInfo[`app_id`], Secret: appInfo[`app_secret`]}, nil
	case lib_define.AppMini:
		return &mini_program.Application{AppID: appInfo[`app_id`], Secret: appInfo[`app_secret`]}, nil
	case lib_define.AppWechatKefu:
		return &wechat_kefu.Application{AppID: appInfo[`app_id`], Secret: appInfo[`app_secret`]}, nil
	case lib_define.FeiShuRobot:
		return &feishu_robot.Application{AppID: appInfo[`app_id`], Secret: appInfo[`app_secret`]}, nil
	case lib_define.DingTalkRobot:
		return &dingtalk_robot.Application{AppID: appInfo[`app_id`], Secret: appInfo[`app_secret`]}, nil
	}
	return nil, errors.New(`app type not support`)
}
