package rpc

import (
	"context"

	"tenkhours/pkg/pb"
)

type CoreClient struct {
	coreClient pb.CoreClient
}

func NewCoreClient(coreClient pb.CoreClient) *CoreClient {
	return &CoreClient{coreClient: coreClient}
}

func (c *CoreClient) CheckPermission(ctx context.Context, profileID, characterID, categoryID *string) (bool, error) {
	req := &pb.CheckPermissionReq{
		ProfileID:   profileID,
		CharacterID: characterID,
		CategoryID:  categoryID,
	}

	res, err := c.coreClient.CheckPermission(ctx, req)
	if err != nil {
		return false, err
	}

	return res.Authorized, nil
}
