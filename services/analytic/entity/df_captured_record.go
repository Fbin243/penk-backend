package entity

type DFCapturedRecord struct {
	ID               string `dataframe:"id"`
	ProfileID        string `dataframe:"profile_id"`
	CharacterID      string `dataframe:"character_id"`
	Year             int    `dataframe:"year"`
	Month            int    `dataframe:"month"`
	Week             int    `dataframe:"week"`
	Day              int    `dataframe:"day"`
	TotalFocusedTime int    `dataframe:"total_focused_time"`
}

type DFCapturedRecordCategory struct {
	ID         string `dataframe:"id"`
	CategoryID string `dataframe:"category_id"`
	Time       int    `dataframe:"time"`
}
