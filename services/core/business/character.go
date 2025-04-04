package business

type CharacterBusiness struct {
	CharacterRepo    ICharacterRepo
	ProfileRepo      IProfileRepo
	GoalRepo         IGoalRepo
	MetricRepo       IMetricRepo
	CategoryRepo     ICategoryRepo
	TimeTrackingRepo ITimeTrackingRepo
	Cache            ICache
}

func NewCharacterBusiness(
	characterRepo ICharacterRepo,
	profileRepo IProfileRepo,
	goalRepo IGoalRepo,
	metricRepo IMetricRepo,
	cateRepo ICategoryRepo,
	timetrackRepo ITimeTrackingRepo,
	cache ICache,
) *CharacterBusiness {
	return &CharacterBusiness{
		CharacterRepo:    characterRepo,
		ProfileRepo:      profileRepo,
		GoalRepo:         goalRepo,
		MetricRepo:       metricRepo,
		CategoryRepo:     cateRepo,
		TimeTrackingRepo: timetrackRepo,
		Cache:            cache,
	}
}
