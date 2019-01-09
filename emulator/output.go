package main

import (
	"fmt"
	"strings"

	"github.com/janhancic/hopper/utils"
)

const outputTpl = `
HLT: halt
MI:  memory address in
RI:  ram in
RO:  ram out
IO:  instruction register out
II:  instruction register in
AI:  A register in
AO:  A register out
EO:  ALU out
SU:  substract
BI:  B register in
OI:  output in
CE:  enable counter
CO:  program counter out
J:   jump
---
FZ: zero flag
FC: carry flag

    HLT[0]                                                             BUS
                                                            [0][0][0][0][0][0][0][0]
                                                             |  |  |  |  |__|__|__|<___J[0]___CO[0]__>[0][0][0][0] CE[0] (PC)           (FZ)[0]  (FC)[0]
                                                             |  |  |  |  |  |  |  |
                       (MAR)  [0][0][0][0]<_________MI[0]____|..|..|..|..|__|__|__|
                                                             |  |  |  |  |  |  |  |
            (RAM) [0][0][0][0][0][0][0][0]<__RO[0]__RI[0]___>|__|__|__|__|__|__|__|
                                                             |  |  |  |  |  |  |  |
                                                             |__|__|__|__|__|__|__|<___AI[0]__AO[0]__>[0][0][0][0][0][0][0][0] (A)
                                                             |  |  |  |  |  |  |  |                    |  |  |  |  |  |  |  |
                                                             |__|__|__|__|__|__|__|<__________EO[0]___[0][0][0][0][0][0][0][0]___SU[0]  (ALU)
                                                             |  |  |  |  |  |  |  |                    |  |  |  |  |  |  |  |
                                                             |__|__|__|__|__|__|__|____BI[0]_________>[0][0][0][0][0][0][0][0] (B)
                                                             |  |  |  |  |  |  |  |
                                                             |__|__|__|__|__|__|__|___________OI[0]__>[0][0][0][0][0][0][0][0]: 123 (OUT)
                                                             |  |  |  |  |  |  |  |
             (IR) [0][0][0][0][0][0][0][0]<__IO[0]__II[0]___>|__|__|__|__|__|__|__|
                                                             |  |  |  |  |  |  |  |`

func replaceSignalWithValue(out string, signalName string, signalState bool) string {
	if signalState {
		return strings.Replace(out, signalName+"[0]", signalName+"[1]", 1)
	} else {
		return out
	}
}

func displayState() {
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
	utils.ClearScreen()
	fmt.Println(output)
}
