package composer

import (
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/analytic/business"
	mongorepo "tenkhours/services/analytic/repo/mongo"
)

type Composer struct {
	AnalyticBiz business.IAnalyticBusiness
}

var composer *Composer

func GetComposer() *Composer {
	if composer != nil {
		return composer
	}

	// Databases
	db := mongodb.GetDBManager().DB

	// Repositories
	timetrackingRepo := mongorepo.NewTimeTrackingRepo(db)

	// Business
	analyticBiz := business.NewAnalyticBusiness(timetrackingRepo)

	return &Composer{
		AnalyticBiz: analyticBiz,
	}
}
