# Copyright 2016 yubo. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
.PHONY: clean run

ifeq ($(strip $(OUTPUT_DIR)),)
  OUTPUT_DIR = .
endif

all: $(OUTPUT_DIR)/bin/falcon $(OUTPUT_DIR)/bin/agent

$(OUTPUT_DIR)/bin/falcon: specs/git.go *.go */*.go cmd/falcon/*.go conf/*.go conf/yyparse.go
	go build -o "${OUTPUT_DIR}/bin/falcon" cmd/falcon/*

$(OUTPUT_DIR)/bin/agent: specs/git.go agent/*.go cmd/agent/*.go specs/*.go
	go build -o "${OUTPUT_DIR}/bin/agent" cmd/agent/*

conf/yyparse.go: conf/yyparse.y
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

prepare: specs/git.go conf/yyparse.go
	go get ./...

tools:
	go get -u github.com/tcnksm/ghr

test:
	go test ./...

coverage:
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

