/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"github.com/astaxie/beego/context"
)

type PluginHooks struct {
	RateLimitsAccess func(*context.Context) bool
	GetTagHostCnt    func(tag, query string, deep bool) (int64, error)
	GetTagHost       func(tag, query string, deep bool,
		limit, offset int) (ret []*TagHostApiGet, err error)
}

var (
	Hooks PluginHooks
)
