/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

/*
 * go test -bench ./timer_test.go
 * BenchmarkTimerSysNow-4          50000000                 24.8 ns/op
 * BenchmarkTimerNow-4             1000000000               2.88 ns/op
 */
package service

import (
	"testing"
	"time"
)

var (
	testTimer *timerModule
)

func init() {
	testTimer = &timerModule{}
	testTimer.prestart(nil)
	testTimer.start(nil)
}

func BenchmarkTimerSysNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Now().Unix()
	}
}

func BenchmarkTimerNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		now()
	}
}
