package composer

import "tenkhours/services/notification/business"

type Composer struct {
	NotificationBiz business.INotificationBusiness
}

var composer *Composer

func GetComposer() *Composer {
	if composer != nil {
		return composer
	}

	return &Composer{
		NotificationBiz: business.NewNotificationBusiness(),
	}
}
