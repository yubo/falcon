/*
 * Copyright 2016 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package auth

import (
	"fmt"

	"github.com/astaxie/beego/context"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/ctrl/api/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOauth2 "google.golang.org/api/oauth2/v1"
)

const (
	GOOGLE_NAME = "google"
)

type googleAuth struct {
	config oauth2.Config
}

func init() {
	models.RegisterAuth(GOOGLE_NAME, &googleAuth{})
}

func (p *googleAuth) Init(conf *falcon.ConfCtrl) error {
	p.config = oauth2.Config{
		Endpoint:     google.Endpoint,
		Scopes:       []string{googleOauth2.PlusMeScope, googleOauth2.UserinfoEmailScope},
		ClientID:     conf.Ctrl.Str(falcon.C_GOOGLE_CLIENT_ID),
		ClientSecret: conf.Ctrl.Str(falcon.C_GOOGLE_CLIENT_SECRET),
		RedirectURL:  conf.Ctrl.Str(falcon.C_GOOGLE_REDIRECT_URL),
	}
	return nil
}

func (p *googleAuth) Verify(_c interface{}) (bool, string, error) {
	return false, "", models.EPERM
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

	uuid = fmt.Sprintf("%s@%s", user.Email, GOOGLE_NAME)
	return
}
