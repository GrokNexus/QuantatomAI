package ipc

import (
	"context"
	"fmt"

	"github.com/apache/arrow/go/v15/arrow/flight"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive" // Ultra Diamond: Resilience
	"time"
)

// RecordReader is an interface for reading Arrow records.
// It abstracts *flight.Reader to allow for mocking in tests.
type RecordReader interface {
	Next() bool
	Record() arrow.Record
	Err() error
	Release()
}

// Client is the interface for communicating with the Atom Engine (Rust).
type Client interface {
	// GetCalculation executes a plan on the engine and returns a stream of Arrow records.
	GetCalculation(ctx context.Context, planID string) (RecordReader, error)
	Close() error
}

// FlightClient implements the Client interface using Apache Arrow Flight.
type FlightClient struct {
	client flight.Client
}

// NewFlightClient creates a connection to the Rust Engine.
func NewFlightClient(addr string) (*FlightClient, error) {
	// Connect to the Flight Server (Rust)
	// Ultra Diamond: KeepAlive Configuration
	// Prevents load balancers (AWS ALB / K8s Service) from killing idle connections.
	kac := keepalive.ClientParameters{
		Time:                10 * time.Second, // Send pings every 10 seconds if no activity
		Timeout:             time.Second,      // Wait 1 second for ping ack before considering dead
		PermitWithoutStream: true,             // Send pings even without active streams
	}

	c, err := flight.NewClientWithMiddleware(
		addr, 
		nil, 
		nil, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(kac),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create flight client: %w", err)
	}

	return &FlightClient{
		client: c,
	}, nil
}

// GetCalculation sends a ticket (PlanID) to the engine and streams back the results.
func (c *FlightClient) GetCalculation(ctx context.Context, planID string) (RecordReader, error) {
	// Create the Ticket
	ticket := &flight.Ticket{
		Ticket: []byte(planID), // In real impl, this might be a serialized Protobuf Plan
	}

	// execute DoGet
	stream, err := c.client.DoGet(ctx, ticket)
	if err != nil {
		return nil, fmt.Errorf("flight DoGet failed: %w", err)
	}

	// Create a Reader from the stream
	rdr, err := flight.NewRecordReader(stream)
	if err != nil {
		return nil, fmt.Errorf("failed to create record reader: %w", err)
	}

	return rdr, nil
}

func (c *FlightClient) Close() error {
	return c.client.Close()
}
