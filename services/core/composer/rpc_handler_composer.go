package composer

import (
	"tenkhours/pkg/db"
	analyticsRepo "tenkhours/services/analytics/repo"
	"tenkhours/services/core/business"
	"tenkhours/services/core/repo"
	"tenkhours/services/core/rpc"
	fishRepo "tenkhours/services/currency/repo"
)

func ComposeRPCHandler() *rpc.RPCHandler {
	// Init dependencies and perform DI manually
	mongodb := db.GetDBManager().DB
	redisClient := db.GetRedisClient()
	profilesRepo := repo.NewProfilesRepo(mongodb, redisClient)
	charactersRepo := repo.NewCharactersRepo(mongodb)
	fishRepo := fishRepo.NewFishRepo(mongodb)
	goalsRepo := repo.NewGoalsRepo(mongodb)

	// TODO: Temporary inject analyticsRepos into profilesBiz for deleting related data
	capturedRepordsRepo := analyticsRepo.NewCapturedRecordsRepo(mongodb)
	snapshotsRepo := repo.NewSnapshotsRepo(mongodb)

	profilesBiz := business.NewProfilesBusiness(profilesRepo, fishRepo, charactersRepo, capturedRepordsRepo, snapshotsRepo, redisClient)
	charactersBiz := business.NewCharactersBusiness(charactersRepo, profilesRepo, goalsRepo)

	return rpc.NewRPCHandler(profilesBiz, charactersBiz)
}
