package entity

import "tenkhours/pkg/db/base"

type Profile struct {
	*base.BaseEntity   `                                    bson:",inline"`
	Name               string `json:"name,omitempty"               bson:"name"`
	Email              string `json:"email,omitempty"              bson:"email"`
	FirebaseUID        string `json:"firebaseUID,omitempty"        bson:"firebase_uid"`
	ImageURL           string `json:"imageURL,omitempty"           bson:"image_url"`
	CurrentCharacterID string `json:"currentCharacterID,omitempty" bson:"current_character_id"`
}
