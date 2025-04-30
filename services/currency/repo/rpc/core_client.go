package rpc

import "tenkhours/proto/pb/core"

type CoreClient struct {
	core.CoreClient
}

func NewCoreClient(coreClient core.CoreClient) *CoreClient {
	return &CoreClient{CoreClient: coreClient}
}
