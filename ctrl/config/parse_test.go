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
	FALCON_DSN		= falcon:1234@tcp(localhost:3306)/falcon?loc=Local&parseTime=true;
	INDEX_DSN		= falcon:1234@tcp(localhost:3306)/idx?loc=Local&parseTime=true;
	ALARM_DSN		= falcon:1234@tcp(localhost:3306)/alarm?loc=Local&parseTime=true;
	AGENT_ADDR		= unix:./var/agent.rpc;
	TRANSFER_ADDR		= unix:./var/transfer.rpc;
	SERVICE_ADDR		= unix:./var/service.rpc;
	ALARM_ADDR		= unix:./var/alarm.rpc;
	CTRL_HTTP_ADDR		= ":8005";


//	disabled;
	MasterMode		true;	//default:	true;
	miMode			false;	//default:	false;
	DevMode			true;	//default:	false;
	debug			1;
	admin			test,service-ocean@misso;
	dsn			${FALCON_DSN}; //falcon:1234@tcp(localhost:3306)/falcon?loc=Local&parseTime=true;

	idxDsn			${INDEX_DSN};
	alarmDsn		${ALARM_DSN};
	etcdendpoints		${FALCON_ETCD};
	transferAddr		${TRANSFER_ADDR};
	callTimeout		5000;
	leasekey		"/open-falcon/ctrl/online/test.dev";
	leasettl		20;
	dbMaxIdle		30;
	dbMaxConn		30;
//	dbSchema		"${GOPATH}/src/github.com/yubo/falcon/scripts/db_schema/03_falcon.sql";
	HttpAddr		${CTRL_HTTP_ADDR};
	EnableDocs		true;
//	SessionGcMaxLifeTime	86400;
//	SessionCookieLifeTime	86400;
	AuthModule		ldap,misso,github,google;
	CacheModule		"host,role,system,tag,user";
	LdapAddr		localhost:389;
	LdapBaseDn		dc=xiaomi,dc=com;
	LdapBindDn		cn=admin,dc=xiaomi,dc=com;
	LdapBindPwd		123456;
	LdapFilter		"(&(objectClass=posixAccount)(cn=%s))";
	MissoRedirectURL	http://auth.dev.pt.xiaomi.com/v1.0/auth/callback/misso;
	GithubClientID		"0c6eb7247bb4bc7ca16a";
	GithubClientSecret	"7c75c029907af4f398a0e6338fcf9680c1138f64";
	GithubRedirectURL	http://auth.dev.pt.xiaomi.com/v1.0/auth/callback/github;
	GoogleClientID		"781171109477-10tu51e8bs1s677na46oct6hdefpntpu.apps.googleusercontent.com";
	GoogleClientSecret	xpEoBFqkmI3KVN9pHt2VW-eN;
	GoogleRedirectURL	http://auth.dev.pt.xiaomi.com/v1.0/auth/callback/google;
	PluginAlarm		true;
	miNornsUrl		http://norns.dev/api/v1/tagstring/cop.xiaomi/hostinfos;
	miNornsInterval		5;
//	rlLimit			10;
//	rlAccuracy		5;
//	rlGcTimeout		60000;
//	rlGcInterval		1000;
//	metric	{
// 		include	./docs/etc/metric_names;
//	};
	agent	{
		leasekey	"/open-falcon/agent/online/test.falcon";
		leasettl	20;
		interval	5;
		apiAddr		${AGENT_ADDR};
		httpAddr	${AGENT_HTTP_ADDR};
		ifacePrefix	eth,em;
		workerProcesses	3;
		callTimeout	5000;
//		upstream	127.0.0.1:1234,127.0.0.1:1235;
		upstream	${TRANSFER_ADDR};
		emuTplDir	./var/tpl;	
		plugins		"sys,emulator";
	};
};
`)

func TestParse(t *testing.T) {
	ret := Parse(text, "test", 0)
	glog.V(4).Infof("%s", ret)
}
