package opcodes

// TODO: This package needs to be rewritten, right now it's tied too much to all users of it.

// OpCodeDescriptor describe a Hopper OPCODE. It holds information for the Hopper assembler to know
// how to translate Hopper Assembly to OPCODEs.
type OpCodeDescriptor struct {
	// The string (or ASM) version of this OPCODE.
	Mnemonic string
	// The OPCODE. This will not be populated in the opCodes, only in the assembling stage.
	OpCode byte
	// The function that will execute the OPCODE.
	Executor func(operand byte) (exitVM bool, incrementPC bool)
}

// OpCodes holds all supported Hopper operation codes.
// The key is the opcode's binary representation.
var OpCodes = map[byte]*OpCodeDescriptor{
	0: &OpCodeDescriptor{
		Mnemonic: "NOP",
	},
	1: &OpCodeDescriptor{
		Mnemonic: "LDA",
	},
	2: &OpCodeDescriptor{
		Mnemonic: "ADD",
	},
	3: &OpCodeDescriptor{
		Mnemonic: "SUB",
	},
	4: &OpCodeDescriptor{
		Mnemonic: "STR",
	},
	5: &OpCodeDescriptor{
		Mnemonic: "LDI",
	},
	6: &OpCodeDescriptor{
		Mnemonic: "JMP",
	},
	7: &OpCodeDescriptor{
		Mnemonic: "JC",
	},
	8: &OpCodeDescriptor{
		Mnemonic: "JZ",
	},
	// Other OP codes are reserved.
	14: &OpCodeDescriptor{
		Mnemonic: "OUT",
	},
	15: &OpCodeDescriptor{
		Mnemonic: "HLT",
	},
}
