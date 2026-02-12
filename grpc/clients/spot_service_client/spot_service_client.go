package spotserviceclient

import (
	"context"
	"fmt"
	"time"

	"github.com/HomidWay/microservice-hw-proto/pb/spotinstrumentservice"
	"google.golang.org/grpc"
)

type SpotServiceClient struct {
	rpcClient spotinstrumentservice.SpotInstrumentServiceClient
}

func NewSpotServiceClient(host string, port int) (*SpotServiceClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to userservice: %w", err)
	}

	client := spotinstrumentservice.NewSpotInstrumentServiceClient(conn)

	return &SpotServiceClient{rpcClient: client}, nil
}

func (s *SpotServiceClient) GetAvailableMarkets(session *sessionmanager.Session) ([]marketmanager.Market, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := s.rpcClient.ViewMarkets(ctx, &spotinstrumentservice.ViewMarketsRequest{
		UserSessionId: session.SessionID(),
		UserId:        session.User().ID(),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch available markets: %w", err)
	}

	markets := make([]marketmanager.Market, len(response.Markets))

	for i, market := range response.Markets {

		var deleteTime *time.Time
		if market.DeletedAt != nil {
			delTime := market.DeletedAt.AsTime()
			deleteTime = &delTime
		}

		markets[i] = marketmanager.NewMarket(
			market.Id,
			market.Active,
			deleteTime,
		)
	}

	return markets, nil
}
