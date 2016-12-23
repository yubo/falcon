/*
 * Copyright 2015 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package mmap

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

var file = filepath.Join(os.TempDir(), "mmap_test")

func TestMmap(t *testing.T) {
	defer os.Remove(file)
	f, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
	f.WriteString("1234")
	f.Close()
	{
		f, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
		m, err := Mmapfd(f, 0, 0, -1, -1)
		if err != nil {
			log.Fatalln(err)
		}
		m.Dump()
		m.WriteString("6")
		m.WriteString("6")
		m.Dump()
		m.Unmap()
		f.Close()
	}
	//flush
	//syscall.Syscall(syscall.SYS_MSYNC, dh.Data, uintptr(dh.Len), syscall.MS_SYNC)
	//unmap
	//syscall.Syscall(syscall.SYS_MUNMAP, dh.Data, uintptr(dh.Len), 0)
	//mmap, err := Map(f, RDWR, 0)
}
