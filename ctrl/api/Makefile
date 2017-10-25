all: api/swagger/swagger.json  
	WORKDIR=$(shell pwd)

.PHONY: db dev doc web

doc:     
	bee generate docs

dev:     
	bee run
