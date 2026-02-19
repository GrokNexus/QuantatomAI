package orchestration

import (
	"context"
	"fmt"
	"io"

	"connectrpc.com/connect"
	"quantatomai/grid-service/pkg/grid/v1" // Assuming protoc gen works
	"quantatomai/grid-service/pkg/grid/v1/gridv1connect"
	mdfv1 "quantatomai/grid-service/pkg/mdf/v1"
	"quantatomai/grid-service/pkg/ipc"
)

type GridQueryServiceHandler struct {
	flightClient ipc.Client
}

func NewGridQueryServiceHandler(client ipc.Client) *GridQueryServiceHandler {
	return &GridQueryServiceHandler{
		flightClient: client,
	}
}

// QueryGrid implements the RPC method.
func (h *GridQueryServiceHandler) QueryGrid(
	ctx context.Context, 
	req *connect.Request[gridv1.QueryGridRequest], 
	stream *connect.ServerStream[gridv1.GridChunk],
) error {
	// 1. Resolve Plan ID (Mock for now)
	// In real life: Look up View ID -> Get AtomScript Formula -> Submit to Engine
	planID := "PLAN-" + req.Msg.ViewId

	// 2. Call Rust Engine via Arrow Flight (Layer 3.2)
	reader, err := h.flightClient.GetCalculation(ctx, planID)
	if err != nil {
		return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to call engine: %w", err))
	}
	defer reader.Release()

	// 3. Stream Results (Ultra Diamond: Zero Copy)
	for reader.Next() {
		record := reader.Record()
		
		// Serialize Arrow RecordBatch to Bytes (Zero-Copy Intent)
		bytes, release, err := ipc.SerializeRecord(record)
		if err != nil {
			return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to serialize record: %w", err))
		}

		chunk := &gridv1.GridChunk{
			Data: &gridv1.GridChunk_ArrowRecordBatch{
				ArrowRecordBatch: bytes,
			},
		}

		if err := stream.Send(chunk); err != nil {
			release()
			return err
		}
		release()
	}

	if reader.Err() != nil {
		return connect.NewError(connect.CodeInternal, reader.Err())
	}

	return nil
}
