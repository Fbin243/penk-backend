package business

type TaskBusiness struct {
	permBiz          IPermissionBusiness
	taskRepo         ITaskRepo
	taskSessionRepo  ITaskSessionRepo
	timetrackingRepo ITimeTrackingRepo
}

func NewTaskBusiness(
	permBiz IPermissionBusiness,
	taskRepo ITaskRepo,
	taskSessionRepo ITaskSessionRepo,
	timetrackingRepo ITimeTrackingRepo,
) *TaskBusiness {
	return &TaskBusiness{
		permBiz:          permBiz,
		taskRepo:         taskRepo,
		taskSessionRepo:  taskSessionRepo,
		timetrackingRepo: timetrackingRepo,
	}
}
