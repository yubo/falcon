/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package auth

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/httplib"
	"github.com/yubo/falcon/lib/core"
	"github.com/yubo/falcon/modules/ctrl/api/models"
)

const (
	MISSO_NAME = "misso"
)

type missoAuth struct {
	RedirectURL string

	CookieSecretKey string
	missoAuthDomain string
	BrokerName      string
	SecretKey       string
	Credential      string
}

func init() {
	models.RegisterAuth(MISSO_NAME, &missoAuth{})
}

func (p *missoAuth) Init(conf *core.Configer) error {
	p.RedirectURL = conf.GetStr("redirect_url")
	p.CookieSecretKey = "secret-key-for-encrypt-cookie"
	p.missoAuthDomain = "http://sso.pt.xiaomi.com"
	p.BrokerName = "test"
	p.SecretKey = "test"
	return nil
}

func (p *missoAuth) Verify(c interface{}) (bool, string, error) {
	return false, "", core.EPERM
}

func (p *missoAuth) AuthorizeUrl(c interface{}) string {
	ctx := c.(*context.Context)

	p.Credential, _ = p.GenerateCredential()
	ctx.SetCookie("broker_cookie", p.Credential)

	v := url.Values{}
	v.Set("callback", p.RedirectURL)

	url, err := p.GetLoginUrl()
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s&%s", url, v.Encode())
}

func (p *missoAuth) LoginCb(c interface{}) (uuid string, err error) {
	ctx := c.(*context.Context)

	user_name, result := ctx.GetSecureCookie(p.CookieSecretKey, "user_name")
	broker_cookie := ctx.GetCookie("broker_cookie")

	//If can get user_name from cookie, user have logined
	if result == true {
		uuid = fmt.Sprintf("%s@%s", user_name, MISSO_NAME)
		return
	} else {
		if broker_cookie == "" {
			//cannot get broker_cookie, may first open in browser, or use service_account
			//try to get user_name from sso, may be login use service account
			authorization := ctx.Input.Header("Authorization")
			if authorization != "" {
				remote_ip := models.GetIPAdress(ctx.Input.Context.Request)
				user_name, err = p.GetServiceUser(authorization, remote_ip)
				if err == nil && user_name != "" {
					uuid = fmt.Sprintf("%s@%s", user_name, MISSO_NAME)
					return
				}
			}
		} else {
			p.Credential = broker_cookie
			if user_name, _ = p.GetUser(); user_name != "" {
				uuid = fmt.Sprintf("%s@%s", user_name, MISSO_NAME)
				ctx.SetSecureCookie(p.CookieSecretKey,
					"user_name", user_name)
				return
			}
		}
	}
	err = core.ErrLogin
	return
}

func (p *missoAuth) LogoutCb(c interface{}) {
	ctx := c.(*context.Context)
	ctx.SetCookie("user_name", "", -1)
}

/***********************
 * from sso_client.go
 ***********************/
func (p *missoAuth) GetUser() (string, error) {
	url := fmt.Sprintf("%s/login/broker/%s/broker_cookies/%s/user",
		p.missoAuthDomain, p.BrokerName, p.Credential)
	resp, err := httplib.Get(url).String()
	var resp_js map[string]string
	err = json.Unmarshal([]byte(resp), &resp_js)
	return resp_js["user_name"], err
}

func (p *missoAuth) GetServiceUser(authorization, user_ip string) (string,
	error) {

	auth_len := strings.Split(authorization, ";")
	if len(auth_len) != 3 {
		return "", fmt.Errorf("authorization wrong")
	}
	url := fmt.Sprintf("%s/mias/api/user_ip/%s/auth/%s/username",
		p.missoAuthDomain, user_ip, authorization)
	resp, err := httplib.Get(url).String()
	if err != nil {
		return "", err
	}
	resp_js := make(map[string]string)
	err = json.Unmarshal([]byte(resp), &resp_js)
	if err != nil {
		return "", fmt.Errorf(resp)
	}

	return resp_js["user_name"], nil
}

func (p *missoAuth) IsLogin() (bool, error) {
	url := fmt.Sprintf("%s/login/broker/%s/broker_cookies/%s/check",
		p.missoAuthDomain, p.BrokerName, p.Credential)
	resp, err := httplib.Get(url).String()
	if err != nil {
		return false, err
	}

	return resp == "1", nil
}

func (p *missoAuth) GetLogoutUrl() string {
	return fmt.Sprintf("%s/login/logout?broker_name=%s",
		p.missoAuthDomain, p.BrokerName)
}

func (p *missoAuth) GetLoginUrl() (string, error) {
	if p.Credential == "" {
		_, err := p.GenerateCredential()
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%s/login?broker_cookies=%s",
		p.missoAuthDomain, p.Credential), nil
}

func (p *missoAuth) GenerateCredential() (string, error) {
	url := p.missoAuthDomain + "/login/broker_cookies"
	req := httplib.Post(url)
	req.Param("broker_name", p.BrokerName)
	req.Param("secret_key", p.SecretKey)
	resp, err := req.SetTimeout(3*time.Second,
		3*time.Second).String()
	p.Credential = resp
	return resp, err
}
