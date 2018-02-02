.PHONY: clean parse doc deploy start vendor stats update dev

EXEC_OUTPUT_PATH=dist/opt/falcon/sbin
SWAGGER_DIR=dist/opt/falcon/html/docs/swagger

MODULES=falcon
SUBMODULES=ctrl agent transfer service alarm
PBFILES=$(shell find . -name "*.proto" -type f -not -path "./cmd*")
DOCFILES=$(SWAGGER_DIR)/ctrl.swagger.json \
	$(SWAGGER_DIR)/transfer.swagger.json \
	$(SWAGGER_DIR)/service.swagger.json \
	$(SWAGGER_DIR)/agent.swagger.json \
	$(SWAGGER_DIR)/alarm.swagger.json
GOFILES=$(shell find . -name "*.go" -type f -not -path "./cmd*") \
	transfer/transfer.pb.go transfer/transfer.pb.gw.go \
	lib/tsdb/tsdb.pb.go \
	service/service.pb.go service/service.pb.gw.go \
	alarm/alarm.pb.go alarm/alarm.pb.gw.go \
	parse/parse.go $(SUBMODULES:%=%/config/parse.go) \
	gitlog.go
#DEPENDS=dist $(GOFILES) $(DOCFILES)
DEPENDS=dist $(GOFILES)

all: $(EXEC_OUTPUT_PATH)/falcon $(EXEC_OUTPUT_PATH)/agent

$(EXEC_OUTPUT_PATH)/falcon: $(DEPENDS) cmd/falcon/*.go
	export GOPATH=$(PWD)/gopath && go build -o $@ ./cmd/falcon
	#echo > $@

$(EXEC_OUTPUT_PATH)/agent: $(DEPENDS) cmd/agent/*.go
	export GOPATH=$(PWD)/gopath && go build -o $@ ./cmd/agent

$(SWAGGER_DIR)/ctrl.swagger.json: docs/ctrl.swagger.json
	cp -f $< $@

$(SWAGGER_DIR)/transfer.swagger.json: transfer/transfer.swagger.json
	cp -f $< $@

$(SWAGGER_DIR)/service.swagger.json: service/service.swagger.json
	cp -f $< $@

$(SWAGGER_DIR)/agent.swagger.json: agent/agent.swagger.json
	cp -f $< $@

$(SWAGGER_DIR)/alarm.swagger.json: alarm/alarm.swagger.json
	cp -f $< $@

$(SWAGGER_DIR)/api_reference.md: docs/api_reference.md
	cp -f $< $@

gitlog.go:
	./scripts/git.sh

dist:
	mkdir -p $(EXEC_OUTPUT_PATH) dist/etc dist/opt/falcon/var \
	dist/opt/falcon/run dist/opt/falcon/log && \
	cp -a docs/html dist/opt/falcon/ && \
	cp -a docs/emu_tpl dist/opt/falcon/var/ && \
	cp -a docs/samples/init.d dist/etc/ && \
	cp -a docs/samples/falcon dist/etc/ && \
	cp -a docs/samples/nginx dist/etc/

clean:
	rm -rf ./dist


install: $(EXEC_OUTPUT_PATH)/falcon
	./scripts/install.sh

deploy: $(DEPENDS)
	cd dist && ../scripts/deploy.sh

start:
	rm -f /tmp/falcon/*.sock; $(EXEC_OUTPUT_PATH)/falcon -config ./docs/samples/falcon/falcon.example.conf start 2>&1

dev:
	rm -f /opt/falcon/*.sock; $(EXEC_OUTPUT_PATH)/falcon -config ./docs/samples/falcon/falcon.dev01.conf start 2>&1

reload:
	$(EXEC_OUTPUT_PATH)/falcon -config ./docs/etc/falcon.example.conf reload 2>&1

usr2:
	cat ./falcon.pid | xargs kill -USR2

parse:
	$(EXEC_OUTPUT_PATH)/falcon -config ./docs/etc/falcon.example.conf parse 2>&1

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
	$(EXEC_OUTPUT_PATH)/falcon stats -config ./docs/samples/falcon/falcon.example.conf 

update:
	git submodule update --recursive --init

#include ./scripts/falcon.mk
