package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/janhancic/hopper/utils"
)

const outputTpl = `
    HLT[0]                                             (BUS)[0][0][0][0][0][0][0][0]
                                                             |  |  |  |  |__|__|__|<___J[0]___CO[0]__>(PC)[0][0][0][0] CE[0]            (FZ)[0]  (FC)[0]
                                                             |  |  |  |  |  |  |  |
                         (MAR)[0][0][0][0]<_________MI[0]____|..|..|..|..|__|__|__|
                                                             |  |  |  |  |  |  |  |
             (RAM)[0][0][0][0][0][0][0][0]<__RO[0]__RI[0]___>|__|__|__|__|__|__|__|
                                                             |  |  |  |  |  |  |  |
                                                             |__|__|__|__|__|__|__|<___AI[0]___AO[0]___>(A)[0][0][0][0][0][0][0][0]
                                                             |  |  |  |  |  |  |  |                         |  |  |  |  |  |  |  |
                                                             |__|__|__|__|__|__|__|<___EO[0]__________(ALU)[0][0][0][0][0][0][0][0]___SU[0]
                                                             |  |  |  |  |  |  |  |                         |  |  |  |  |  |  |  |
                                                             |__|__|__|__|__|__|__|____BI[0]___________>(B)[0][0][0][0][0][0][0][0]
                                                             |  |  |  |  |  |  |  |
                                                             |__|__|__|__|__|__|__|____OI[0]__>(OUT)[0][0][0][0][0][0][0][0]: $OUT 
                                                             |  |  |  |  |  |  |  |
              (IR)[0][0][0][0][0][0][0][0]<__IO[0]__II[0]___>|__|__|__|__|__|__|__|
                                                             |  |  |  |  |  |  |  |
`

func replaceSignalWithValue(out string, signalName string, signalState bool) string {
	if signalState {
		return strings.Replace(out, signalName+"[0]", signalName+"[1]", 1)
	}
	return out
}

func replaceBitRegisterWithValue(out string, registerName string, registerState bool) string {
	if registerState {
		return strings.Replace(out, "("+registerName+")[0]", "("+registerName+")[1]", 1)
	}
	return out
}

func replaceByteRegisterWithValue(out string, registerName string, registerValue byte) string {
	if registerValue == 0 {
		return out
	}

	valueAsStrings := strings.Split(utils.ByteToString(registerValue), "")
	return strings.Replace(
		out,
		"("+registerName+")[0][0][0][0][0][0][0][0]",
		fmt.Sprintf(
			"(%s)[%s][%s][%s][%s][%s][%s][%s][%s]",
			registerName,
			valueAsStrings[0],
			valueAsStrings[1],
			valueAsStrings[2],
			valueAsStrings[3],
			valueAsStrings[4],
			valueAsStrings[5],
			valueAsStrings[6],
			valueAsStrings[7],
		),
		1,
	)
}

func replaceNibbleRegisterWithValue(out string, registerName string, registerValue byte) string {
	if registerValue == 0 {
		return out
	}

	valueAsStrings := strings.Split(utils.ByteToNibble(registerValue), "")
	return strings.Replace(
		out,
		"("+registerName+")[0][0][0][0]",
		fmt.Sprintf(
			"(%s)[%s][%s][%s][%s]",
			registerName,
			valueAsStrings[0],
			valueAsStrings[1],
			valueAsStrings[2],
			valueAsStrings[3],
		),
		1,
	)
}

func displayState() {
	// Output all signal indicators as SGN[1] if they are on
	output := replaceSignalWithValue(outputTpl, "HLT", HLT)
	output = replaceSignalWithValue(output, "MI", MI)
	output = replaceSignalWithValue(output, "RI", RI)
	output = replaceSignalWithValue(output, "RO", RO)
	output = replaceSignalWithValue(output, "IO", IO)
	output = replaceSignalWithValue(output, "II", II)
	output = replaceSignalWithValue(output, "AI", AI)
	output = replaceSignalWithValue(output, "AO", AO)
	output = replaceSignalWithValue(output, "EO", EO)
	output = replaceSignalWithValue(output, "SU", SU)
	output = replaceSignalWithValue(output, "BI", BI)
	output = replaceSignalWithValue(output, "OI", OI)
	output = replaceSignalWithValue(output, "CE", CE)
	output = replaceSignalWithValue(output, "CO", CO)
	output = replaceSignalWithValue(output, "J", J)
	// Output all single bit registers as (R)[1] if they are set
	output = replaceBitRegisterWithValue(output, "FZ", FZ)
	output = replaceBitRegisterWithValue(output, "FC", FC)
	// Output all byte registers as (R)[1][0][1][0][1][0][1][0] according to their value
	output = replaceByteRegisterWithValue(output, "BUS", BUS) // also used for bus
	output = replaceByteRegisterWithValue(output, "A", RegisterA)
	output = replaceByteRegisterWithValue(output, "B", RegisterB)
	output = replaceByteRegisterWithValue(output, "ALU", ALU)
	output = replaceByteRegisterWithValue(output, "IR", IR)
	output = replaceByteRegisterWithValue(output, "OUT", RegisterOut)
	// Output all nibble registers as (R)[1][0][1][0] according to their value
	output = replaceNibbleRegisterWithValue(output, "MAR", MAR)
	output = replaceNibbleRegisterWithValue(output, "PC", PC)
	// Output the contents of OUT register in decimal
	// TODO: this should be modeled more like the binary to decimal hardware decoder
	output = strings.Replace(output, "$OUT", strconv.Itoa(int(RegisterOut)), 1)

	// TODO: output RAM

	utils.ClearScreen()
	fmt.Println(output)
}
