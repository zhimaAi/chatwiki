// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"bytes"
	"chatwiki/internal/pkg/lib_define"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"sort"
	"strings"
)

func Sha1(data string) string {
	ctx := sha1.New()
	ctx.Write([]byte(data))
	return hex.EncodeToString(ctx.Sum(nil))
}

func VerifySignature(signature, timestamp, nonce string, argv ...string) bool {
	sl := []string{lib_define.SignToken, timestamp, nonce}
	for _, s := range argv {
		sl = append(sl, s)
	}
	sort.Strings(sl)
	ss := strings.Join(sl, ``)
	return Sha1(ss) == signature
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func BytesToUint32(bs4 []byte) uint32 {
	if len(bs4) > 4 {
		bs4 = bs4[:4]
	}
	return binary.BigEndian.Uint32(bs4)
}

func Uint32ToBytes(i uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, i)
	return buf
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func AesEncrypt(origData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, len(key))
	blockMode := cipher.NewCBCEncrypter(block, key[:aes.BlockSize])
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return encrypted, nil
}

func GenerateSignature(timestamp, nonce string, argv ...string) string {
	sl := []string{lib_define.SignToken, timestamp, nonce}
	for _, s := range argv {
		sl = append(sl, s)
	}
	sort.Strings(sl)
	ss := strings.Join(sl, ``)
	return Sha1(ss)
}
