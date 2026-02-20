package gridv1

import (
	mdfv1 "quantatomai/grid-service/pkg/mdf/v1"
)

type QueryGridRequest struct {
	ViewId           string   `protobuf:"bytes,1,opt,name=view_id,json=viewId,proto3" json:"view_id,omitempty"`
	Dimensions       []string `protobuf:"bytes,2,rep,name=dimensions,proto3" json:"dimensions,omitempty"`
	FilterExpression string   `protobuf:"bytes,3,opt,name=filter_expression,json=filterExpression,proto3" json:"filter_expression,omitempty"`
}

type GridChunk struct {
	Data isGridChunk_Data `protobuf_oneof:"data"`
}

type GridChunk_CellsV1 struct {
	CellsV1 *Molecules `protobuf:"bytes,1,opt,name=cells_v1,json=cellsV1,proto3,oneof"`
}

type GridChunk_ArrowRecordBatch struct {
	ArrowRecordBatch []byte `protobuf:"bytes,2,opt,name=arrow_record_batch,json=arrowRecordBatch,proto3,oneof"`
}

func (*GridChunk_CellsV1) isGridChunk_Data()          {}
func (*GridChunk_ArrowRecordBatch) isGridChunk_Data() {}

type isGridChunk_Data interface {
	isGridChunk_Data()
}

type Molecules struct {
	List []*mdfv1.Molecule `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
}

func (m *QueryGridRequest) Reset()         { *m = QueryGridRequest{} }
func (m *QueryGridRequest) String() string { return "QueryGridRequest" }
func (*QueryGridRequest) ProtoMessage()    {}

func (m *GridChunk) Reset()         { *m = GridChunk{} }
func (m *GridChunk) String() string { return "GridChunk" }
func (*GridChunk) ProtoMessage()    {}

func (m *Molecules) Reset()         { *m = Molecules{} }
func (m *Molecules) String() string { return "Molecules" }
func (*Molecules) ProtoMessage()    {}
