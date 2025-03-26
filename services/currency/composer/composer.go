package composer

import (
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/currency/business"
	mongorepo "tenkhours/services/currency/repo/mongo"

	"google.golang.org/grpc"
)

type Composer struct {
	DailyRewardRepo business.IRewardRepo
	CoreClient      business.ICoreClient
	CoreClientConn  *grpc.ClientConn
	RewardBiz       business.IRewardBusiness
}

var composer *Composer

func GetComposer() *Composer {
	if composer != nil {
		return composer
	}

	// Database
	db := mongodb.GetDBManager().DB

	// Repository
	rewardRepo := mongorepo.NewRewardRepo(db)

	// RPC Client
	coreClient, coreClientConn := ComposeCoreClient()

	// Business
	rewardBiz := business.NewRewardBusiness(rewardRepo, coreClient)

	return &Composer{
		CoreClient:     coreClient,
		CoreClientConn: coreClientConn,
		RewardBiz:      rewardBiz,
	}
}
