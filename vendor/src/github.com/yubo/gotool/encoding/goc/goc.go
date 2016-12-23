package goc

import (
	"io"
	"reflect"
	"strconv"
)

func Read(r io.Reader, data interface{}) error {
	return nil
}

func Write(w io.Writer, data interface{}) error {
	return nil
}

func Size(v interface{}) int {
	return dataSize(reflect.Indirect(reflect.ValueOf(v)))
}

// dataSize returns the number of bytes the actual data represented by v occupies in memory.
// For compound structures, it sums the sizes of the elements. Thus, for instance, for a slice
// it returns the length of the slice times the element size and does not count the memory
// occupied by the header. If the type of v is not acceptable, dataSize returns -1.
func dataSize(v reflect.Value) int {
	if v.Kind() == reflect.Slice {
		if s := sizeof(v.Type().Elem(), reflect.StructTag("")); s >= 0 {
			return s * v.Len()
		}
		return -1
	}
	return sizeof(v.Type(), reflect.StructTag(""))
}

// sizeof returns the size >= 0 of variables for the given type or -1 if the type is not acceptable.
func sizeof(t reflect.Type, tag reflect.StructTag) int {
	switch t.Kind() {
	case reflect.Array:
		if s := sizeof(t.Elem(), tag); s >= 0 {
			return s * t.Len()
		}

	case reflect.Struct:
		if tag.Get("type") == "union" {
			max := 0
			for i, n := 0, t.NumField(); i < n; i++ {
				s := sizeof(t.Field(i).Type, t.Field(i).Tag)
				if s < 0 {
					return -1
				}
				if s > max {
					max = s
				}
			}
			return max
		} else {
			sum := 0
			for i, n := 0, t.NumField(); i < n; i++ {
				s := sizeof(t.Field(i).Type, t.Field(i).Tag)
				if s < 0 {
					return -1
				}
				sum += s
			}
			return sum
		}

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return int(t.Size())

	case reflect.String:
		if s, err := strconv.Atoi(tag.Get("size")); err == nil && s > 0 {
			return s
		} else {
			return -1
		}
	}

	return -1
}
