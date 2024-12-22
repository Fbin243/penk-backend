package composer

import (
	"tenkhours/pkg/db"
	coreRepo "tenkhours/services/core/repo"
	"tenkhours/services/currency/business"
	"tenkhours/services/currency/graph"
	"tenkhours/services/currency/repo"
)

func ComposeGraphQLResolver() *graph.Resolver {
	// Init dependencies and perform DI manually
	mongodb := db.GetDBManager().DB
	fishRepo := repo.NewFishRepo(mongodb)
	redisClient := db.GetRedisClient()
	profilesRepo := coreRepo.NewProfilesRepo(mongodb, redisClient)
	charactersRepo := coreRepo.NewCharactersRepo(mongodb)
	fishBiz := business.NewFishBusiness(fishRepo, charactersRepo, profilesRepo, redisClient)

	return &graph.Resolver{
		FishBusiness: fishBiz,
	}
}
