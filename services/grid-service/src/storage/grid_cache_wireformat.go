package storage

import (
	"encoding/json"
)

type WireFormatType string

const (
	WireFormatJSON       WireFormatType = "json"
	WireFormatFlatBuffer WireFormatType = "flatbuffer"
	WireFormatCapnProto  WireFormatType = "capnproto"
)

type WireFormatCodec interface {
	EncodeGrid(grid interface{}) ([]byte, error)
	DecodeGrid(data []byte, target interface{}) error
	Format() WireFormatType
}

type JSONWireFormatCodec struct{}

func NewJSONWireFormatCodec() *JSONWireFormatCodec { return &JSONWireFormatCodec{} }

func (c *JSONWireFormatCodec) EncodeGrid(grid interface{}) ([]byte, error) {
	return json.Marshal(grid)
}

func (c *JSONWireFormatCodec) DecodeGrid(data []byte, target interface{}) error {
	return json.Unmarshal(data, target)
}

func (c *JSONWireFormatCodec) Format() WireFormatType { return WireFormatJSON }

type CapnProtoWireFormatCodec struct{}

func NewCapnProtoWireFormatCodec() *CapnProtoWireFormatCodec {
	return &CapnProtoWireFormatCodec{}
}

func (c *CapnProtoWireFormatCodec) EncodeGrid(grid interface{}) ([]byte, error) {
	return nil, ErrWireFormatNotImplemented("capnproto EncodeGrid")
}

func (c *CapnProtoWireFormatCodec) DecodeGrid(data []byte, target interface{}) error {
	return ErrWireFormatNotImplemented("capnproto DecodeGrid")
}

func (c *CapnProtoWireFormatCodec) Format() WireFormatType { return WireFormatCapnProto }

type ErrWireFormatNotImplemented string

func (e ErrWireFormatNotImplemented) Error() string {
	return "wire format not implemented: " + string(e)
}
