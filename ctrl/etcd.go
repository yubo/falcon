/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ctrl

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/ctrl/config"
)

var (
	etcdCli *EtcdCliModule
)

type EtcdCliModule struct {
	cli *falcon.EtcdCli
}

func EtcdPuts(kvs map[string]string) error {
	if etcdCli == nil {
		return falcon.EPERM
	}
	return etcdCli.cli.Puts(kvs)
}

func EtcdGetPrefix(key string) (resp *clientv3.GetResponse, err error) {
	if etcdCli == nil {
		return nil, falcon.EPERM
	}
	return etcdCli.cli.GetPrefix(key)
}

func EtcdGet(key string) (string, error) {
	if etcdCli == nil {
		return "", falcon.EPERM
	}
	return etcdCli.cli.Get(key)
}

func (p *EtcdCliModule) PreStart(c *config.ConfCtrl) error {

	p.cli = falcon.NewEtcdCli(c.Ctrl)
	p.cli.Prestart()
	etcdCli = p
	return nil
}

func (p *EtcdCliModule) Start(c *config.ConfCtrl) error {
	p.cli.Start()
	return nil
}

func (p *EtcdCliModule) Stop(c *config.ConfCtrl) error {
	p.cli.Stop()
	return nil
}

func (p *EtcdCliModule) Reload(old, c *config.ConfCtrl) error {
	p.cli.Reload(c.Ctrl)
	return p.Start(c)
}
