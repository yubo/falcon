/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package parse

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/golang/glog"
	_ "github.com/yubo/falcon/agent/modules"
	_ "github.com/yubo/falcon/alarm/modules"
	_ "github.com/yubo/falcon/ctrl/modules"
	_ "github.com/yubo/falcon/service/modules"
	_ "github.com/yubo/falcon/transfer/modules"
)

type Test struct {
	in  string
	out string
}

func init() {
	flag.Lookup("logtostderr").Value.Set("true")
	//flag.Lookup("v").Value.Set("5")
}

func TestExprText(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	goroot := os.Getenv("GOROOT")

	for i, test := range []Test{
		{"hello,{world}", "hello,{world}"},
		{"${GOPATH}", gopath},
		{"gopath:${GOPATH}", fmt.Sprintf("gopath:%s", gopath)},
		{"gopath:${GOPATH};goroot:${GOROOT}",
			fmt.Sprintf("gopath:%s;goroot:%s", gopath, goroot)},
	} {
		out := exprText([]byte(test.in))
		if out != test.out {
			t.Errorf(`#%d: exprText("%s")="%s"; want "%s"`,
				i, test.in, out, test.out)
		}
	}
}

var text = []byte(`
log	stdout	4;
pidFile	./falcon.pid;
//root	${GOPATH}/src/github.com/open-falcon/falcon/ctrl;

// 配置顺序决定了启动顺序，模块之间的依赖关系会限制启动顺序
// 可以设置本地的环境变量
// TRANSFER_ADDR = "8080";
// TRANSFER_ADDR = 0.0.0.0:8080;
FALCON_DSN    = falcon:1234@tcp(localhost:3306)/falcon?loc=Local&parseTime=true;
INDEX_DSN     = falcon:1234@tcp(localhost:3306)/idx?loc=Local&parseTime=true;
ALARM_DSN     = falcon:1234@tcp(localhost:3306)/alarm?loc=Local&parseTime=true;
AGENT_ADDR	= unix:./var/agent.rpc;
TRANSFER_ADDR	= unix:./var/transfer.rpc;
SERVICE_ADDR	= unix:./var/service.rpc;
ALARM_ADDR	= unix:./var/alarm.rpc;

AGENT_HTTP_ADDR		= ":8001";
TRANSFER_HTTP_ADDR	= ":8002";
SERVICE_HTTP_ADDR	= ":8003";
ALARM_HTTP_ADDR		= ":8004";
CTRL_HTTP_ADDR		= 127.0.0.1:8005;

agent {
	debug;
	burstSize	16;	// client put burst size to remote service
	leasekey	"/open-falcon/agent/online/test.falcon";
	leasettl	20;
	host		localhost;
	interval	5;
	apiAddr		${AGENT_ADDR};
	httpAddr	${AGENT_HTTP_ADDR};
	ifacePrefix	eth,em;
	workerProcesses	3;
	callTimeout	5000;
//	upstream	stdout;
//	upstream	127.0.0.1:1234,127.0.0.1:1235;
	upstream	${TRANSFER_ADDR};
	emuTplDir	./var/tpl;	
	plugins		"emulator";
};
`)

func TestParseText(t *testing.T) {
	ret := ParseText(text, "test", 0)
	glog.V(4).Infof("%s", ret)
}

func TestParse(t *testing.T) {
	ret := Parse("../docs/etc/falcon.example.conf")
	glog.V(4).Infof("%s", ret)
}
