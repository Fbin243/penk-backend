package rpc

import (
	"context"

	"tenkhours/pkg/pb"
	"tenkhours/services/analytic/business"
)

type AnalyticHandler struct {
	pb.UnimplementedAnalyticServer
	analyticBusiness business.IAnalyticBusiness
}

func NewAnalyticHandler(analyticBiz business.IAnalyticBusiness) *AnalyticHandler {
	return &AnalyticHandler{
		analyticBusiness: analyticBiz,
	}
}

func (hdl *AnalyticHandler) DeleteCapturedRecords(ctx context.Context, req *pb.DeleteCapturedRecordsReq) (*pb.DeleteCapturedRecordsResp, error) {
	res := &pb.DeleteCapturedRecordsResp{Success: false}

	err := hdl.analyticBusiness.DeleteCapturedRecords(ctx, req.ProfileID)
	if err != nil {
		return res, err
	}

	res.Success = true
	return res, nil
}
