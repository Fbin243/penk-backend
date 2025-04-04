package rpc

import (
	"context"
	"errors"
	"time"

	"tenkhours/proto/pb/common"
	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"

	"github.com/samber/lo"
)

func (hdl *CoreHandler) CreateTimeTracking(ctx context.Context, req *core.CreateTimeTrackingReq) (*core.TimeTracking, error) {
	input := &entity.TimeTrackingInput{
		CategoryID:  req.CategoryId,
		ReferenceID: req.ReferenceId,
		StartTime:   time.Unix(req.StartTime, 0),
	}

	if req.ReferenceType != nil {
		input.ReferenceType = lo.ToPtr(entity.EntityType(req.ReferenceType.String()))
	}

	timetracking, err := hdl.timetrackingBiz.CreateTimeTracking(ctx, input)
	if err != nil {
		return nil, err
	}

	resp := &core.TimeTracking{
		Id:            timetracking.ID,
		CharacterId:   timetracking.CharacterID,
		CategoryId:    timetracking.CategoryID,
		ReferenceId:   timetracking.ReferenceID,
		ReferenceType: req.ReferenceType,
		StartTime:     timetracking.StartTime.Unix(),
		EndTime:       nil,
	}

	return resp, nil
}

func (hdl *CoreHandler) UpdateTimeTracking(ctx context.Context, req *common.EmptyReq) (*core.TimeTracking, error) {
	timetracking, err := hdl.timetrackingBiz.UpdateTimeTracking(ctx)
	if err != nil {
		return nil, err
	}

	var referenceType *core.ReferenceType
	if timetracking.ReferenceType != nil {
		referenceType = lo.ToPtr(core.ReferenceType(core.ReferenceType_value[string(*timetracking.ReferenceType)]))
	}

	resp := &core.TimeTracking{
		Id:            timetracking.ID,
		CharacterId:   timetracking.CharacterID,
		CategoryId:    timetracking.CategoryID,
		ReferenceId:   timetracking.ReferenceID,
		ReferenceType: referenceType,
		StartTime:     timetracking.StartTime.Unix(),
		EndTime:       nil,
	}

	return resp, nil
}

// TODO: @Fbin243 implements the following methods later
func (hdl *CoreHandler) GetCurrencoreimeTracking(context.Context, *common.EmptyReq) (*core.TimeTracking, error) {
	return nil, errors.New("not implemented")
}

func (hdl *CoreHandler) GecoreotalCurrencoreimeTracking(context.Context, *core.TotalTimeTrackingReq) (*core.TotalTimeTrackingResp, error) {
	return nil, errors.New("not implemented")
}
