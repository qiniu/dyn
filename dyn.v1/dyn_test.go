package dyn

import (
	"testing"

	"qiniupkg.com/x/log.v7"
	"qiniupkg.com/x/ts.v7"
)

func init() {
	log.SetOutputLevel(0)
}

// ----------------------------------------------------------

type M map[string]interface{}
type M2 map[string]interface{}

func TestFuncGet(t *testing.T) {

	c1 := 1
	data := M{
		"a": 1,
		"b": func() interface{} {
			return map[string]interface{}{
				"c": func() interface{} {
					c1++
					return c1
				},
			}
		},
		"d": nil,
		"e": "123",
		"f": M2{
			"g": 456,
		},
	}

	_, ok := Get(nil, "a")
	if ok {
		ts.Fatal(t, "a?")
	}

	cc, ok := Get(data, "b.c")
	if !ok || cc == nil {
		ts.Fatal(t, "b.c:", cc, ok)
	}

	c, ok := GetInt(data, "b.c")
	if !ok || c != 2 {
		ts.Fatal(t, "b.c != 2", c, ok)
	}

	c2, ok := GetFloat(data, "b.c")
	if !ok || c2 != 3 {
		ts.Fatal(t, "b.c != 3")
	}

	_, ok = GetFloat(data, "d")
	if ok {
		ts.Fatal(t, "d is float?")
	}

	_, ok = GetString(data, "d")
	if ok {
		ts.Fatal(t, "a is string?")
	}
}

// ----------------------------------------------------------

func TestMapGet(t *testing.T) {

	data := map[string]interface{}{
		"a": 1,
		"b": map[string]interface{}{
			"c": 2,
		},
		"d": nil,
	}

	_, ok := Get(nil, "a")
	if ok {
		ts.Fatal(t, "a?")
	}

	c, ok := GetInt(data, "b.c")
	if !ok || c != 2 {
		ts.Fatal(t, "b.c != 2")
	}

	c2, ok := GetFloat(data, "b.c")
	if !ok || c2 != 2 {
		ts.Fatal(t, "b.c != 2")
	}

	_, ok = GetFloat(data, "d")
	if ok {
		ts.Fatal(t, "d is float?")
	}

	_, ok = GetString(data, "d")
	if ok {
		ts.Fatal(t, "a is string?")
	}
}

// ----------------------------------------------------------

type bar struct {
	C int `json:"c"`
}

type foo struct {
	A int         `json:"a"`
	B bar         `json:"b"`
	D interface{} `json:"d"`
}

func TestStructGet(t *testing.T) {

	data := &foo{
		A: 1,
		B: bar{
			C: 2,
		},
		D: nil,
	}

	_, ok := Get(nil, "a")
	if ok {
		ts.Fatal(t, "a?")
	}

	c, ok := GetInt(data, "b.c")
	if !ok || c != 2 {
		ts.Fatal(t, "b.c != 2, b.c =", c, ok)
	}

	c2, ok := GetFloat(data, "b.c")
	if !ok || c2 != 2 {
		ts.Fatal(t, "b.c != 2")
	}

	_, ok = GetFloat(data, "d")
	if ok {
		ts.Fatal(t, "d is float?")
	}

	_, ok = GetString(data, "d")
	if ok {
		ts.Fatal(t, "a is string?")
	}
}

// ----------------------------------------------------------

func TestArrayGet(t *testing.T) {

	data := []interface{}{
		5,
		"a",
		map[string]interface{}{
			"c": 2,
		},
		[]int{1, 3},
	}

	c, ok := GetInt(data, "0")
	if !ok || c != 5 {
		ts.Fatal(t, "0 != 5")
	}

	c1, ok := GetString(data, "1")
	if !ok || c1 != "a" {
		ts.Fatal(t, "1 != a")
	}

	c2, ok := GetInt(data, "2.c")
	if !ok || c2 != 2 {
		ts.Fatal(t, "2.c != 2")
	}

	c3, ok := GetInt(data, "3.1")
	if !ok || c3 != 3 {
		ts.Fatal(t, "3.1 != 3")
	}
}

// ----------------------------------------------------------
