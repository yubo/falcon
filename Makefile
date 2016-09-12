# Copyright 2016 yubo. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
.PHONY: clean run

ifeq ($(strip $(OUTPUT_DIR)),)
  OUTPUT_DIR = .
endif

all: $(OUTPUT_DIR)/bin/agent $(OUTPUT_DIR)/bin/handoff $(OUTPUT_DIR)/bin/backend $(OUTPUT_DIR)/bin/falcon


$(OUTPUT_DIR)/bin/backend: git.go *.go specs/*.go backend/*.go cmd/backend/*.go backend/*.c backend/*.h
	go build -gcflags "-N -l" -o "${OUTPUT_DIR}/bin/backend" cmd/backend/*

$(OUTPUT_DIR)/bin/falcon: git.go *.go */*.go cmd/falcon/*.go
	go build -o "${OUTPUT_DIR}/bin/falcon" cmd/falcon/*

$(OUTPUT_DIR)/bin/handoff: git.go *.go specs/*.go handoff/*.go cmd/handoff/*.go
	go build -o "${OUTPUT_DIR}/bin/handoff" cmd/handoff/*

$(OUTPUT_DIR)/bin/agent: git.go *.go specs/*.go agent/*.go agent/plugin/*.go cmd/agent/*.go
	go build -o "${OUTPUT_DIR}/bin/agent" cmd/agent/*

git.go:
	/bin/sh scripts/git.sh

clean:
	rm -rf bin/* git.go

run:
	./bin/backend -config ./etc/backend.conf -logtostderr -v 3

prepare: git.go
	go get ./...

tools:
	go get github.com/tcnksm/ghr

test:
	go test ./...

compile:
	mkdir -p "${OUTPUT_DIR}/bin"
	make

targz:
	mkdir -p ${OUTPUT_DIR}/dist
	cd ${OUTPUT_DIR}/bin; tar czvf ../dist/falcon_single_linux_amd64.tar.gz ./falcon ; tar czvf ../dist/falcon_multi_linux_amd64.tar.gz ./backend ./handoff ./agent

shasums:
	cd ${OUTPUT_DIR}/dist; shasum * > ./SHASUMS

release:
	ghr --delete --prerelease -u yubo -r falcon pre-release ${OUTPUT_DIR}/dist
#./bin/agent -config ./etc/agent.conf -logtostderr -v 3
#./bin/backend -config ./etc/backend.conf -logtostderr -v 3

