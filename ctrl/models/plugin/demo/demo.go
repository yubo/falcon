/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package plugins

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/yubo/falcon/ctrl/models"
)

func init() {
	models.RegisterPlugin(start)
}

func start() error {
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
