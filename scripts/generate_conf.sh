#!/bin/bash

declare -A confs
confs=(
    [%%PLUGIN_HTTP%%]=plugin.falcon.srv

    [%%SMS_HTTP%%]=0.0.0.0:8000
    [%%SMS_MI%%]=10.108.163.63
    [%%SMS_OTHER%%]=58.56.119.22:7822
    [%%SMS_ACCOUNT%%]=922010
    [%%SMS_PASS%%]=Qhzh2j
    [%%SMS_EXTNO%%]=106907895678
    [%%SENDER_NAME%%]=falcon
    [%%SENDER_TOKEN%%]=c1646e480

    [%%SMTP_ADDR%%]=10.101.12.50:25
    [%%SMTP_USER%%]=falcon@xiaomi.com
    [%%SMTP_PASS%%]=

    [%%AGENT_PORT%%]=1988
    [%%AGENT_HTTP%%]=0.0.0.0:1988

    [%%ETCD_HTTP%%]=0.0.0.0:2379
    [%%CTRL_PORT%%]=8001
    [%%CTRL_HTTP%%]=0.0.0.0:8001
    [%%AGGREGATOR_HTTP%%]=0.0.0.0:6055
    [%%ALARM_HTTP%%]=0.0.0.0:9912
    [%%NODATA_HTTP%%]=0.0.0.0:6090
    [%%PING_HTTP%%]=0.0.0.0:9981
    [%%TASK_HTTP%%]=0.0.0.0:8002
    [%%HBS_RPC%%]=0.0.0.0:8422
    [%%HBS_HTTP%%]=0.0.0.0:6031
    [%%SENDER_HTTP%%]=0.0.0.0:6066
    [%%LDAP_ADDR%%]=127.0.0.1:389
    [%%LDAP_PASS%%]=123456
    
    [%%GRAPH_HTTP%%]=0.0.0.0:6071
    [%%GRAPH_RPC%%]=0.0.0.0:6070
    [%%GRAPH_GRPC%%]=0.0.0.0:6072
    [%%GRAPH_RRD%%]="/home/work/data/6070"

    [%%JUDGE_HTTP%%]=0.0.0.0:6081
    [%%JUDGE_RPC%%]=0.0.0.0:6080
    [%%JUDGE_GRPC%%]=0.0.0.0:6082

    [%%TRANSFER_HTTP%%]=0.0.0.0:6060
    [%%TRANSFER_RPC%%]=0.0.0.0:8433
    [%%TRANSFER_GRPC%%]=0.0.0.0:8434
    [%%TRANSFER_SOCKET%%]=0.0.0.0:4444

    [%%PORTAL_HTTP%%]=0.0.0.0:9980
    [%%LINKER_HTTP%%]=0.0.0.0:3030
    [%%LINKS_HTTP%%]=l.a.xiaomi.com

    [%%SLINKER_HTTP%%]=0.0.0.0:8080
    [%%SLINKER_URL%%]=127.0.0.1:8080
    [%%SLINKER_DB_NAME%%]=slinker
    [%%SLINKER_DB_USER%%]=root
    [%%SLINKER_DB_PASS%%]=
    [%%SLINKER_DB_HOST%%]=127.0.0.1
    [%%SLINKER_DB_PORT%%]=3306

    [%%DB_GRAPH%%]="falcon:54ca48ac2@tcp(127.0.0.1:3306)/graph"
    [%%DB_ALARM%%]="falcon:54ca48ac2@tcp(127.0.0.1:3306)/event"
    [%%DB_FALCON%%]="falcon:54ca48ac2@tcp(127.0.0.1:3306)/falcon"
    [%%DB_LINKER%%]="falcon:54ca48ac2@tcp(127.0.0.1:3306)/linker"
    [%%DB_HOSTINFO%]="falcon:54ca48ac2@tcp(127.0.0.1:3306)/hostinfo"
    [%%DB_SMS_GATEWAY%]="falcon:54ca48ac2@tcp(127.0.0.1:3306)/sms_gateway"
    [%%REDIS_HTTP%%]=127.0.0.1:6379

    [%%NORNS_HTTP%%]=127.0.0.1:6060
    [%%XPHAROS_HTTP%%]=xpharos.pt.xiaomi.com
)

configurer() {
    for i in "${!confs[@]}"
    do
        search=$i
        replace=${confs[$i]}
        # Note the "" after -i, needed in OS X
        find ./dist/etc -type f -exec sed -i "s!${search}!${replace}!g" {} \;
    done
}
configurer
