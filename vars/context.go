package vars

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/qiniu/dyn/dyn"
	"github.com/qiniu/dyn/text"
	"github.com/qiniu/x/log"

	"github.com/qiniu/dyn/proto"
)

const (
	Fmttype_Invalid = proto.Fmttype_Invalid
	Fmttype_Json    = proto.Fmttype_Json
	Fmttype_Form    = proto.Fmttype_Form
	Fmttype_Text    = proto.Fmttype_Text
	Fmttype_Jsonstr = proto.Fmttype_Jsonstr
)

var (
	ErrNotVar           = errors.New("assign to a non variable expression")
	ErrUnmatchedValue   = errors.New("unmatched value")
	ErrSliceLenNotEqual = errors.New("slice length not equal")
)

// ---------------------------------------------------------------------------

type Context struct {
	vars map[string]interface{}
}

func New() *Context {

	vars := make(map[string]interface{})
	return &Context{
		vars: vars,
	}
}

// ---------------------------------------------------------------------------

func (p *Context) GetVars() map[string]interface{} {

	return p.vars
}

func (p *Context) GetVar(key string) (v interface{}, ok bool) {

	return dyn.Get(p.vars, key)
}

func (p *Context) DeleteVar(key string) {

	v := p.vars
	parts := strings.Split(key, ".")
	ilast := len(parts) - 1

	for i := 0; i < ilast; i++ {
		v1, ok1 := v[parts[i]]
		if !ok1 {
			return
		}
		v2, ok2 := v1.(map[string]interface{})
		if !ok2 || v2 == nil {
			return
		}
		v = v2
	}
	delete(v, parts[ilast])
}

func (p *Context) MatchVar(key string, vreal interface{}) (err error) {

	var v interface{}
	var ok bool

	parts := strings.Split(key, ".")
	ilast := len(parts) - 1

	v = p.vars
	for i, part := range parts {
		v1, ok1 := v.(map[string]interface{})
		if !ok1 || v1 == nil {
			return ErrUnmatchedValue
		}
		if v, ok = v1[part]; ok {
			continue
		}
		if i == ilast {
			v1[part] = vreal
			return nil
		}
		v = make(map[string]interface{})
		v1[part] = v
	}

	return p.Match(v, vreal)
}

// ---------------------------------------------------------------------------

func (p *Context) Match(vexp, vreal interface{}) (err error) {

	return p.doMatch(vexp, vreal, "")
}

func (p *Context) doMatch(vexp, vreal interface{}, field string) (err error) {

	switch v := vexp.(type) {
	case map[string]interface{}:
		if v2, ok := vreal.(map[string]interface{}); ok {
			sfield := field
			if sfield != "" {
				sfield += "."
			}
			for sk, sv := range v {
				if sv2, ok2 := v2[sk]; ok2 {
					err = p.doMatch(sv, sv2, sfield+sk)
					if err == nil {
						continue
					}
				} else {
					err = errors.New("field not found: `" + sfield + sk + "`")
				}
				return
			}
			return nil
		}
		log.Debug("Match map object failed:", vexp, vreal)

	case proto.Var:
		err2 := p.MatchVar(v.Key, vreal)
		if err2 != nil {
			if field == "" {
				return ErrUnmatchedValue
			}
			return errors.New("match field `" + field + "` failed: " + err2.Error())
		}
		return nil

	case []interface{}:
		v2 := reflect.ValueOf(vreal)
		if v2.Kind() == reflect.Slice {
			if len(v) != v2.Len() {
				err = ErrSliceLenNotEqual
				return
			}
			sfield := field
			if sfield != "" {
				sfield += "."
			}
			for i, sv := range v {
				err = p.doMatch(sv, v2.Index(i).Interface(), sfield+strconv.Itoa(i))
				if err != nil {
					return
				}
			}
			return nil
		}
		log.Debug("Match slice object failed:", vexp, vreal)

	default:
		if reflect.DeepEqual(vexp, vreal) {
			return nil
		}
		log.Debug("Match value failed:", vexp, vreal)
	}

	if field == "" {
		return ErrUnmatchedValue
	}
	return errors.New("unmatched field: `" + field + "`")
}

// ---------------------------------------------------------------------------

func (p *Context) Let(vexp, vreal interface{}) (err error) {

	if v, ok := vexp.(proto.Var); ok {
		p.DeleteVar(v.Key)
		return p.Match(vexp, vreal)
	}
	return ErrNotVar
}

// ---------------------------------------------------------------------------

func (p *Context) Subst(vexp interface{}, ft int) (vres interface{}, err error) {

	switch v := vexp.(type) {
	case map[string]interface{}:
		vres1 := make(map[string]interface{})
		for sk, sv := range v {
			vres2, err2 := p.Subst(sv, Fmttype_Invalid)
			if err2 != nil {
				return nil, err2
			}
			vres1[sk] = vres2
		}
		return vres1, nil

	case string:
		if ft == Fmttype_Invalid {
			return vexp, nil
		}
		vres, err = p.SubstText(v, ft)
		return

	case proto.Var:
		vres1, ok1 := dyn.Get(p.vars, v.Key)
		if !ok1 {
			return nil, errors.New("var `" + v.Key + "` not found")
		}
		return vres1, nil

	case []interface{}:
		n := len(v)
		vres1 := make([]interface{}, n)
		for i, sv := range v {
			vres2, err2 := p.Subst(sv, Fmttype_Invalid)
			if err2 != nil {
				return nil, err2
			}
			vres1[i] = vres2
		}
		return vres1, nil

	default:
		return vexp, nil
	}
}

func (p *Context) SubstText(exprvar string, ft int) (v string, err error) {

	return text.Subst(exprvar, p.vars, ft, true)
}

// ---------------------------------------------------------------------------
