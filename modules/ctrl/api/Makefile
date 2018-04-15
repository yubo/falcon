all: api/swagger/swagger.json  

.PHONY: dev doc

doc:
	bee generate docs

dev:     
	bee run -gendoc=true -downdoc=true
