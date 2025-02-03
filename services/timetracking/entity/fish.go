package entity

type Fish struct {
	ProfileID string `json:"profileID,omitempty"`
	Gold      int32  `json:"gold"                bson:"gold"`
	Normal    int32  `json:"normal"              bson:"normal"`
}

type CatchFishResult struct {
	FishType string `json:"fishType"`
	Number   int32  `json:"number"`
}
