package orchestration

import (
	"context"
	"testing"
	"errors"

	"quantatomai/grid-service/pkg/grid/v1"
	"quantatomai/grid-service/pkg/ipc"
	"connectrpc.com/connect"
	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockIPCClient
type MockIPCClient struct {
	mock.Mock
}

func (m *MockIPCClient) GetCalculation(ctx context.Context, planID string) (ipc.RecordReader, error) {
	args := m.Called(ctx, planID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(ipc.RecordReader), args.Error(1)
}

func (m *MockIPCClient) Close() error {
	return nil
}

// MockRecordReader
type MockRecordReader struct {
	mock.Mock
}

func (m *MockRecordReader) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRecordReader) Record() arrow.Record {
	args := m.Called()
	if rec := args.Get(0); rec != nil {
		return rec.(arrow.Record)
	}
	return nil
}

func (m *MockRecordReader) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRecordReader) Release() {
	m.Called()
}

// MockServerStream
type MockServerStream struct {
	mock.Mock
	connect.ServerStream[gridv1.GridChunk] // Embed to satisfy interface
}

func (m *MockServerStream) Send(msg *gridv1.GridChunk) error {
	args := m.Called(msg)
	return args.Error(0)
}

func TestQueryGrid_Streaming(t *testing.T) {
	// 1. Setup Arrow Data
	pool := memory.NewGoAllocator()
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "col1", Type: arrow.PrimitiveTypes.Int64},
		},
		nil,
	)
	b := array.NewRecordBuilder(pool, schema)
	defer b.Release()

	b.Field(0).(*array.Int64Builder).AppendValues([]int64{1, 2, 3}, nil)
	rec := b.NewRecord()
	// Note: We don't defer release here because the Mock will return it, 
	// and the service might Release it? serialized bytes? 
	// Actually RecordReader ownership: usually Reader owns the record returned by Record() until Next() is called again.
	// We'll let the test clean it up.
	defer rec.Release() 

	// 2. Setup Mocks
	mockClient := new(MockIPCClient)
	mockReader := new(MockRecordReader)
	mockStream := new(MockServerStream)

	ctx := context.Background()
	req := connect.NewRequest(&gridv1.QueryGridRequest{ViewId: "123"})

	// Expectation: GetCalculation called
	mockClient.On("GetCalculation", ctx, "PLAN-123").Return(mockReader, nil)

	// Expectation: Reader.Next() -> True (1st record)
	mockReader.On("Next").Return(true).Once()
	// Expectation: Reader.Record() -> returns our record
	mockReader.On("Record").Return(rec).Once()
	
	// Expectation: Stream.Send() -> called with serialized bytes
	// We verify the bytes are not empty
	mockStream.On("Send", mock.AnythingOfType("*grid.GridChunk")).Return(nil).Run(func(args mock.Arguments) {
		chunk := args.Get(0).(*gridv1.GridChunk)
		data := chunk.Data.(*gridv1.GridChunk_ArrowRecordBatch).ArrowRecordBatch
		assert.NotEmpty(t, data)
		// Ideally we would deserialize looking for [1,2,3] but for now checking non-empty is good progress.
	}).Once()

	// Expectation: Reader.Next() -> False (End of stream)
	mockReader.On("Next").Return(false).Once()

	// Expectation: Reader.Err() -> nil
	mockReader.On("Err").Return(nil)

	// Expectation: Reader.Release() -> called
	mockReader.On("Release").Return()

	// 3. Execute
	handler := NewGridQueryServiceHandler(mockClient)
	err := handler.QueryGrid(ctx, req, mockStream)

	// 4. Verify
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
	mockReader.AssertExpectations(t)
	mockStream.AssertExpectations(t)
}
