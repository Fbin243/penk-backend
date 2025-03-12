package composer

import (
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/services/timetracking/business"
	redisrepo "tenkhours/services/timetracking/repo/redis"
	"tenkhours/services/timetracking/repo/rpc"

	"google.golang.org/grpc"
)

type Composer struct {
	redisRepo          *redisrepo.RedisRepo
	coreClient         *rpc.CoreClient
	notificationClient *rpc.NotificationClient
	currencyClient     *rpc.CurrencyClient
	coreConn           *grpc.ClientConn
	notificationConn   *grpc.ClientConn
	currencyConn       *grpc.ClientConn
	timetrackingBiz    *business.TimeTrackingBusiness
}

var composer *Composer

func GetComposer() *Composer {
	if composer != nil {
		return composer
	}

	redisClient := rdb.GetRedisClient()
	redisRepo := redisrepo.NewRedisRepo(redisClient)
	coreClient, coreConn := ComposeCoreClient()
	notiClient, notiConn := ComposeNotificationClient()
	currencyClient, currencyConn := ComposeCurrencyClient()
	timetrackingBiz := business.NewTimeTrackingsBusiness(coreClient, currencyClient, notiClient, redisRepo)

	return &Composer{
		redisRepo:          redisRepo,
		coreClient:         coreClient,
		notificationClient: notiClient,
		currencyClient:     currencyClient,
		coreConn:           coreConn,
		notificationConn:   notiConn,
		currencyConn:       currencyConn,
		timetrackingBiz:    timetrackingBiz,
	}
}

func (c *Composer) Close() {
	c.coreConn.Close()
	c.currencyConn.Close()
	c.notificationConn.Close()
}
