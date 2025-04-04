package business

type TimeTrackingBusiness struct {
	permBiz          IPermissionBusiness
	notiClient       INotificationClient
	cache            ICache
	habitRepo        IHabitRepo
	habitLogRepo     IHabitLogRepo
	timetrackingRepo ITimeTrackingRepo
}

func NewTimeTrackingBusiness(
	permissionBusiness IPermissionBusiness,
	notiClient INotificationClient,
	cache ICache,
	habitRepo IHabitRepo,
	habitLogRepo IHabitLogRepo,
	timetrackingRepo ITimeTrackingRepo,
) *TimeTrackingBusiness {
	return &TimeTrackingBusiness{
		permBiz:          permissionBusiness,
		notiClient:       notiClient,
		cache:            cache,
		habitRepo:        habitRepo,
		habitLogRepo:     habitLogRepo,
		timetrackingRepo: timetrackingRepo,
	}
}
