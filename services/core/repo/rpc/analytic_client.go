package rpc

import (
	"context"
	"fmt"

	"tenkhours/pkg/pb"
)

type AnalyticClient struct {
	pb.AnalyticClient
}

func NewAnalyticClient(analyticClient pb.AnalyticClient) *AnalyticClient {
	return &AnalyticClient{AnalyticClient: analyticClient}
}

func (c *AnalyticClient) DeleteCapturedRecords(ctx context.Context, profileID string) error {
	req := &pb.DeleteCapturedRecordsReq{
		ProfileID: profileID,
	}

	res, err := c.AnalyticClient.DeleteCapturedRecords(ctx, req)
	if err != nil && !res.Success {
		return fmt.Errorf("failed to create fish: %v", err)
	}

	return nil
}
