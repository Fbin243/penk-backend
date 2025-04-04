package business

type ProfileBusiness struct {
	permBiz          IPermissionBusiness
	ProfileRepo      IProfileRepo
	CharacterRepo    ICharacterRepo
	CategoryRepo     ICategoryRepo
	MetricRepo       IMetricRepo
	GoalRepo         IGoalRepo
	HabitRepo        IHabitRepo
	TimeTrackingRepo ITimeTrackingRepo
	CurrencyClient   ICurrencyClient
	Cache            ICache
}

func NewProfileBusiness(
	permBiz IPermissionBusiness,
	profileRepo IProfileRepo,
	characterRepo ICharacterRepo,
	categoryRepo ICategoryRepo,
	metricRepo IMetricRepo,
	goalRepo IGoalRepo,
	habitRepo IHabitRepo,
	timeTrackingRepo ITimeTrackingRepo,
	currencyClient ICurrencyClient,
	cache ICache,
) *ProfileBusiness {
	return &ProfileBusiness{
		permBiz:          permBiz,
		ProfileRepo:      profileRepo,
		CharacterRepo:    characterRepo,
		CategoryRepo:     categoryRepo,
		MetricRepo:       metricRepo,
		GoalRepo:         goalRepo,
		HabitRepo:        habitRepo,
		TimeTrackingRepo: timeTrackingRepo,
		CurrencyClient:   currencyClient,
		Cache:            cache,
	}
}
