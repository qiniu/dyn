package vars

import (
	"errors"
	"reflect"
	"testing"

	"qiniupkg.com/x/jsonutil.v7"
)

// ---------------------------------------------------------------------------

type caseMatchVar struct {
	Key         string
	RealVal     string
	Err         error
	OldDom      string
	ExpectedDom string
}

func TestMatchVar(t *testing.T) {

	cases := []caseMatchVar{
		{
			Key: `a.b`,
			RealVal: `{"value": 1}`,
			Err: nil,
			OldDom: `{}`,
			ExpectedDom: `{"a": {"b": {"value": 1}}}`,
		},
		{
			Key: `a.b`,
			RealVal: `{"value": 1}`,
			Err: errors.New("unmatched field: `value`"),
			OldDom: `{"a": {"b": {"value": 2}}}`,
			ExpectedDom: `{"a": {"b": {"value": 2}}}`,
		},
		{
			Key: `a.b`,
			RealVal: `{"value": 1}`,
			Err: ErrUnmatchedValue,
			OldDom: `{"a": "123"}`,
			ExpectedDom: `{"a": "123"}`,
		},
		{
			Key: `a.b`,
			RealVal: `{"value": {"37": "64"}}`,
			Err: nil,
			OldDom: `{"a": {}, "c": 1.2}`,
			ExpectedDom: `{"a": {"b": {"value": {"37": "64"}}}, "c": 1.2}`,
		},
	}

	for _, c := range cases {
		ctx := New()
		err := jsonutil.Unmarshal(c.OldDom, &ctx.vars)
		if err != nil {
			t.Fatal("jsonutil.Unmarshal OldDom failed:", c.OldDom, err)
		}
		var vreal interface{}
		err = jsonutil.Unmarshal(c.RealVal, &vreal)
		if err != nil {
			t.Fatal("jsonutil.Unmarshal RealVal failed:", c.RealVal, err)
		}
		err = ctx.MatchVar(c.Key, vreal)
		if err != c.Err {
			if !(err != nil && c.Err != nil && err.Error() == c.Err.Error()) {
				t.Fatal("MatchVar unexpected error:", c, err)
			}
		}
		var vdom interface{}
		err = jsonutil.Unmarshal(c.ExpectedDom, &vdom)
		if err != nil {
			t.Fatal("jsonutil.Unmarshal ExpectedDom failed:", c.ExpectedDom, err)
		}
		if !reflect.DeepEqual(ctx.vars, vdom) {
			t.Fatal("MatchVar unexpected dom:", c, vdom)
		}
	}
}

// ---------------------------------------------------------------------------

