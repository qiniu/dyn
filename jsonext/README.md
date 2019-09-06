qiniupkg.com/x/dyn/jsonext.v1
============

这个包扩展了 json 的文法，增加变量的支持。变量的形式有两种：

* `$(...)`
* `${...}`

样例：

```
{
	"a": $(a),
	"b": {
		"c": ${b.c},
		"e": $(e)
	},
	"f": ${f}
}
```

这段 json 文本进行 Unmarshal 后变量的数据类型成为：

* qiniupkg.com/x/dyn/proto.v1.Var

而 Var 类型的变量 Marshal 后又变回 $(...) 字符串。像上面的 json 文本进行 Unmarshal 再 Marshal 得到的结果会是：

```
{"a":$(a),"b":{"c":$(b.c),"e":$(e)},"f":$(f)}
```

