package main

import (
	"github.com/gboddin/vmtest/gbocomputer"
	"github.com/gboddin/vmtest/system"
	"io"
	"os"
)

// ResetVector Arbitrary decision for our computer architecture
const ResetVector uint16 = 0x8000

func main() {
	computer := system.NewComputer(ResetVector, gbocomputer.StdOutput{})

	romFile, err := os.Open("test.rom")
	if err != nil {
		panic(err)
	}
	romBytes, err := io.ReadAll(romFile)
	if err != nil {
		panic(err)
	}
	// Load ROM
	computer.LoadRom(ResetVector, romBytes)
	for {
		computer.Step()
	}
}
