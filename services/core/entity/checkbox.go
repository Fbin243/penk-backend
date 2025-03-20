package entity

type Checkbox struct {
	ID    string `json:"id"    bson:"id"`
	Name  string `json:"name"  bson:"name"`
	Value bool   `json:"value" bson:"value"`
}
