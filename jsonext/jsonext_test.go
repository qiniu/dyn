package jsonext

import (
	"fmt"
	"reflect"
	"testing"
)

func TestExt(t *testing.T) {

	var v interface{}

	err := Unmarshal([]byte(`
		{"a": ${a}, "b": $(b)}
	`), &v)
	if err != nil {
		t.Fatal("Unmarshal failed:", err)
	}

	b, _ := Marshal(v)
	fmt.Println("jsonext:", string(b))

	var v2 interface{}
	err = Unmarshal(b, &v2)
	if err != nil {
		t.Fatal("Unmarshal v2 failed:", err)
	}

	if !reflect.DeepEqual(v, v2) {
		t.Fatal("v != v2 -", v, v2)
	}
}

