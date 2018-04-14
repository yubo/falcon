.PHONY: test

test: parse.go
	go test -v -args -v 3 -logtostderr true

%.go: %.y
	goyacc -o $@ $<
