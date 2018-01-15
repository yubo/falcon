/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package alarm

import (
	"unsafe"

	"github.com/yubo/gotool/list"
)

func list_entry(l *list.ListHead) *eventEntry {
	return (*eventEntry)(unsafe.Pointer((uintptr(unsafe.Pointer(l)) -
		unsafe.Offsetof(((*eventEntry)(nil)).list))))
}
