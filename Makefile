.PHONY: clean parse doc deploy run vendor

all: dist/bin/falcon

MODULES=falcon
SUBMODULES=ctrl agent transfer backend
PBFILES=$(shell find . -name "*.proto" -type f -not -path "./cmd*")
DOCFILES=dist/etc/falcon.example.conf \
	dist/etc/metric_names \
	dist/html/ctrl.swagger.json \
	dist/html/transfer.swagger.json \
	dist/html/backend.swagger.json \
	dist/html/api_reference.md
GOFILES=$(shell find . -name "*.go" -type f -not -path "./cmd*") \
	transfer/transfer.pb.go transfer/transfer.pb.gw.go \
	backend/backend.pb.go backend/backend.pb.gw.go \
	falcon.pb.go \
	parse/parse.go $(SUBMODULES:%=%/parse/parse.go) \
	gitlog.go
DEPENDS=dist $(GOFILES) $(DOCFILES)
TARGETS=dist dist/bin/falcon $(DOCFILES)

dist/bin/falcon: $(DEPENDS) cmd/falcon/*.go
	go build -o $@ ./cmd/falcon

dist/html/ctrl.swagger.json: docs/ctrl.swagger.json
	cp -f $< $@

dist/html/transfer.swagger.json: transfer/transfer.swagger.json
	cp -f $< $@

dist/html/backend.swagger.json: backend/backend.swagger.json
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

run:
	./dist/bin/falcon -config ./docs/etc/agent.example.conf -logtostderr -v 4 start 2>&1

reload:
	./dist/bin/falcon -config ./docs/etc/falcon.conf -logtostderr -v 4 reload 2>&1

parse: $(TARGETS)
	./dist/bin/falcon -config ./docs/etc/falcon.example.conf -logtostderr -v 4 parse 2>&1

coverage: $(DEPENDS)
	./scripts/test_coverage.sh
	curl -s https://codecov.io/bash | bash

pb: $(PBFILES)
	find -name \*.pb.go

vendor:
	./scripts/vendor.sh

doc:
	./scripts/generate_doc.sh

include ./scripts/falcon.mk
