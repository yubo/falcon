all: api/swagger/swagger.json  
	WORKDIR=$(shell pwd)

.PHONY: db dev doc web

doc:     
	cd ${GOPATH}/src/github.com/open-falcon/falcon/ctrl/api && \
		echo `pwd` && \
	bee generate docs

dev:     
	cd ${GOPATH}/src/github.com/open-falcon/falcon/ctrl/api && \
		bee run -gendoc=true -downdoc=true
