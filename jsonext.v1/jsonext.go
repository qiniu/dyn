package jsonext

import (
	"reflect"
	"unsafe"

	. "qiniupkg.com/dyn/proto.v1"
)

// ----------------------------------------------------------

var varType = reflect.TypeOf(Var{})

func encodeVar(e *encodeState, v reflect.Value, quoted bool) {

	key := v.Interface().(Var).Key
	e.WriteByte('$')
	e.WriteByte('(')
	e.WriteString(key)
	e.WriteByte(')')
}

// ----------------------------------------------------------

func UnmarshalString(data string, v interface{}) error {

	sh := *(*reflect.StringHeader)(unsafe.Pointer(&data))
	arr := (*[1<<30]byte)(unsafe.Pointer(sh.Data))
	return Unmarshal(arr[:sh.Len], v)
}

func MarshalToString(v interface{}) (text string, err error) {

	b, err := Marshal(v)
	if err != nil {
		return
	}
	return string(b), nil
}

func MarshalIndentToString(v interface{}, prefix, indent string) (text string, err error) {

	b, err := MarshalIndent(v, prefix, indent)
	if err != nil {
		return
	}
	return string(b), nil
}

// ----------------------------------------------------------

