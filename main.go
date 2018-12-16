package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Operations use the first 4 bits as the Op Code. The last 4 bits are operation specific.
const (
	// A no operation.
	opNoop byte = 0
	// Loads contents of ADDR to A.
	// 0 0 0 1 A D D R
	opLoda byte = 1
	// Loads contents of ADDR to B.
	// 0 0 1 0 A D D R
	opLodb byte = 2
	// Increment value of ADDR.
	// 0 0 1 1 A D D R
	opIncr byte = 3
	// Subtracts B from A and stores in ADDR (ram[ADDR]=A-B).
	// 0 1 0 0 A D D R
	opSubt byte = 4
	// Move PC to ADDR if flag register is flagPositive.
	// 0 1 0 1 A D D R
	opJmpp byte = 5
	// Halt the VM.
	// 0 1 1 0 _ _ _ _
	opHalt byte = 6
)

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

func updateFlagRegister(result byte) {
	if result == 0 {
		flagRegister = flagZero
	} else if isBitSet(result, 8) {
		flagRegister = flagNegative
	} else {
		flagRegister = flagPositive
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
		opArgument := clearMsb(instruction)

		doExit := false
		incrementPc := true
		switch opCode {
		case opLoda:
			registerA = ram[opArgument]
			updateFlagRegister(registerA)

		case opLodb:
			registerB = ram[opArgument]
			updateFlagRegister(registerB)

		case opIncr:
			ram[opArgument]++
			updateFlagRegister(0)

		case opSubt:
			result := registerA - registerB
			ram[opArgument] = result
			updateFlagRegister(result)

		case opJmpp:
			if flagRegister == flagPositive {
				incrementPc = false
				pc = opArgument
				updateFlagRegister(pc)
			}

		case opHalt:
			doExit = true

		default:
			panic("Unknown instruction")
		}

		// JMP ops can manually set the pc
		if incrementPc {
			pc++
		}

		if doExit {
			break
		}
	}

	clearScreen()
	printState()
}

/*
TODO:
- control clock with a command line (maybe even have a manual step)
- show contents of all registers in binary and hex mode
- show ram +/- a couple of locations around the current PC counter
- show the actual command in text form of the command the PC counter is pointing to
*/

func clearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}
