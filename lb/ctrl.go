/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package lb

import (
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
)

func (p *Lb) scanCtrl() error {

	client, err := rpcDial(p.Params.CtrlAddr,
		time.Duration(p.Params.ConnTimeout)*time.Millisecond)
	if err != nil {
		return err
	}
	defer client.Close()

	bs := []specs.Backend{}
	err = netRpcCall(client, "CTRL.ListBackend", specs.Null{}, &bs,
		time.Duration(p.Params.CallTimeout)*time.Millisecond)
	if err != nil {
		return err
	}
	p.Backends = bs

	return nil
}

func (p *Lb) ctrlStart() {
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

func (p *Lb) ctrlStop() {
}
