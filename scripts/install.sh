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

mkdir -p /opt/falcon
tar cf -  -C ./dist/ . | tar xvf - -C /opt/falcon

if [ -f /etc/lsb-release ]; then
	ln -s /opt/falcon/docs/samples/init.d/falcon.lsb.init /etc/init.d/falcon
	#force_install_file ./docs/samples/init.d/falcon.lsb.init /etc/init.d/falcon 0755
	update-rc.d falcon defaults 80
elif [ -f /etc/redhat-release ]; then
	force_install_file ./docs/samples/init.d/falcon.redhat.init /etc/init.d/falcon 0755
	chkconfig --add falcon 
	chkconfig falcon on
fi


