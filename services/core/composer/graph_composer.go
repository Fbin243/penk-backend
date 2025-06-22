package composer

import (
	"tenkhours/services/core/transport/graph"
)

func ComposeGraphQLResolver() *graph.Resolver {
	composer := GetComposer()
	return &graph.Resolver{
		ProfileBusiness:      composer.ProfileBiz,
		CharacterBusiness:    composer.CharacterBiz,
		GoalBusiness:         composer.GoalBiz,
		HabitBusiness:        composer.HabitBiz,
		MetricBusiness:       composer.MetricBiz,
		CategoryBusiness:     composer.CategoryBiz,
		TimeTrackingBusiness: composer.TimeTrackingBiz,
		TaskBusiness:         composer.TaskBiz,

		CharacterRepo:    composer.CharacterRepo,
		MetricRepo:       composer.MetricRepo,
		CategoryRepo:     composer.CategoryRepo,
		TimeTrackingRepo: composer.TimeTrackingRepo,
		HabitRepo:        composer.HabitRepo,
		HabitLogRepo:     composer.HabitLogRepo,
		TaskRepo:         composer.TaskRepo,
		TaskSessionRepo:  composer.TaskSessionRepo,
		GoalRepo:         composer.GoalRepo,
	}
}
