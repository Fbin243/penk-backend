package business

import (
	"context"
	"fmt"
)

// AddDeviceToken adds a new device token to the user's profile
func (biz *NotificationBusiness) RegisterDeviceToken(ctx context.Context, profileID, token, deviceID, platform string) (bool, error) {
	err := biz.DevicesTokenRepo.UpsertDeviceToken(ctx, profileID, token, deviceID, platform)
	if err != nil {
		return false, fmt.Errorf("failed to add Devices Token: %v", err)
	}

	return true, nil
}

// RemoveDeviceToken removes a device token from the user's profile
func (biz *NotificationBusiness) RemoveDeviceToken(ctx context.Context, profileID, token string) (bool, error) {
	err := biz.DevicesTokenRepo.RemoveDeviceToken(ctx, profileID, token)
	if err != nil {
		return false, fmt.Errorf("failed to remove Devices Token: %v", err)
	}

	return true, nil
}
