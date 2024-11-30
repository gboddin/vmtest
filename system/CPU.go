package system

import (
	"encoding/binary"
	"log"
	"os"
)

type CPU struct {
	ProgramCounter uint16
	StackPointer   uint8
	Accumulator    uint8
	IndexRegisterX uint8
	IndexRegisterY uint8
	Status         uint8
	Memory         *Memory
	Bus            *Memory
}

func (cpu *CPU) SetFlag(flag uint8) {
	cpu.Status |= flag
}

func (cpu *CPU) UnsetFlag(flag uint8) {
	cpu.Status &= ^flag
}

func (cpu *CPU) GetFlag(flag uint8) bool {
	return (cpu.Status & flag) != 0
}

func (cpu *CPU) Debug() {
	//log.Printf("Accumulator: %#02x, RX: %#02x, RY: %#02x, Status: %08b", cpu.Accumulator, cpu.IndexRegisterX, cpu.IndexRegisterY, cpu.Status)
	//time.Sleep(100 * time.Millisecond)
}

func (cpu *CPU) Step() {
	//log.Printf("CPU PC: %#04x", cpu.ProgramCounter)
	defer cpu.Debug()
	opcode := cpu.Memory.Read(cpu.ProgramCounter)
	//log.Printf("OPCode: %#02x", opcode)
	// Do our thing
	switch opcode {
	case OP_INX:
		cpu.IndexRegisterX++
	case OP_INY:
		cpu.IndexRegisterY++
	case OP_LDA_BYTE:
		cpu.ProgramCounter++
		cpu.Accumulator = cpu.Memory.Read(cpu.ProgramCounter)
	case OP_LDA_ABS_X:
		cpu.ProgramCounter++
		absAddr := binary.LittleEndian.Uint16(cpu.Memory.content[cpu.ProgramCounter:])
		cpu.ProgramCounter += 2
		cpu.Accumulator = cpu.Memory.Read(absAddr + uint16(cpu.IndexRegisterX))
		return
	case OP_STA_ABS:
		cpu.ProgramCounter++
		absAddr := binary.LittleEndian.Uint16(cpu.Memory.content[cpu.ProgramCounter:])
		cpu.ProgramCounter += 2
		cpu.Memory.Write(absAddr, cpu.Accumulator)
		return
	case OP_LDX_BYTE:
		cpu.ProgramCounter++
		cpu.IndexRegisterX = cpu.Memory.Read(cpu.ProgramCounter)
	case OP_LDY_BYTE:
		cpu.ProgramCounter++
		cpu.IndexRegisterY = cpu.Memory.Read(cpu.ProgramCounter)
	case OP_JMP_ABS:
		cpu.ProgramCounter++
		cpu.ProgramCounter = binary.LittleEndian.Uint16(cpu.Memory.content[cpu.ProgramCounter:])
		// Returning right away, we don't need to touch PC anymore
		return
	case OP_CPX:
		cpu.ProgramCounter++
		toCompare := cpu.Memory.Read(cpu.ProgramCounter)
		if cpu.IndexRegisterX < toCompare {
			cpu.UnsetFlag(ZeroFlag)
			cpu.UnsetFlag(CarryFlag)
		}
		if cpu.IndexRegisterX > toCompare {
			cpu.UnsetFlag(ZeroFlag)
			cpu.SetFlag(CarryFlag)
		}
		if cpu.IndexRegisterX == toCompare {
			cpu.SetFlag(ZeroFlag)
			cpu.SetFlag(CarryFlag)
			cpu.UnsetFlag(Negative)
		}
	case OP_CMP:
		cpu.ProgramCounter++
		toCompare := cpu.Memory.Read(cpu.ProgramCounter)
		if cpu.Accumulator < toCompare {
			cpu.UnsetFlag(ZeroFlag)
			cpu.UnsetFlag(CarryFlag)
		}
		if cpu.Accumulator > toCompare {
			cpu.UnsetFlag(ZeroFlag)
			cpu.SetFlag(CarryFlag)
		}
		if cpu.Accumulator == toCompare {
			cpu.SetFlag(ZeroFlag)
			cpu.SetFlag(CarryFlag)
			cpu.UnsetFlag(Negative)
		}
	case OP_BNE:
		cpu.ProgramCounter++
		if cpu.GetFlag(ZeroFlag) {
			cpu.ProgramCounter++
			return
		}
		offset := int8(cpu.Memory.Read(cpu.ProgramCounter)) + 1
		cpu.ProgramCounter = uint16(int16(cpu.ProgramCounter) + int16(offset))
		return
	case OP_BEQ:
		cpu.ProgramCounter++
		if !cpu.GetFlag(ZeroFlag) {
			cpu.ProgramCounter++
			return
		}
		offset := int8(cpu.Memory.Read(cpu.ProgramCounter)) + 1
		cpu.ProgramCounter = uint16(int16(cpu.ProgramCounter) + int16(offset))
		return
	case OP_HALT:
		os.Exit(0)
	default:
		log.Fatalf("Unknown opcode: %#02x\n", opcode)
	}
	// Default increase PC
	cpu.ProgramCounter++
}

func (cpu *CPU) Init(mem *Memory) {
	cpu.Memory = mem
}

func (cpu *CPU) Reset() {
	cpu.ProgramCounter = binary.LittleEndian.Uint16(cpu.Memory.content[0xfffc:])
}

func (cpu *CPU) SetResetVector(address uint16) {
	binary.LittleEndian.PutUint16(cpu.Memory.content[0xfffc:], address)
}

const (
	CarryFlag uint8 = 1 << iota
	ZeroFlag
	InterruptDisable
	DecimalMode
	BreakCommand
	Unused1
	Overflow
	Negative
)

const (
	OP_INX       uint8 = 0xe8
	OP_INY             = 0xc8
	OP_JMP_ABS         = 0x4c
	OP_LDA_BYTE        = 0xa9
	OP_LDX_BYTE        = 0xa2
	OP_LDY_BYTE        = 0xa0
	OP_LDA_ABS_X       = 0xbd
	OP_CPX             = 0xe0
	OP_CMP             = 0xc9
	OP_BNE             = 0xd0
	OP_BEQ             = 0xf0
	OP_STA_ABS         = 0x8d
	OP_JMP_REL         = 0x6c
	OP_HALT            = 0x00
)
