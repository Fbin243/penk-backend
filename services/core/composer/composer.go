package composer

import (
	mongodb "tenkhours/pkg/db/mongo"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/services/core/business"
	mongorepo "tenkhours/services/core/repo/mongo"
	redisrepo "tenkhours/services/core/repo/redis"

	"google.golang.org/grpc"
)

type Composer struct {
	ProfileBiz    business.IProfileBusiness
	CharacaterBiz business.ICharacterBusiness
	GoalBiz       business.IGoalBusiness
	TemplateBiz   business.ITemplateBusiness
	CharacterRepo business.ICharacterRepo
	CurrencyConn  *grpc.ClientConn
	AnalyticConn  *grpc.ClientConn
}

var composer *Composer

func GetComposer() *Composer {
	if composer != nil {
		return composer
	}

	// Databases
	mongodb := mongodb.GetDBManager().DB
	redisClient := rdb.GetRedisClient()

	// Repositories
	profileRepo := mongorepo.NewProfileRepo(mongodb)
	characterRepo := mongorepo.NewCharacterRepo(mongodb)
	goalRepo := mongorepo.NewGoalRepo(mongodb)
	templateRepo := mongorepo.NewTemplateRepo(mongodb)
	templateTopicRepo := mongorepo.NewTemplateTopicRepo(mongodb)
	redisRepo := redisrepo.NewRedisRepo(redisClient)

	// RPC Clients
	currencyClient, currencyConn := ComposeCurrencyClient()
	analyticClient, analyticConn := ComposeAnalyticClient()

	// Business
	profileBiz := business.NewProfileBusiness(profileRepo, characterRepo, currencyClient, analyticClient, redisRepo)
	characterBiz := business.NewCharacterBusiness(characterRepo, profileRepo, goalRepo)
	goalBiz := business.NewGoalBusiness(goalRepo, characterRepo)
	templateBiz := business.NewTemplateBusiness(templateRepo, templateTopicRepo)

	return &Composer{
		ProfileBiz:    profileBiz,
		CharacaterBiz: characterBiz,
		GoalBiz:       goalBiz,
		TemplateBiz:   templateBiz,
		CharacterRepo: characterRepo,
		CurrencyConn:  currencyConn,
		AnalyticConn:  analyticConn,
	}
}
