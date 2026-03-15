package orchestration

import (
	"context"
	gridv1 "quantatomai/grid-service/pkg/grid/v1"
	"quantatomai/grid-service/pkg/ipc"
	"testing"

	"connectrpc.com/connect"
	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
)

type benchmarkIPCClient struct {
	record arrow.Record
}

func (c *benchmarkIPCClient) GetCalculation(_ context.Context, _ string) (ipc.RecordReader, error) {
	return &benchmarkRecordReader{record: c.record}, nil
}

func (c *benchmarkIPCClient) Close() error {
	return nil
}

type benchmarkRecordReader struct {
	record arrow.Record
	read   bool
}

func (r *benchmarkRecordReader) Next() bool {
	if r.read {
		return false
	}

	r.read = true
	return true
}

func (r *benchmarkRecordReader) Record() arrow.Record {
	return r.record
}

func (r *benchmarkRecordReader) Err() error {
	return nil
}

func (r *benchmarkRecordReader) Release() {}

type benchmarkGridChunkSender struct {
	bytesSent int
}

func (s *benchmarkGridChunkSender) Send(msg *gridv1.GridChunk) error {
	if batch, ok := msg.Data.(*gridv1.GridChunk_ArrowRecordBatch); ok {
		s.bytesSent += len(batch.ArrowRecordBatch)
	}

	return nil
}

func BenchmarkGridQueryServiceSingleRecord(b *testing.B) {
	record := newBenchmarkRecord(b, 256)
	defer record.Release()

	handler := NewGridQueryServiceHandler(&benchmarkIPCClient{record: record})
	req := connect.NewRequest(&gridv1.QueryGridRequest{ViewId: "bench-single"})
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sender := &benchmarkGridChunkSender{}
		if err := handler.queryGrid(ctx, req, sender); err != nil {
			b.Fatalf("queryGrid failed: %v", err)
		}
	}
}

func BenchmarkGridQueryServiceLargeRecord(b *testing.B) {
	record := newBenchmarkRecord(b, 8192)
	defer record.Release()

	handler := NewGridQueryServiceHandler(&benchmarkIPCClient{record: record})
	req := connect.NewRequest(&gridv1.QueryGridRequest{ViewId: "bench-large"})
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sender := &benchmarkGridChunkSender{}
		if err := handler.queryGrid(ctx, req, sender); err != nil {
			b.Fatalf("queryGrid failed: %v", err)
		}
	}
}

func newBenchmarkRecord(b *testing.B, rowCount int) arrow.Record {
	b.Helper()

	pool := memory.NewGoAllocator()
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "account_id", Type: arrow.PrimitiveTypes.Int64},
		{Name: "value", Type: arrow.PrimitiveTypes.Float64},
	}, nil)

	builder := array.NewRecordBuilder(pool, schema)
	b.Cleanup(builder.Release)

	accountBuilder := builder.Field(0).(*array.Int64Builder)
	valueBuilder := builder.Field(1).(*array.Float64Builder)

	for index := 0; index < rowCount; index++ {
		accountBuilder.Append(int64(index + 1))
		valueBuilder.Append(float64(index) * 1.25)
	}

	return builder.NewRecord()
}
