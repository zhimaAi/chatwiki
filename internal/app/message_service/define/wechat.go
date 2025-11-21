// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package define

const PassiveReplyTextTemplate = `<xml>
	<ToUserName><![CDATA[:ToUserName]]></ToUserName>
	<FromUserName><![CDATA[:FromUserName]]></FromUserName>
	<CreateTime>:CreateTime</CreateTime>
	<MsgType><![CDATA[text]]></MsgType>
	<Content><![CDATA[:Content]]></Content>
</xml>`

const PassiveReplyImageTemplate = `<xml>
	<ToUserName><![CDATA[:ToUserName]]></ToUserName>
	<FromUserName><![CDATA[:FromUserName]]></FromUserName>
	<CreateTime>:CreateTime</CreateTime>
	<MsgType><![CDATA[image]]></MsgType>
	<Image>
		<MediaId><![CDATA[:MediaId]]></MediaId>
	</Image>
</xml>`
