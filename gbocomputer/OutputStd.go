package gbocomputer

import (
	"os"
)

type StdOutput struct {
}

func (s StdOutput) GetMemoryRange() (start uint16, end uint16) {
	return 0xc000, 0xc000
}

func (s StdOutput) Read(address uint16) (byte, bool) {
	return 0, false
}

func (s StdOutput) Write(address uint16, data byte) bool {
	if address != 0xc000 {
		return false
	}
	_, err := os.Stdout.Write([]byte{data})
	if err != nil {
		panic(err)
	}
	return true
}
