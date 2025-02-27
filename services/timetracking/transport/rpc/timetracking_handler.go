package rpc

import (
	"context"
	"errors"
	"time"

	"tenkhours/proto/pb/common"
	tt "tenkhours/proto/pb/timetracking"
	"tenkhours/services/timetracking/business"

	"github.com/samber/lo"
)

type TimetrackingHandler struct {
	tt.TimeTrackingServiceServer
	timetrackingBiz business.ITimeTrackingBusiness
}

func NewtimetrackingHandler(timetrackingBiz business.ITimeTrackingBusiness) *TimetrackingHandler {
	return &TimetrackingHandler{
		timetrackingBiz: timetrackingBiz,
	}
}

func (hdl *TimetrackingHandler) CreateTimeTracking(ctx context.Context, req *tt.CreateTimeTrackingRequest) (*tt.TimeTracking, error) {
	timetracking, err := hdl.timetrackingBiz.CreateTimeTracking(ctx, req.CharacterId, req.CategoryId, time.Unix(req.StartTime, 0))
	if err != nil {
		return nil, err
	}

	resp := &tt.TimeTracking{
		Id:          timetracking.ID,
		CharacterId: timetracking.CharacterID,
		CategoryId:  timetracking.CategoryID,
		StartTime:   timetracking.StartTime.Unix(),
		EndTime:     nil,
	}

	return resp, nil
}

func (hdl *TimetrackingHandler) UpdateTimeTracking(ctx context.Context, req *common.EmptyReq) (*tt.TimeTrackingWithFish, error) {
	timetracking, fish, err := hdl.timetrackingBiz.UpdateTimeTracking(ctx)
	if err != nil {
		return nil, err
	}

	resp := &tt.TimeTrackingWithFish{
		TimeTracking: &tt.TimeTracking{
			Id:          timetracking.ID,
			CharacterId: timetracking.CharacterID,
			CategoryId:  timetracking.CategoryID,
			StartTime:   timetracking.StartTime.Unix(),
			EndTime:     lo.ToPtr(timetracking.EndTime.Unix()),
		},
		Normal: fish.Normal,
		Gold:   fish.Gold,
	}

	return resp, nil
}

// TODO: @Fbin243 implements the following methods later
func (hdl *TimetrackingHandler) GetCurrentTimeTracking(context.Context, *common.EmptyReq) (*tt.TimeTracking, error) {
	return nil, errors.New("not implemented")
}

func (hdl *TimetrackingHandler) GetTotalCurrentTimeTracking(context.Context, *tt.TotalTimeTrackingRequest) (*tt.TotalTimeTrackingResponse, error) {
	return nil, errors.New("not implemented")
}
