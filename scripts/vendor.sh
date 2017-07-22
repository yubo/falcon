#!/bin/bash
OD=${PWD}
FD=${PWD}/cmd/falcon
VD=${PWD}/cmd/vendor

if ! [[ "$0" =~ "scripts/vendor.sh" ]]; then
	echo "must be run from repository root"
	exit 255
fi


cd ${FD}

rm -rf ${VD}
mkdir ${VD}
govendor add +external
cd ${VD}/github.com/yubo
rm -rf falcon
ln -s ../../../.. falcon
cd ${OD}

