all: ctrl swagger/swagger.json
WORKDIR=$(shell pwd)
VENDOR=${WORKDIR}/../../vendor

.PHONY: db dev doc

ctrl: *.go */*.go
	go build

swagger/swagger.json: routers/router.go
	bee generate docs

doc:
	bee generate docs

#cd ${WORKDIR} &&\
#echo ${WORKDIR} && \ 
# export GOPATH=${VENDOR} && cd ${WORKDIR} && \
#	bee run -gendoc=true -downdoc=true
#	bee run

dev:
	bee run -gendoc=true -downdoc=true

db:
	mysql -uroot -p12341234 falcon < ./scripts/falcon.sql
	mysql -uroot -p12341234 falcon_test < ./scripts/falcon.sql

env:
	@echo export GOPATH=${VENDOR}

