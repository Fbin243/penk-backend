package rpc

import (
	"context"

	"tenkhours/proto/pb/core"
	"tenkhours/services/core/business"
)

type CoreHandler struct {
	core.UnimplementedCoreServer
	profileBiz   business.IProfileBusiness
	characterBiz business.ICharacterBusiness
	goalBiz      business.IGoalBusiness
}

func NewCoreHandler(profilesBusiness business.IProfileBusiness, charactersBusiness business.ICharacterBusiness, goalBiz business.IGoalBusiness) *CoreHandler {
	return &CoreHandler{
		profileBiz:   profilesBusiness,
		characterBiz: charactersBusiness,
		goalBiz:      goalBiz,
	}
}

func (hdl *CoreHandler) IntrospectToken(ctx context.Context, req *core.IntrospectReq) (*core.IntrospectResp, error) {
	resp := &core.IntrospectResp{Success: false}

	authSession, err := hdl.profileBiz.IntrospectToken(ctx, req.Token)
	if err != nil {
		return resp, err
	}

	resp.Success = true
	resp.ProfileId = authSession.ProfileID

	return resp, nil
}

func (hdl *CoreHandler) CheckPermission(ctx context.Context, req *core.CheckPermissionReq) (*core.CheckPermissionResp, error) {
	resp := &core.CheckPermissionResp{Authorized: false}

	err := hdl.profileBiz.CheckPermission(ctx, req.ProfileId, req.CharacterId, req.CategoryId)
	if err != nil {
		return resp, err
	}

	resp.Authorized = true

	return resp, nil
}
