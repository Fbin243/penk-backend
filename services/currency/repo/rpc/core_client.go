package rpc

import (
	"context"
	"fmt"

	"tenkhours/pkg/pb"
	"tenkhours/services/currency/entity"
)

type CoreClient struct {
	pb.CoreClient
}

func NewCoreClient(coreClient pb.CoreClient) *CoreClient {
	return &CoreClient{CoreClient: coreClient}
}

func (c *CoreClient) BuyItem(ctx context.Context, profileID, characterID, metricID *string, item entity.ItemType, amount int32) error {
	req := &pb.BuyItemReq{
		ProfileID:   profileID,
		CharacterID: characterID,
		MetricID:    metricID,
		ItemType:    pb.BuyItemReq_ItemType(item),
		Amount:      amount,
	}

	res, err := c.CoreClient.BuyItem(ctx, req)
	if !res.Success || err != nil {
		return fmt.Errorf("failed to buy item: %v", err)
	}

	return err
}
