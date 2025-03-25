package composer

import (
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/analytic/business"
	"tenkhours/services/analytic/repo/rpc"
	tt "tenkhours/services/timetracking/repo/mongo"

	"google.golang.org/grpc"
)

type Composer struct {
	AnalyticBiz    business.IAnalyticBusiness
	CoreClient     *rpc.CoreClient
	CoreClientConn *grpc.ClientConn
}

var composer *Composer

func GetComposer() *Composer {
	if composer != nil {
		return composer
	}

	// Databases
	db := mongodb.GetDBManager().DB

	// Repositories
	timetrackingRepo := tt.NewTimeTrackingRepo(db)

	// RPC Clients
	coreClient, conn := ComposeCoreClient()

	// Business
	analyticBiz := business.NewAnalyticBusiness(coreClient, timetrackingRepo)

	return &Composer{
		AnalyticBiz:    analyticBiz,
		CoreClient:     coreClient,
		CoreClientConn: conn,
	}
}
