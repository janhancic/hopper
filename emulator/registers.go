package main

// RegisterA is the general purpose A register
var RegisterA byte = 56

// RegisterB is a register used by the computer for arithmetic operations
var RegisterB byte = 89

// ALU holds either the sum of A+B or the difference between A-B
var ALU byte = 213

// MAR holds the current memory addres
var MAR byte = 13

// PC is the program counter register
var PC byte = 25

// IR is the instructions register. MSB holds the OP code and the LSB holds the operand.
// Only LSB is connected to the BUS.
var IR byte = 255

// RegisterOut serves as the output register
var RegisterOut byte = 133

// FZ flag zero register
var FZ bool

// FC flag carry register
var FC bool
