package composer

import (
	"tenkhours/pkg/db"
	analyticsRepo "tenkhours/services/analytics/repo"
	"tenkhours/services/core/business"
	"tenkhours/services/core/graph"
	"tenkhours/services/core/repo"
	fishRepo "tenkhours/services/currency/repo"
)

func ComposeGraphQLResolver() *graph.Resolver {
	// Init dependencies and perform DI manually
	mongodb := db.GetDBManager().DB
	redisClient := db.GetRedisClient()
	profilesRepo := repo.NewProfilesRepo(mongodb, redisClient)
	charactersRepo := repo.NewCharactersRepo(mongodb)
	fishRepo := fishRepo.NewFishRepo(mongodb)
	goalsRepo := repo.NewGoalsRepo(mongodb)
	templatesRepo := repo.NewTemplatesRepo(mongodb)
	templateCategoriesRepo := repo.NewTemplateCategoriesRepo(mongodb)

	// TODO: Temporary inject analyticsRepos into profilesBiz for deleting related data
	capturedRepordsRepo := analyticsRepo.NewCapturedRecordsRepo(mongodb)
	snapshotsRepo := analyticsRepo.NewSnapshotsRepo(mongodb)

	profilesBiz := business.NewProfilesBusiness(profilesRepo, fishRepo, charactersRepo, capturedRepordsRepo, snapshotsRepo, redisClient)
	charactersBiz := business.NewCharactersBusiness(charactersRepo, profilesRepo, goalsRepo)
	goalsBiz := business.NewGoalsBusiness(goalsRepo, charactersRepo)
	templatesBiz := business.NewTemplatesBusiness(templatesRepo, templateCategoriesRepo)

	return &graph.Resolver{
		ProfilesBusiness:   profilesBiz,
		CharactersBusiness: charactersBiz,
		GoalsBusiness:      goalsBiz,
		TemplatesBusiness:  templatesBiz,
	}
}
