package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Useful for clearing the op code from an instruction.
const msbMask = 0xF0 // 11110000

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

func clearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func printUsageAndDie() {
	fmt.Println("Usage: hopper run program.hop.bin or hopper assemble program.hop")
	os.Exit(1)
}
