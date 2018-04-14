/*
 * Copyright 2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package emulator

import (
	"fmt"
	"math"
	"testing"
)

func TestEmu(t *testing.T) {
	tpl, err := readTpl("./sample.tpl")
	if err != nil {
		t.Fatal(err)
	}
	// x1
	for i := 0; i < tpl.n; i++ {
		v := tpl.emuValue(int64(i * tpl.interval))
		if math.Abs(v-tpl.v[i]) > 0.1 {
			t.Fatalf("tpl.emuvalue(%d * %d) got %f, but want %f",
				i, tpl.interval, v, tpl.v[i])
		} else {
			fmt.Printf("%d, %f\n", i, v)
		}
	}
	// x2
	for i := 0; i < tpl.n; i++ {
		v := tpl.emuValue(int64(i * tpl.interval * 2))
		fmt.Printf("%d, %f\n", i, v)
	}
	// x0.5
	for i := 0; i < tpl.n; i++ {
		v := tpl.emuValue(int64(i * tpl.interval / 2))
		fmt.Printf("%d, %f\n", i, v)
	}

}
