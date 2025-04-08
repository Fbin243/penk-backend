package business

type CategoryBusiness struct {
	permBiz          IPermissionBusiness
	cateRepo         ICategoryRepo
	metricRepo       IMetricRepo
	timetrackingRepo ITimeTrackingRepo
	habitRepo        IHabitRepo
	taskRepo         ITaskRepo
}

func NewCategoryBusiness(
	permBiz IPermissionBusiness,
	cateRepo ICategoryRepo,
	metricRepo IMetricRepo,
	timetrackingRepo ITimeTrackingRepo,
	habitRepo IHabitRepo,
	taskRepo ITaskRepo,
) *CategoryBusiness {
	return &CategoryBusiness{
		permBiz:          permBiz,
		cateRepo:         cateRepo,
		metricRepo:       metricRepo,
		timetrackingRepo: timetrackingRepo,
		habitRepo:        habitRepo,
		taskRepo:         taskRepo,
	}
}
