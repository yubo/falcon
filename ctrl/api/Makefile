all: falcon-api swagger/swagger.json
WORKDIR=$(shell pwd)
VENDOR=${WORKDIR}/../../vendor

.PHONY: db dev doc

falcon-api: *.go */*.go
	go build -o $@

swagger/swagger.json: routers/router.go controllers/*.go
	bee generate docs

doc:
	bee generate docs

dev:
	bee run -gendoc=true -downdoc=true

run: api
	./api -logtostderr -v=5 start


db:
	mysql -uroot -p12341234 yubo_falcon < ../scripts/falcon.sql
	mysql -uroot -p12341234 yubo_falcon_test < ../scripts/falcon.sql

