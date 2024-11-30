package system

type MemoryMappedDevice interface {
	GetMemoryRange() (start uint16, end uint16)
	Read(address uint16) (byte, bool)
	Write(address uint16, data byte) bool
}
