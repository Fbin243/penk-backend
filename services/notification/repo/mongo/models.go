package mongorepo

import (
	mongodb "tenkhours/pkg/db/mongo"
)

type DevicesToken struct {
	*mongodb.BaseEntity `                           bson:",inline"`
	ProfileID           string  `json:"profileID,omitempty" bson:"profile_id"`
	Tokens              []Token `json:"tokens,omitempty"    bson:"tokens"`
}

type Token struct {
	DeviceID string `json:"deviceID,omitempty" bson:"device_id"`
	Token    string `json:"token,omitempty"    bson:"token"`
	Platform string `json:"platform,omitempty" bson:"platform"`
	CreateAt string `json:"createAt,omitempty" bson:"create_at"`
}
