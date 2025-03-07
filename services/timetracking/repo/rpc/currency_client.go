package rpc

import (
	"context"
	"fmt"
	"log"

	"tenkhours/proto/pb/currency"
	"tenkhours/services/timetracking/entity"
)

type CurrencyClient struct {
	currencyClient currency.CurrencyClient
}

func NewCurrencyClient(currencyClient currency.CurrencyClient) *CurrencyClient {
	return &CurrencyClient{currencyClient: currencyClient}
}

func (c *CurrencyClient) CatchFish(ctx context.Context) (*entity.CatchFishResult, error) {
	log.Print("Send request to Currency to catch fish ...")
	req := &currency.CatchFishReq{}

	res, err := c.currencyClient.CatchFish(ctx, req)
	if err != nil {
		return nil, err
	}

	fmt.Println("res.FishType: ", res.FishType.String())

	return &entity.CatchFishResult{
		FishType: res.FishType.String(),
		Number:   res.Number,
	}, nil
}

func (c *CurrencyClient) UpdateFish(ctx context.Context, fish *entity.Fish) error {
	req := &currency.UpdateFishReq{
		ProfileId: fish.ProfileID,
		Gold:      fish.Gold,
		Normal:    fish.Normal,
	}

	_, err := c.currencyClient.UpdateFish(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
