# Hopper Emulator
This is currently work in progress. The goal is to create an emulator that will run the Hopper
binaries and produce the same results as running them in the VM would.
The difference between this and the VM is that in the emulator tries to mimic how the actual
hardware operates. For example it doesn't execute the whole instruction in one clock cycle, but
rather uses micro operations to perform one instructions across multiple clock cycles.

## Running
Use binaries assembled using the Hopper assembler.

```shell
go build
./emulator path/to/program.hop.bin
```

You can also pass the clock frequency (in Hz) with which the VM will run. This will make the VM run
automatically using the specified frequency.
