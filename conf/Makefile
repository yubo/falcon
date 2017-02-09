# Copyright 2016 yubo. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
.PHONY: clean test

all: yyparse.go

yyparse.go: yyparse.y
	go tool yacc -o yyparse.go yyparse.y

test: *.go
	go test

clean:
	rm -f yyparse.go *.output

