//go:build !windows
// +build !windows

// File: services/grid-service/src/projection/offheap_arena.go
package projection

import (
	"os"
	"sync"
	"syscall"
	"unsafe"
)

type OffHeapArena struct {
	ptr    unsafe.Pointer
	length int
	cap    int
	file   *os.File
}

func (a *OffHeapArena) Bytes() []byte {
	if a.ptr == nil || a.length == 0 {
		return nil
	}
	hdr := &struct {
		Data unsafe.Pointer
		Len  int
		Cap  int
	}{a.ptr, a.length, a.cap}
	return *(*[]byte)(unsafe.Pointer(hdr))
}

func (a *OffHeapArena) Float64Slice() []float64 {
	if a.ptr == nil || a.length == 0 {
		return nil
	}
	count := a.length / 8
	return unsafe.Slice((*float64)(a.ptr), count)
}

func (a *OffHeapArena) Close() error {
	if a.ptr == nil {
		return nil
	}
	err := syscall.Munmap(a.Bytes())
	a.ptr = nil
	a.length = 0
	a.cap = 0
	return err
}

type ArenaManager struct {
	mu     sync.Mutex
	pool   []*OffHeapArena
	maxCap int
}

func NewArenaManager(maxCap int) *ArenaManager {
	return &ArenaManager{
		maxCap: maxCap,
		pool:   make([]*OffHeapArena, 0, 16),
	}
}

func (m *ArenaManager) Acquire(capacity int) (*OffHeapArena, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, a := range m.pool {
		if a.cap >= capacity {
			m.pool[i] = m.pool[len(m.pool)-1]
			m.pool = m.pool[:len(m.pool)-1]
			a.length = 0
			return a, nil
		}
	}

	aligned := alignToPage(capacity)
	data, err := syscall.Mmap(
		-1,
		0,
		aligned,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_ANON|syscall.MAP_PRIVATE,
	)
	if err != nil {
		return nil, err
	}

	return &OffHeapArena{
		ptr:    unsafe.Pointer(&data[0]),
		length: 0,
		cap:    aligned,
	}, nil
}

func (m *ArenaManager) Release(a *OffHeapArena) {
	if a == nil || a.ptr == nil {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if a.cap <= m.maxCap {
		a.length = 0
		m.pool = append(m.pool, a)
	} else {
		_ = a.Close()
	}
}

func (a *OffHeapArena) WriteBytes(src []byte) {
	if len(src) > a.cap {
		panic("offheap arena overflow")
	}
	dst := a.Bytes()
	copy(dst, src)
	a.length = len(src)
}

func (a *OffHeapArena) WriteFloat64Slice(src []float64) {
	byteLen := len(src) * 8
	if byteLen > a.cap {
		panic("offheap arena overflow")
	}
	dst := a.Bytes()
	srcBytes := unsafe.Slice((*byte)(unsafe.Pointer(&src[0])), byteLen)
	copy(dst, srcBytes)
	a.length = byteLen
}

func alignToPage(n int) int {
	page := os.Getpagesize()
	if n%page == 0 {
		return n
	}
	return ((n / page) + 1) * page
}
