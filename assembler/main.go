package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/janhancic/hopper/opcodes"
)

func assemble(hopAsm []string) []byte {
	mnemonicCodes := map[string]opcodes.OpCodeDescriptor{} // reverse lookup map from mnemonic->OPCODE
	for opCode, opCodeDesc := range opcodes.OpCodes {
		opCodeDesc.OpCode = opCode
		mnemonicCodes[opCodeDesc.Mnemonic] = opCodeDesc
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
				machineInstruction = opCodeDesc.OpCode << 4 // move the OPCODE to MSB
			} else {
				opOperandAsInt, _ := strconv.ParseUint(asmOperand, 10, 8)
				opOperandAsByte := byte(opOperandAsInt)

				machineInstruction = opCodeDesc.OpCode << 4 // move the OPCODE to MSB
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
	if len(os.Args) != 2 {
		printUsageAndDie()
	}

	hopContents, _ := ioutil.ReadFile(os.Args[1])
	fmt.Println("Assembling ...")
	machineCode := assemble(strings.Split(string(hopContents), "\n"))
	ioutil.WriteFile(os.Args[1]+".bin", machineCode, 0644)
	fmt.Println("Done!")
}

func printUsageAndDie() {
	fmt.Println("Usage: assembler path/to/program.hop")
	os.Exit(1)
}
