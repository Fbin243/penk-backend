package entity

import "time"

type StatAnalyticFilter struct {
	CharacterID      string            `json:"characterID"`
	StartTime        *time.Time        `json:"startTime,omitempty"`
	EndTime          *time.Time        `json:"endTime,omitempty"`
	AnalyticSections []AnalyticSection `json:"analyticSections"`
}

type CapturedRecord struct {
	CharacterID string `bson:"character_id" dataframe:"character_id"`
	CategoryID  string `bson:"category_id"  dataframe:"category_id"`
	Year        int    `bson:"year"         dataframe:"year"`
	Date        string `bson:"date"         dataframe:"date"`
	Month       int    `bson:"month"        dataframe:"month"`
	Week        int    `bson:"week"         dataframe:"week"`
	Day         int    `bson:"day"          dataframe:"day"`
	Time        int    `bson:"time"         dataframe:"time"`
}
