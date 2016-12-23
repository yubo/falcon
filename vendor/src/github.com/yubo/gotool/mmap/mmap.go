package mmap

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

type Mmap struct {
	buf      []byte
	buf_addr uintptr
	buf_len  int
	off      int // read at &buf[off], write at &buf[len(buf)]
}

func (m *Mmap) Write(p []byte) (n int, err error) {
	if m.off >= len(m.buf) {
		// Buffer is empty, reset to recover space.
		m.off = 0
		if len(p) == 0 {
			return
		}
		return 0, io.EOF
	}
	n = copy(m.buf[m.off:], p)
	m.off += n
	return
}

func (m *Mmap) WriteString(p string) (n int, err error) {
	if m.off >= len(m.buf) {
		// Buffer is empty, reset to recover space.
		m.off = 0
		if len(p) == 0 {
			return
		}
		return 0, io.EOF
	}
	n = copy(m.buf[m.off:], p)
	m.off += n
	return
}

func (m *Mmap) Read(p []byte) (n int, err error) {
	if m.off >= len(m.buf) {
		// Buffer is empty, reset to recover space.
		m.off = 0
		if len(p) == 0 {
			return
		}
		return 0, io.EOF
	}
	n = copy(p, m.buf[m.off:])
	m.off += n
	return
}

func (m *Mmap) String() string {
	return fmt.Sprintf("off:%d size:%d data:%x", m.off, m.buf_len, m.buf)
}

func (m *Mmap) Reset() {
	m.off = 0
}

func (m *Mmap) Seek(off int) error {
	if off < 0 {
		return errors.New("off < 0")
	} else if off >= m.buf_len {
		return errors.New("off > buf_len")
	} else {
		m.off = off
		return nil
	}
}

func (m *Mmap) Dump() {
	log.Println(m)
}

func (m *Mmap) Lock() error {
	if _, _, errno := syscall.Syscall(syscall.SYS_MLOCK, m.buf_addr, uintptr(m.buf_len), 0); errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}

func (m *Mmap) Unlock() error {
	if _, _, errno := syscall.Syscall(syscall.SYS_MUNLOCK, m.buf_addr, uintptr(m.buf_len), 0); errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}

func (m *Mmap) Flush() error {
	if _, _, errno := syscall.Syscall(syscall.SYS_MSYNC, m.buf_addr, uintptr(m.buf_len), syscall.MS_SYNC); errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}

func (m *Mmap) Unmap() error {
	if _, _, errno := syscall.Syscall(syscall.SYS_MUNMAP, m.buf_addr, uintptr(m.buf_len), 0); errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}

func Mmapfd(fd *os.File, offset, length int, prot, flags int) (*Mmap, error) {
	var fi os.FileInfo
	var m Mmap
	var err error

	if fi, err = fd.Stat(); err != nil {
		return nil, err
	}
	if fi, err = fd.Stat(); err != nil {
		return nil, err
	}
	size := int(fi.Size())
	if offset < 0 {
		if offset = size + offset; offset < 0 {
			return nil, errors.New("abs(offset) > file_size")
		}
	} else if offset > size {
		return nil, errors.New("abs(offset) > file_size")

	}
	if length <= 0 {
		length = size - offset
	}

	if prot < 0 {
		prot = syscall.PROT_WRITE | syscall.PROT_READ
	}

	if flags < 0 {
		flags = syscall.MAP_SHARED
	}
	log.Println(offset, length)
	if m.buf, err = syscall.Mmap(int(fd.Fd()), int64(offset), int(length), syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED); err != nil {
		return nil, err
	}
	dh := (*reflect.SliceHeader)(unsafe.Pointer(&m.buf))
	m.buf_addr = dh.Data
	m.buf_len = dh.Len
	return &m, nil
}

func Mmapf(filename string, offset, length int, prot, flags int) (fd *os.File, mmap *Mmap, err error) {
	if fd, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644); err != nil {
		return nil, nil, err
	}
	if mmap, err = Mmapfd(fd, offset, length, prot, flags); err != nil {
		fd.Close()
		return nil, nil, err
	}
	return fd, mmap, nil
}
