package system

type Computer struct {
	ResetVector uint16
	CPU         *CPU
	Memory      *Memory
	Devices     []MemoryMappedDevice
}

func NewComputer(resetVector uint16, devices ...MemoryMappedDevice) *Computer {
	comp := &Computer{
		ResetVector: resetVector,
		Devices:     devices,
		CPU:         &CPU{},
		Memory: &Memory{
			Devices: make(map[uint16]MemoryMappedDevice),
		},
	}
	for _, device := range devices {
		comp.Memory.MapDevice(device)
	}
	return comp.Init()
}

func (computer *Computer) Init() *Computer {
	// Give CPU access to memory
	computer.CPU.Init(computer.Memory)
	// Set the ResetVector for the CPU in RAM
	computer.CPU.SetResetVector(computer.ResetVector)
	// Trigger CPU reset
	computer.CPU.Reset()
	return computer
}

func (computer *Computer) Step() {
	computer.CPU.Step()
	return
}

func (computer *Computer) LoadRom(addr uint16, data []byte) {
	computer.Memory.Load(addr, data)
}
