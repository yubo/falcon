/*
 * Copyright 2016 yubo. All rights reserved.
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

	"github.com/astaxie/beego/context"
	"github.com/yubo/falcon/ctrl/api/models"
	"github.com/yubo/falcon/ctrl/config"
	"github.com/yubo/falcon/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const (
	GITHUB_BASE_URL = "https://api.github.com"
	GITHUB_NAME     = "github"
)

type githubAuth struct {
	config oauth2.Config
}

func init() {
	models.RegisterAuth(GITHUB_NAME, &githubAuth{})
}

func (p *githubAuth) Init(conf *config.ConfCtrl) error {
	p.config = oauth2.Config{
		Endpoint:     github.Endpoint,
		Scopes:       []string{"user:email"},
		ClientID:     conf.Ctrl.Str(utils.C_GITHUB_CLIENT_ID),
		ClientSecret: conf.Ctrl.Str(utils.C_GITHUB_CLIENT_SECRET),
		RedirectURL:  conf.Ctrl.Str(utils.C_GITHUB_REDIRECT_URL),
	}
	return nil
}

func (p *githubAuth) Verify(c interface{}) (bool, string, error) {
	return false, "", utils.EPERM
}

func (p *githubAuth) AuthorizeUrl(c interface{}) string {
	ctx := c.(*context.Context)

	v := url.Values{}
	v.Set("cb", ctx.Input.Query("cb"))

	conf := p.config
	conf.RedirectURL = fmt.Sprintf("%s?%s", conf.RedirectURL, v.Encode())
	return conf.AuthCodeURL(models.RandString(8))
}

func (p *githubAuth) LoginCb(c interface{}) (uuid string, err error) {
	var user githubUser

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

	uuid = fmt.Sprintf("%s@%s", user.Login, GITHUB_NAME)
	return
}

func (p *githubAuth) LogoutCb(c interface{}) {
}

type githubUser struct {
	Name  string `json:"name"`
	Login string `json:"login"`
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// user queries the GitHub API for profile information using the provided client. The HTTP
// client is expected to be constructed by the golang.org/x/oauth2 package, which inserts
// a bearer token as part of the request.
func (c *githubAuth) user(r *http.Request, client *http.Client) (githubUser, error) {
	var u githubUser
	req, err := http.NewRequest("GET", GITHUB_BASE_URL+"/user", nil)
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
