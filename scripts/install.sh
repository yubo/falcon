#!/bin/sh
mkdir -p /opt/falcon/{bin,etc,swagger}
cp -af ./bin/* /opt/falcon/bin/
cp -f ./src/falcon/etc/falcon.example.conf /opt/falcon/etc/falcon.conf
cp -f ./src/falcon/etc/metric_names /opt/falcon/etc/
cp -f ./Docs/ctrl/swagger.json /opt/falcon/swagger/swagger.json


