package flag

import (
	"reflect"
	"testing"

	"github.com/qiniu/dyn/vars"
)

// ---------------------------------------------------------------------------

type retArgs struct {
	Ret  Cmd `cmd:"ret"`
	Code int `arg:"code,opt"` // opt: 可选参数
}

type hostArgs struct {
	HostCmd Cmd    `cmd:"host"`
	Host    string `arg:"host - eg. api.qiniu.com"`
	Portal  string `arg:"portal - eg. <ip>:<port>"`
}

type base64Args struct {
	Base64      Cmd    `cmd:"base64"`
	StdEncoding bool   `flag:"std - use standard base64 encoding. default is urlsafe base64 encoding."`
	Fdecode     bool   `flag:"d - to decode data. default is to encode data."`
	Data        string `arg:"data"`
}

// ---------------------------------------------------------------------------

type caseParseArgs struct {
	argsType reflect.Type
	cmd      []string
	args     interface{}
	err      error
}

func TestParse(t *testing.T) {

	cases := []caseParseArgs{
		{
			argsType: reflect.TypeOf((*retArgs)(nil)),
			cmd:      []string{"ret", "200"},
			args:     &retArgs{Cmd{}, 200},
			err:      nil,
		},
		{
			argsType: reflect.TypeOf((*retArgs)(nil)),
			cmd:      []string{"ret", "$(code)"},
			args:     &retArgs{Cmd{}, 200},
			err:      nil,
		},
		{
			argsType: reflect.TypeOf((*retArgs)(nil)),
			cmd:      []string{"ret"},
			args:     &retArgs{Cmd{}, 0},
			err:      nil,
		},
		{
			argsType: reflect.TypeOf((*hostArgs)(nil)),
			cmd:      []string{"host", "api.qiniu.com", "192.168.1.10:8888"},
			args:     &hostArgs{Cmd{}, "api.qiniu.com", "192.168.1.10:8888"},
			err:      nil,
		},
		{
			argsType: reflect.TypeOf((*base64Args)(nil)),
			cmd:      []string{"base64", "-std", "hello"},
			args:     &base64Args{StdEncoding: true, Fdecode: false, Data: "hello"},
			err:      nil,
		},
	}

	ctx := vars.New()
	setVar(t, ctx, "code", 200)

	for _, c := range cases {
		args, err := Parse(ctx, c.argsType, c.cmd)
		if err != c.err {
			t.Fatal("Parse unexpected error:", err, c)
		}
		if !reflect.DeepEqual(args.Interface(), c.args) {
			t.Fatal("Parse unexpected args:", args.Interface(), c)
		}
	}
}

func setVar(t *testing.T, ctx *vars.Context, varName string, obj interface{}) {

	ctx.DeleteVar(varName)
	err := ctx.MatchVar(varName, obj)
	if err != nil {
		t.Fatal("setVar Match failed:", err)
	}
}

// ---------------------------------------------------------------------------
