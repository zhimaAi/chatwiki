// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/pkg/lib_define"
	"errors"

	"github.com/zhimaAi/go_tools/tool"
)

func IsSafeMode(signature, encryptType string) bool {
	return len(signature) > 0 && encryptType == `aes`
}

func MsgDecrypt(content, msgSignature, nonce, timestamp string) (string, error) {
	if !VerifySignature(msgSignature, timestamp, nonce, content) {
		return "", errors.New(`signature verification failure`)
	}
	key, err := tool.Base64Decode(lib_define.AesKey + `=`)
	if err != nil {
		return "", err
	}
	content, err = tool.Base64Decode(content)
	if err != nil {
		return "", err
	}
	bs, err := AesDecrypt([]byte(content), []byte(key))
	if err != nil {
		return "", err
	}
	if len(bs) < 20 {
		return "", errors.New(`aes error:` + string(bs))
	}
	length := int(BytesToUint32(bs[16:20]))
	if len(bs) < 16+4+length {
		return "", errors.New(`length error:` + string(bs))
	}
	return string(bs)[16+4 : 16+4+length], nil
}

func MsgEncrypt(xmlStr string) (string, error) {
	key, err := tool.Base64Decode(lib_define.AesKey + `=`)
	if err != nil {
		return "", err
	}
	xmlStr = tool.Random(16) + string(Uint32ToBytes(uint32(len(xmlStr)))) + xmlStr + lib_define.OpenAppid
	bs, err := AesEncrypt([]byte(xmlStr), []byte(key))
	if err != nil {
		return "", err
	}
	encrypt := tool.Base64Encode(string(bs))
	timestamp, nonce := tool.Time2String(), tool.Random(10)
	sign := GenerateSignature(timestamp, nonce, encrypt)
	result := `<xml>
	<Encrypt><![CDATA[` + encrypt + `]]></Encrypt>
	<MsgSignature><![CDATA[` + sign + `]]></MsgSignature>
	<TimeStamp>` + timestamp + `</TimeStamp>
	<Nonce><![CDATA[` + nonce + `]]></Nonce>
</xml>`
	return result, nil
}
