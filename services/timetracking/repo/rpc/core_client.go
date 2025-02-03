package rpc

import (
	"context"
	"log"

	"tenkhours/pkg/pb"
)

type CoreClient struct {
	coreClient pb.CoreClient
}

func NewCoreClient(coreClient pb.CoreClient) *CoreClient {
	return &CoreClient{coreClient: coreClient}
}

func (c *CoreClient) UpdateTimeInCharacter(ctx context.Context, characterID string, metricID *string, time int32) (bool, error) {
	log.Print("Send request to Core to update time in character ...")
	req := &pb.UpdateTimeReq{
		CharacterID: characterID,
		MetricID:    metricID,
		Time:        time,
	}

	res, err := c.coreClient.UpdateTimeInCharacter(ctx, req)
	if err != nil {
		return false, err
	}

	return res.Success, nil
}

func (c *CoreClient) CheckPermission(ctx context.Context, profileID, characterID string, metricID *string) (bool, error) {
	log.Print("Send request to Core to check permission ...")
	req := &pb.CheckPermissionReq{
		ProfileID:   profileID,
		CharacterID: characterID,
		MetricID:    metricID,
	}

	res, err := c.coreClient.CheckPermission(ctx, req)
	if err != nil {
		return false, err
	}

	return res.Authorized, nil
}
