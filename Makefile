# Copyright 2016 yubo. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
.PHONY: clean run
MODULES=agent ctrl falcon
VENDOR=$(shell pwd)/vendor
VERSION=$(shell cat VERSION)
TARGETS=$(MODULES:%=bin/%)
ORG_PATH=github.com/yubo/falcon
WORKDIR=${VENDOR}/src/${ORG_PATH}

ifeq ($(strip $(OUTPUT_DIR)),)
  OUTPUT_DIR = .
endif

all: $(TARGETS)

$(TARGETS): conf/yyparse.go specs/git.go
	export GOPATH=${VENDOR} && cd ${WORKDIR} &&\
	go build -o ${OUTPUT_DIR}/$@ ./$(subst bin/,cmd/,$@)

conf/yyparse.go: conf/yyparse.y
	export GOPATH=${VENDOR} && cd ${WORKDIR} &&\
	go tool yacc -o conf/yyparse.go conf/yyparse.y

specs/git.go: scripts/git.sh
	/bin/sh scripts/git.sh

clean:
	rm -rf bin/* specs/git.go conf/yyparse.go conf/*.output

run:
	./bin/falcon -config ./etc/falcon.conf -logtostderr -v 4 start 2>&1

parse:
	./bin/falcon -config ./etc/falcon.conf -logtostderr -v 3 parse 2>&1
#./bin/agent -http=false -v 4 start 2>&1

tools:
	go get -u github.com/tcnksm/ghr

test:
	export GOPATH=${VENDOR} && cd ${WORKDIR} &&\
	go test ./...

coverage: conf/yyparse.go specs/git.go
	export GOPATH=${VENDOR} && cd ${WORKDIR} &&\
	./scripts/test_coverage.sh
	curl -s https://codecov.io/bash | bash

compile:
	mkdir -p "${OUTPUT_DIR}/bin"
	make

targz:
	mkdir -p ${OUTPUT_DIR}/dist; rm -rf ${OUTPUT_DIR}/dist/*
	cd ${OUTPUT_DIR}/bin; tar czvf ../dist/falcon_Linux_amd64.tar.gz ./falcon; tar czvf ../dist/agent_Linux_amd64.tar.gz ./agent

shasums:
	cd ${OUTPUT_DIR}/dist; shasum * > ./SHASUMS

release:
	ghr --delete --prerelease -u yubo -r falcon pre-release ${OUTPUT_DIR}/dist

