package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

func PKCS7Padding(ciphertext []byte) []byte {
	padding := aes.BlockSize - len(ciphertext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

//AESDecode aes解碼
func AESDecode(data string, aeskey string) string {
	key, _ := hex.DecodeString(aeskey)
	iv, _ := hex.DecodeString(aeskey)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	decodedata, err := hex.DecodeString(data)

	mode.CryptBlocks(decodedata, decodedata)

	decodedata = PKCS7UnPadding(decodedata)
	retstr := ""
	retstr = fmt.Sprintf("%s", decodedata)
	return retstr
}
