package entity

import "tenkhours/pkg/db/base"

type DevicesToken struct {
	*base.BaseEntity `                            bson:",inline"`
	ProfileID        string  `json:"profile_id,omitempty" bson:"profile_id"`
	Tokens           []Token `json:"tokens,omitempty"     bson:"tokens"`
}

type Token struct {
	DeviceID string `json:"device_id,omitempty" bson:"device_id"`
	Token    string `json:"token,omitempty"     bson:"token"`
	Platform string `json:"platform,omitempty"  bson:"platform"`
	CreateAt string `json:"create_at,omitempty" bson:"create_at"`
}
