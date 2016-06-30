# 
# Copyright 2016 Xiaomi Corporation. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
# 
# Authors:    Yu Bo <yubo@xiaomi.com>
# 
#all: git.go
all: build/storage
#build/falcon

build/storage: cmd/storage.go git.go *.go storage/*.go
	go build -o build/storage cmd/storage.go

build/falcon: cmd/falcon.go git.go *.go
	go build -o build/falcon cmd/falcon.go

git.go:
	/bin/sh scripts/git.sh
