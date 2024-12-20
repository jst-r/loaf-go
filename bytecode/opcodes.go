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
	OpPrint
	OpPop
	OpDefineGlobal
	OpGetGlobal
	OpSetGlobal
	OpGetLocal
	OpSetLocal
	OpJumpIfFalse
	OpJump
	OpLoop
	OpReturn
)
