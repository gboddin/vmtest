package system

import (
	"encoding/hex"
	"fmt"
)

type Memory struct {
	content [64 * 1024]byte
	Devices map[uint16]MemoryMappedDevice
}

func NewMemory() *Memory {
	return &Memory{
		Devices: make(map[uint16]MemoryMappedDevice),
	}
}

func (mem *Memory) MapDevice(device MemoryMappedDevice) {
	start, end := device.GetMemoryRange()
	for i := start; i <= end; i++ {
		mem.Devices[i] = device
	}
}

func (mem *Memory) HexDump() {
	fmt.Println(hex.Dump(mem.content[:]))
}

func (mem *Memory) Write(addr uint16, data uint8) {
	if device, found := mem.Devices[addr]; found {
		device.Write(addr, data)
	}
	mem.content[addr] = data
}

func (mem *Memory) Read(addr uint16) uint8 {
	for _, h := range mem.Devices {
		if data, hooked := h.Read(addr); hooked {
			return data
		}
	}
	return mem.content[addr]
}

func (mem *Memory) Load(addr uint16, data []byte) {
	copy(mem.content[addr:], data)
}
