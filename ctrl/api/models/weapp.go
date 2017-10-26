/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

/*
 * encryptedData AES 128 CBC PKCS#7
 * https://github.com/go-web/tokenizer/blob/master/pkcs7.go
 * https://golang.org/pkg/crypto/cipher/#example_NewCBCDecrypter
 */
package models

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang/glog"
)

type OpenidGet struct {
	Openid  string `json:"openid"`
	Session string `json:"session"`
}

type WxSession struct {
	Session string `json:"session_key"`
	Expires int    `json:"expires_in"`
	Openid  string `json:"openid"`
}

var (
	ErrInvalidBlockSize    = errors.New("invalid blocksize")
	ErrInvalidPKCS7Data    = errors.New("invalid PKCS7 data (empty or not padded)")
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

func DesDecryption(key, iv, ciphertext []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, ErrInvalidBlockSize
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, ErrInvalidBlockSize
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(ciphertext, ciphertext)

	return pkcs7Unpad(ciphertext, aes.BlockSize)

}

// pkcs7Unpad validates and unpads data from the given bytes slice.
// The returned value will be 1 to n bytes smaller depending on the
// amount of padding, where n is the block size.
func pkcs7Unpad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	if len(b)%blocksize != 0 {
		return nil, ErrInvalidPKCS7Padding
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return b[:len(b)-n], nil
}

type WeappUserInfo_t struct {
	OpenId    string `json:"openId"`
	NickName  string `json:"nickName"`
	Gender    int32  `json:"gender"`
	Language  string `json:"language"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarUrl string `json:"avatarUrl"`
	Watermark struct {
		Appid     string `json:"appid"`
		Timestamp int64  `json:"timestamp"`
	} `json:"watermark"`
}

func WeappLogin(code, encrypt_data, iv string) (map[string]interface{}, error) {
	//wxappid, wxappsecret
	var userinfo WeappUserInfo_t
	var resp = make(map[string]interface{})

	session, err := WeappSession(code)
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("session %#v", session)

	keyB, _ := base64.StdEncoding.DecodeString(session.Session)
	ivB, _ := base64.StdEncoding.DecodeString(iv)
	dataB, _ := base64.StdEncoding.DecodeString(encrypt_data)

	text, err := DesDecryption(keyB, ivB, dataB)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(text, &userinfo)
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("userinfo %s", userinfo)

	resp["session"] = map[string]string{"id": session.Openid, "skey": RandString(32)}

	return resp, err
}

// response: '{"session_key":"AmBppkx\/wIcgLBfesd\/mig==","expires_in":7200,"openid":"owU0h0eb8mzxHIiSdtjlIbyUjV3U"}'
//
func WeappSession(code string) (*WxSession, error) {
	var session WxSession
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", wxappid, wxappsecret, code)

	err := getJson(url, &session, 30*time.Second)

	return &session, err
}

func WeappOpenid(code string) (*OpenidGet, error) {
	session, err := WeappSession(code)
	if err != nil {
		return nil, err
	}

	ret := &OpenidGet{Openid: session.Openid, Session: RandString(32)}

	return ret, err
}
