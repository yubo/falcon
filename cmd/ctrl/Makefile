all: api/swagger/swagger.json  
	WORKDIR=$(shell pwd)

.PHONY: db dev doc web

doc:     
	cd ${GOPATH}/src/github.com/yubo/falcon/ctrl/api && \
		echo `pwd` && \
	bee generate docs

dev:     
	bee run -gendoc=true -downdoc=true
