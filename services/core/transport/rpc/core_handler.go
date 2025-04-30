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
	metricBiz       business.IMetricBusiness
	categoryBiz     business.ICategoryBusiness
	habitBiz        business.IHabitBusiness
}

func NewCoreHandler(
	profilesBusiness business.IProfileBusiness,
	charactersBusiness business.ICharacterBusiness,
	goalBiz business.IGoalBusiness,
	timetrackingBiz business.ITimeTrackingBusiness,
	taskBiz business.ITaskBusiness,
	metricBiz business.IMetricBusiness,
	categoryBiz business.ICategoryBusiness,
	habitBiz business.IHabitBusiness,
) *CoreHandler {
	return &CoreHandler{
		profileBiz:      profilesBusiness,
		characterBiz:    charactersBusiness,
		goalBiz:         goalBiz,
		timetrackingBiz: timetrackingBiz,
		taskBiz:         taskBiz,
		metricBiz:       metricBiz,
		categoryBiz:     categoryBiz,
		habitBiz:        habitBiz,
	}
}

func (hdl *CoreHandler) IntrospectUser(ctx context.Context, req *core.IntrospectReq) (*core.IntrospectResp, error) {
	resp := &core.IntrospectResp{Success: false}

	authSession, err := hdl.profileBiz.IntrospectUser(ctx, req.Token, req.UserId, req.DeviceId)
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
