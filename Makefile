.PHONY: clean parse doc deploy start vendor stats update dev

all: dist/sbin/falcon dist/sbin/agent

MODULES=falcon
SUBMODULES=ctrl agent transfer service alarm
PBFILES=$(shell find . -name "*.proto" -type f -not -path "./cmd*")
DOCFILES=dist/html/doc/ctrl.swagger.json \
	dist/html/doc/transfer.swagger.json \
	dist/html/doc/service.swagger.json \
	dist/html/doc/agent.swagger.json \
	dist/html/doc/alarm.swagger.json \
	dist/html/doc/api_reference.md
GOFILES=$(shell find . -name "*.go" -type f -not -path "./cmd*") \
	transfer/transfer.pb.go transfer/transfer.pb.gw.go \
	lib/tsdb/tsdb.pb.go \
	service/service.pb.go service/service.pb.gw.go \
	alarm/alarm.pb.go alarm/alarm.pb.gw.go \
	parse/parse.go $(SUBMODULES:%=%/config/parse.go) \
	gitlog.go
DEPENDS=dist $(GOFILES) $(DOCFILES)

dist/sbin/falcon: $(DEPENDS) cmd/falcon/*.go
	export GOPATH=$(PWD)/gopath && go build -o $@ ./cmd/falcon

dist/sbin/agent: $(DEPENDS) cmd/agent/*.go
	export GOPATH=$(PWD)/gopath && go build -o $@ ./cmd/agent

dist/html/doc/ctrl.swagger.json: docs/ctrl.swagger.json
	cp -f $< $@

dist/html/doc/transfer.swagger.json: transfer/transfer.swagger.json
	cp -f $< $@

dist/html/doc/service.swagger.json: service/service.swagger.json
	cp -f $< $@

dist/html/doc/agent.swagger.json: agent/agent.swagger.json
	cp -f $< $@

dist/html/doc/alarm.swagger.json: alarm/alarm.swagger.json
	cp -f $< $@

dist/html/doc/api_reference.md: docs/api_reference.md
	cp -f $< $@

dist/etc/doc/falcon.example.conf: docs/samples/falcon/falcon.example.conf
	cp -f $< $@

gitlog.go:
	./scripts/git.sh

dist:
	mkdir -p dist/sbin dist/etc && \
	cp -a docs/html dist/ && \
	cp -a docs/samples/init.d dist/etc/ && \
	cp -a docs/samples/default dist/etc/ && \
	cp -a docs/samples/falcon dist/etc/ && \
	cp -a docs/samples/nginx dist/etc/

clean:
	rm -rf ./dist


install: dist/sbin/falcon
	./scripts/install.sh

deploy: $(DEPENDS)
	cd dist && ../scripts/deploy.sh

start:
	rm -f /tmp/falcon/*.sock; ./dist/sbin/falcon start -config ./docs/samples/falcon/falcon.example.conf 2>&1

dev:
	rm -f /opt/falcon/*.sock; ./dist/sbin/falcon start -config ./docs/samples/falcon/falcon.dev01.conf 2>&1

reload:
	./dist/sbin/falcon reload -config ./docs/etc/falcon.example.conf 2>&1

usr2:
	cat ./falcon.pid | xargs kill -USR2

parse: dist/sbin/falcon
	./dist/sbin/falcon parse -config ./docs/etc/falcon.example.conf 2>&1

coverage: $(DEPENDS)
	./scripts/test_coverage.sh
	curl -s https://codecov.io/bash | bash

pb: $(PBFILES)
	find -name \*.pb.go

vendor:
	./scripts/vendor.sh

doc:
	./scripts/generate_doc.sh

stats:
	./dist/sbin/falcon stats -config ./docs/etc/falcon.example.conf 

update:
	git submodule update --recursive --init

include ./scripts/falcon.mk
