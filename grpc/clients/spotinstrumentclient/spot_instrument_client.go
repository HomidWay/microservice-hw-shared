package spotinstrumentclient

import (
	"context"
	"fmt"
	"time"

	"github.com/HomidWay/microservice-hw-proto/pb/spotinstrumentservice"
	"google.golang.org/grpc"
)

type SpotInstrumentClient struct {
	rpcClient spotinstrumentservice.SpotInstrumentServiceClient
}

func NewSpotInstrumentClient(rpcHost string, port int, options ...grpc.DialOption) (*SpotInstrumentClient, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", rpcHost, port), options...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SpotInstrument RPC server: %w", err)
	}

	client := spotinstrumentservice.NewSpotInstrumentServiceClient(conn)

	return &SpotInstrumentClient{rpcClient: client}, nil
}

func (si SpotInstrumentClient) ViewMarkets(ctx context.Context, userID string) ([]Market, error) {

	request := &spotinstrumentservice.ViewMarketsRequest{
		UserId: userID,
	}

	marketResponse, err := si.rpcClient.ViewMarkets(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to ViewMarkets with error: %s", err.Error())
	}

	markets := make([]Market, 0, len(marketResponse.Markets))

	for _, m := range marketResponse.Markets {

		var deletedAt *time.Time
		if m.DeletedAt != nil {
			deletedTime := m.DeletedAt.AsTime()
			deletedAt = &deletedTime
		}

		market := NewMarket(m.Id, m.Active, deletedAt)

		markets = append(markets, market)
	}

	return markets, nil
}
