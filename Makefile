# Copyright 2016 yubo. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
.PHONY: clean run
MODULES=agent falcon
VERSION=$(shell cat VERSION)
TARGETS=$(MODULES:%=dist/%)
ORG_PATH=github.com/yubo/falcon
WORKDIR=${shell pwd}

ifeq ($(strip $(OUTPUT_DIR)),)
  OUTPUT_DIR = .
endif

all: $(TARGETS)

$(TARGETS): conf/yyparse.go git.go
	export GOPATH=$(WORKDIR)/gopath:$$GOPATH &&\
	go build -o ${OUTPUT_DIR}/$@ github.com/yubo/falcon/$(subst dist/,cmd/,$@)

conf/yyparse.go: conf/yyparse.y
	export GOPATH=${WORKDIR}/gopath:$$GOPATH &&\
	go tool yacc -o conf/yyparse.go conf/yyparse.y

git.go: scripts/git.sh
	/bin/sh scripts/git.sh

clean:
	rm -rf dist/* git.go conf/yyparse.go y.output

run:
	./dist/falcon -config ./etc/falcon.conf -logtostderr -v 4 start 2>&1

parse:
	./dist/falcon -config ./etc/falcon.conf -logtostderr -v 3 parse 2>&1
#./bin/agent -http=false -v 4 start 2>&1

tools:
	go get -u github.com/tcnksm/ghr &&\
	go get -u github.com/kardianos/govendor

test:
	export GOPATH=${WORKDIR}/gopath:$$GOPATH &&\
	go test ./...

coverage: conf/yyparse.go git.go
	export GOPATH=${WORKDIR}/gopath:$$GOPATH &&\
	./scripts/test_coverage.sh
	curl -s https://codecov.io/bash | bash

compile:
	mkdir -p "${OUTPUT_DIR}/dist"
	make

targz:
	mkdir -p ${OUTPUT_DIR}/dist; rm -rf ${OUTPUT_DIR}/dist/*
	cd ${OUTPUT_DIR}/bin; tar czvf ../dist/falcon_Linux_amd64.tar.gz ./falcon; tar czvf ../dist/agent_Linux_amd64.tar.gz ./agent

shasums:
	cd ${OUTPUT_DIR}/dist; shasum * > ./SHASUMS

release:
	ghr --delete --prerelease -u yubo -r falcon pre-release ${OUTPUT_DIR}/dist

