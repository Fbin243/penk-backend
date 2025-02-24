package rpc

import (
	"context"

	"tenkhours/proto/pb/core"
)

type CoreClient struct {
	core.CoreClient
}

func NewCoreClient(coreClient core.CoreClient) *CoreClient {
	return &CoreClient{CoreClient: coreClient}
}

func (c *CoreClient) CheckPermission(ctx context.Context, profileID, characterID, categoryID *string) (bool, error) {
	req := &core.CheckPermissionReq{
		ProfileId:   profileID,
		CharacterId: characterID,
		CategoryId:  categoryID,
	}

	res, err := c.CoreClient.CheckPermission(context.TODO(), req)
	if err != nil {
		return false, err
	}

	return res.Authorized, nil
}
