/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
)

func (p *Backend) scanCtrl() error {

	client, err := dial(p.Conf.Params.CtrlAddr, p.Conf.Params.ConnTimeout)
	if err != nil {
		return err
	}
	defer client.Close()

	m := specs.Migrate{}
	err = netRpcCall(client, "CTRL.ListMigrate", specs.Null{}, &m,
		time.Duration(p.Conf.Params.CallTimeout)*time.Millisecond)
	if err != nil {
		return err
	}
	p.Conf.Migrate = m

	return nil
}

func (p *Backend) ctrlStart() {
	if err := p.scanCtrl(); err != nil {
		glog.Exitf(MODULE_NAME+"getBackends failed %s", err.Error())
	}
	ticker := time.NewTicker(time.Second * CTRL_STEP).C
	go func() {
		for {
			select {
			case _, ok := <-p.running:
				if !ok {
					return
				}
			case <-ticker:
				p.scanCtrl()
			}
		}
	}()
}

func (p *Backend) ctrlStop() {
}
