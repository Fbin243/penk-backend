package business

type MetricBusiness struct {
	permBiz    IPermissionBusiness
	metricRepo IMetricRepo
	cateRepo   ICategoryRepo
}

func NewMetricBusiness(
	permBiz IPermissionBusiness,
	metricRepo IMetricRepo,
	cateRepo ICategoryRepo,
) *MetricBusiness {
	return &MetricBusiness{
		permBiz:    permBiz,
		metricRepo: metricRepo,
		cateRepo:   cateRepo,
	}
}
