package mongorepo

import (
	"time"

	mongodb "tenkhours/pkg/db/mongo"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CapturedRecord struct {
	OID              primitive.ObjectID           `json:"id"                      bson:"_id,omitempty"`
	Timestamp        time.Time                    `json:"timestamp"               bson:"timestamp,omitempty"`
	TotalFocusedTime int32                        `json:"totalFocusedTime"        bson:"total_focused_time,omitempty"`
	Categories       []CapturedRecordCategory     `json:"categories,omitempty"    bson:"categories,omitempty"`
	TimeTrackings    []CapturedRecordTimeTracking `json:"timeTrackings,omitempty" bson:"time_trackings,omitempty"`
	Metadata         CapturedRecordMetadata       `json:"metadata"                bson:"metadata,omitempty"`
}

func (c *CapturedRecord) ID(id string) {
	c.OID = mongodb.ToObjectID(id)
}

type CapturedRecordCategory struct {
	OID  primitive.ObjectID `json:"id"   bson:"_id,omitempty"`
	Time int32              `json:"time" bson:"time,omitempty"`
}

func (c *CapturedRecordCategory) ID(id string) {
	c.OID = mongodb.ToObjectID(id)
}

type CapturedRecordMetadata struct {
	CharacterOID primitive.ObjectID `json:"characterID" bson:"character_id,omitempty"`
	ProfileOID   primitive.ObjectID `json:"profileID"   bson:"profile_id,omitempty"`
}

func (c *CapturedRecordMetadata) ProfileID(id string) {
	c.ProfileOID = mongodb.ToObjectID(id)
}

func (c *CapturedRecordMetadata) CharacterID(id string) {
	c.CharacterOID = mongodb.ToObjectID(id)
}

type CapturedRecordTimeTracking struct {
	CategoryOID *primitive.ObjectID `json:"categoryID,omitempty" bson:"category_id,omitempty"`
	Time        int32               `json:"time"                 bson:"time,omitempty"`
	StartTime   time.Time           `json:"startTime"            bson:"start_time,omitempty"`
	EndTime     time.Time           `json:"endTime"              bson:"end_time,omitempty"`
}

func (c *CapturedRecordTimeTracking) CategoryID(id *string) {
	if id != nil {
		c.CategoryOID = lo.ToPtr(mongodb.ToObjectID(*id))
	}
}
