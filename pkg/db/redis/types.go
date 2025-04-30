package rdb

type AuthSession struct {
	FirebaseUID        string `json:"firebaseUID"`
	ProfileID          string `json:"profileID"`
	CurrentCharacterID string `json:"currentCharacterID"`
	DeviceID           string `json:"deviceID"`
}
