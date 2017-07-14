#!/bin/bash
src=${GOPATH}/src/github.com/open-falcon/falcon/ctrl/api/swagger/swagger.json

cd ${GOPATH}/src/github.com/open-falcon/falcon/ctrl/api && \
mkdir -p dist/html/swagger && \
bee generate docs && \
cp -f $src ${GOPATH}/src/github.com/open-falcon/falcon/Docs/swagger.json
