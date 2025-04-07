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
	ProfileBiz      business.IProfileBusiness
	CharacterBiz    business.ICharacterBusiness
	GoalBiz         business.IGoalBusiness
	CategoryBiz     business.ICategoryBusiness
	MetricBiz       business.IMetricBusiness
	HabitBiz        business.IHabitBusiness
	TimeTrackingBiz business.ITimeTrackingBusiness

	CharacterRepo    business.ICharacterRepo
	CategoryRepo     business.ICategoryRepo
	MetricRepo       business.IMetricRepo
	TimeTrackingRepo business.ITimeTrackingRepo
	HabitRepo        business.IHabitRepo
	HabitLogRepo     business.IHabitLogRepo

	CurrencyConn *grpc.ClientConn
	AnalyticConn *grpc.ClientConn
	NotiConn     *grpc.ClientConn
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
	timetrackingRepo := mongorepo.NewTimeTrackingRepo(mongodb)
	habitRepo := mongorepo.NewHabitRepo(mongodb)
	habitLogRepo := mongorepo.NewHabitLogRepo(mongodb)

	// RPC Clients
	currencyClient, currencyConn := ComposeCurrencyClient()
	notiClient, notiConn := ComposeNotificationClient()

	// Business
	permBiz := business.NewPermissionBusiness(profileRepo, characterRepo, categoryRepo, metricRepo, goalRepo, habitRepo, timetrackingRepo)
	profileBiz := business.NewProfileBusiness(permBiz, profileRepo, characterRepo, categoryRepo, metricRepo, goalRepo, habitRepo, habitLogRepo, timetrackingRepo, currencyClient, redisRepo)
	characterBiz := business.NewCharacterBusiness(characterRepo, profileRepo, goalRepo, metricRepo, categoryRepo, timetrackingRepo, redisRepo)
	goalBiz := business.NewGoalBusiness(permBiz, goalRepo, characterRepo, categoryRepo, metricRepo)
	catgoryBiz := business.NewCategoryBusiness(permBiz, categoryRepo, metricRepo, timetrackingRepo)
	metricBiz := business.NewMetricBusiness(permBiz, metricRepo, categoryRepo)
	habitBiz := business.NewHabitBusiness(permBiz, habitRepo, habitLogRepo, categoryRepo, timetrackingRepo)
	timetrackingBiz := business.NewTimeTrackingBusiness(permBiz, notiClient, redisRepo, habitRepo, habitLogRepo, timetrackingRepo)

	return &Composer{
		ProfileBiz:      profileBiz,
		CharacterBiz:    characterBiz,
		GoalBiz:         goalBiz,
		CategoryBiz:     catgoryBiz,
		MetricBiz:       metricBiz,
		HabitBiz:        habitBiz,
		TimeTrackingBiz: timetrackingBiz,

		CharacterRepo:    characterRepo,
		CategoryRepo:     categoryRepo,
		MetricRepo:       metricRepo,
		TimeTrackingRepo: timetrackingRepo,
		HabitRepo:        habitRepo,
		HabitLogRepo:     habitLogRepo,

		CurrencyConn: currencyConn,
		AnalyticConn: notiConn,
		NotiConn:     notiConn,
	}
}

func (c *Composer) Close() {
	c.AnalyticConn.Close()
	c.CurrencyConn.Close()
	c.NotiConn.Close()
}
