/*
 * Copyright 2017 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ratelimits

import (
	"fmt"
	"testing"
	"time"
)

func test(qps, accuracy, ts uint32) {
	stats := make(map[bool]int)
	x, err := New(qps, accuracy)
	if err != nil {
		fmt.Println(err)
		return
	}
	done := make(chan struct{}, 1)
	go func(done chan struct{}, x *RateLimits) {
		for {
			select {
			case <-done:
				fmt.Printf("true:%fHz false:%fHz\n",
					float64(stats[true])/float64(ts),
					float64(stats[false])/float64(ts))
				return
			default:
				ok := x.Update("test")
				stats[ok]++
				time.Sleep(time.Nanosecond)
			}
		}
	}(done, x)
	time.Sleep(time.Second * time.Duration(ts))
	done <- struct{}{}
}

func testGc() {
	x, _ := New(10, 1)

	err := x.GcStart(time.Second, time.Millisecond*10)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 200; i++ {
		x.Update(fmt.Sprintf("i_%d", i))
		time.Sleep(time.Millisecond * 10)
	}

	for i := 0; i < 10; i++ {
		fmt.Printf("rl.Len() %d\n", x.Len())
		time.Sleep(time.Millisecond * 100)
	}

	x.GcStop()
	time.Sleep(100 * time.Millisecond) // let cancelation propagate
}

func TestAll(t *testing.T) {
	test(10, 1, 2)
	test(10, 2, 2)
	test(10, 4, 2)
	test(10, 8, 2)
	test(100, 1, 2)
	test(100, 2, 2)
	test(100, 4, 2)
	test(100, 8, 2)
	test(10000, 1, 1)
	test(10000, 2, 1)
	test(10000, 4, 1)
	test(10000, 8, 1)
	testGc()
}
