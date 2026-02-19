package ipc

import (
	"sync"
)

// bufferPool recycles bytes.Buffers to reduce GC pressure.
var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// SerializeRecord serializes a record and returns the bytes plus a release function.
// The caller MUST call the release function when done with the bytes to return the buffer to the pool.
func SerializeRecord(record arrow.Record) ([]byte, func(), error) {
	// 1. Get buffer from pool
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()

	// 2. Create IPC Writer
	writer := ipc.NewWriter(buf, ipc.WithSchema(record.Schema()), ipc.WithAllocator(memory.DefaultAllocator))
	
	// 3. Write
	if err := writer.Write(record); err != nil {
		writer.Close()
		bufferPool.Put(buf) // Return on error
		return nil, nil, fmt.Errorf("failed to write record: %w", err)
	}

	if err := writer.Close(); err != nil {
		bufferPool.Put(buf)
		return nil, nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// 4. Return bytes and the cleanup closure
	// Note: buf.Bytes() is valid until Reset() or write.
	// The caller must process these bytes (e.g., stream.Send) before calling release.
	return buf.Bytes(), func() { bufferPool.Put(buf) }, nil
}
