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
	SERVICE_ADDR	= unix:./var/service.rpc;
	FALCON_DSN	= falcon:1234@tcp(localhost:3306)/falcon?loc=Local&parseTime=true;

	debug;
	leasekey		"/open-falcon/service/online/test.falcon";
	leasettl		20;
	apiAddr			${SERVICE_ADDR};
	httpAddr		${SERVICE_ADDR}.gw;
	dsn			${FALCON_DSN};
	idxDsn			${INDEX_DSN};
	dbMaxIdle		4;
	callTimeout		5000;	// 请求超时时间
	workerProcesses		2;	// 数据迁移时连接目标服务器的并发数量
	hdisk			"./var/data/hd01";
	confInterval		600;	// 同步配置间隔时间
	judgeInterval		60;	// 事件触发器扫描间隔时间
	shardIds		"0,1,2,3,4,5,6,7,8,9";
	lastTsthreshold		30;		// 最近活动时间阈值，超过的数据项会被忽略
	indexTimeout		3600*26;
	indexUpdateCycleTime	3600*24;
};
`)

func TestParse(t *testing.T) {
	ret := Parse(text, "test", 0)
	glog.V(4).Infof("%s", ret)
}
