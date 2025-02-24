package rpc

import (
	"context"

	"tenkhours/proto/pb/currency"
	"tenkhours/services/currency/business"
	"tenkhours/services/currency/entity"
)

type CurrencyHandler struct {
	currency.UnimplementedCurrencyServer
	currencyBusiness business.ICurrencyBusiness
}

func NewCurrencyHandler(currencyBiz business.ICurrencyBusiness) *CurrencyHandler {
	return &CurrencyHandler{
		currencyBusiness: currencyBiz,
	}
}

func (h *CurrencyHandler) CreateFish(ctx context.Context, req *currency.CreateFishReq) (*currency.CreateFishResp, error) {
	res := &currency.CreateFishResp{Success: false}

	_, err := h.currencyBusiness.CreateFish(ctx, req.ProfileId)
	if err != nil {
		return res, err
	}

	res.Success = true
	return res, nil
}

func (hdl *CurrencyHandler) CatchFish(ctx context.Context, req *currency.CatchFishReq) (*currency.CatchFishResp, error) {
	res := &currency.CatchFishResp{FishType: currency.CatchFishResp_None, Number: 0}

	catchFishResult, err := hdl.currencyBusiness.CatchFish(ctx)
	if err != nil {
		return res, err
	}

	res.FishType = currency.CatchFishResp_FishType(currency.CatchFishResp_FishType_value[string(catchFishResult.FishType)])
	res.Number = catchFishResult.Number
	return res, nil
}

func (hdl *CurrencyHandler) UpdateFish(ctx context.Context, req *currency.UpdateFishReq) (*currency.UpdateFishResp, error) {
	res := &currency.UpdateFishResp{FishId: ""}
	fish, err := hdl.currencyBusiness.UpdateFish(ctx, &entity.Fish{
		ProfileID: req.ProfileId,
		Gold:      req.Gold,
		Normal:    req.Normal,
	})
	if err != nil {
		return res, err
	}

	res.FishId = fish.ID
	return res, nil
}
