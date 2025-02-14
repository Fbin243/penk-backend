package rpc

import (
	"tenkhours/pkg/pb"
)

type CoreClient struct {
	pb.CoreClient
}

func NewCoreClient(coreClient pb.CoreClient) *CoreClient {
	return &CoreClient{CoreClient: coreClient}
}
