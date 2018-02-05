.PHONY: clean parse doc deploy start vendor stats update dev gen

EXEC_OUTPUT_PATH=dist/bin
SWAGGER_DIR=dist/var/html/docs/swagger

MODULES=falcon
GOFILES=$(shell find . -name "*.go" -type f -not -path "./cmd*") gitlog.go
DEPENDS=dist $(GOFILES) $(DOCFILES)
#DEPENDS=dist $(GOFILES)

all: $(EXEC_OUTPUT_PATH)/falcon $(EXEC_OUTPUT_PATH)/agent
	for name in `find -maxdepth 2 -name \*.swagger.json`; \
	do \
		cp -f "$$name" $(SWAGGER_DIR)/$$(basename "$$name"); \
	done

$(EXEC_OUTPUT_PATH)/falcon: $(DEPENDS) cmd/falcon/*.go
	export GOPATH=$(PWD)/gopath && go build -o $@ ./cmd/falcon
	#echo > $@

$(EXEC_OUTPUT_PATH)/agent: $(DEPENDS) cmd/agent/*.go
	export GOPATH=$(PWD)/gopath && go build -o $@ ./cmd/agent
	#echo > $@

gitlog.go:
	./scripts/git.sh

dist:
	mkdir -p $(EXEC_OUTPUT_PATH) dist/var \
	dist/etc dist/run dist/log && \
	cp -a docs dist/ && \
	cp -a var/html/ui dist/var/html && \
	cp -a var/html/docs dist/var/html/ && \
	cp -a var/emu_tpl dist/var/ && \
	cp -a docs/samples/falcon/* dist/etc/

clean:
	rm -rf ./dist


install: $(EXEC_OUTPUT_PATH)/falcon
	./scripts/install.sh

deploy: $(DEPENDS)
	cd dist && ../scripts/deploy.sh

start:
	$(EXEC_OUTPUT_PATH)/falcon -config ./etc/falcon.example.conf start 2>&1

reload:
	$(EXEC_OUTPUT_PATH)/falcon -config ./etc/falcon.example.conf reload 2>&1

usr2:
	cat ./falcon.pid | xargs kill -USR2

parse:
	$(EXEC_OUTPUT_PATH)/falcon -config ./etc/falcon.example.conf parse 2>&1

coverage: $(DEPENDS)
	./scripts/test_coverage.sh
	curl -s https://codecov.io/bash | bash

vendor:
	./scripts/vendor.sh

doc:
	./scripts/generate_doc.sh

gen:
	make -f scripts/generate.mk

stats:
	$(EXEC_OUTPUT_PATH)/falcon  -config ./etc/falcon.example.conf stats

update:
	git submodule update --recursive --init
