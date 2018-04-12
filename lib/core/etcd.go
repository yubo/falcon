/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

/*
 * export ETCDCTL_API=3
 * etcdctl get --prefix /openfalcon
 */
package core

import (
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/golang/glog"
	"golang.org/x/net/context"
)

const (
	ETCD_CLIENT_MAX_RETRY       = 2
	ETCD_CLIENT_REQUEST_TIMEOUT = 3 * time.Second
)

// just for falcon-lite(graph/transfer)
type EtcdCliConfig struct {
	Endpoints  []string `json:"endpoints"`
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	CertFile   string   `json:"cert_file"`
	KeyFile    string   `json:"key_file"`
	CaFile     string   `json:"ca_file"`
	LeaseKey   string   `json:"lease_key"`
	LeaseValue string   `json:"lease_value"`
	LeaseTtl   int64    `json:"lease_ttl"`
}

type EtcdCli struct {
	Conf *EtcdCliConfig

	// runtime
	enable  bool
	leaseid clientv3.LeaseID
	config  clientv3.Config
	client  *clientv3.Client
	cancel  context.CancelFunc
}

func NewEtcdCli(c *EtcdCliConfig) *EtcdCli {
	return &EtcdCli{Conf: c}
}

func (p *EtcdCli) Prestart() {
	c := p.Conf
	if len(c.Endpoints) == 0 || c.Endpoints[0] == "" {
		return
	}

	p.config = clientv3.Config{
		Endpoints:   c.Endpoints,
		DialTimeout: 3 * time.Second,
		Username:    c.Username,
		Password:    c.Password,
	}

	if c.CertFile != "" && c.KeyFile != "" {
		tlsInfo := transport.TLSInfo{
			CertFile:      c.CertFile,
			KeyFile:       c.KeyFile,
			TrustedCAFile: c.CaFile,
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			glog.Infof("etcd ClientConfig() error %s", err.Error())
			return
		}
		p.config.TLS = tlsConfig
	}

	if err := p.connection(); err != nil {
		return
	}

	p.enable = true
	return
}

func (p *EtcdCli) connection() (err error) {
	p.client, err = clientv3.New(p.config)
	if err != nil {
		glog.Infof("etcd New() error %s", err.Error())
	}
	return
}

func (p *EtcdCli) reconnection() error {

	p.client.Close()
	if err := p.connection(); err != nil {
		return err
	}

	// start keepalive
	p.Start()

	return nil
}

func (p *EtcdCli) Get(key string) (string, error) {
	if !p.enable {
		return "", ErrUnsupported
	}

	for i := 0; i < ETCD_CLIENT_MAX_RETRY; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		resp, err := p.client.Get(ctx, key)
		cancel()
		if err == nil {
			if len(resp.Kvs) > 0 {
				return string(resp.Kvs[0].Value), nil
			}
			return "", ErrEmpty
		}
		glog.Infof("etcd Get(%s) error %s", key, err.Error())

		p.reconnection()
	}
	return "", clientv3.ErrNoAvailableEndpoints
}

func (p *EtcdCli) GetPrefix(key string) (resp *clientv3.GetResponse, err error) {
	if !p.enable {
		return nil, ErrUnsupported
	}

	for i := 0; i < ETCD_CLIENT_MAX_RETRY; i++ {
		ctx, cancel := context.WithTimeout(context.Background(),
			ETCD_CLIENT_REQUEST_TIMEOUT)
		resp, err = p.client.Get(ctx, key, clientv3.WithPrefix(),
			clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
		cancel()
		if err == nil {
			return
		}
		p.reconnection()
	}
	return
}

func (p *EtcdCli) Put(key, value string) (err error) {
	if !p.enable {
		return ErrUnsupported
	}

	for i := 0; i < ETCD_CLIENT_MAX_RETRY; i++ {
		ctx, cancel := context.WithTimeout(context.Background(),
			ETCD_CLIENT_REQUEST_TIMEOUT)
		_, err = p.client.Put(ctx, key, value)
		cancel()
		if err == nil {
			return
		}
		p.reconnection()
	}
	return
}

func (p *EtcdCli) Puts(kvs map[string]string) (err error) {
	if !p.enable {
		return ErrUnsupported
	}

	i := 0
	ops := make([]clientv3.Op, len(kvs))
	for k, v := range kvs {
		ops[i] = clientv3.OpPut(k, v)
		i++
	}

	for i := 0; i < ETCD_CLIENT_MAX_RETRY; i++ {
		ctx, cancel := context.WithTimeout(context.Background(),
			ETCD_CLIENT_REQUEST_TIMEOUT)
		_, err = clientv3.NewKV(p.client).Txn(ctx).If().Then(ops...).Commit()
		cancel()
		if err == nil {
			return
		}
		p.reconnection()
	}
	return
}

func (p *EtcdCli) Start() error {
	var ctx context.Context

	if !p.enable {
		glog.V(3).Infof("etcd client disabled")
		return nil
	}
	ctx, p.cancel = context.WithCancel(context.Background())

	resp, err := p.client.Grant(ctx, p.Conf.LeaseTtl)
	if err != nil {
		glog.Infof("etcd Grant() error %s", err.Error())
		p.enable = false
		return err
	}

	p.leaseid = resp.ID

	_, err = p.client.Put(ctx, p.Conf.LeaseKey, p.Conf.LeaseValue,
		clientv3.WithLease(p.leaseid))
	if err != nil {
		glog.Infof("etcd put with lease error %s", err.Error())
		p.enable = false
		return err
	}

	// the key will be kept forever
	_, err = p.client.KeepAlive(ctx, p.leaseid)
	if err != nil {
		glog.Infof("etcd keepalive error %s", err.Error())
		p.enable = false
		return err
	}
	glog.V(3).Infof("etcd keepalive success")
	return nil
}

func (p *EtcdCli) Stop() {
	if !p.enable {
		return
	}
	// cancel background routine
	p.cancel()

	p.client.Revoke(context.Background(), p.leaseid)

	p.client.Close()
}

func (p *EtcdCli) Reload(c *EtcdCliConfig) {
	/* TODO
	p.Stop()
	p.Prestart()
	p.Start()
	*/
	return
}
