/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package agent

import (
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
)

func (p *Agent) scanCtrl() error {

	client, err := rpcDial(p.Conf.Params.CtrlAddr,
		time.Duration(p.Conf.Params.ConnTimeout)*time.Millisecond)
	if err != nil {
		return err
	}
	defer client.Close()

	resp := []string{}
	err = netRpcCall(client, "CTRL.ListLb", falcon.Null{}, &resp,
		time.Duration(p.Conf.Params.CallTimeout)*time.Millisecond)
	if err != nil {
		return err
	}

	p.Conf.Upstreams = resp
	return nil
}

func (p *Agent) ctrlStart() {
	if err := p.scanCtrl(); err != nil {
		glog.Exitf(MODULE_NAME+"getLbs failed %s", err.Error())
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

func (p *Agent) ctrlStop() {
}
