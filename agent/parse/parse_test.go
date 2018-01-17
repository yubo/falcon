/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package parse

import (
	"testing"

	"github.com/golang/glog"
)

var text = []byte(`
{
	debug;
	leasekey	"/open-falcon/agent/online/test.falcon";
	leasettl	20;
	host		localhost;
	http		on;
	interval	5;
	httpAddr	0.0.0.0:7007;
	ifacePrefix	eth,em;
	workerProcesses	3;
	callTimeout	60*(3+1);
	plugins     	sys;
//	upstream	stdout;
//	upstream	"127.0.0.1:1234,127.0.0.1:1235";
	upstream        stdout;
	emuTplDir	./var/tpl;	
	plugins		"sys,emulator";
};
`)

func TestParse(t *testing.T) {
	ret := Parse(text, "test", 0)
	glog.V(4).Infof("%s", ret)
}
