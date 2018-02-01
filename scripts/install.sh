#!/bin/bash
export LC_ALL=en_US.UTF-8

if [[ $EUID -ne 0 ]]; then
	echo "This script must be run as root" 
	exit 1
fi


install_file()
{
	src=$1
	des=$2
	mod=$3
	if [ ! -f $des ]; then
		cp -f $src $des
		chmod $mod $des
	fi
} 

force_install_file()
{
	src=$1
	des=$2
	mod=$3
	cp -f $src $des
	chmod $mod $des
}

tar cf -  -C ./dist/ --exclude etc/init.d . | tar xvf - -C /

if [ -f /etc/lsb-release ]; then
	force_install_file ./dist/etc/init.d/falcon.lsb.init /etc/init.d/falcon 0755
	update-rc.d falcon defaults 80
elif [ -f /etc/redhat-release ]; then
	force_install_file ./dist/etc/init.d/falcon.redhat.init /etc/init.d/falcon 0755
	chkconfig --add falcon 
	chkconfig falcon on
fi


