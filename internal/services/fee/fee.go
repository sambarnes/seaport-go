package fee

import (
	"context"

	"github.com/bufbuild/connect-go"

	v1 "github.com/sambarnes/seaport-go/pkg/proto/seaport/v1"
)

// Service holds state for the fee.Service
type Service struct {
	// TODO: all in-memory state for now
}

// GetFeeStructure returns the Service structure.
func (s *Service) GetFeeStructure(_ context.Context, _ *connect.Request[v1.Empty]) (
	*connect.Response[v1.FeeStructure], error,
) {
	// TODO: placeholder fee structure for now
	fee := &v1.FeeStructure{
		Maker: &v1.TradeFees{
			NotionalBps: 0,
			PremiumBps:  0,
			SpotBps:     0,
			Flat:        0,
		},
		Taker: &v1.TradeFees{
			NotionalBps: 0,
			PremiumBps:  0,
			SpotBps:     0,
			Flat:        0,
		},
		ClearWriteNotionalBps:    0,
		ClearRedeemedNotionalBps: 0,
		ClearExerciseNotionalBps: 0,
		Address: &v1.H160{
			Hi: &v1.H128{Hi: 17, Lo: 38},
			Lo: 420,
		},
	}

	res := connect.NewResponse(fee)
	res.Header().Set("Seaport-Version", "v1")
	return res, nil
}
