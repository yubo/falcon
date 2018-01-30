.PHONY: clean parse doc deploy start vendor stats

all: dist/bin/falcon

MODULES=falcon
SUBMODULES=ctrl agent transfer service alarm
PBFILES=$(shell find . -name "*.proto" -type f -not -path "./cmd*")
DOCFILES=dist/etc/falcon.example.conf \
	dist/etc/metric_names \
	dist/html/ctrl.swagger.json \
	dist/html/transfer.swagger.json \
	dist/html/service.swagger.json \
	dist/html/agent.swagger.json \
	dist/html/alarm.swagger.json \
	dist/html/api_reference.md
GOFILES=$(shell find . -name "*.go" -type f -not -path "./cmd*") \
	transfer/transfer.pb.go transfer/transfer.pb.gw.go \
	lib/tsdb/tsdb.pb.go \
	service/service.pb.go service/service.pb.gw.go \
	alarm/alarm.pb.go alarm/alarm.pb.gw.go \
	parse/parse.go $(SUBMODULES:%=%/config/parse.go) \
	gitlog.go
DEPENDS=dist $(GOFILES) $(DOCFILES)
TARGETS=dist dist/bin/falcon $(DOCFILES)

dist/bin/falcon: $(DEPENDS) cmd/falcon/*.go
	go build -o $@ ./cmd/falcon

dist/html/ctrl.swagger.json: docs/ctrl.swagger.json
	cp -f $< $@

dist/html/transfer.swagger.json: transfer/transfer.swagger.json
	cp -f $< $@

dist/html/service.swagger.json: service/service.swagger.json
	cp -f $< $@

dist/html/agent.swagger.json: agent/agent.swagger.json
	cp -f $< $@

dist/html/alarm.swagger.json: alarm/alarm.swagger.json
	cp -f $< $@

dist/html/api_reference.md: docs/api_reference.md
	cp -f $< $@

gitlog.go:
	./scripts/git.sh

dist:
	mkdir -p dist/bin dist/etc dist/html

dist/etc/falcon.example.conf: docs/etc/falcon.example.conf dist
	cp docs/etc/falcon.example.conf $@

dist/etc/metric_names: dist
	cp docs/etc/metric_names dist/etc/metric_names

clean:
	rm -rf ./dist


install: $(TARGETS)
	./scripts/install.sh

deploy: $(DEPENDS)
	cd dist && ../scripts/deploy.sh

start:
	rm -f ./var/*.rpc; ./dist/bin/falcon start -config ./docs/etc/falcon.example.conf 2>&1

reload:
	./dist/bin/falcon reload -config ./docs/etc/falcon.example.conf 2>&1

usr2:
	cat ./falcon.pid | xargs kill -USR2

parse: $(TARGETS)
	./dist/bin/falcon parse -config ./docs/etc/falcon.example.conf 2>&1

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
	./dist/bin/falcon stats -config ./docs/etc/falcon.example.conf 

include ./scripts/falcon.mk
