all: api/swagger/swagger.json  

WORKDIR=${GOPATH}/src/github.com/yubo/falcon/ctrl/api

.PHONY: db dev doc web

doc:
	cd ${WORKDIR} && \
	bee generate docs

dev:     
	cd ${WORKDIR} && \
	bee run -gendoc=true -downdoc=true
