package business

import "context"

type INotificationBusiness interface {
	AddEmailToWaitlist(ctx context.Context, email string) error
}
