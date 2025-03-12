package entity

type SendNotiReq struct {
	ProfileID string `json:"profileID,omitempty"`
	DeviceID  string `json:"deviceID,omitempty"`
	Title     string `json:"title,omitempty"`
	Body      string `json:"body,omitempty"`
}
