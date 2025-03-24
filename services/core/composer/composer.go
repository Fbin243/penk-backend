package composer

import (
	mongodb "tenkhours/pkg/db/mongo"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/services/core/business"
	mongorepo "tenkhours/services/core/repo/mongo"
	redisrepo "tenkhours/services/core/repo/redis"
	timetrackrepo "tenkhours/services/timetracking/repo/mongo"

	"google.golang.org/grpc"
)

type Composer struct {
	ProfileBiz       business.IProfileBusiness
	CharacaterBiz    business.ICharacterBusiness
	GoalBiz          business.IGoalBusiness
	CategoryBiz      business.ICategoryBusiness
	MetricBiz        business.IMetricBusiness
	CharacterRepo    business.ICharacterRepo
	CategoryRepo     business.ICategoryRepo
	MetricRepo       business.IMetricRepo
	TimeTrackingRepo business.ITimeTrackingRepo
	CurrencyConn     *grpc.ClientConn
	AnalyticConn     *grpc.ClientConn
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
	redisRepo := redisrepo.NewRedisRepo(redisClient)
	categoryRepo := mongorepo.NewCategoryRepo(mongodb)
	metricRepo := mongorepo.NewMetricRepo(mongodb)
	timetrackingRepo := timetrackrepo.NewTimeTrackingRepo(mongodb)

	// RPC Clients
	currencyClient, currencyConn := ComposeCurrencyClient()

	// Business
	profileBiz := business.NewProfileBusiness(profileRepo, characterRepo, categoryRepo, metricRepo, goalRepo, timetrackingRepo, currencyClient, redisRepo)
	characterBiz := business.NewCharacterBusiness(characterRepo, profileRepo, goalRepo, metricRepo, categoryRepo, timetrackingRepo)
	goalBiz := business.NewGoalBusiness(goalRepo, characterRepo, categoryRepo, metricRepo)
	catgoryBiz := business.NewCategoryBusiness(categoryRepo, characterRepo, metricRepo, timetrackingRepo)
	metricBiz := business.NewMetricBusiness(metricRepo, characterRepo, categoryRepo)

	return &Composer{
		ProfileBiz:       profileBiz,
		CharacaterBiz:    characterBiz,
		GoalBiz:          goalBiz,
		CategoryBiz:      catgoryBiz,
		MetricBiz:        metricBiz,
		CharacterRepo:    characterRepo,
		CategoryRepo:     categoryRepo,
		MetricRepo:       metricRepo,
		TimeTrackingRepo: timetrackingRepo,
		CurrencyConn:     currencyConn,
	}
}
