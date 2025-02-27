package composer

import (
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/services/timetracking/business"
	redisrepo "tenkhours/services/timetracking/repo/redis"
	"tenkhours/services/timetracking/repo/rpc"

	"google.golang.org/grpc"
)

type Composer struct {
	redisRepo       *redisrepo.RedisRepo
	coreClient      *rpc.CoreClient
	currencyClient  *rpc.CurrencyClient
	coreConn        *grpc.ClientConn
	currencyConn    *grpc.ClientConn
	timetrackingBiz business.ITimeTrackingBusiness
}

var composer *Composer

func GetComposer() *Composer {
	if composer != nil {
		return composer
	}

	redisClient := rdb.GetRedisClient()
	redisRepo := redisrepo.NewRedisRepo(redisClient)
	coreClient, coreConn := ComposeCoreClient()
	currencyClient, currencyConn := ComposeCurrencyClient()
	timetrackingBiz := business.NewTimeTrackingsBusiness(coreClient, currencyClient, redisRepo)

	return &Composer{
		redisRepo:       redisRepo,
		coreClient:      coreClient,
		currencyClient:  currencyClient,
		coreConn:        coreConn,
		currencyConn:    currencyConn,
		timetrackingBiz: timetrackingBiz,
	}
}

func (c *Composer) Close() {
	c.coreConn.Close()
	c.currencyConn.Close()
}
