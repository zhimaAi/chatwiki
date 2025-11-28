package lib_define

// 飞书消息回调
type FeishuMsgEvent struct {
	Schema    string `json:"schema"`
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
	Header    struct {
		EventId    string `json:"event_id"`
		Token      string `json:"token"`
		CreateTime string `json:"create_time"`
		EventType  string `json:"event_type"`
		TenantKey  string `json:"tenant_key"`
		AppId      string `json:"app_id"`
	} `json:"header"`
	Event struct {
		Message struct {
			ChatId     string `json:"chat_id"`
			ChatType   string `json:"chat_type"`
			Content    string `json:"content"`
			CreateTime string `json:"create_time"`
			Mentions   []struct {
				Id struct {
					OpenId  string      `json:"open_id"`
					UnionId string      `json:"union_id"`
					UserId  interface{} `json:"user_id"`
				} `json:"id"`
				Key       string `json:"key"`
				Name      string `json:"name"`
				TenantKey string `json:"tenant_key"`
			} `json:"mentions"`
			MessageId   string `json:"message_id"`
			MessageType string `json:"message_type"`
			UpdateTime  string `json:"update_time"`
		} `json:"message"`
		Sender struct {
			SenderId struct {
				OpenId  string      `json:"open_id"`
				UnionId string      `json:"union_id"`
				UserId  interface{} `json:"user_id"`
			} `json:"sender_id"`
			SenderType string `json:"sender_type"`
			TenantKey  string `json:"tenant_key"`
		} `json:"sender"`
		ChatID                string `json:"chat_id"`
		LastMessageCreateTime string `json:"last_message_create_time"`
		LastMessageID         string `json:"last_message_id"`
		OperatorID            struct {
			OpenID  string      `json:"open_id"`
			UnionID string      `json:"union_id"`
			UserID  interface{} `json:"user_id"`
		} `json:"operator_id"`
	} `json:"event"`
}

type DingtalkImgContent struct {
	PictureDownloadCode string `json:"pictureDownloadCode"`
	DownloadCode        string `json:"downloadCode"`
}

type DingtalkMsgEvent struct {
	SenderPlatform string `json:"senderPlatform"`
	ConversationId string `json:"conversationId"`
	AtUsers        []struct {
		DingtalkId string `json:"dingtalkId"`
	} `json:"atUsers"`
	ChatbotCorpId             string `json:"chatbotCorpId"`
	ChatbotUserId             string `json:"chatbotUserId"`
	OpenThreadId              string `json:"openThreadId"`
	MsgId                     string `json:"msgId"`
	SenderNick                string `json:"senderNick"`
	IsAdmin                   bool   `json:"isAdmin"`
	SenderStaffId             string `json:"senderStaffId"`
	SessionWebhookExpiredTime int64  `json:"sessionWebhookExpiredTime"`
	CreateAt                  int64  `json:"createAt"`
	SenderCorpId              string `json:"senderCorpId"`
	ConversationType          string `json:"conversationType"`
	SenderId                  string `json:"senderId"`
	ConversationTitle         string `json:"conversationTitle"`
	IsInAtList                bool   `json:"isInAtList"`
	SessionWebhook            string `json:"sessionWebhook"`
	Text                      struct {
		Content string `json:"content"`
	} `json:"text"`
	Content   DingtalkImgContent `json:"content"`
	RobotCode string             `json:"robotCode"`
	Msgtype   string             `json:"msgtype"`
}
