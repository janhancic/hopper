# Hopper Assembler
A very simplistic assembler that does a direct translation of mnemonics and their operands into
binary.

## Running
```shell
go build
./assemble path/to/program.hop
```

This will generate a `program.hop.bin` in the same folder as the `.hop` file. You can run the
assembled binary with the vm or the emulator programs.
