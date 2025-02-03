package rpc

import (
	"context"
	"fmt"

	"tenkhours/pkg/pb"
)

type CurrencyClient struct {
	pb.CurrencyClient
}

func NewCurrencyClient(currencyClient pb.CurrencyClient) *CurrencyClient {
	return &CurrencyClient{CurrencyClient: currencyClient}
}

func (c *CurrencyClient) CreateFish(ctx context.Context, profileID string) error {
	req := &pb.CreateFishReq{
		ProfileID: profileID,
	}

	res, err := c.CurrencyClient.CreateFish(ctx, req)
	if err != nil && !res.Success {
		return fmt.Errorf("failed to create fish: %v", err)
	}

	return nil
}
