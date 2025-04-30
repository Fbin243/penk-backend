package composer

import (
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/services/currency/business"
	mongorepo "tenkhours/services/currency/repo/mongo"
)

type Composer struct {
	DailyRewardRepo business.IRewardRepo
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

	// Business
	rewardBiz := business.NewRewardBusiness(rewardRepo)

	return &Composer{
		RewardBiz: rewardBiz,
	}
}
