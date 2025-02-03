package rpc

import (
	"context"
	"fmt"
	"log"

	"tenkhours/pkg/pb"
	"tenkhours/services/timetracking/entity"
)

type CurrencyClient struct {
	currencyClient pb.CurrencyClient
}

func NewCurrencyClient(currencyClient pb.CurrencyClient) *CurrencyClient {
	return &CurrencyClient{currencyClient: currencyClient}
}

func (c *CurrencyClient) CatchFish(ctx context.Context) (*entity.CatchFishResult, error) {
	log.Print("Send request to Currency to catch fish ...")
	req := &pb.CatchFishReq{}

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
	req := &pb.UpdateFishReq{
		ProfileID: fish.ProfileID,
		Gold:      fish.Gold,
		Normal:    fish.Normal,
	}

	_, err := c.currencyClient.UpdateFish(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
