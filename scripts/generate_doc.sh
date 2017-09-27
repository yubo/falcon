#!/bin/bash

export GOPATH=${PWD}/gopath
src=${GOPATH}/src/github.com/yubo/falcon/ctrl/api/swagger/swagger.json
dst=${GOPATH}/src/github.com/yubo/falcon/docs/ctrl.swagger.json

cd ${GOPATH}/src/github.com/yubo/falcon/ctrl/api && \
mkdir -p $(dirname $dst) && \
bee generate docs && \
cp -f $src $dst
cd -

## api_reference.md
tempdir=tt
rm -rf ${tempdir}
mkdir -p ${tempdir}/{falcon,transfer,backend}
cd  ${tempdir}
cp -f ../falcon.proto falcon/falcon.proto
cp -f ../transfer/transfer.proto transfer/transfer.proto
cp -f ../backend/backend.proto backend/backend.proto
protodoc --directories="falcon=service_message,transfer=service_message,backend=service_message" \
	--title="falcon API Reference" \
	--output="api_reference.md" \
	--disclaimer="This is a generated documentation. Please read the proto files for more."
sed -i 's/falcon\/falcon\.proto/falcon.proto/g' api_reference.md
cd -
cp -f ${tempdir}/api_reference.md ${GOPATH}/src/github.com/yubo/falcon/docs/api_reference.md
rm -rf ${tempdir}
