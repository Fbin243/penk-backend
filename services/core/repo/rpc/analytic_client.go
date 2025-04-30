package rpc

import (
	"context"
	"fmt"

	"tenkhours/proto/pb/analytic"
)

type AnalyticClient struct {
	analytic.AnalyticClient
}

func NewAnalyticClient(analyticClient analytic.AnalyticClient) *AnalyticClient {
	return &AnalyticClient{AnalyticClient: analyticClient}
}

func (c *AnalyticClient) DeleteCapturedRecords(ctx context.Context, profileID string) error {
	req := &analytic.DeleteCapturedRecordsReq{
		ProfileId: profileID,
	}

	res, err := c.AnalyticClient.DeleteCapturedRecords(ctx, req)
	if err != nil && !res.Success {
		return fmt.Errorf("failed to create fish: %v", err)
	}

	return nil
}
