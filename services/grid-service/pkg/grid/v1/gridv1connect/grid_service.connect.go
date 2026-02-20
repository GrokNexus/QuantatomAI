package gridv1connect

import (
	"context"
	gridv1 "quantatomai/grid-service/pkg/grid/v1"

	"connectrpc.com/connect"
)

// GridQueryServiceName is the fully-qualified name of the GridQueryService service.
const GridQueryServiceName = "quantatomai.grid.v1.GridQueryService"

// GridQueryServiceClient is a client for the quantatomai.grid.v1.GridQueryService service.
type GridQueryServiceClient interface {
	QueryGrid(context.Context, *connect.Request[gridv1.QueryGridRequest]) (*connect.ServerStreamForClient[gridv1.GridChunk], error)
}

// GridQueryServiceHandler is an implementation of the quantatomai.grid.v1.GridQueryService service.
type GridQueryServiceHandler interface {
	QueryGrid(context.Context, *connect.Request[gridv1.QueryGridRequest], *connect.ServerStream[gridv1.GridChunk]) error
}
