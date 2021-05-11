package model

type Struct struct {
	Name    string // Name of the struct
	Fields  []*Field
	Methods []*Method
	Abbrv   string
}

// Method defines a single method of a Struct
type Method struct {
	Abbrv            string
	Struct           string
	RemoteStructName string
	RemoteMethodName string
	OutputStruct     *Argument
	Name             string      // field name
	Arguments        []*Argument // function arguments
	ReturnValues     []*Field    // return types
}

type Argument struct {
	Name    string // type name e.g. ChannelID
	Type    string // like string, int etc.
	Pointer bool
}

type Field struct {
	Import  string
	Name    string
	Type    string
	Pointer bool
}
