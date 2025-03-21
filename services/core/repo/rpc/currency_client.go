package rpc

import (
	"context"
	"fmt"

	"tenkhours/proto/pb/currency"
)

type CurrencyClient struct {
	currency.CurrencyClient
}

func NewCurrencyClient(currencyClient currency.CurrencyClient) *CurrencyClient {
	return &CurrencyClient{CurrencyClient: currencyClient}
}

func (c *CurrencyClient) CreateFish(ctx context.Context, profileID string) error {
	req := &currency.CreateFishReq{
		ProfileId: profileID,
	}

	res, err := c.CurrencyClient.CreateFish(ctx, req)
	if err != nil && !res.Success {
		return fmt.Errorf("failed to create fish: %v", err)
	}

	return nil
}

func (c *CurrencyClient) DeleteFish(ctx context.Context, profileID string) error {
	req := &currency.DeleteFishReq{
		ProfileId: profileID,
	}

	res, err := c.CurrencyClient.DeleteFish(ctx, req)
	if err != nil && !res.Success {
		return fmt.Errorf("failed to delete fish: %v", err)
	}

	return nil
}
