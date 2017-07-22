#!/bin/sh
mkdir -p /opt/falcon/bin
mkdir -p /opt/falcon/etc
mkdir -p /opt/falcon/html
cp -af ./dist/bin/* /opt/falcon/bin/
cp -af ./dist/etc/* /opt/falcon/etc/
cp -af ./dist/html/* /opt/falcon/html/


