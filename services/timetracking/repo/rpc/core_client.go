package rpc

import (
	"context"

	"tenkhours/proto/pb/core"
)

type CoreClient struct {
	coreClient core.CoreClient
}

func NewCoreClient(coreClient core.CoreClient) *CoreClient {
	return &CoreClient{coreClient: coreClient}
}

func (c *CoreClient) CheckPermission(ctx context.Context, profileID, characterID, categoryID *string) (bool, error) {
	req := &core.CheckPermissionReq{
		ProfileId:   profileID,
		CharacterId: characterID,
		CategoryId:  categoryID,
	}

	res, err := c.coreClient.CheckPermission(ctx, req)
	if err != nil {
		return false, err
	}

	return res.Authorized, nil
}
