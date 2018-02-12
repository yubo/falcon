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
	TSDB_DIR	= "var/tsdb";

	debug;
	leasekey	"/open-falcon/service/online/test.falcon";
	leasettl	20;
	apiAddr		${SERVICE_ADDR};
	httpAddr	${SERVICE_ADDR}.gw;
	dsn		${FALCON_DSN};
	idxDsn		${INDEX_DSN};
	dbMaxIdle	4;
	callTimeout	5000;	// 请求超时时间
	confInterval	600;	// 同步配置间隔时间
	judgeInterval	60;	// 事件触发器扫描间隔时间
	shardNum	10;
	shardIds	"0,1,2,3,4,5,6,7,8,9";
	cacheTimeout	24*3600;
	rrdTimeout	31*24*3600;
	tsdbBucketNum	13;
	tsdbBucketSize	2*3600;	// time (second)
	//tsdbTimeout	26*3600;// tsdbBucketNum * tsdbBucketSize
	tsdbDir		${TSDB_DIR};
};
`)

func TestParse(t *testing.T) {
	ret := Parse(text, "test", 0)
	glog.V(4).Infof("%s", ret)
}
