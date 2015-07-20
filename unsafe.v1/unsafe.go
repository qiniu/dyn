package unsafe

import (
	"reflect"
	"unsafe"
)

// ----------------------------------------------------------

func ToBytes(data string) []byte {

	sh := *(*reflect.StringHeader)(unsafe.Pointer(&data))
	arr := (*[1<<30]byte)(unsafe.Pointer(sh.Data))
	return arr[:sh.Len]
}

func ToString(data []byte) string {

	sh := *(*reflect.SliceHeader)(unsafe.Pointer(&data))
	ret := reflect.StringHeader{Data: sh.Data, Len: sh.Len}
	return *(*string)(unsafe.Pointer(&ret))
}

// ----------------------------------------------------------

