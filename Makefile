.PHONY: clean parse doc deploy run

MODULES=falcon
SUBMODULES=ctrl agent transfer backend
TARGETS=$(MODULES:%=dist/bin/%)

DOCFILES=dist/html/swagger/swagger.json dist/etc/falcon.example.conf
YYFILES=parse/parse.go $(SUBMODULES:%=%/parse/parse.go)
DEPENDS=$(DOCFILES) $(YYFILES) dist gitlog.go



all: $(TARGETS) dist dist/etc/falcon.example.conf dist/etc/falcon.example.conf

#VENDOR=$(shell readlink -f ./gopath)
VENDOR=$(GOPATH)


%.go: %.y
	goyacc -o $@ $<


$(TARGETS): $(YYFILES) dist gitlog.go
	export GOPATH=${VENDOR} && \
	go build -o $@ ./cmd/$(subst dist/bin/,,$@)


dist/html/swagger/swagger.json:
	cp docs/swagger.json $@

gitlog.go:
	./scripts/git.sh

doc:
	./scripts/generate_doc.sh

dist:
	mkdir -p dist/bin dist/etc dist/html/swagger

dist/etc/falcon.example.conf: docs/etc/falcon.example.conf dist
	cp docs/etc/falcon.example.conf $@

dist/etc/metric_names: docs/etc/metric_names dist
	cp docs/etc/metric_names dist/etc/metric_names

clean:
	rm -rf ./dist gitlog.go $(YYFILES)

install: $(DEPENDS)
	./scripts/install.sh

deploy: $(DEPENDS)
	cd dist && ../scripts/deploy.sh

run:
	./dist/bin/falcon -config ./etc/agent.example.conf -logtostderr -v 4 start 2>&1

reload:
	./dist/bin/falcon -config ./etc/falcon.conf -logtostderr -v 4 reload 2>&1

parse: $(TARGETS)
	./dist/bin/falcon -config ./etc/falcon.conf -logtostderr -v 4 parse 2>&1

coverage: $(DEPENDS)
	export GOPATH=${WORKDIR}/gopath:$$GOPATH &&\
	./scripts/test_coverage.sh
	curl -s https://codecov.io/bash | bash
