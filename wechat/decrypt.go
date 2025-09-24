package wechat

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"log"
)

func (s *Session) DecryptTo(encryptedData, iv string, v any) error {

	if Debug() {
		log.Println("sessionKey: ", s.Key)
		log.Println("encryptedData: ", encryptedData)
		log.Println("iv: ", iv)
	}

	data, err := s.DecryptData(encryptedData, iv)
	if err != nil {
		if Debug() {
			log.Println("decrypt err: ", err.Error())
		}
		return err
	}

	if Debug() {
		log.Println("decryptedData: ", string(data))
	}

	return jsoniter.Unmarshal(data, &v)
}

// DecryptData 解密用户数据
func (s *Session) DecryptData(encryptedData, iv string) ([]byte, error) {

	if len(s.Key) != 24 {
		return nil, fmt.Errorf("illegal AES key")
	}
	key, err := base64.StdEncoding.DecodeString(s.Key)
	if err != nil {
		return nil, err
	}

	if len(iv) != 24 {
		return nil, fmt.Errorf("illegal IV")
	}
	aesIV, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}

	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("illegal buffer")
	}

	cipher.NewCBCDecrypter(block, aesIV).
		CryptBlocks(data, data)

	return bytes.Trim(data, "\x00\x0e"), nil // 去除多余填充
}
