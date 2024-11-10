package bytecode

const (
	OpConstant uint8 = iota
	OpNil
	OpTrue
	OpFalse
	OpNegate
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpReturn
)
