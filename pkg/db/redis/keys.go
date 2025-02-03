package rdb

const (
	AuthSessionKey    = "auth_session_"
	CapturedRecordKey = "captured_record_"
	FishKey           = "fish_"
	TimeTrackingKey   = "time_tracking_"
)

func GetAuthSessionKey(firebaseUID string) string {
	return AuthSessionKey + firebaseUID
}

func GetCapturedRecordKey(profileID string) string {
	return CapturedRecordKey + profileID
}

func GetFishKey(profileID string) string {
	return FishKey + profileID
}

func GetTimeTrackingKey(profileID string) string {
	return TimeTrackingKey + profileID
}
