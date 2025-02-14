package mongorepo

import (
	"time"

	"tenkhours/services/core/entity"

	mongodb "tenkhours/pkg/db/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	*mongodb.BaseEntity `                                    bson:",inline"`
	Name                string             `json:"name,omitempty"               bson:"name"`
	Email               string             `json:"email,omitempty"              bson:"email"`
	FirebaseUID         string             `json:"firebaseUID,omitempty"        bson:"firebase_uid"`
	ImageURL            string             `json:"imageURL,omitempty"           bson:"image_url"`
	CurrentCharacterOID primitive.ObjectID `json:"currentCharacterID,omitempty" bson:"current_character_id,omitempty"`
	AvailableSnapshots  int32              `json:"availableSnapshots,omitempty" bson:"available_snapshots"`
	AutoSnapshot        bool               `json:"autoSnapshot,omitempty"       bson:"auto_snapshot"`
}

func (p *Profile) CurrentCharacterID(id string) {
	p.CurrentCharacterOID = mongodb.ToObjectID(id)
}

type Character struct {
	*mongodb.BaseEntity `                                  bson:",inline"`
	ProfileOID          primitive.ObjectID `json:"profileID,omitempty"        bson:"profile_id"`
	Name                string             `json:"name,omitempty"             bson:"name"`
	Gender              bool               `json:"gender,omitempty"           bson:"gender"`
	Tags                []string           `json:"tags,omitempty"             bson:"tags"`
	TotalFocusedTime    int32              `json:"totalFocusedTime,omitempty" bson:"total_focused_time"`
	Categories          []Category         `json:"categories,omitempty"       bson:"categories"`
	Metrics             []Metric           `json:"metrics,omitempty"          bson:"metrics"`
	Vision              Vision             `json:"vision,omitempty"           bson:"vision"`
}

func (p *Character) ProfileID(id string) {
	p.ProfileOID = mongodb.ToObjectID(id)
}

type Vision struct {
	Name        string `json:"name,omitempty"        bson:"name"`
	Description string `json:"description,omitempty" bson:"description"`
}

type Category struct {
	OID         primitive.ObjectID `json:"id,omitempty"          bson:"_id"`
	Name        string             `json:"name,omitempty"        bson:"name"`
	Description string             `json:"description,omitempty" bson:"description"`
	Time        int32              `json:"time,omitempty"        bson:"time"`
	Style       CategoryStyle      `json:"style,omitempty"       bson:"style"`
}

func (m *Category) ID(id string) {
	m.OID = mongodb.ToObjectID(id)
}

type CategoryStyle struct {
	Color string `json:"color,omitempty" bson:"color"`
	Icon  string `json:"icon,omitempty"  bson:"icon"`
}

type Metric struct {
	OID         primitive.ObjectID `json:"id,omitempty"          bson:"_id"`
	CategoryOID primitive.ObjectID `json:"category_id,omitempty" bson:"category_id,omitempty"`
	Name        string             `json:"name,omitempty"        bson:"name"`
	Value       float64            `json:"value,omitempty"       bson:"value"`
	Unit        string             `json:"unit,omitempty"        bson:"unit"`
}

func (m *Metric) ID(id string) {
	m.OID = mongodb.ToObjectID(id)
}

func (m *Metric) CategoryID(id *string) {
	if id != nil {
		m.CategoryOID = mongodb.ToObjectID(*id)
	}
}

type Goal struct {
	*mongodb.BaseEntity `                   bson:",inline"`
	CharacterOID        primitive.ObjectID      `json:"characterID" bson:"character_id"`
	Name                string                  `json:"name"        bson:"name"`
	Description         string                  `json:"description" bson:"description"`
	StartDate           time.Time               `json:"startDate"   bson:"start_date"`
	EndDate             time.Time               `json:"endDate"     bson:"end_date"`
	Status              entity.GoalFinishStatus `json:"status"      bson:"status"`
	Target              []Category              `json:"target"      bson:"target"`
}

func (g *Goal) CharacterID(id string) {
	g.CharacterOID = mongodb.ToObjectID(id)
}

type Template struct {
	*mongodb.BaseEntity `                             bson:",inline"`
	Name                string             `json:"name,omitempty"        bson:"name"`
	Description         string             `json:"description,omitempty" bson:"description"`
	TopicOID            primitive.ObjectID `json:"topicID,omitempty"     bson:"topic_id"`
	Style               TemplateStyle      `json:"style,omitempty"       bson:"style"`
	Categories          []TemplateCategory `json:"categories,omitempty"  bson:"categories"`
}

func (t *Template) TopicID(id string) {
	t.TopicOID = mongodb.ToObjectID(id)
}

type TemplateCategory struct {
	Name        string           `json:"name,omitempty"        bson:"name"`
	Description string           `json:"description,omitempty" bson:"description"`
	Style       CategoryStyle    `json:"style,omitempty"       bson:"style"`
	Metrics     []TemplateMetric `json:"metrics,omitempty"     bson:"metrics"`
}

type TemplateMetric struct {
	Name  string  `json:"name,omitempty"  bson:"name"`
	Value float64 `json:"value,omitempty" bson:"value"`
	Unit  string  `json:"unit,omitempty"  bson:"unit"`
}

type TemplateStyle struct {
	Color string `json:"color,omitempty" bson:"color"`
	Icon  string `json:"icon,omitempty"  bson:"icon"`
}

type TemplateTopic struct {
	*mongodb.BaseEntity `                             bson:",inline"`
	Name                string `json:"name,omitempty"        bson:"name"`
	Description         string `json:"description,omitempty" bson:"description"`
}
