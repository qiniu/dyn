package text

import (
	"strings"
	"syscall"
	"testing"

	"qiniupkg.com/x/errors.v7"
	"qiniupkg.com/x/ts.v7"
)

func execOld(exprvar string, data interface{}, ft int) (val string, err error) {

	for {
		pos := strings.Index(exprvar, "$")
		if pos < 0 {
			return val + exprvar, nil
		}

		val += exprvar[:pos]
		if pos+2 < len(exprvar) {
			start := exprvar[pos+1]
			switch start {
			case '(', '{':
				end := ")"
				if start == '{' {
					end = "}"
				}
				exprleft := exprvar[pos+2:]
				pos2 := strings.Index(exprleft, end)
				if pos2 >= 0 {
					key2 := exprleft[:pos2]
					val2, err2 := GetAsString(data, key2, ft, false)
					if err2 != nil {
						err = errors.Info(err2, "expr.Subst - GetAsString failed", key2).Detail(err2)
						return
					}
					exprvar = exprleft[pos2+1:]
					val += val2
					continue
				}
			}
		}
		if pos+1 < len(exprvar) {
			if exprvar[pos+1] == '$' {
				exprvar = exprvar[pos+2:]
				val += "$"
				continue
			}
		}
		err = errors.Info(syscall.EINVAL, "expr.Subst - invalid expr", exprvar[pos:])
		break
	}
	return
}

var c1 = 1

var data = map[string]interface{}{
	"a": 1,
	"b": func() interface{} {
		return map[string]interface{}{
			"c": func() interface{} {
				c1++
				return c1
			},
		}
	},
	"d": "12\\&3",
}

type testCase struct {
	expr            string
	result          string
	resultForm      string
	resultText      string
}

var cases = []testCase{
	{"abc", "abc", "abc", "abc"},
	{"abc $$ aaa", "abc $ aaa", "abc $ aaa", "abc $ aaa"},
	{"abc $(b.c) aaa", "abc 2 aaa", "abc 3 aaa", "abc 4 aaa"},
	{"abc $(d) aaa", "abc \"12\\\\&3\" aaa", "abc 12%5C%263 aaa", "abc 12\\&3 aaa"},
	{"abc \"$(d)\" aaa", "abc \"12\\\\&3\" aaa", "abc \"12%5C%263\" aaa", "abc \"12\\&3\" aaa"},
	{"$$", "$", "$", "$"},
	{"$(b.c) aaa", "5 aaa", "6 aaa", "7 aaa"},
	{"abc $(d)", "abc \"12\\\\&3\"", "abc 12%5C%263", "abc 12\\&3"},
	{"abc=${e}&dee", "abc=null&dee", "abc=&dee", "abc=&dee"},
}

func TestExpr(t *testing.T) {

	for _, c := range cases {
		result, err := Subst(c.expr, data, Fmttype_Json, false)
		if err != nil || result != c.result {
			ts.Fatal(t, "Subst failed:", c.expr, "-", err, result, c.result)
		}
		result, err = Subst(c.expr, data, Fmttype_Form, false)
		if err != nil || result != c.resultForm {
			ts.Fatal(t, "Subst failed:", c.expr, "-", err, result, c.resultForm)
		}
		result, err = Subst(c.expr, data, Fmttype_Text, false)
		if err != nil || result != c.resultText {
			ts.Fatal(t, "Subst failed:", c.expr, "-", err, result, c.resultText)
		}
	}
	_, err := Subst("abc=${e}&dee", data, Fmttype_Form, true)
	if !(err != nil && err.Error() == "dyn.Get key `e` not found") {
		ts.Fatal(t, "Subst failed:", err)
	}
}

