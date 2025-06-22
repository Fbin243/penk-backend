package rpc

import (
	"context"

	"tenkhours/proto/pb/common"
	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"
)

func (hdl *CoreHandler) UpsertMetric(ctx context.Context, req *core.MetricInput) (*core.Metric, error) {
	metricInput, err := Map[core.MetricInput, entity.MetricInput](req, UnixTimeConverter)
	if err != nil {
		return nil, err
	}

	metric, err := hdl.metricBiz.Upsert(ctx, metricInput)
	if err != nil {
		return nil, err
	}

	return Map[entity.Metric, core.Metric](metric, UnixTimeConverter)
}

func (hdl *CoreHandler) DeleteMetric(ctx context.Context, req *common.IdReq) (*common.IdResp, error) {
	metric, err := hdl.metricBiz.Delete(ctx, req.Id)

	return &common.IdResp{
		Id: metric.ID,
	}, err
}
