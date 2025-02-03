package composer

import (
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/currency/business"
	mongorepo "tenkhours/services/currency/repo/mongo"

	"google.golang.org/grpc"
)

type Composer struct {
	FishRepo       business.IFishRepo
	CoreClient     business.ICoreClient
	CoreClientConn *grpc.ClientConn
	CurrencyBiz    business.ICurrencyBusiness
}

var composer *Composer

func GetComposer() *Composer {
	if composer != nil {
		return composer
	}

	// Database
	db := mongodb.GetDBManager().DB

	// Repository
	fishRepo := mongorepo.NewFishRepo(db)

	// RPC Client
	coreClient, coreClientConn := ComposeCoreClient()

	// Business
	currencyBiz := business.NewCurrencyBusiness(fishRepo, coreClient)

	return &Composer{
		FishRepo:       fishRepo,
		CoreClient:     coreClient,
		CoreClientConn: coreClientConn,
		CurrencyBiz:    currencyBiz,
	}
}
