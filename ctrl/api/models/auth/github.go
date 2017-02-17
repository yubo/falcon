/*
 * Copyright 2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/yubo/falcon/ctrl/api/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const (
	baseURL = "https://api.github.com"
)

type githubAuth struct {
	models.AuthModule
	config oauth2.Config
}

var (
	_github = &githubAuth{
		AuthModule: models.AuthModule{Name: "github"},
		config: oauth2.Config{
			Endpoint: github.Endpoint,
			Scopes:   []string{"user:email"},
		},
	}
)

func init() {
	models.RegisterAuth(_github)
}

func (p *githubAuth) PreStart() error {
	if p.AuthModule.Prestarted {
		return models.ErrRePreStart
	}
	p.AuthModule.Prestarted = true
	p.config.ClientID = beego.AppConfig.String("githubclientid")
	p.config.ClientSecret = beego.AppConfig.String("githubclientsecret")
	p.config.RedirectURL = beego.AppConfig.String("githubredirecturl")
	beego.Debug("clientid", p.config.ClientID)
	return nil
}

func (p *githubAuth) AuthorizeUrl(c interface{}) string {
	ctx := c.(*context.Context)

	v := url.Values{}
	v.Set("cb", ctx.Input.Query("cb"))

	conf := p.config
	conf.RedirectURL = fmt.Sprintf("%s?%s", conf.RedirectURL, v.Encode())
	return conf.AuthCodeURL(models.RandString(8))
}

func (p *githubAuth) CallBack(c interface{}) (uuid string, err error) {
	var user user

	r := c.(*context.Context).Request
	q := r.URL.Query()

	if errType := q.Get("error"); errType != "" {
		err = fmt.Errorf("%s:%s", errType, q.Get("error_description"))
		return
	}

	ctx := r.Context()

	token, err := p.config.Exchange(ctx, q.Get("code"))
	if err != nil {
		err = fmt.Errorf("github: failed to get token: %v", err)
		return
	}

	client := p.config.Client(ctx, token)

	user, err = p.user(r, client)
	if err != nil {
		err = fmt.Errorf("github: get user: %v", err)
		return
	}

	if user.Name == "" {
		err = fmt.Errorf("github: get empty username")
		return
	}
	uuid = fmt.Sprintf("%s@%s", user.Name, p.GetName())
	return
}

type user struct {
	Name  string `json:"name"`
	Login string `json:"login"`
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// user queries the GitHub API for profile information using the provided client. The HTTP
// client is expected to be constructed by the golang.org/x/oauth2 package, which inserts
// a bearer token as part of the request.
func (c *githubAuth) user(r *http.Request, client *http.Client) (user, error) {
	var u user
	req, err := http.NewRequest("GET", baseURL+"/user", nil)
	if err != nil {
		return u, fmt.Errorf("github: new req: %v", err)
	}
	req = req.WithContext(r.Context())
	resp, err := client.Do(req)
	if err != nil {
		return u, fmt.Errorf("github: get URL %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return u, fmt.Errorf("github: read body: %v", err)
		}
		return u, fmt.Errorf("%s: %s", resp.Status, body)
	}

	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode response: %v", err)
	}
	return u, nil
}
