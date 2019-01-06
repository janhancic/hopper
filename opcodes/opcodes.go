// Package opcodes defines all the supported operation codes in a Hopper architecture.
package opcodes

// OpCode represents a Hopper supported operation code.
type OpCode byte

// NOP  does nothing. It has no operands. 0000
var NOP OpCode

// LDA stores the value addressed by its operand into register A. 0001
var LDA OpCode = 1

// ADD adds the value addressed by its operand and the value in register A together. It stores the
// result in register A. 0010
var ADD OpCode = 2

// SUB subtracts the value addressed by its operand from the value in register A. It stores the
// result in register A. 0011
var SUB OpCode = 3

// STR stores the value addressed by its operand in register A. 0100
var STR OpCode = 4

// LDI stores the operand in register A. 0101
var LDI OpCode = 5

// JMP jumps to the address specified in the operand. 0110
var JMP OpCode = 6

// JC jumps to the address specified in the operand if the result of the last operation resulted in
// a carry bit. 0111
var JC OpCode = 7

// JZ jumps to the address specified in the operand if the result of the last operation was 0. 1000
var JZ OpCode = 8

// OUT puts the value of register A into the Out register. 1110
var OUT OpCode = 14

// HLT halts the execution of the computer. 1111
var HLT OpCode = 15

// OpCodeMnemonics is a convenience map of OpCode->mnemonic for easy lookups.
var OpCodeMnemonics = map[OpCode]string{
	NOP: "NOP",
	LDA: "LDA",
	ADD: "ADD",
	SUB: "SUB",
	STR: "STR",
	LDI: "LDI",
	JMP: "JMP",
	JC:  "JC",
	JZ:  "JZ",
	OUT: "OUT",
	HLT: "HLT",
}

// MnemonicOpCodes is a convenience map of mnemonic->OpCode for easy lookups.
var MnemonicOpCodes = map[string]OpCode{
	"NOP": NOP,
	"LDA": LDA,
	"ADD": ADD,
	"SUB": SUB,
	"STR": STR,
	"LDI": LDI,
	"JMP": JMP,
	"JC":  JC,
	"JZ":  JZ,
	"OUT": OUT,
	"HLT": HLT,
}
