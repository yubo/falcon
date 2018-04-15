/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"golang.org/x/net/context"

	"github.com/yubo/falcon/lib/core"
	"github.com/yubo/falcon/modules/transfer"
	"google.golang.org/grpc"
)

type clientModule struct {
	conn *grpc.ClientConn
}

func (p *clientModule) PreStart(ctrl *Ctrl) error {
	return nil
}

func (p *clientModule) Start(ctrl *Ctrl) error {
	if core.AddrIsDisable(ctrl.Conf.TransferAddr) {
		return nil
	}

	conn, _, err := core.DialRr(context.Background(), ctrl.Conf.TransferAddr, true)
	if err != nil {
		return err
	}
	ctrl.transferCli = transfer.NewTransferClient(conn)

	return nil
}

func (p *clientModule) Stop(ctrl *Ctrl) error {
	if ctrl.transferCli != nil {
		p.conn.Close()
	}
	ctrl.transferCli = nil
	return nil
}

func (p *clientModule) Reload(ctrl *Ctrl) error {
	// TODO
	return nil
}
