/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package config

import (
	"testing"

	"github.com/golang/glog"
)

var text = []byte(`
{
	debug;
	leasekey	"/open-falcon/transfer/online/test.falcon";
	leasettl	20;
	apiAddr		"unix:./var/transfer.rpc";
	httpAddr	"unix:./var/transfer.rpc.gw";
	workerProcesses	2;	// upstream 的并发连接数
	burstSize	2;	// client put burst size to remote service
	callTimeout	5000;	// 请求超时时间
	shardMap {	// 后端服务
		0 "unix:./var/service.rpc";
		1 "unix:./var/service.rpc";
		2 "unix:./var/service.rpc";
		3 "unix:./var/service.rpc";
		4 "unix:./var/service.rpc";
		5 "unix:./var/service.rpc";
		6 "unix:./var/service.rpc";
		7 "unix:./var/service.rpc";
		8 "unix:./var/service.rpc";
		9 "unix:./var/service.rpc";
	};
};
`)

func TestParse(t *testing.T) {
	ret := Parse(text, "test", 0)
	glog.V(4).Infof("%s", ret)
}
