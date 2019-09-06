package proto

// ----------------------------------------------------------

const (
	Fmttype_Invalid = -1
	Fmttype_Json    = 1
	Fmttype_Form    = 2
	Fmttype_Text    = 3
	Fmttype_Jsonstr = 4 // 在json的字符串内
)

// ----------------------------------------------------------

type Var struct {
	Key string
}

// ----------------------------------------------------------
