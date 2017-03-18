/*
 * Copyright 2016 2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package plugins

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/ctrl"
)

func init() {
	ctrl.RegisterPrestart(start)
}

func start(conf *falcon.ConfCtrl) error {
	ticker := time.NewTicker(time.Second * 60)
	go func() {
		for {
			select {
			case <-ticker.C:
				beego.Debug("demo")
			}
		}
	}()
	return nil
}
