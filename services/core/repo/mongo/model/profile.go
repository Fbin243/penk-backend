package mongomodel

import (
	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	*mongodb.BaseEntity `                                    bson:",inline"`
	Name                string             `json:"name,omitempty"               bson:"name"`
	Email               string             `json:"email,omitempty"              bson:"email"`
	FirebaseUID         string             `json:"firebaseUID,omitempty"        bson:"firebase_uid"`
	ImageURL            string             `json:"imageURL,omitempty"           bson:"image_url"`
	CurrentCharacterOID primitive.ObjectID `json:"currentCharacterID,omitempty" bson:"current_character_id"`
}

func (p *Profile) CurrentCharacterID(id string) {
	p.CurrentCharacterOID = mongodb.ToObjectID(id)
}
