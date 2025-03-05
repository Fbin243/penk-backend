package mongorepo

import (
	"time"

	"tenkhours/services/core/entity"

	mongodb "tenkhours/pkg/db/mongo"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	*mongodb.BaseEntity `                                    bson:",inline"`
	Name                string              `json:"name,omitempty"               bson:"name"`
	Email               string              `json:"email,omitempty"              bson:"email"`
	FirebaseUID         string              `json:"firebaseUID,omitempty"        bson:"firebase_uid"`
	ImageURL            string              `json:"imageURL,omitempty"           bson:"image_url"`
	CurrentCharacterOID *primitive.ObjectID `json:"currentCharacterID,omitempty" bson:"current_character_id,omitempty"`
}

func (p *Profile) CurrentCharacterID(id *string) {
	if id != nil {
		p.CurrentCharacterOID = lo.ToPtr(mongodb.ToObjectID(*id))
	}
}

type Character struct {
	*mongodb.BaseEntity `                            bson:",inline"`
	ProfileOID          primitive.ObjectID `json:"profileID,omitempty"  bson:"profile_id"`
	Name                string             `json:"name,omitempty"       bson:"name"`
	Gender              bool               `json:"gender,omitempty"     bson:"gender"`
	Tags                []string           `json:"tags,omitempty"       bson:"tags"`
	Categories          []Category         `json:"categories,omitempty" bson:"categories"`
	Metrics             []Metric           `json:"metrics,omitempty"    bson:"metrics"`
}

func (p *Character) ProfileID(id string) {
	p.ProfileOID = mongodb.ToObjectID(id)
}

type Category struct {
	OID         primitive.ObjectID `json:"id,omitempty"          bson:"_id"`
	Name        string             `json:"name,omitempty"        bson:"name"`
	Description string             `json:"description,omitempty" bson:"description"`
	Style       CategoryStyle      `json:"style,omitempty"       bson:"style"`
	Metrics     []Metric           `json:"metrics,omitempty"     bson:"metrics"`
}

func (m *Category) ID(id string) {
	m.OID = mongodb.ToObjectID(id)
}

type CategoryStyle struct {
	Color string `json:"color,omitempty" bson:"color"`
	Icon  string `json:"icon,omitempty"  bson:"icon"`
}

type Metric struct {
	OID   primitive.ObjectID `json:"id,omitempty"    bson:"_id"`
	Name  string             `json:"name,omitempty"  bson:"name"`
	Value float64            `json:"value,omitempty" bson:"value"`
	Unit  string             `json:"unit,omitempty"  bson:"unit"`
}

func (m *Metric) ID(id string) {
	m.OID = mongodb.ToObjectID(id)
}

type Goal struct {
	*mongodb.BaseEntity `                   bson:",inline"`
	CharacterOID        primitive.ObjectID      `json:"characterID" bson:"character_id"`
	Name                string                  `json:"name"        bson:"name"`
	Description         string                  `json:"description" bson:"description"`
	StartTime           time.Time               `json:"startTime"   bson:"start_time"`
	EndTime             time.Time               `json:"endTime"     bson:"end_time"`
	Status              entity.GoalFinishStatus `json:"status"      bson:"status"`
	Target              GoalTarget              `json:"target"      bson:"target"`
}

func (g *Goal) CharacterID(id string) {
	g.CharacterOID = mongodb.ToObjectID(id)
}

type GoalTarget struct {
	Categories []GoalCategory `json:"categories" bson:"categories"`
	Metrics    []GoalMetric   `json:"metrics"    bson:"metrics"`
	Checkboxes []Checkbox     `json:"checkboxes" bson:"checkboxes"`
}

type GoalCategory struct {
	*Category `               bson:",inline"`
	Metrics   []GoalMetric `json:"metrics" bson:"metrics"`
}

type GoalMetric struct {
	*Metric     `                   bson:",inline"`
	Condition   entity.MetricCondition `json:"condition"   bson:"condition"`
	TargetValue *float64               `json:"targetValue" bson:"target_value,omitempty"`
	RangeValue  *entity.Range          `json:"rangeValue"  bson:"range_value,omitempty"`
}

type Checkbox struct {
	OID   primitive.ObjectID `json:"id"    bson:"id"`
	Name  string             `json:"name"  bson:"name"`
	Value bool               `json:"value" bson:"value"`
}

func (c *Checkbox) ID(id string) {
	c.OID = mongodb.ToObjectID(id)
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
