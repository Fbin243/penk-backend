package rpc

import (
	"context"

	"tenkhours/pkg/pb"
	"tenkhours/services/currency/business"
	"tenkhours/services/currency/entity"
)

type CurrencyHandler struct {
	pb.UnimplementedCurrencyServer
	currencyBusiness business.ICurrencyBusiness
}

func NewCurrencyHandler(currencyBiz business.ICurrencyBusiness) *CurrencyHandler {
	return &CurrencyHandler{
		currencyBusiness: currencyBiz,
	}
}

func (h *CurrencyHandler) CreateFish(ctx context.Context, req *pb.CreateFishReq) (*pb.CreateFishResp, error) {
	res := &pb.CreateFishResp{Success: false}

	_, err := h.currencyBusiness.CreateFish(ctx, req.ProfileID)
	if err != nil {
		return res, err
	}

	res.Success = true
	return res, nil
}

func (hdl *CurrencyHandler) CatchFish(ctx context.Context, req *pb.CatchFishReq) (*pb.CatchFishResp, error) {
	res := &pb.CatchFishResp{FishType: pb.CatchFishResp_None, Number: 0}

	catchFishResult, err := hdl.currencyBusiness.CatchFish(ctx)
	if err != nil {
		return res, err
	}

	res.FishType = pb.CatchFishResp_FishType(pb.CatchFishResp_FishType_value[string(catchFishResult.FishType)])
	res.Number = catchFishResult.Number
	return res, nil
}

func (hdl *CurrencyHandler) UpdateFish(ctx context.Context, req *pb.UpdateFishReq) (*pb.UpdateFishResp, error) {
	res := &pb.UpdateFishResp{FishID: ""}
	fish, err := hdl.currencyBusiness.UpdateFish(ctx, &entity.Fish{
		ProfileID: req.ProfileID,
		Gold:      req.Gold,
		Normal:    req.Normal,
	})
	if err != nil {
		return res, err
	}

	res.FishID = fish.ID
	return res, nil
}
