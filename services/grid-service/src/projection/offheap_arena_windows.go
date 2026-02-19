//go:build windows
// +build windows

// File: services/grid-service/src/projection/offheap_arena_windows.go
package projection

import "unsafe"

type OffHeapArena struct {
	data []byte
}

func (a *OffHeapArena) Bytes() []byte {
	return a.data
}

func (a *OffHeapArena) Float64Slice() []float64 {
	if len(a.data) == 0 {
		return nil
	}
	count := len(a.data) / 8
	return unsafe.Slice((*float64)(unsafe.Pointer(&a.data[0])), count)
}

func (a *OffHeapArena) Close() error {
	a.data = nil
	return nil
}

type ArenaManager struct {
	maxCap int
}

func NewArenaManager(maxCap int) *ArenaManager {
	return &ArenaManager{maxCap: maxCap}
}

func (m *ArenaManager) Acquire(capacity int) (*OffHeapArena, error) {
	return &OffHeapArena{
		data: make([]byte, capacity),
	}, nil
}

func (m *ArenaManager) Release(a *OffHeapArena) {
	// no-op
}

func (a *OffHeapArena) WriteBytes(src []byte) {
	if len(src) > cap(a.data) {
		panic("offheap arena overflow (windows stub)")
	}
	a.data = a.data[:len(src)]
	copy(a.data, src)
}

func (a *OffHeapArena) WriteFloat64Slice(src []float64) {
	byteLen := len(src) * 8
	if byteLen > cap(a.data) {
		panic("offheap arena overflow (windows stub)")
	}
	a.data = a.data[:byteLen]
	srcBytes := unsafe.Slice((*byte)(unsafe.Pointer(&src[0])), byteLen)
	copy(a.data, srcBytes)
}
