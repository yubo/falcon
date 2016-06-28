#all: git.go
all: build/storage
#build/falcon

build/storage: cmd/storage.go git.go *.go storage/*.go
	go build -o build/storage cmd/storage.go

build/falcon: cmd/falcon.go git.go *.go
	go build -o build/falcon cmd/falcon.go

git.go:
	/bin/sh scripts/git.sh
