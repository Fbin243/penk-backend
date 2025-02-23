package composer

import (
	mongodb "tenkhours/pkg/db/mongo"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/services/analytic/business"
	mongorepo "tenkhours/services/analytic/repo/mongo"
	"tenkhours/services/analytic/repo/redis"
	"tenkhours/services/analytic/repo/rpc"

	"google.golang.org/grpc"
)

type Composer struct {
	CapturedRecordRepo *mongorepo.CapturedRecordRepo
	AnalyticBiz        business.IAnalyticBusiness
	CoreClient         *rpc.CoreClient
	CoreClientConn     *grpc.ClientConn
}

var composer *Composer

func GetComposer() *Composer {
	if composer != nil {
		return composer
	}

	// Databases
	db := mongodb.GetDBManager().DB
	redisClient := rdb.GetRedisClient()

	// Repositories
	capturedRecordRepo := mongorepo.NewCapturedRecordRepo(db)
	redisRepo := redis.NewRedisRepo(redisClient)

	// RPC Clients
	coreClient, conn := ComposeCoreClient()

	// Business
	analyticBiz := business.NewAnalyticBusiness(capturedRecordRepo, coreClient, redisRepo)

	return &Composer{
		CapturedRecordRepo: capturedRecordRepo,
		AnalyticBiz:        analyticBiz,
		CoreClient:         coreClient,
		CoreClientConn:     conn,
	}
}
