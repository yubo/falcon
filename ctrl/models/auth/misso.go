/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package auth

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/yubo/falcon/ctrl/controllers"
	"github.com/yubo/falcon/ctrl/models"
)

type missoAuth struct {
	models.AuthModule
	callback string

	CookieSecretKey string
	missoAuthDomain string
	BrokerName      string
	SecretKey       string
	Credential      string
}

var (
	_misso = &missoAuth{
		AuthModule:      models.AuthModule{Name: "misso"},
		CookieSecretKey: "secret-key-for-encrypt-cookie",
		missoAuthDomain: "http://sso.pt.xiaomi.com",
		BrokerName:      "test",
		SecretKey:       "test",
	}
)

func init() {
	models.RegisterAuth(_misso)
}

func (p *missoAuth) LoginHtml(_c interface{}) string {
	c := _c.(*controllers.AuthController)
	ctx := c.Ctx

	//If not login, at first generate a key from sso system
	p.Credential, _ = p.GenerateCredential()
	ctx.SetCookie("broker_cookie", p.Credential)
	base_url := ctx.Input.Context.Request.Host

	login_url, _ := p.GetLoginUrl()
	return fmt.Sprintf("<a href='%s&callback=http://%s/pub/auth/callback/"+
		"%s'>mioss</a>", login_url, base_url, p.Name)
}

func (p *missoAuth) CallBack(_c interface{}) {
	var (
		uuid string
	)
	c := _c.(*controllers.AuthController)
	ctx := c.Ctx

	remote_ips, ok := ctx.Input.Context.Request.Header["X-Forward-For"]
	remote_ip := ""
	if ok && len(remote_ips) > 0 {
		remote_ip = remote_ips[0]
	} else {
		remote_ip = ctx.Input.IP()
	}

	user_name, result := ctx.GetSecureCookie(p.CookieSecretKey, "user_name")
	broker_cookie := ctx.GetCookie("broker_cookie")

	//If can get user_name from cookie, user have logined
	if result == true {
		uuid = fmt.Sprintf("%s@%s", user_name, p.Name)
		goto success_out
	} else {
		if broker_cookie == "" {
			//cannot get broker_cookie, may first open in browser, or use service_account
			//try to get user_name from sso, may be login use service account
			authorization := ctx.Input.Header("Authorization")
			if authorization != "" {
				user_name, err := p.GetServiceUser(authorization, remote_ip)
				if err == nil && user_name != "" {
					uuid = fmt.Sprintf("%s@%s", user_name, p.Name)
					goto success_out
				}
			}
		} else {
			p.Credential = broker_cookie
			if user_name, _ = p.GetUser(); user_name != "" {
				uuid = fmt.Sprintf("%s@%s", user_name, p.Name)
				ctx.SetSecureCookie(p.CookieSecretKey,
					"user_name", user_name)
				goto success_out
			}
		}
	}

	ctx.Redirect(302, "/pub/auth/login")
	return
success_out:
	beego.Info("[Get user's user_name by broker_cookie][user_name:",
		user_name, "remote_ip:", ctx.Input.IP())
	//user_name
	c.Access(uuid)
	ctx.Redirect(302, "/")
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
	var resp_js map[string]string
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
