/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ipc

/*
#include "msg.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// MSGMAX Maximum size of a message text, in bytes (default value: 8192 bytes). /proc/sys/kernel/msgmax
// MSGMNB Maximum number of bytes that can be held in a message queue (default value: 16384 bytes). /proc/sys/kernel/msgmnb

func MsgGet(key, mode int) int {
	return int(C.msgGet(C.key_t(key), C.int(mode|C.IPC_CREAT)))
}

func MsgSnd(msqid int, msg string, msgtyp, msgflg int) int {
	msgp := C.CString(msg)
	defer C.free(unsafe.Pointer(msgp))

	return int(C.msgSnd(C.int(msqid), unsafe.Pointer(msgp), C.size_t(len(msg)),
		C.int(msgtyp), C.int(msgflg)))
}

func MsgRcv(msgid int, msgtyp, msgflg int) (string, error) {
	var msg *C.char

	len := int(C.msgRcv(C.int(msgid), C.int(msgtyp), C.int(msgflg), &msg))
	if len < 0 {
		return "", fmt.Errorf("msgrcv error")
	}
	defer C.free(unsafe.Pointer(msg))
	return C.GoStringN(msg, C.int(len)), nil
}
