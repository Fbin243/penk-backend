package business

type ReminderBusiness struct {
	reminderRepo IReminderRepo
	permBiz      IPermissionBusiness
}

func NewReminderBusiness(reminderRepo IReminderRepo, permBiz IPermissionBusiness) *ReminderBusiness {
	return &ReminderBusiness{
		reminderRepo: reminderRepo,
		permBiz:      permBiz,
	}
}
