package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// opCodeDescriptor describe a Hopper OPCODE. It holds information for the Hopper assembler to know
// how to translate Hopper Assembly to OPCODEs.
type opCodeDescriptor struct {
	// The string (or ASM) version of this OPCODE.
	asm string
	// Indicates wether this OPCODE has any operands.
	noOperand bool
	// The function that will execute the OPCODE.
	executor func(operand byte) (result byte, exitVM bool, incrementPC bool)
}

// opCodes holds all supported OPCODEs in Hopper VM.
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
TODO: Possible instructions (* already implemented)
NOOP*
LOAD A ADDR*
LOAD B ADDR*
SAVE A ADDR
SAVE B ADDR
ADD  ADDR ; Add ADDR=A+B
SUB  ADDR ; Subtract ADDR=A-B*
JMP  ADDR ; Jump to ADDR
JMPZ ADDR ; Jump to ADDR if flagRegister is zero
JMPP ADDR ; Jump to ADDR if flagRegister is positive*
JMPN ADDR ; Jump to ADDR if flagRegister is negative
INC  ADDR*
HALT*
*/

// Flags for the flagRegister.
const (
	flagZero     = 0
	flagPositive = 1
	flagNegative = 2
)

// Useful for clearing the op code from an instruction.
const msbMask = 0xF0 // 11110000

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

// convert a byte to a binary string representation
func byteToString(b byte) string {
	return fmt.Sprintf("%08b", b)
}

// convert half a byte to a binary string representation
func byteToNibble(b byte) string {
	return fmt.Sprintf("%04b", b)
}

// convert a string containing binary into a byte
func stringToByte(s string) byte {
	b, err := strconv.ParseUint(s, 2, 8)
	if err != nil {
		panic(fmt.Sprintf("Invalid binary sequence in string '%v': %v", s, err))
	}
	return byte(b)
}

// determines if the n-th bit (from the right) is set
func isBitSet(b byte, n uint8) bool {
	return (b & (1 << n)) > 0
}

// sets the first 4 (MSB) bits to 0
func clearMsb(b byte) byte {
	return b &^ msbMask
}

// prints the current state of the VM
func printState() {
	fmt.Printf("Register A:    %s (%d)\n", byteToString(registerA), registerA)
	fmt.Printf("Register B:    %s (%d)\n", byteToString(registerB), registerB)
	fmt.Printf("Flag Register: %s (%d)\n", byteToString(flagRegister), flagRegister)
	fmt.Printf("PC:            %s (%d)\n", byteToString(pc), pc)
	fmt.Println("RAM:")
	for addr, val := range ram {
		pcIndicator := ""
		if byte(addr) == pc {
			pcIndicator = " <---"
		}
		fmt.Printf(
			"%02d: %s: %s%s\n",
			addr+1,
			byteToNibble(byte(addr)),
			byteToString(val),
			pcIndicator,
		)
	}
}

func main() {
	// the count_to_three.hop
	ram[0] = stringToByte("00111100")
	ram[1] = stringToByte("00011101")
	ram[2] = stringToByte("00101100")
	ram[3] = stringToByte("01001110")
	ram[4] = stringToByte("01010000")
	ram[5] = stringToByte("01100000")
	ram[6] = stringToByte("00000000")
	ram[7] = stringToByte("00000000")
	ram[8] = stringToByte("00000000")
	ram[9] = stringToByte("00000000")
	ram[10] = stringToByte("00000000")
	ram[11] = stringToByte("00000000")
	ram[12] = stringToByte("00000000")
	ram[13] = stringToByte("00000011")
	ram[14] = stringToByte("00000000")
	ram[15] = stringToByte("00000000")

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

/*
TODO:
- control clock with a command line (maybe even have a manual step)
- show the actual command in text form of the command the PC counter is pointing to
*/

func clearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}
