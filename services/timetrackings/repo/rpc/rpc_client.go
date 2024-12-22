package rpc

import (
	"context"
	"log"
	"tenkhours/pkg/pb"
)

type RPCClient struct {
	coreClient pb.CoreClient
}

func NewRPCClient(coreClient pb.CoreClient) *RPCClient {
	return &RPCClient{coreClient: coreClient}
}

func (c *RPCClient) UpdateTimeInCharacter(ctx context.Context, characterID string, metricID string, time int32) (bool, error) {
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

func (c *RPCClient) CheckPermission(ctx context.Context, profileID, characterID, metricID string) (bool, error) {
	log.Print("Send request to Core to check permission ...")
	log.Printf("Check permission for profileID: %s, characterID: %s, metricID: %s", profileID, characterID, metricID)
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
