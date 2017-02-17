/*
 * Copyright 2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package auth

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/yubo/falcon/ctrl/api/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOauth2 "google.golang.org/api/oauth2/v1"
)

const (
	baseURL = "https://api.github.com"
)

type googleAuth struct {
	models.AuthModule
	config oauth2.Config
}

var (
	_google = &googleAuth{
		AuthModule: models.AuthModule{Name: "google"},
		config: oauth2.Config{
			Endpoint: google.Endpoint,
			Scopes:   []string{googleOauth2.PlusMeScope, googleOauth2.UserinfoEmailScope},
		},
	}
)

func init() {
	models.RegisterAuth(_google)
}

func (p *googleAuth) PreStart() error {
	if p.AuthModule.Prestarted {
		return models.ErrRePreStart
	}
	p.AuthModule.Prestarted = true
	p.config.ClientID = beego.AppConfig.String("googleclientid")
	p.config.ClientSecret = beego.AppConfig.String("googleclientsecret")
	p.config.RedirectURL = beego.AppConfig.String("googleredirecturl")
	beego.Debug("clientid", p.config.ClientID)
	return nil
}

func (p *googleAuth) AuthorizeUrl(c interface{}) string {
	// not support cb param

	conf := p.config
	return conf.AuthCodeURL(models.RandString(8))
}

func (p *googleAuth) CallBack(c interface{}) (uuid string, err error) {
	r := c.(*context.Context).Request
	q := r.URL.Query()

	if errType := q.Get("error"); errType != "" {
		err = fmt.Errorf("%s:%s", errType, q.Get("error_description"))
		return
	}

	ctx := r.Context()

	token, err := p.config.Exchange(ctx, q.Get("code"))
	if err != nil {
		err = fmt.Errorf("google: failed to get token: %v", err)
		return
	}

	client := p.config.Client(ctx, token)

	svc, err := googleOauth2.New(client)
	if err != nil {
		err = fmt.Errorf("google: get user: %v", err)
		return
	}
	user, err := googleOauth2.NewUserinfoService(svc).V2.Me.Get().Do()

	if !*user.VerifiedEmail {
		err = fmt.Errorf("google: email not verified")
	}

	uuid = fmt.Sprintf("%s@%s", user.Email, p.GetName())
	return
}
