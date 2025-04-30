package composer

import (
	mongodb "tenkhours/pkg/db/mongo"
	rdb "tenkhours/pkg/db/redis"
	"tenkhours/services/core/business"
	mongorepo "tenkhours/services/core/repo/mongo"
	redisrepo "tenkhours/services/core/repo/redis"
	rewardrepo "tenkhours/services/currency/repo/mongo"

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
	TaskBiz         business.ITaskBusiness

	CharacterRepo    business.ICharacterRepo
	CategoryRepo     business.ICategoryRepo
	MetricRepo       business.IMetricRepo
	TimeTrackingRepo business.ITimeTrackingRepo
	HabitRepo        business.IHabitRepo
	HabitLogRepo     business.IHabitLogRepo
	TaskRepo         business.ITaskRepo
	TaskSessionRepo  business.ITaskSessionRepo
	GoalRepo         business.IGoalRepo

	CurrencyConn *grpc.ClientConn
	AnalyticConn *grpc.ClientConn
	NotiConn     *grpc.ClientConn

	RewardRepo business.IRewardRepo
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
	taskRepo := mongorepo.NewTaskRepo(mongodb)
	taskSessionRepo := mongorepo.NewTaskSessionRepo(mongodb)
	rewardRepo := rewardrepo.NewRewardRepo(mongodb)

	// RPC Clients
	currencyClient, currencyConn := ComposeCurrencyClient()
	notiClient, notiConn := ComposeNotificationClient()

	// Business
	permBiz := business.NewPermissionBusiness(profileRepo, characterRepo, categoryRepo, metricRepo, goalRepo, habitRepo, timetrackingRepo, taskRepo)
	profileBiz := business.NewProfileBusiness(permBiz, profileRepo, characterRepo, categoryRepo, metricRepo, goalRepo, habitRepo, habitLogRepo, timetrackingRepo, taskRepo, taskSessionRepo, currencyClient, redisRepo, rewardRepo)
	characterBiz := business.NewCharacterBusiness(characterRepo, profileRepo, goalRepo, metricRepo, categoryRepo, timetrackingRepo, habitRepo, habitLogRepo, taskRepo, taskSessionRepo, redisRepo)
	goalBiz := business.NewGoalBusiness(permBiz, goalRepo, characterRepo, categoryRepo, metricRepo)
	catgoryBiz := business.NewCategoryBusiness(permBiz, categoryRepo, metricRepo, timetrackingRepo, habitRepo, taskRepo)
	metricBiz := business.NewMetricBusiness(permBiz, metricRepo, categoryRepo)
	habitBiz := business.NewHabitBusiness(permBiz, habitRepo, habitLogRepo, categoryRepo, timetrackingRepo)
	timetrackingBiz := business.NewTimeTrackingBusiness(permBiz, notiClient, habitRepo, habitLogRepo, timetrackingRepo)
	taskBiz := business.NewTaskBusiness(permBiz, taskRepo, taskSessionRepo, timetrackingRepo)

	return &Composer{
		ProfileBiz:      profileBiz,
		CharacterBiz:    characterBiz,
		GoalBiz:         goalBiz,
		CategoryBiz:     catgoryBiz,
		MetricBiz:       metricBiz,
		HabitBiz:        habitBiz,
		TimeTrackingBiz: timetrackingBiz,
		TaskBiz:         taskBiz,

		CharacterRepo:    characterRepo,
		CategoryRepo:     categoryRepo,
		MetricRepo:       metricRepo,
		TimeTrackingRepo: timetrackingRepo,
		HabitRepo:        habitRepo,
		HabitLogRepo:     habitLogRepo,
		TaskRepo:         taskRepo,
		TaskSessionRepo:  taskSessionRepo,
		GoalRepo:         goalRepo,

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
