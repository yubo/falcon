#!/bin/sh
VENDOR=$PWD/cmd/vendor
DIR=$PWD/cmd/$1

if [ ! -d $DIR ]; then
	echo "$DIR dir not found"
	exit 1
fi

if [ ! -d $VENDOR ]; then
	echo "cmd/vendor not found"
	exit 1
fi


cd $DIR
mkdir vendor
govendor add +external
rm -rf vendor/github.com/open-falcon/falcon
tar cf - vendor | tar xf - -C $VENDOR/..
rm -rf vendor
