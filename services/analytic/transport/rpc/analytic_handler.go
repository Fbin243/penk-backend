package rpc

import (
	"context"

	"tenkhours/proto/pb/analytic"
	"tenkhours/services/analytic/business"
)

type AnalyticHandler struct {
	analytic.UnimplementedAnalyticServer
	analyticBusiness business.IAnalyticBusiness
}

func NewAnalyticHandler(analyticBiz business.IAnalyticBusiness) *AnalyticHandler {
	return &AnalyticHandler{
		analyticBusiness: analyticBiz,
	}
}

func (hdl *AnalyticHandler) DeleteCapturedRecords(ctx context.Context, req *analytic.DeleteCapturedRecordsReq) (*analytic.DeleteCapturedRecordsResp, error) {
	res := &analytic.DeleteCapturedRecordsResp{Success: false}

	err := hdl.analyticBusiness.DeleteCapturedRecords(ctx, req.ProfileId)
	if err != nil {
		return res, err
	}

	res.Success = true
	return res, nil
}
