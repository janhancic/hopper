# Hopper
A simplistic 8-bit Virtual Machine implemented in Go.

This is a learning exercise for me to learn how a basic VM works.

## Architecture
The VM implements the [Ben Eater](https://eater.net/)'s [8bit computer](https://eater.net/8bit) 
architecture. It can run the same exact programs that Ben's breadboard computer can.

## Running
Clone this repository and `cd` into its folder.

To assemble a .hop file (Hopper Assembly) into an executable run `go run *.go assemble program.hop`.

Then run the program with `go run *.go run program.hop.bin`. 

You can also pass the clock frequency (in Hz) with which the VM will run. This will make the program
advance automatically.

## Notes on code
A lot of the code is very simplistic and lacks error handling. This is on purpose. The project is
not meant to be used for anything serious. It is purely a learning tool. The subject here is VMs,
not proper file error handling etc.

I've also tried to limit the number of files so that the code is more approachable for someone just
having a curious glance at it. My goal was for someone to be able to read `main.go` and understand
how everything works.
