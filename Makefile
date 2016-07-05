# Copyright 2016 yubo. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
.PHONY: clean run

all: bin/agent bin/handoff bin/storage bin/falcon

bin/storage: git.go *.go specs/*.go storage/*.go cmd/storage/*.go
	go build -gcflags "-N -l" -o bin/storage cmd/storage/*

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
	./bin/agent -config ./etc/agent.conf -logtostderr -v 3
	#./bin/storage -config ./etc/storage.conf -logtostderr -v 3
