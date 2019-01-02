package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// opCodeDescriptor describe a Hopper OPCODE. It holds information for the Hopper assembler to know
// how to translate Hopper Assembly to OPCODEs.
type opCodeDescriptor struct {
	// The string (or ASM) version of this OPCODE.
	mnemonic string
	// The OPCODE. This will not be populated in the opCodes, only in the assembling stage.
	opCode byte
	// The function that will execute the OPCODE.
	executor func(operand byte) (exitVM bool, incrementPC bool)
}

// opCodes holds all supported OPCODEs in Hopper VM. The key is the OPCODE's binary representation.
var opCodes = map[byte]opCodeDescriptor{
	0: opCodeDescriptor{
		mnemonic: "NOP",
		executor: func(operand byte) (bool, bool) { return false, true },
	},
	1: opCodeDescriptor{
		mnemonic: "LDA",
		executor: func(operand byte) (exitVM bool, incrementPC bool) {
			registerA = ram[operand]
			return false, true
		},
	},
	2: opCodeDescriptor{
		mnemonic: "ADD",
		executor: func(operand byte) (exitVM bool, incrementPC bool) {
			result, carry := byteAdder(registerA, ram[operand])

			registerA = result
			flagCarryRegister = carry
			flagZeroRegister = result == 0

			return false, true
		},
	},
	3: opCodeDescriptor{
		mnemonic: "SUB",
		executor: func(operand byte) (exitVM bool, incrementPC bool) {
			result, carry := byteSubtractor(registerA, ram[operand])

			registerA = result
			flagCarryRegister = carry
			flagZeroRegister = result == 0

			return false, true
		},
	},
	4: opCodeDescriptor{
		mnemonic: "STR",
		executor: func(operand byte) (exitVM bool, incrementPC bool) {
			ram[operand] = registerA
			return false, true
		},
	},
	5: opCodeDescriptor{
		mnemonic: "LDI",
		executor: func(operand byte) (exitVM bool, incrementPC bool) {
			registerA = operand
			return false, true
		},
	},
	6: opCodeDescriptor{
		mnemonic: "JMP",
		executor: func(operand byte) (exitVM bool, incrementPC bool) {
			pc = operand
			return false, false
		},
	},
	7: opCodeDescriptor{
		mnemonic: "JC",
		executor: func(operand byte) (exitVM bool, incrementPC bool) {
			if flagCarryRegister {
				pc = operand
				return false, false
			}
			return false, true
		},
	},
	8: opCodeDescriptor{
		mnemonic: "JZ",
		executor: func(operand byte) (exitVM bool, incrementPC bool) {
			if flagZeroRegister {
				pc = operand
				return false, false
			}
			return false, true
		},
	},
	// Other OP codes are reserved.
	14: opCodeDescriptor{
		mnemonic: "OUT",
		executor: func(_ byte) (exitVM bool, incrementPC bool) {
			registerOut = registerA
			return false, true
		},
	},
	15: opCodeDescriptor{
		mnemonic: "HLT",
		executor: func(_ byte) (exitVM bool, incrementPC bool) { return true, true },
	},
}

var (
	// General purpose register. In the real computer there is also a register B, but we don't need
	// it in the VM.
	registerA byte
	// The output register.
	registerOut byte
	// Indicates the result of the last operation.
	flagZeroRegister  bool
	flagCarryRegister bool
	// The working memory for the VM.
	ram [16]byte
	// Program counter.
	pc byte
)

func runVM(stepDelay float64) {
	var sleepTime time.Duration
	if stepDelay != 0 {
		parsedSleepTime, _ := time.ParseDuration(fmt.Sprintf("%vs", stepDelay))
		sleepTime = parsedSleepTime
	}
	for {
		clearScreen()
		printState()
		if stepDelay == 0 {
			fmt.Printf("Press enter to advance program.")
			fmt.Scanln()
		} else {
			time.Sleep(sleepTime)
		}

		instruction := ram[pc]
		// Op codes are defined in the first 4 bits, shifting the instruction by 4 bits to the right
		// gives us the op code only without arguments.
		opCode := instruction >> 4
		opOperand := clearMsb(instruction)

		exitVM, incrementPC := opCodes[opCode].executor(opOperand)

		if incrementPC {
			pc++
		}

		if exitVM {
			break
		}
	}

	clearScreen()
	printState()
}

func assemble(hopAsm []string) []byte {
	mnemonicCodes := map[string]opCodeDescriptor{} // reverse lookup map from mnemonic->OPCODE
	for opCode, opCodeDesc := range opCodes {
		opCodeDesc.opCode = opCode
		mnemonicCodes[opCodeDesc.mnemonic] = opCodeDesc
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

		if opCodeDesc, isOpCode := mnemonicCodes[asmCode]; isOpCode {
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
		clockSpeed := float64(0)
		if len(os.Args) == 4 {
			hz, err := strconv.ParseFloat(os.Args[3], 64)
			if err != nil {
				panic(err)
			}
			clockSpeed = 1 / hz
			fmt.Println(clockSpeed)
		}
		runVM(clockSpeed)
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
- show the actual command in text form of the command the PC counter is pointing to
*/
