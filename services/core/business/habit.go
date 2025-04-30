package business

type HabitBusiness struct {
	permBiz          IPermissionBusiness
	habitRepo        IHabitRepo
	habitLogRepo     IHabitLogRepo
	cateRepo         ICategoryRepo
	timetrackingRepo ITimeTrackingRepo
}

func NewHabitBusiness(
	permBiz IPermissionBusiness,
	habitRepo IHabitRepo,
	habitLogRepo IHabitLogRepo,
	cateRepo ICategoryRepo,
	timetrackingRepo ITimeTrackingRepo,
) *HabitBusiness {
	return &HabitBusiness{
		permBiz:          permBiz,
		habitRepo:        habitRepo,
		habitLogRepo:     habitLogRepo,
		cateRepo:         cateRepo,
		timetrackingRepo: timetrackingRepo,
	}
}
