// Copyright 2016 yubo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipc

import (
	"fmt"
	"testing"
	"time"
)

func TestMsgSnd(t *testing.T) {
	mode := 0666
	key := 0x1234

	go func() {
		msqid := MsgGet(key, mode)
		msgtype := 1
		for i := 0; i < 10; i++ {
			MsgSnd(msqid, fmt.Sprintf("%d", i), msgtype, 0)
			time.Sleep(time.Millisecond * 100)
		}
		MsgSnd(msqid, ".", msgtype, 0)
	}()

	{
		msqid := MsgGet(key, mode)
		msgtype := 1
		for {
			msg, err := MsgRcv(msqid, msgtype, 0)
			if err != nil {
				t.Error(err)
			} else {
				fmt.Println(msg)
				if msg == "." {
					return
				}
			}
		}
	}

}
