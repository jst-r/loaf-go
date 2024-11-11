package bytecode

const (
	OpConstant uint8 = iota
	OpNil
	OpTrue
	OpFalse
	OpNot
	OpNegate
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpEqual
	OpGreater
	OpLess
	OpReturn
)
