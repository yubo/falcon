all: api/swagger/swagger.json

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
	mysql -uroot -p12341234 falcon < ./scripts/falcon.sql
	mysql -uroot -p12341234 falcon_test < ./scripts/falcon.sql

pkg:
	tar --exclude=./ui/node_modules -czvf /tmp/ui.tar.gz . 
