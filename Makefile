.PHONY: clean run vendor
MODULES=falcon
VERSION=$(shell cat VERSION)
TARGETS=$(MODULES:%=bin/%)
WORKDIR=${shell pwd}

ifeq ($(strip $(OUTPUT_DIR)),)
  OUTPUT_DIR = .
endif

all: $(TARGETS)

$(TARGETS): yyparse.go git.go *.go */*.go */*/*.go
	export GOPATH=$(WORKDIR)/gopath:$$GOPATH &&\
	go build -o ${OUTPUT_DIR}/$@ github.com/yubo/falcon/$(subst bin/,cmd/,$@)

git.go: scripts/git.sh
	/bin/sh scripts/git.sh

clean:
	rm -rf bin/* git.go conf/yyparse.go y.output

run:
	rm -rf var/*.sock; ./bin/falcon -config ./etc/falcon.conf -logtostderr -v 4 start 2>&1

reload:
	./bin/falcon -config ./etc/falcon.conf -logtostderr -v 4 reload 2>&1

parse:
	./bin/falcon -config ./etc/falcon.conf -logtostderr -v 3 parse 2>&1
#./bin/agent -http=false -v 4 start 2>&1

tools:
	go get -u github.com/tcnksm/ghr &&\
	go get -u github.com/kardianos/govendor

yyparse.go: yyparse.y
	go tool yacc -o yyparse.go yyparse.y

coverage: yyparse.go git.go
	export GOPATH=${WORKDIR}/gopath:$$GOPATH &&\
	./scripts/test_coverage.sh
	curl -s https://codecov.io/bash | bash

test: yyparse.go
	go test -logtostderr

build:
	mkdir -p "${OUTPUT_DIR}/bin"
	make

targz:
	mkdir -p ${OUTPUT_DIR}/bin; rm -rf ${OUTPUT_DIR}/bin/*
	cd ${OUTPUT_DIR}/bin; tar czvf ../bin/falcon_Linux_amd64.tar.gz ./falcon; tar czvf ../bin/agent_Linux_amd64.tar.gz ./agent

shasums:
	cd ${OUTPUT_DIR}/bin; shasum * > ./SHASUMS

release:
	ghr --delete --prerelease -u yubo -r falcon pre-release ${OUTPUT_DIR}/bin

vendor:
	cd cmd && govendor add +external

