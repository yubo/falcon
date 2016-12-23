package goc

import (
	"reflect"
	"testing"
)

type t2_t struct {
	i32 int32
	i64 int64
	str string `size:"5"`
}

type t1_t struct {
	i32 int32
	i64 int64
	f64 float64
	str string `size:"5"`
	t2  t2_t   `type:"union"`
}

var src = []byte{1, 2, 3, 4, 5, 6, 7, 8}
var res = []int32{0x01020304, 0x05060708}

func checkResult(t *testing.T, dir string, err error, have, want interface{}) {
	if err != nil {
		t.Errorf("%v : %v", dir, err)
		return
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("%v :\n\thave %+v\n\twant %+v", dir, have, want)
	}
}

func TestRead(t *testing.T) {
	//t.Errorf("hello world")
}

func checkSize(t *testing.T, value interface{}, size int) {
	if size != Size(value) {
		t.Errorf("%+v size %d, want %d", value, Size(value), size)
	}
}
func TestSize(t *testing.T) {
	checkSize(t, int8(1), 1)
	checkSize(t, uint8(1), 1)
	checkSize(t, int16(1), 2)
	checkSize(t, uint16(1), 2)
	checkSize(t, int32(1), 4)
	checkSize(t, uint32(1), 4)
	checkSize(t, int64(1), 8)
	checkSize(t, uint64(1), 8)
	checkSize(t, float32(1), 4)
	checkSize(t, float64(1), 8)
	checkSize(t, t2_t{}, 17)
	checkSize(t, t1_t{}, 33)
}
