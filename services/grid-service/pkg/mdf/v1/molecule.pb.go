package mdfv1

type Molecule struct {
	Coordinate []int64 `protobuf:"varint,1,rep,packed,name=coordinate,proto3" json:"coordinate,omitempty"`
	ScenarioId int64   `protobuf:"varint,2,opt,name=scenario_id,json=scenarioId,proto3" json:"scenario_id,omitempty"`
	Value      float64 `protobuf:"fixed64,3,opt,name=value,proto3" json:"value,omitempty"`
	State      int32   `protobuf:"varint,4,opt,name=state,proto3" json:"state,omitempty"` // 0: RAW, 1: CALC, 2: CONSOL, 3: ADJUST
	Tag        string  `protobuf:"bytes,5,opt,name=tag,proto3" json:"tag,omitempty"`
	TraceId    string  `protobuf:"bytes,6,opt,name=trace_id,json=traceId,proto3" json:"trace_id,omitempty"`
}

func (m *Molecule) Reset()         { *m = Molecule{} }
func (m *Molecule) String() string { return "Molecule" }
func (*Molecule) ProtoMessage()    {}
