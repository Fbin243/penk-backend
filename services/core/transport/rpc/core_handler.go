package rpc

import (
	"context"

	"tenkhours/proto/pb/core"
	"tenkhours/services/core/business"
)

type CoreHandler struct {
	core.UnimplementedCoreServer
	profileBiz      business.IProfileBusiness
	characterBiz    business.ICharacterBusiness
	goalBiz         business.IGoalBusiness
	timetrackingBiz business.ITimeTrackingBusiness
	taskBiz         business.ITaskBusiness
}

func NewCoreHandler(
	profilesBusiness business.IProfileBusiness,
	charactersBusiness business.ICharacterBusiness,
	goalBiz business.IGoalBusiness,
	timetrackingBiz business.ITimeTrackingBusiness,
	taskBiz business.ITaskBusiness,
) *CoreHandler {
	return &CoreHandler{
		profileBiz:      profilesBusiness,
		characterBiz:    charactersBusiness,
		goalBiz:         goalBiz,
		timetrackingBiz: timetrackingBiz,
		taskBiz:         taskBiz,
	}
}

func (hdl *CoreHandler) IntrospectToken(ctx context.Context, req *core.IntrospectReq) (*core.IntrospectResp, error) {
	resp := &core.IntrospectResp{Success: false}

	authSession, err := hdl.profileBiz.IntrospectToken(ctx, req.Token, req.DeviceId)
	if err != nil {
		return resp, err
	}

	resp.Success = true
	resp.ProfileId = authSession.ProfileID
	resp.DeviceId = authSession.DeviceID
	resp.CurrentCharacterId = authSession.CurrentCharacterID
	resp.FirebaseUid = authSession.FirebaseUID

	return resp, nil
}
