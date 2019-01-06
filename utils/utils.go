package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Useful for clearing the op code from an instruction.
const msbMask = 0xF0 // 11110000

// convert a byte to a binary string representation
func ByteToString(b byte) string {
	return fmt.Sprintf("%08b", b)
}

// convert half a byte to a binary string representation
func ByteToNibble(b byte) string {
	return fmt.Sprintf("%04b", b)
}

// convert a string containing binary into a byte
func StringToByte(s string) byte {
	b, err := strconv.ParseUint(s, 2, 8)
	if err != nil {
		panic(fmt.Sprintf("Invalid binary sequence in string '%v': %v", s, err))
	}
	return byte(b)
}

// determines if the n-th bit (from the right) is set
func IsBitSet(b byte, n uint8) bool {
	return (b & (1 << n)) > 0
}

// sets the first 4 (MSB) bits to 0
func ClearMsb(b byte) byte {
	return b &^ msbMask
}

// adds two bytes together and returns the result and if an overflow occurred (carry)
func ByteAdder(a, b byte) (result byte, carry bool) {
	c := a + b
	if (c > a) == (b > 0) {
		return c, false
	}
	return c, true
}

// subtracts two bytes and returns the result and if an overflow occurred (carry)
func ByteSubtractor(a, b byte) (result byte, carry bool) {
	c := a - b
	if (c < a) == (b > 0) {
		return c, true
	}
	return c, false
}

func ClearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}
