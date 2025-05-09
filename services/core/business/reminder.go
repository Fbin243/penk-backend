package business

type ReminderBusiness struct {
	reminderRepo  IReminderRepo
	permBiz       IPermissionBusiness
	reminderCache IReminderCache
}

func NewReminderBusiness(reminderRepo IReminderRepo, permBiz IPermissionBusiness, reminderCache IReminderCache) *ReminderBusiness {
	return &ReminderBusiness{
		reminderRepo:  reminderRepo,
		permBiz:       permBiz,
		reminderCache: reminderCache,
	}
}
