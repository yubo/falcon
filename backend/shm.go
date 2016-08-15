/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

/*
#include "shm.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

//void *shmat(int shmid, const void *shmaddr, int shmflg);
func Shmat(shmid int, shmaddr uintptr,
	shmflg int) (uintptr, error) {
	addr := uintptr(unsafe.Pointer(C.shmat(C.int(shmid),
		unsafe.Pointer(shmaddr), C.int(shmflg))))
	if addr == uintptr(unsafe.Pointer(C.err_addr)) {
		return 0, errors.New("shmat error")
	}
	return addr, nil
}

//int shmget(key_t key, size_t size, int shmflg);
//return shm_id, segment_size, error
func Shmget(key, size, shmflg int) (int, int, error) {
	var (
		buf C.struct_shmid_ds
	)

	shm_id := C.shmget(C.key_t(key), C.size_t(size),
		C.int(shmflg))
	if shm_id == -1 {
		return -1, 0, errors.New("shmget error")
	}

	if C.shmctl(shm_id, C.IPC_STAT, &buf) == -1 {
		buf.shm_segsz = 0
	}

	return int(shm_id), int(buf.shm_segsz), nil
}

func Shmdt(p uintptr) error {
	if int(C.shmdt(unsafe.Pointer(p))) == -1 {
		return errors.New("detach error")
	}
	return nil
}

func Shmrm(shmid int) error {
	if int(C.shmctl(C.int(shmid), C.IPC_RMID,
		(*C.struct_shmid_ds)(unsafe.Pointer(uintptr(0))))) == -1 {
		return errors.New("shmrm error")
	}
	return nil
}
