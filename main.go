package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// opCodeDescriptor describe a Hopper OPCODE. It holds information for the Hopper assembler to know
// how to translate Hopper Assembly to OPCODEs.
type opCodeDescriptor struct {
	// The string (or ASM) version of this OPCODE.
	asm string
	// The OPCODE. This will not be populated in the opCodes, only in the assembling stage.
	opCode byte
	// Indicates wether this OPCODE has any operands.
	noOperand bool
	// The function that will execute the OPCODE.
	executor func(operand byte) (result byte, exitVM bool, incrementPC bool)
}

// opCodes holds all supported OPCODEs in Hopper VM. The key is the OPCODE's binary representation.
var opCodes = map[byte]opCodeDescriptor{
	0: opCodeDescriptor{
		asm:       "NOOP",
		noOperand: true,
		executor:  func(operand byte) (byte, bool, bool) { return 0, false, true },
	},
	1: opCodeDescriptor{
		asm: "LODA",
		executor: func(operand byte) (result byte, exitVM bool, incrementPC bool) {
			registerA = ram[operand]
			return registerA, false, true
		},
	},
	2: opCodeDescriptor{
		asm: "LODB",
		executor: func(operand byte) (result byte, exitVM bool, incrementPC bool) {
			registerB = ram[operand]
			return registerB, false, true
		},
	},
	3: opCodeDescriptor{
		asm: "INCR",
		executor: func(operand byte) (result byte, exitVM bool, incrementPC bool) {
			ram[operand]++
			return 0, false, true
		},
	},
	4: opCodeDescriptor{
		asm: "SUBT",
		executor: func(operand byte) (result byte, exitVM bool, incrementPC bool) {
			result = registerA - registerB
			ram[operand] = result
			return result, false, true
		},
	},
	5: opCodeDescriptor{
		asm: "JMPP",
		executor: func(operand byte) (result byte, exitVM bool, incrementPC bool) {
			if flagRegister == flagPositive {
				pc = operand
				return 0, false, false
			}
			return 0, false, true
		},
	},
	6: opCodeDescriptor{
		asm:       "HALT",
		noOperand: true,
		executor:  func(_ byte) (result byte, exitVM bool, incrementPC bool) { return 0, true, false },
	},
}

/*
TODO: Possible future instructions
SAVE A ADDR
SAVE B ADDR
ADD  ADDR ; Add ADDR=A+B
JMP  ADDR ; Jump to ADDR
JMPZ ADDR ; Jump to ADDR if flagRegister is zero
JMPN ADDR ; Jump to ADDR if flagRegister is negative
*/

// Flags for the flagRegister.
const (
	flagZero     = 0
	flagPositive = 1
	flagNegative = 2
)

var (
	// Two general purpose registers.
	registerA byte
	registerB byte
	// Indicates the result of the last operation.
	flagRegister byte
	// The working memory for the VM.
	ram [16]byte
	// Program counter.
	pc byte
)

func runVM() {
	for {
		clearScreen()
		printState()
		fmt.Printf("Press enter to advance program.")
		fmt.Scanln()

		instruction := ram[pc]
		// Op codes are defined in the first 4 bits, shifting the instruction by 4 bits to the right
		// gives us the op code only without arguments.
		opCode := instruction >> 4
		opOperand := clearMsb(instruction)

		result, exitVM, incrementPC := opCodes[opCode].executor(opOperand)
		if result == 0 {
			flagRegister = flagZero
		} else if isBitSet(result, 8) {
			flagRegister = flagNegative
		} else {
			flagRegister = flagPositive
		}

		if exitVM {
			break
		}

		if incrementPC {
			pc++
		}
	}

	clearScreen()
	printState()
}

func assemble(hopAsm []string) []byte {
	asmCodes := map[string]opCodeDescriptor{} // reverse lookup map from ASMCODE->OPCODE
	for opCode, opCodeDesc := range opCodes {
		opCodeDesc.opCode = opCode
		asmCodes[opCodeDesc.asm] = opCodeDesc
	}

	machineCode := make([]byte, len(hopAsm))
	for idx, asmInstruction := range hopAsm {
		var machineInstruction byte
		commentPos := strings.Index(asmInstruction, ";")
		if commentPos != -1 {
			asmInstruction = strings.Trim(asmInstruction[0:commentPos], " ")
		}

		asmCode := asmInstruction
		asmOperand := ""
		spacePos := strings.Index(asmInstruction, " ")
		if spacePos != -1 {
			asmCode = asmInstruction[0:spacePos]
			asmOperand = asmInstruction[spacePos+1:]
		}

		if opCodeDesc, isOpCode := asmCodes[asmCode]; isOpCode {
			if asmOperand == "" {
				machineInstruction = opCodeDesc.opCode << 4 // move the OPCODE to MSB
			} else {
				opOperandAsInt, _ := strconv.ParseUint(asmOperand, 10, 8)
				opOperandAsByte := byte(opOperandAsInt)

				machineInstruction = opCodeDesc.opCode << 4 // move the OPCODE to MSB
				machineInstruction ^= opOperandAsByte       // set the LSB to the operand value
			}
		} else {
			// treat as memory value, convert to byte
			instructionAsByte, _ := strconv.ParseUint(asmInstruction, 10, 8)
			machineInstruction = byte(instructionAsByte)
		}

		machineCode[idx] = machineInstruction
	}

	return machineCode
}

func main() {
	if len(os.Args) < 3 {
		printUsageAndDie()
	}

	if os.Args[1] == "run" {
		binContents, _ := ioutil.ReadFile(os.Args[2])
		copy(ram[:], binContents)
		runVM()
	} else if os.Args[1] == "assemble" {
		hopContents, _ := ioutil.ReadFile(os.Args[2])
		machineCode := assemble(strings.Split(string(hopContents), "\n"))
		ioutil.WriteFile(os.Args[2]+".bin", machineCode, 0644)
	} else {
		printUsageAndDie()
	}
}

/*
TODO:
- control clock with a command line (maybe even have a manual step)
- show the actual command in text form of the command the PC counter is pointing to
*/
