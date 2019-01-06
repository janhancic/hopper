package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/janhancic/hopper/opcodes"
	"github.com/janhancic/hopper/utils"
)

func init() {
	// NOP
	opcodes.OpCodes[0].Executor = func(operand byte) (bool, bool) { return false, true }
	// LDA
	opcodes.OpCodes[1].Executor = func(operand byte) (exitVM bool, incrementPC bool) {
		registerA = ram[operand]
		return false, true
	}
	// ADD
	opcodes.OpCodes[2].Executor = func(operand byte) (exitVM bool, incrementPC bool) {
		result, carry := utils.ByteAdder(registerA, ram[operand])

		registerA = result
		flagCarryRegister = carry
		flagZeroRegister = result == 0

		return false, true
	}
	// SUB
	opcodes.OpCodes[3].Executor = func(operand byte) (exitVM bool, incrementPC bool) {
		result, carry := utils.ByteSubtractor(registerA, ram[operand])

		registerA = result
		flagCarryRegister = carry
		flagZeroRegister = result == 0

		return false, true
	}
	// STR
	opcodes.OpCodes[4].Executor = func(operand byte) (exitVM bool, incrementPC bool) {
		ram[operand] = registerA
		return false, true
	}
	// LDI
	opcodes.OpCodes[5].Executor = func(operand byte) (exitVM bool, incrementPC bool) {
		registerA = operand
		return false, true
	}
	// JMP
	opcodes.OpCodes[6].Executor = func(operand byte) (exitVM bool, incrementPC bool) {
		pc = operand
		return false, false
	}
	// JC
	opcodes.OpCodes[7].Executor = func(operand byte) (exitVM bool, incrementPC bool) {
		if flagCarryRegister {
			pc = operand
			return false, false
		}
		return false, true
	}
	// JZ
	opcodes.OpCodes[8].Executor = func(operand byte) (exitVM bool, incrementPC bool) {
		if flagZeroRegister {
			pc = operand
			return false, false
		}
		return false, true
	}
	// OUT
	opcodes.OpCodes[14].Executor = func(_ byte) (exitVM bool, incrementPC bool) {
		registerOut = registerA
		return false, true
	}
	// HLT
	opcodes.OpCodes[15].Executor = func(_ byte) (exitVM bool, incrementPC bool) { return true, true }
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
		utils.ClearScreen()
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
		opOperand := utils.ClearMsb(instruction)

		exitVM, incrementPC := opcodes.OpCodes[opCode].Executor(opOperand)

		if incrementPC {
			pc++
		}

		if exitVM {
			break
		}
	}

	utils.ClearScreen()
	printState()
}

func main() {
	if len(os.Args) < 2 {
		printUsageAndDie()
	}

	binContents, _ := ioutil.ReadFile(os.Args[1])
	copy(ram[:], binContents)
	clockSpeed := float64(0)
	if len(os.Args) == 3 {
		hz, err := strconv.ParseFloat(os.Args[2], 64)
		if err != nil {
			panic(err)
		}
		clockSpeed = 1 / hz
		fmt.Println(clockSpeed)
	}
	runVM(clockSpeed)

}

// prints the current state of the VM
func printState() {
	fmt.Printf("Register A:    %s (%d)\n", utils.ByteToString(registerA), registerA)
	fmt.Printf("Register Out:  %s (%d)\n", utils.ByteToString(registerOut), registerOut)
	fmt.Printf("Flag Zero:     %v\n", flagZeroRegister)
	fmt.Printf("Flag Carry:    %v\n", flagCarryRegister)
	fmt.Printf("PC:            %s (%d)\n", utils.ByteToNibble(pc), pc)
	fmt.Println("RAM:")
	for addr, val := range ram {
		pcIndicator := ""
		if byte(addr) == pc {
			pcIndicator = " <---"
		}
		fmt.Printf(
			"%02d: %s: %s%s\n",
			addr,
			utils.ByteToNibble(byte(addr)),
			utils.ByteToString(val),
			pcIndicator,
		)
	}
}

func printUsageAndDie() {
	fmt.Println("Usage: vm path/to/program.hop.bin or vm path/to/program.hop.bin 10")
	os.Exit(1)
}
