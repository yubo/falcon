/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/transfer"
	"golang.org/x/net/context"
)

var (
	Client *ClientModule
)

type ClientModule struct {
	ctx    context.Context
	cancel context.CancelFunc

	client      transfer.TransferClient
	callTimeout int
}

func (p *ClientModule) PreStart(ctrl *Ctrl) error {
	Client = p
	return nil
}

func (p *ClientModule) Start(ctrl *Ctrl) error {
	conf := &ctrl.Conf.Ctrl
	p.callTimeout, _ = conf.Int(C_CALL_TIMEOUT)

	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.worker(conf.Str(C_TRANSFER_ADDR))

	return nil
}

func (p *ClientModule) Stop(ctrl *Ctrl) error {
	p.cancel()
	return nil
}

func (p *ClientModule) Reload(ctrl *Ctrl) error {
	return nil
}

func (p *ClientModule) worker(transferAddr string) error {

	go func() {
		conn, _, err := falcon.DialRr(p.ctx, transferAddr, true)
		if err != nil {
			return
		}
		defer conn.Close()

		p.client = transfer.NewTransferClient(conn)
		select {
		case <-p.ctx.Done():
			return
		}
	}()
	return nil
}

func GetDps(req *transfer.GetRequest) (*transfer.GetResponse, error) {

	statsInc(ST_GET_ITERS, 1)
	statsInc(ST_GET_DPS, len(req.Keys))

	glog.V(6).Infof("%s tx get %v", MODULE_NAME, len(req.Keys))

	ctx, _ := context.WithTimeout(context.Background(),
		time.Duration(Client.callTimeout)*time.Millisecond)
	resp, err := Client.client.Get(ctx, req)
	if err != nil {
		statsInc(ST_GET_ITERS_ERR, 1)
	} else {
		statsInc(ST_GET_DPS_ERR, int(len(req.Keys)-len(resp.Data)))
	}
	return resp, err
}
