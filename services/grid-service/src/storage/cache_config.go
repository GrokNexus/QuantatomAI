package storage

import "time"

type CacheConfig struct {
	WireCodec              WireFormatCodec
	WireFormat             WireFormatType
	EnableChunkedHydration bool
	DefaultTTL             time.Duration
	XFetchProbability      float64
}

func NewDefaultCacheConfig() *CacheConfig {
	return &CacheConfig{
		WireFormat:             WireFormatJSON,
		WireCodec:              NewJSONWireFormatCodec(),
		EnableChunkedHydration: false,
		DefaultTTL:             5 * time.Minute,
		XFetchProbability:      0.05,
	}
}

// BuildWireCodec builds a codec based on WireFormat if WireCodec is nil.
func (c *CacheConfig) BuildWireCodec() WireFormatCodec {
	if c.WireCodec != nil {
		return c.WireCodec
	}
	switch c.WireFormat {
	case WireFormatFlatBuffer:
		return NewFlatBufferWireFormatCodec()
	case WireFormatCapnProto:
		return NewCapnProtoWireFormatCodec()
	default:
		return NewJSONWireFormatCodec()
	}
}
