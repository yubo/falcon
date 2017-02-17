all: api/swagger/swagger.json
WORKDIR=$(shell pwd)
VENDOR=${WORKDIR}/../vendor

.PHONY: db dev doc web

doc:
	make -C api doc

web:
	make -C ui build && \
	rm -rf api/static && \
	cp -a ui/dist api/static

dev:
	make -C api dev

db:
	mysql -uroot -p123456 yubo_falcon < ./scripts/falcon.sql
	mysql -uroot -p123456 yubo_falcon_test < ./scripts/falcon.sql

env:
	@echo export GOPATH=${VENDOR}

pkg:
	tar --exclude=./ui/node_modules -czvf /tmp/ui.tar.gz . 
