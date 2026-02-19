package event

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy" // Ultra Diamond: Snappy Compression
)

// KafkaBus implements the Bus interface using segmentio/kafka-go.
// It is compatible with Redpanda, Apache Kafka, etc.
type KafkaBus struct {
	writer *kafka.Writer
	reader *kafka.Reader // Optional, initialized on Subscribe
	brokers []string
}

// NewKafkaBus creates a new Bus connected to the specified brokers.
func NewKafkaBus(brokers []string, topic string) *KafkaBus {
	// Configure the Writer for high throughput/low latency
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		// Optimized batch settings for Redpanda
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
		Async:        true, // Non-blocking write
		
		// Ultra Diamond: Compression Enabled (Save Bandwidth)
		Compression: kafka.Snappy, 

		// Ultra Diamond: Error Logger (DLQ Simulation)
		// If write fails asynchronously, we log heavily so external monitoring picks it up.
		ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			// In production, this might write to a "dead_letter.log" file or specific error topic
			log.Printf("[KAFKA-ERROR] "+msg, args...)
		}),
	}

	return &KafkaBus{
		writer:  w,
		brokers: brokers,
	}
}

// Publish sends an event asynchronously.
func (kb *KafkaBus) Publish(ctx context.Context, event *AtomEventGo) error {
	// Convert AtomEventGo to Kafka Message
	// In a real implementation, we would Serialize to FlatBuffers here using the generated code.
	// For now, we wrap the raw payload.
	
	msg := kafka.Message{
		Key:   []byte(event.TenantID), // Key by Tenant for partitioning
		Value: event.Payload,          // The FlatBuffer bytes
		Time:  time.Unix(0, event.Timestamp*int64(time.Nanosecond)),
		Headers: []kafka.Header{
			{Key: "trace_id", Value: []byte(event.TraceID)},
			{Key: "type", Value: []byte(fmt.Sprintf("%d", event.Type))},
		},
	}

	// Because Async is true, this puts it in the buffer and returns nil immediately
	// If the buffer is full, it might block depending on config, but context cancels it.
	return kb.writer.WriteMessages(ctx, msg)
}

// Subscribe returns a channel of events.
// Note: This implementation creates a new Reader for the topic.
func (kb *KafkaBus) Subscribe(ctx context.Context, topic string) (<-chan *AtomEventGo, error) {
	// Config reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   kb.brokers,
		Topic:     topic,
		GroupID:   "grid-service-group",
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	
	kb.reader = r
	out := make(chan *AtomEventGo, 100)

	go func() {
		defer close(out)
		defer r.Close()
		
		for {
			m, err := r.ReadMessage(ctx)
			if err != nil {
				// Context canceled or reader closed
				return
			}
			
			// Deserialize (Simulated)
			// In real impl: Use generated FlatBuffers code to read message.
			evt := &AtomEventGo{
				Payload: m.Value,
				TenantID: string(m.Key),
				// Extract headers...
			}
			
			select {
			case out <- evt:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, nil
}

// Close closes the writer and reader interactions.
func (kb *KafkaBus) Close() error {
	if kb.writer != nil {
		return kb.writer.Close()
	}
	return nil
}
