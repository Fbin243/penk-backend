package business

type PermissionBusiness struct {
	ProfileRepo      IProfileRepo
	CharacterRepo    ICharacterRepo
	CategoryRepo     ICategoryRepo
	MetricRepo       IMetricRepo
	GoalRepo         IGoalRepo
	HabitRepo        IHabitRepo
	TimeTrackingRepo ITimeTrackingRepo
	TaskRepo         ITaskRepo
}

func NewPermissionBusiness(
	profileRepo IProfileRepo,
	characterRepo ICharacterRepo,
	categoryRepo ICategoryRepo,
	metricRepo IMetricRepo,
	goalRepo IGoalRepo,
	habitRepo IHabitRepo,
	timeTrackingRepo ITimeTrackingRepo,
	taskRepo ITaskRepo,
) *PermissionBusiness {
	return &PermissionBusiness{
		ProfileRepo:      profileRepo,
		CharacterRepo:    characterRepo,
		CategoryRepo:     categoryRepo,
		MetricRepo:       metricRepo,
		GoalRepo:         goalRepo,
		HabitRepo:        habitRepo,
		TimeTrackingRepo: timeTrackingRepo,
		TaskRepo:         taskRepo,
	}
}
