package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Operations use the first 4 bits as the Op Code. The last 4 bits are operation specific.
const (
	// 8 7 6 5 4 3 2 1
	// 0 0 0 0 _ _ _ _
	// Adds values of Register A and Register B together and stores them in Register B.
	// TODO: use LSBs to control which register the sum goes in
	opAdd byte = 0
	// 8 7 6 5 4 3 2 1
	// 0 0 0 1 V A L U
	// Loads V A L U into register A.
	opLoadA byte = 1
	// 8 7 6 5 4 3 2 1
	// 0 0 1 0 V A L U
	// Loads V A L U into register B.
	opLoadB byte = 2
	// 8 7 6 5 4 3 2 1
	// 0 0 1 1 _ _ _ _
	// Halt the program.
	opHalt byte = 3
)

// Useful for clearing the op code from an instruction.
const msbMask = 0xF0 // 11110000

var (
	// TODO: use array for registers
	registerA byte
	registerB byte

	ram [255]byte

	pc byte
)

// convert a byte to a binary string representation
func byteToString(b byte) string {
	return fmt.Sprintf("%08b", b)
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
	fmt.Printf("Register A: %s (%d)\n", byteToString(registerA), registerA)
	fmt.Printf("Register B: %s (%d)\n", byteToString(registerB), registerB)
	fmt.Printf("PC:         %s (%d)\n", byteToString(pc), pc)
}

func main() {
	// sample program that does 2 sums
	ram[0] = stringToByte("00010011") // LOAD A 0011 (3)
	ram[1] = stringToByte("00100100") // LOAD B 0100 (4)
	ram[2] = stringToByte("00000000") // SUM
	ram[3] = stringToByte("00010001") // LOAD A 0001 (1)
	ram[4] = stringToByte("00000000") // SUM // should be 8
	ram[5] = stringToByte("00110000") // HALT

	for {
		clearScreen()
		fmt.Println("Current state:")
		printState()
		fmt.Printf("Press enter to advance program.")
		fmt.Scanln()

		instruction := ram[pc]
		// Op codes are defined in the first 4 bits, shifting the instruction by 4 bits to the right
		// gives us the op code only without arguments.
		opCode := instruction >> 4

		doExit := false
		switch opCode {
		case opLoadA:
			registerA = clearMsb(instruction)

		case opLoadB:
			registerB = clearMsb(instruction)

		case opAdd:
			registerB = registerA + registerB

		case opHalt:
			doExit = true

		default:
			panic("Unknown instruction")
		}

		pc++

		if doExit {
			break
		}
	}

	clearScreen()
	fmt.Println("Program ended. End state:")
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
