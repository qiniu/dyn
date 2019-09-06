package text

import (
	"net/url"
	"strconv"
	"strings"
	"syscall"
	"unicode/utf8"

	"github.com/qiniu/dyn/dyn"
	"github.com/qiniu/dyn/text/internal/encoding/json"
	"github.com/qiniu/x/errors"
)

const (
	Fmttype_Json    = 1
	Fmttype_Form    = 2
	Fmttype_Text    = 3
	Fmttype_Jsonstr = 4 // 在json的字符串内
)

// ----------------------------------------------------------

func AsJsonString(data interface{}) (val string, err error) {

retry:
	switch v := data.(type) {
	case func() interface{}:
		data = v()
		goto retry

	default:
		v2, ok2 := dyn.Int(data)
		if ok2 {
			return strconv.FormatInt(v2, 10), nil
		}
		v3, err3 := json.Marshal(data)
		if err3 != nil {
			return "", err3
		}
		val = string(v3)
	}
	return
}

func AsJsonstrString(data interface{}) (val string, err error) {

retry:
	switch v := data.(type) {
	case func() interface{}:
		data = v()
		goto retry

	default:
		v2, ok2 := dyn.Int(data)
		if ok2 {
			return strconv.FormatInt(v2, 10), nil
		}
		v3, err3 := json.Marshal(data)
		if err3 != nil {
			return "", err3
		}
		if v3[0] == '"' {
			v3 = v3[1 : len(v3)-1]
		}
		val = string(v3)
	}
	return
}

func AsQueryString(data interface{}) (val string, err error) {

	val, err = AsTextString(data)
	if err == nil {
		val = url.QueryEscape(val)
	}
	return
}

func AsTextString(data interface{}) (val string, err error) {

retry:
	switch v := data.(type) {
	case string:
		val = v

	case func() interface{}:
		data = v()
		goto retry

	default:
		if data == nil {
			return "", nil
		}
		v2, ok2 := dyn.Int(data)
		if ok2 {
			return strconv.FormatInt(v2, 10), nil
		}
		v3, err3 := json.Marshal(data)
		if err3 != nil {
			return "", err3
		}
		val = string(v3)
	}
	return
}

func AsString(data interface{}, ft int) (val string, err error) {

	switch ft {
	case Fmttype_Json:
		return AsJsonString(data)
	case Fmttype_Form:
		return AsQueryString(data)
	case Fmttype_Text:
		return AsTextString(data)
	case Fmttype_Jsonstr:
		return AsJsonstrString(data)
	}
	return "", syscall.EINVAL
}

func GetAsString(data interface{}, key string, ft int, failIfNotExists bool) (val string, err error) {

	v, ok := dyn.Get(data, key)
	if !ok {
		if failIfNotExists {
			return "", errors.New("dyn.Get key `" + key + "` not found")
		}
		return AsString(nil, ft)
	}
	return AsString(v, ft)
}

// ----------------------------------------------------------

func decodeVar(
	b []byte, exprvar string, pos int,
	data interface{}, ft, instring int, failIfNotExists bool) ([]byte, int, error) {

	if ft == Fmttype_Json && instring != 0 {
		ft = Fmttype_Jsonstr
	}

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
				val2, err2 := GetAsString(data, key2, ft, failIfNotExists)
				if err2 != nil {
					return nil, 0, errors.Info(err2, "expr.Exec - GetAsString failed", key2).Detail(err2)
				}
				return append(b, val2...), len(exprvar) - len(exprleft[pos2+1:]), nil
			}
		}
	}
	if pos+1 < len(exprvar) {
		if exprvar[pos+1] == '$' {
			return append(b, '$'), pos + 2, nil
		}
	}
	return nil, 0, errors.Info(syscall.EINVAL, "expr.Exec - invalid expr", exprvar[pos:])
}

// ----------------------------------------------------------

func Subst(exprvar string, data interface{}, ft int, failIfNotExists bool) (v string, err error) {

	var b []byte

	instring := 0
	start, pos := 0, 0
	end := len(exprvar)
	for pos < end {
		ch, size := utf8.DecodeRuneInString(exprvar[pos:])
		switch ch {
		case '"':
			instring ^= 1
		case '\\':
			pos += size
			if pos < end {
				_, size = utf8.DecodeRuneInString(exprvar[pos:])
			} else {
				size = 0
			}
		case '$':
			if b == nil {
				b = make([]byte, 0, len(exprvar))
			}
			b = append(b, exprvar[start:pos]...)
			b, pos, err = decodeVar(b, exprvar, pos, data, ft, instring, failIfNotExists)
			if err != nil {
				return
			}
			start, size = pos, 0
		}
		pos += size
	}

	if start < pos {
		if b == nil {
			return exprvar, nil
		}
		b = append(b, exprvar[start:pos]...)
	}
	return string(b), nil
}

// ----------------------------------------------------------
