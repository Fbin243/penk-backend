package rdb

const (
	AuthSessionKey  = "auth_session_"
	FishKey         = "fish_"
	TimeTrackingKey = "time_tracking_"
	DeviceTokenKey  = "device_token_"
)

func GetAuthSessionKey(firebaseUID string) string {
	return AuthSessionKey + firebaseUID
}

func GetFishKey(profileID string) string {
	return FishKey + profileID
}

func GetTimeTrackingKey(profileID string) string {
	return TimeTrackingKey + profileID
}

func GetDeviceTokenKey(profileID string) string {
	return DeviceTokenKey + profileID
}
