package rfq

import (
	"context"

	"github.com/bufbuild/connect-go"

	v1 "github.com/sambarnes/seaport-go/pkg/proto/seaport/v1"
)

// Service holds state for the rfq.Service.
type Service struct {
	// TODO: all in-memory state for now
}

// WebTaker requests quotes from makers via a single QuoteRequest message and
// receives a stream of QuoteResponse messages (for use by gRPC-web clients).
func (s *Service) WebTaker(
	_ context.Context,
	_ *connect.Request[v1.QuoteRequest],
	_ *connect.ServerStream[v1.QuoteResponse],
) error {
	return nil
}

// Taker requests quotes from makers via a stream of QuoteRequest messages and
// receives a stream of QuoteResponse messages.
func (s *Service) Taker(
	_ context.Context,
	_ *connect.BidiStream[v1.QuoteRequest, v1.QuoteResponse],
) error {
	return nil
}

// Maker sends quotes to takers via a stream of QuoteResponse messages and
// receives a stream of QuoteRequest messages.
func (s *Service) Maker(
	_ context.Context,
	_ *connect.BidiStream[v1.QuoteResponse, v1.QuoteRequest],
) error {
	return nil
}
