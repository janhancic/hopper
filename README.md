# Hopper
A collection of learning projects that I'm working on with the goal of acquiring a better
understanding of how computers work at the lowest level.

My ultimate plan is to build an 8bit breadboard computer. Most of the things here are heavily
inspired by [Ben Eater](https://eater.net/)'s [8bit computer](https://eater.net/8bit).

The current plan is roughly as follows:

- build a VM that can run the same programs as Ben's computer - DONE
- build a hardware emulator for Ben's computer - IN PROGRESS
- come up with & implement some sort of improvement to Ben's architecture (most likely expand the
available memory
- design the computer in some sort of circuit design&simulate software (most likely in
[Logisim](http://www.cburch.com/logisim/) or [Digital](https://github.com/hneemann/Digital))
- build the computer on a real breadboard
- explore [PCBs](https://en.wikipedia.org/wiki/Printed_circuit_board) and
[FPGAs](https://en.wikipedia.org/wiki/Field-programmable_gate_array)

Considering I am learning as I go I expect the above to change.

## Goals
As mentioned the goal is to learn as much as I can about how computers work on the lowest level.
Ultimately I want to have a physical computer of my own. It will be based on existing architectures,
but with some improvements or modifications of my own so I can call it my own.

## What is in the repository
This repository contains all the source code I will write for these projects and any other files I
deem interesting and/or necessary for understanding of what's going on.

The main programming language I'm using is [Go](https://golang.org/), that is simply because I'm
working with it most at the moment and I feel comfortable enough using it. There are several
programs in this repository, each is contained in its own folder and contains a README that explains
how to use it. There is an overview of all of them below.

## Programs in this repository

### Hopper Virtual Machine
Contained in the `vm` folder. An 8bit Virtual Machine that virtualises a Hopper computer.

### Hopper Hardware Emulator
Contained in the `emulator` folder. An 8bit emulator that emulates a Hopper computer.

### Hopper Assembler
Contained in the `assembler` folder. A simple assembler that assembles .hop files into binary files
that can be run with the VM or the emulator.

## Hopper Architecture
At the moment Hopper uses the same architecture as [Ben Eater](https://eater.net/)'s
[8bit computer](https://eater.net/8bit). In the future this won't be the case and I will provide a
detailed document explaining it at that point.

## Running
The easiest way to play around with Hopper is to clone this repository and then refer to each
binary's README for instructions on how to run it.

## Notes on code
A lot of the code is very simplistic and lacks error handling. This is on purpose. The project is
not meant to be used for anything serious. It is purely a learning tool.
