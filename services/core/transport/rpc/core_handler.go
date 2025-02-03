package rpc

import (
	"context"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/pb"
	"tenkhours/services/core/business"
	"tenkhours/services/core/entity"
)

type CoreHandler struct {
	pb.UnimplementedCoreServer
	profilesBusiness   business.IProfileBusiness
	charactersBusiness business.ICharacterBusiness
}

func NewCoreHandler(profilesBusiness business.IProfileBusiness, charactersBusiness business.ICharacterBusiness) *CoreHandler {
	return &CoreHandler{
		profilesBusiness:   profilesBusiness,
		charactersBusiness: charactersBusiness,
	}
}

func (hdl *CoreHandler) UpdateTimeInCharacter(ctx context.Context, req *pb.UpdateTimeReq) (*pb.UpdateTimeResp, error) {
	resp := &pb.UpdateTimeResp{Success: false}

	err := hdl.charactersBusiness.UpdateTimeInCharacter(ctx, req.CharacterID, req.MetricID, req.Time)
	if err != nil {
		return resp, err
	}

	resp.Success = true

	return resp, nil
}

func (hdl *CoreHandler) IntrospectProfile(ctx context.Context, req *pb.IntrospectReq) (*pb.IntrospectResp, error) {
	resp := &pb.IntrospectResp{Success: false}

	firebaseProfile := auth.FirebaseProfile{
		UID:     req.FirebaseUID,
		Email:   req.Email,
		Name:    req.Name,
		Picture: req.Picture,
	}

	profile, err := hdl.profilesBusiness.IntrospectProfile(ctx, firebaseProfile)
	if err != nil {
		return resp, err
	}

	resp.Success = true
	resp.ProfileID = profile.ID

	return resp, nil
}

func (hdl *CoreHandler) CheckPermission(ctx context.Context, req *pb.CheckPermissionReq) (*pb.CheckPermissionResp, error) {
	resp := &pb.CheckPermissionResp{Authorized: false}

	err := hdl.profilesBusiness.CheckPermission(ctx, req.ProfileID, req.CharacterID, req.MetricID)
	if err != nil {
		return resp, err
	}

	resp.Authorized = true

	return resp, nil
}

func (hdl *CoreHandler) BuyItem(ctx context.Context, req *pb.BuyItemReq) (*pb.BuyItemRes, error) {
	resp := &pb.BuyItemRes{Success: false}

	err := hdl.profilesBusiness.BuyItem(ctx, req.ProfileID, req.CharacterID, req.MetricID, entity.ItemType(req.ItemType), req.Amount)
	if err != nil {
		return resp, err
	}

	resp.Success = true

	return resp, nil
}
