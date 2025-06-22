package business

type CharacterBusiness struct {
	CharacterRepo    ICharacterRepo
	ProfileRepo      IProfileRepo
	GoalRepo         IGoalRepo
	MetricRepo       IMetricRepo
	CategoryRepo     ICategoryRepo
	TimeTrackingRepo ITimeTrackingRepo
	HabitRepo        IHabitRepo
	HabitLogRepo     IHabitLogRepo
	TaskRepo         ITaskRepo
	TaskSessionRepo  ITaskSessionRepo
	Cache            ICache
}

func NewCharacterBusiness(
	characterRepo ICharacterRepo,
	profileRepo IProfileRepo,
	goalRepo IGoalRepo,
	metricRepo IMetricRepo,
	cateRepo ICategoryRepo,
	timetrackRepo ITimeTrackingRepo,
	habitRepo IHabitRepo,
	habitLogRepo IHabitLogRepo,
	taskRepo ITaskRepo,
	taskSessionRepo ITaskSessionRepo,
	cache ICache,
) *CharacterBusiness {
	return &CharacterBusiness{
		CharacterRepo:    characterRepo,
		ProfileRepo:      profileRepo,
		GoalRepo:         goalRepo,
		MetricRepo:       metricRepo,
		CategoryRepo:     cateRepo,
		TimeTrackingRepo: timetrackRepo,
		HabitRepo:        habitRepo,
		HabitLogRepo:     habitLogRepo,
		TaskRepo:         taskRepo,
		TaskSessionRepo:  taskSessionRepo,
		Cache:            cache,
	}
}
