# Copyright 2016 yubo. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
.PHONY: clean run

all: bin/agent bin/handoff bin/backend bin/falcon

bin/backend: git.go *.go specs/*.go backend/*.go cmd/backend/*.go backend/*.c backend/*.h
	go build -gcflags "-N -l" -o bin/backend cmd/backend/*

bin/falcon: git.go *.go */*.go cmd/falcon/*.go
	go build -o bin/falcon cmd/falcon/*

bin/handoff: git.go *.go specs/*.go handoff/*.go cmd/handoff/*.go
	go build -o bin/handoff cmd/handoff/*

bin/agent: git.go *.go specs/*.go agent/*.go agent/plugin/*.go cmd/agent/*.go
	go build -o bin/agent cmd/agent/*

git.go:
	/bin/sh scripts/git.sh

clean:
	rm -rf bin/* git.go

run:
	./bin/backend -config ./etc/backend.conf -logtostderr -v 3

prepare: git.go
	go get ./...
	go get github.com/tcnksm/ghr

test:
	go test ./...

compile:
	make

targz:
	mkdir -p ${OUTPUT_DIR}/dist
	cd ${OUTPUT_DIR}/bin; tar czvf ../dist/falcon_linux_amd64.tar.gz ./*

shasums:
	cd ${OUTPUT_DIR}/dist; shasum * > ./SHASUMS

release:
	ghr --delete --prerelease -u yubo -r falcon pre-release ${OUTPUT_DIR}/dist
#./bin/agent -config ./etc/agent.conf -logtostderr -v 3
#./bin/backend -config ./etc/backend.conf -logtostderr -v 3

