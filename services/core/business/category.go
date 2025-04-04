package business

type CategoryBusiness struct {
	permBiz          IPermissionBusiness
	cateRepo         ICategoryRepo
	metricRepo       IMetricRepo
	timetrackingRepo ITimeTrackingRepo
}

func NewCategoryBusiness(
	permBiz IPermissionBusiness,
	cateRepo ICategoryRepo,
	metricRepo IMetricRepo,
	timetrackingRepo ITimeTrackingRepo,
) *CategoryBusiness {
	return &CategoryBusiness{
		permBiz:          permBiz,
		cateRepo:         cateRepo,
		metricRepo:       metricRepo,
		timetrackingRepo: timetrackingRepo,
	}
}
