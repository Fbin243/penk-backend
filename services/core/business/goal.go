package business

type GoalBusiness struct {
	permBiz       IPermissionBusiness
	goalRepo      IGoalRepo
	characterRepo ICharacterRepo
	categoryRepo  ICategoryRepo
	metricRepo    IMetricRepo
}

func NewGoalBusiness(
	permBiz IPermissionBusiness,
	goalRepo IGoalRepo,
	characterRepo ICharacterRepo,
	categoryRepo ICategoryRepo,
	metricRepo IMetricRepo,
) *GoalBusiness {
	return &GoalBusiness{
		permBiz:       permBiz,
		goalRepo:      goalRepo,
		characterRepo: characterRepo,
		categoryRepo:  categoryRepo,
		metricRepo:    metricRepo,
	}
}
