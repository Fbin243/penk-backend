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
	ProfileBiz    business.IProfileBusiness
	CharacterBiz  business.ICharacterBusiness
	GoalBiz       business.IGoalBusiness
	CategoryBiz   business.ICategoryBusiness
	MetricBiz     business.IMetricBusiness
	HabitBusiness business.IHabitBusiness

	CharacterRepo    business.ICharacterRepo
	CategoryRepo     business.ICategoryRepo
	MetricRepo       business.IMetricRepo
	TimeTrackingRepo business.ITimeTrackingRepo
	HabitRepo        business.IHabitRepo
	HabitLogRepo     business.IHabitLogRepo

	CurrencyConn *grpc.ClientConn
	AnalyticConn *grpc.ClientConn
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
	habitRepo := mongorepo.NewHabitRepo(mongodb)
	habitLogRepo := mongorepo.NewHabitLogRepo(mongodb)

	// RPC Clients
	currencyClient, currencyConn := ComposeCurrencyClient()

	// Business
	profileBiz := business.NewProfileBusiness(profileRepo, characterRepo, categoryRepo, metricRepo, goalRepo, timetrackingRepo, currencyClient, redisRepo)
	characterBiz := business.NewCharacterBusiness(characterRepo, profileRepo, goalRepo, metricRepo, categoryRepo, timetrackingRepo)
	goalBiz := business.NewGoalBusiness(goalRepo, characterRepo, categoryRepo, metricRepo)
	catgoryBiz := business.NewCategoryBusiness(categoryRepo, characterRepo, metricRepo, timetrackingRepo)
	metricBiz := business.NewMetricBusiness(metricRepo, characterRepo, categoryRepo)
	habitBusiness := business.NewHabitBusiness(habitRepo, habitLogRepo, characterRepo, categoryRepo)

	return &Composer{
		ProfileBiz:       profileBiz,
		CharacterBiz:     characterBiz,
		GoalBiz:          goalBiz,
		CategoryBiz:      catgoryBiz,
		MetricBiz:        metricBiz,
		CharacterRepo:    characterRepo,
		CategoryRepo:     categoryRepo,
		MetricRepo:       metricRepo,
		TimeTrackingRepo: timetrackingRepo,
		HabitRepo:        habitRepo,
		HabitLogRepo:     habitLogRepo,
		HabitBusiness:    habitBusiness,
		CurrencyConn:     currencyConn,
	}
}
