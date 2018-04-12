.PHONY: clean parse doc deploy start vendor stats update dev gen tools

EXEC_OUTPUT_PATH=dist/bin
SWAGGER_DIR=dist/var/html/docs/swagger

MODULES=falcon
GOFILES=$(shell find . -name "*.go" -type f -not -path "./cmd*") gitlog.go
DEPENDS=dist $(GOFILES) $(DOCFILES)
#DEPENDS=dist $(GOFILES)

all: $(EXEC_OUTPUT_PATH)/falcon
	@for name in `find . -maxdepth 2 -name \*.swagger.json`; \
	do \
		cp -f "$$name" $(SWAGGER_DIR)/$$(basename "$$name"); \
	done

$(EXEC_OUTPUT_PATH)/falcon: $(DEPENDS) cmd/falcon/*.go
	go build -o $@ ./cmd/falcon

$(EXEC_OUTPUT_PATH)/agent: $(DEPENDS) cmd/agent/*.go
	@echo > $@
	@export GOPATH=$(PWD)/gopath && go build -o $@ ./cmd/agent

gitlog.go:
	./scripts/git.sh

dist:
	@mkdir -p $(EXEC_OUTPUT_PATH) dist/scripts \
	dist/etc dist/run dist/log dist/var/html/docs && \
	cp -a scripts/db_schema	dist/scripts && \
	cp -a docs		dist/ && \
	cp -a docs/img		dist/var/html/docs/ && \
	cp -a var/html/ui	dist/var/html && \
	cp -a var/html/docs	dist/var/html/ && \
	cp -a var/emu_tpl	dist/var/ && \
	cp -a docs/samples/falcon/* dist/etc/

clean:
	rm -rf ./dist


install: $(EXEC_OUTPUT_PATH)/falcon
	./scripts/install.sh

deploy: $(DEPENDS)
	cd dist && ../scripts/deploy.sh

start:
	$(EXEC_OUTPUT_PATH)/falcon -logtostderr -v 6 -f ./etc/values.example.yaml -config ./etc/falcon.example.yaml start 2>&1

reload:
	$(EXEC_OUTPUT_PATH)/falcon -config ./etc/falcon.example.yaml reload 2>&1

usr2:
	cat ./falcon.pid | xargs kill -USR2

parse:
	$(EXEC_OUTPUT_PATH)/falcon -f ./etc/values.example.yaml -config ./etc/falcon.example.yaml parse 2>&1

coverage: $(DEPENDS)
	export GOPATH=$(PWD)/gopath && ./scripts/test_coverage.sh
	curl -s https://codecov.io/bash | bash

vendor:
	./scripts/vendor.sh

doc:
	./scripts/generate_doc.sh

gen:
	make -f scripts/generate.mk

stats:
	$(EXEC_OUTPUT_PATH)/falcon -f ./etc/values.example.yaml -config ./etc/falcon.example.conf stats

update:
	git submodule update --recursive --init

tools:
	go get -u github.com/tcnksm/ghr

shasums:
	shasum dist/bin/* > ./dist/SHASUMS

release:
	tar czvf falcon_linux_amd64.tar.gz dist && ghr --delete --prerelease -u yubo -r falcon pre-release ./falcon_linux_amd64.tar.gz


