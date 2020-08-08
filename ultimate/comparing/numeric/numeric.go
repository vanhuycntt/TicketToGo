package numeric

type IntType int

type valueType struct {
	IntTypeLevel IntType
	Value        int
}

func New() valueType {
	return valueType{
		IntType(200),
		10,
	}
}
