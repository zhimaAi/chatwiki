// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/beevik/etree"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

type xmlMap map[string]interface{}

func (m *xmlMap) UnmarshalXML(d *xml.Decoder, _ xml.StartElement) error {
	type xmlMapEntry struct {
		XMLName xml.Name
		Value   string `xml:",chardata"`
	}
	*m = xmlMap{}
	for {
		var e xmlMapEntry
		if err := d.Decode(&e); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}

func IsInCharacterRange(r rune) bool {
	return r == 0x09 || r == 0x0A || r == 0x0D ||
		r >= 0x20 && r <= 0xD7FF ||
		r >= 0xE000 && r <= 0xFFFD ||
		r >= 0x10000 && r <= 0x10FFFF
}

func ReplaceInvalidCharacter(xmlStr string) string {
	rs := []rune(xmlStr)
	for i, r := range rs {
		if !IsInCharacterRange(r) {
			rs[i] = 42 // 42 is *
		}
	}
	return string(rs)
}

func ParseMessage(content string) map[string]interface{} {
	content = ReplaceInvalidCharacter(strings.TrimSpace(content))
	data := make(map[string]interface{})
	if len(content) > 0 && content[:1] == `<` {
		data = XmlParse(content)
	} else if len(content) > 0 {
		if err := tool.JsonDecode(content, &data); err != nil {
			logs.Error(err.Error()+`[%s]`, content)
			return make(map[string]interface{})
		}
	}
	return data
}

func GetMessage(body, signature, encryptType, msgSignature, nonce, timestamp string) map[string]interface{} {
	data := ParseMessage(body)
	if IsSafeMode(signature, encryptType) && data["Encrypt"] != nil {
		content, _ := data["Encrypt"].(string)
		raw, err := MsgDecrypt(content, msgSignature, nonce, timestamp)
		if err != nil {
			logs.Error(err.Error())
			if data["MsgId"] != nil {
				delete(data, `Encrypt`)
			} else {
				data = make(map[string]interface{})
			}
		} else {
			data = ParseMessage(raw)
		}
	}
	if len(data) == 0 {
		logs.Debug(body)
	}
	return data
}

func xmlErgodic(element *etree.Element) (map[string]interface{}, bool) {
	data := make(map[string]interface{})
	if element == nil {
		return data, true
	}
	elements := element.FindElements(`./*`)
	if len(elements) == 0 {
		return data, true
	}
	for idx := range elements {
		var text interface{}
		text, scalar := xmlErgodic(elements[idx])
		if scalar {
			text = elements[idx].Text()
		}
		if value, ok := data[elements[idx].Tag]; !ok {
			data[elements[idx].Tag] = text
		} else if values, ok2 := value.([]interface{}); !ok2 {
			data[elements[idx].Tag] = []interface{}{value, text}
		} else {
			data[elements[idx].Tag] = append(values, text)
		}
	}
	return data, false
}

func XmlParse(content string) map[string]interface{} {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(content); err != nil {
		logs.Error(`content:%s,err:%s`, content, err.Error())
		return make(map[string]interface{})
	}
	data, _ := xmlErgodic(doc.FindElement(`xml`))
	return data
}

// 飞书消息解密
func GetFeiShuMessage(encrypt string, key string) (string, error) {
	buf, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", fmt.Errorf("base64StdEncode Error[%v]", err)
	}
	if len(buf) < aes.BlockSize {
		return "", errors.New("cipher  too short")
	}
	keyBs := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(keyBs[:sha256.Size])
	if err != nil {
		return "", fmt.Errorf("AESNewCipher Error[%v]", err)
	}
	iv := buf[:aes.BlockSize]
	buf = buf[aes.BlockSize:]
	if len(buf)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(buf, buf)
	n := strings.Index(string(buf), "{")
	if n == -1 {
		n = 0
	}
	m := strings.LastIndex(string(buf), "}")
	if m == -1 {
		m = len(buf) - 1
	}
	return string(buf[n : m+1]), nil
}
