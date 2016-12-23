all: ctrl swagger/swagger.json
WORKDIR=$(shell readlink -f .)
VENDOR=$(shell readlink -f ${WORKDIR}/../vendor)

.PHONY: db dev

ctrl: *.go */*.go
	go build

swagger/swagger.json: routers/router.go
	bee generate docs

dev:
	export GOPATH=${VENDOR} && cd ${WORKDIR} &&\
	bee run -gendoc=true -downdoc=true

db:
	mysql -uroot -p12341234 < falcon.sql

env:
	@echo export GOPATH=${VENDOR}

