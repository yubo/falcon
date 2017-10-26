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
	"sync"
	"time"

	"github.com/golang/glog"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/yubo/falcon"
)

type WxSession struct {
	Session string `json:"session_key"`
	Expires int64  `json:"expires_in"`
	Openid  string `json:"openid"`
}

type weappTask struct {
	key     string
	action  int32
	states  int32
	expires int64
	data    interface{}
}

type WeappSession struct {
	Key       string
	Expires   int64
	User      *User
	AppUser   *WeappUserInfo_t
	Wxsession string
	Wxexpires int64
	Wxopenid  string
}

type weapp_t struct {
	sync.RWMutex
	sessions map[string]*WeappSession
	tasks    map[string]*weappTask
}

const (
	WEAPP_TASK_BIND_USER = iota
	WEAPP_TASK_UNBIND_USER
	WEAPP_TASK_LOGIN
)

const (
	WEAPP_TASK_PENDING = iota
	WEAPP_TASK_DONE
	WEAPP_TASK_FAIL
)

var (
	ErrInvalidBlockSize    = errors.New("invalid blocksize")
	ErrInvalidPKCS7Data    = errors.New("invalid PKCS7 data (empty or not padded)")
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
	ErrInvalidAppid        = errors.New("invalid appid")
	ErrInvalidTask         = errors.New("invalid Task")
	weapp                  weapp_t
)

func init() {
	weapp.sessions = make(map[string]*WeappSession)
	weapp.tasks = make(map[string]*weappTask)
	ticker := time.NewTicker(time.Second * 5).C
	weapp.AddTask(&weappTask{
		key:     "1234567890",
		action:  WEAPP_TASK_BIND_USER,
		states:  WEAPP_TASK_PENDING,
		data:    int64(1000),
		expires: time.Now().Unix() + 60000,
	})
	go func() {
		for {
			<-ticker
			now := time.Now().Unix()
			for skey, session := range weapp.sessions {
				if session.Expires < now {
					delete(weapp.sessions, skey)
				}
			}
			for skey, task := range weapp.tasks {
				if task.expires < now {
					delete(weapp.tasks, skey)
				}
			}
		}
	}()
}

func (p *weapp_t) AddTask(t *weappTask) error {
	p.Lock()
	defer p.Unlock()
	p.tasks[t.key] = t
	return nil
}
func (p *weapp_t) delTask(key string) (*weappTask, error) {
	p.Lock()
	defer p.Unlock()
	if t, ok := p.tasks[key]; !ok {
		return nil, falcon.ErrNoExits
	} else {
		delete(p.tasks, key)
		return t, nil
	}
}
func (p *weapp_t) getTask(key string) (*weappTask, error) {
	p.RLock()
	defer p.RUnlock()
	if t, ok := p.tasks[key]; !ok {
		return nil, falcon.ErrNoExits
	} else {
		return t, nil
	}
}

func (p *weapp_t) addSession(s *WeappSession) error {
	p.Lock()
	defer p.Unlock()
	p.sessions[s.Key] = s
	return nil
}
func (p *weapp_t) delSession(key string) (*WeappSession, error) {
	p.Lock()
	defer p.Unlock()
	if t, ok := p.sessions[key]; !ok {
		return nil, falcon.ErrNoExits
	} else {
		delete(p.sessions, key)
		return t, nil
	}
}
func (p *weapp_t) getSession(key string) (*WeappSession, error) {
	p.RLock()
	defer p.RUnlock()
	if t, ok := p.sessions[key]; !ok {
		return nil, falcon.ErrNoExits
	} else {
		return t, nil
	}
}

func WeappGetSession(key string) (*WeappSession, error) {
	return weapp.getSession(key)
}

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

func WeappLogin(code, encrypt_data, iv string) (*WeappSession, error) {
	//wxappid, wxappsecret
	var appuser WeappUserInfo_t

	session, err := GetWeappSession(code)
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

	err = json.Unmarshal(text, &appuser)
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("appuser %s", appuser)

	if appuser.Watermark.Appid != wxappid {
		return nil, ErrInvalidAppid
	}

	now := time.Now().Unix()

	s := &WeappSession{
		Key:       RandString(32),
		Expires:   now + 3600,
		AppUser:   &appuser,
		Wxsession: session.Session,
		Wxexpires: now + session.Expires,
		Wxopenid:  session.Openid,
	}

	if user, err := SysOp.GetUserByUuid(s.AppUser.OpenId + "@weapp"); err == nil {
		if user, err = SysOp.GetUser(user.Muid); err == nil {
			s.User = user
		}
	}

	weapp.addSession(s)

	return s, nil
}

// response: '{"session_key":"AmBppkx\/wIcgLBfesd\/mig==","expires_in":7200,"openid":"owU0h0eb8mzxHIiSdtjlIbyUjV3U"}'
//
func GetWeappSession(code string) (*WxSession, error) {
	var session WxSession
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", wxappid, wxappsecret, code)

	err := getJson(url, &session, 30*time.Second)

	return &session, err
}

func WeappOpenid(code string) (string, error) {
	session, err := GetWeappSession(code)
	if err != nil {
		return "", err
	}

	return session.Openid, nil
}

type QrTask struct {
	Key string `json:"key"`
	Img string `json:"img"`
}

func WeappBindQr(uid int64) (*QrTask, error) {
	task := &weappTask{
		key:     RandString(64),
		action:  WEAPP_TASK_BIND_USER,
		states:  WEAPP_TASK_PENDING,
		data:    uid,
		expires: time.Now().Unix() + 60,
	}

	weapp.Lock()
	defer weapp.Unlock()
	weapp.tasks[task.key] = task

	png, err := qrcode.Encode(task.key, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	return &QrTask{Key: task.key, Img: base64.StdEncoding.EncodeToString(png)}, nil
}

func WeappTaskAck(key string, sess *WeappSession) (interface{}, error) {
	task, err := weapp.getTask(key)
	if err != nil {
		return nil, err
	}

	if task.states != WEAPP_TASK_PENDING {
		return nil, ErrInvalidTask
	}

	switch task.action {
	case WEAPP_TASK_BIND_USER:
		if user, err := weappBind(task.data.(int64), sess); err != nil {
			task.states = WEAPP_TASK_FAIL
			task.expires = time.Now().Unix() + 5
			return nil, err
		} else {
			task.states = WEAPP_TASK_DONE
			task.expires = time.Now().Unix() + 5
			return user, err
		}
	case WEAPP_TASK_UNBIND_USER:
	case WEAPP_TASK_LOGIN:
	default:
		return nil, falcon.EINVAL
	}
	return nil, nil
}

func weappBind(uid int64, sess *WeappSession) (user *User, err error) {

	if user, err = GetUser(uid, SysOp.O); err != nil {
		return
	}

	// create user from weapp userinfo
	wuser := &User{
		Uuid:      fmt.Sprintf("%s@weapp", sess.AppUser.OpenId),
		Name:      sess.AppUser.OpenId,
		Cname:     sess.AppUser.NickName,
		AvatarUrl: sess.AppUser.AvatarUrl,
	}

	wuser, err = SysOp.AddUser(wuser)
	if err != nil {
		return
	}

	err = SysOp.BindUser(wuser.Id, user.Id)
	return
}
