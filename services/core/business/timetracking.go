package business

type TimeTrackingBusiness struct {
	permBiz          IPermissionBusiness
	notiClient       INotificationClient
	habitRepo        IHabitRepo
	habitLogRepo     IHabitLogRepo
	taskRepo         ITaskRepo
	timetrackingRepo ITimeTrackingRepo
}

func NewTimeTrackingBusiness(
	permissionBusiness IPermissionBusiness,
	notiClient INotificationClient,
	habitRepo IHabitRepo,
	habitLogRepo IHabitLogRepo,
	timetrackingRepo ITimeTrackingRepo,
) *TimeTrackingBusiness {
	return &TimeTrackingBusiness{
		permBiz:          permissionBusiness,
		notiClient:       notiClient,
		habitRepo:        habitRepo,
		habitLogRepo:     habitLogRepo,
		timetrackingRepo: timetrackingRepo,
	}
}
