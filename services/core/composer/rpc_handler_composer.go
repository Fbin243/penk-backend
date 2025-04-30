package composer

import (
	"tenkhours/services/core/transport/rpc"
)

func ComposeRPCHandler() *rpc.CoreHandler {
	composer = GetComposer()
	return rpc.NewCoreHandler(composer.ProfileBiz, composer.CharacterBiz, composer.GoalBiz, composer.TimeTrackingBiz, composer.TaskBiz, composer.MetricBiz, composer.CategoryBiz, composer.HabitBiz)
}
