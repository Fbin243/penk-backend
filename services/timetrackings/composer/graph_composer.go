package composer

import (
	"tenkhours/pkg/db"
	"tenkhours/services/core/repo"
	"tenkhours/services/timetrackings/business"
	"tenkhours/services/timetrackings/graph"

	fishRepo "tenkhours/services/currency/repo"
	timetrackingsRepo "tenkhours/services/timetrackings/repo"

	"google.golang.org/grpc"
)

func ComposeGraphQLResolver() (*graph.Resolver, *grpc.ClientConn) {
	mongodb := db.GetDBManager().DB
	redisClient := db.GetRedisClient()
	profilesRepo := repo.NewProfilesRepo(mongodb, redisClient)
	timetrackingsRepo := timetrackingsRepo.NewTimeTrackingsRepo(mongodb)
	fishRepo := fishRepo.NewFishRepo(mongodb)
	coreClient, conn := ComposeRPCClient()
	timetrackingsBiz := business.NewTimeTrackingsBusiness(timetrackingsRepo, coreClient, fishRepo, profilesRepo, redisClient)

	return &graph.Resolver{
		TimeTrackingsBusiness: timetrackingsBiz,
	}, conn
}
