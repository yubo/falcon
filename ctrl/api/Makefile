all: falcon-api swagger/swagger.json
WORKDIR=$(shell pwd)
VENDOR=${WORKDIR}/../../vendor

.PHONY: db dev doc

ctrl: *.go */*.go
	go build

falcon-api: *.go */*.go
	go build -o $@

swagger/swagger.json: routers/router.go controllers/*.go
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

run: api
	./api -logtostderr -v=5 start


db:
	mysql -uroot -p123456 yubo_falcon < ../scripts/falcon.sql
	mysql -uroot -p123456 yubo_falcon_test < ../scripts/falcon.sql

env:
	@echo export GOPATH=${VENDOR}

