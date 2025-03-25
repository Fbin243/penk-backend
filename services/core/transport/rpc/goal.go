package rpc

import (
	"context"
	"fmt"

	"tenkhours/pkg/utils"
	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"

	"github.com/jinzhu/copier"
)

func (hdl *CoreHandler) UpsertGoal(ctx context.Context, req *core.GoalInput) (*core.Goal, error) {
	resp := &core.Goal{}

	goalInput := entity.GoalInput{}

	// Mapping rpc to entity
	copier.CopyWithOption(goalInput, req, copier.Option{
		IgnoreEmpty: true,
		Converters: []copier.TypeConverter{
			{
				SrcType: core.MetricCondition(0),
				DstType: entity.MetricCondition(fmt.Sprint(0)),
				Fn: func(src any) (any, error) {
					return entity.MetricCondition(
						core.MetricCondition_name[int32(src.(core.MetricCondition).Number())],
					), nil
				},
			},
		},
	})
	goalInput.EndTime = utils.UnixToTime(req.EndTime)
	goalInput.StartTime = utils.UnixToTime(req.StartTime)

	// Call business logic
	goal, err := hdl.goalBiz.UpsertGoal(ctx, goalInput)
	if err != nil {
		return resp, err
	}

	// Mapping entity to rpc
	copier.Copy(resp, &goal)
	resp.Status = core.GoalStatus(core.GoalStatus_value[string(goal.Status)])
	resp.CreatedAt = goal.CreatedAt.Unix()
	resp.UpdatedAt = goal.UpdatedAt.Unix()
	resp.StartTime = goal.StartTime.Unix()
	resp.EndTime = goal.EndTime.Unix()
	copier.Copy(&resp.Metrics, &goal.Metrics)
	copier.Copy(&resp.Checkboxes, &goal.Checkboxes)

	return resp, nil
}
