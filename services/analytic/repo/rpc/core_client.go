package rpc

import (
	"context"

	"tenkhours/pkg/pb"
)

type CoreClient struct {
	pb.CoreClient
}

func NewCoreClient(coreClient pb.CoreClient) *CoreClient {
	return &CoreClient{CoreClient: coreClient}
}

func (c *CoreClient) CheckPermission(ctx context.Context, profileID, characterID string, metricID *string) (bool, error) {
	req := &pb.CheckPermissionReq{
		ProfileID:   profileID,
		CharacterID: characterID,
		MetricID:    metricID,
	}

	res, err := c.CoreClient.CheckPermission(context.TODO(), req)
	if err != nil {
		return false, err
	}

	return res.Authorized, nil
}
