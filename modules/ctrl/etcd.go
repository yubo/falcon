/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"github.com/yubo/falcon/lib/core"
)

type etcdCliModule struct {
}

func (p *etcdCliModule) PreStart(ctrl *Ctrl) error {
	return nil
}

func (p *etcdCliModule) Start(ctrl *Ctrl) error {
	cli := core.NewEtcdCli(ctrl.Conf.EtcdClient)
	cli.Prestart()
	cli.Start()
	ctrl.etcdCli = cli
	return nil
}

func (p *etcdCliModule) Stop(ctrl *Ctrl) error {
	ctrl.etcdCli.Stop()
	return nil
}

func (p *etcdCliModule) Reload(ctrl *Ctrl) error {
	return nil

	/* TODO
	p.Stop(ctrl)
	time.Sleep(time.Second)
	return p.Start(ctrl)
	*/
}
