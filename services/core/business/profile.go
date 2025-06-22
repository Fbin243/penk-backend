package business

type ProfileBusiness struct {
	permBiz          IPermissionBusiness
	ProfileRepo      IProfileRepo
	CharacterRepo    ICharacterRepo
	CategoryRepo     ICategoryRepo
	MetricRepo       IMetricRepo
	GoalRepo         IGoalRepo
	HabitRepo        IHabitRepo
	HabitLogRepo     IHabitLogRepo
	TimeTrackingRepo ITimeTrackingRepo
	TaskRepo         ITaskRepo
	TaskSessionRepo  ITaskSessionRepo
	CurrencyClient   ICurrencyClient
	Cache            ICache
	RewardRepo       IRewardRepo
}

func NewProfileBusiness(
	permBiz IPermissionBusiness,
	profileRepo IProfileRepo,
	characterRepo ICharacterRepo,
	categoryRepo ICategoryRepo,
	metricRepo IMetricRepo,
	goalRepo IGoalRepo,
	habitRepo IHabitRepo,
	habitLogRepo IHabitLogRepo,
	timeTrackingRepo ITimeTrackingRepo,
	taskRepo ITaskRepo,
	taskSessionRepo ITaskSessionRepo,
	currencyClient ICurrencyClient,
	cache ICache,
	rewardRepo IRewardRepo,
) *ProfileBusiness {
	return &ProfileBusiness{
		permBiz:          permBiz,
		ProfileRepo:      profileRepo,
		CharacterRepo:    characterRepo,
		CategoryRepo:     categoryRepo,
		MetricRepo:       metricRepo,
		GoalRepo:         goalRepo,
		HabitRepo:        habitRepo,
		HabitLogRepo:     habitLogRepo,
		TimeTrackingRepo: timeTrackingRepo,
		TaskRepo:         taskRepo,
		TaskSessionRepo:  taskSessionRepo,
		CurrencyClient:   currencyClient,
		Cache:            cache,
		RewardRepo:       rewardRepo,
	}
}
