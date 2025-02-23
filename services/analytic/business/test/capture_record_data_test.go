package business_test

import (
	"tenkhours/pkg/utils"
	"tenkhours/services/analytic/entity"
)

var capturedRecords = []entity.CapturedRecord{
	{
		Timestamp: utils.ParseTime("2025-01-10T14:30:00.000Z"),
		Metadata: entity.CapturedRecordMetadata{
			CharacterID: "character_id_1",
			ProfileID:   "profile_id",
		},
		Categories: []entity.CapturedRecordCategory{
			{ID: "category_id_1", Time: 120},
			{ID: "category_id_2", Time: 180},
		},
		TimeTrackings: []entity.CapturedRecordTimeTracking{
			{
				Time:      120,
				StartTime: utils.ParseTime("2025-01-10T14:30:00.000Z"),
				EndTime:   utils.ParseTime("2025-01-10T14:32:00.000Z"),
			},
			{
				Time:      180,
				StartTime: utils.ParseTime("2025-02-15T09:45:00.000Z"),
				EndTime:   utils.ParseTime("2025-02-15T09:50:00.000Z"),
			},
		},
		TotalFocusedTime: 300,
		ID:               "captured_record_id_1",
	},
	{
		Timestamp: utils.ParseTime("2024-03-12T08:15:00.000Z"),
		Metadata: entity.CapturedRecordMetadata{
			CharacterID: "character_id_2",
			ProfileID:   "profile_id",
		},
		Categories: []entity.CapturedRecordCategory{
			{ID: "category_id_1", Time: 200},
			{ID: "category_id_2", Time: 150},
		},
		TimeTrackings: []entity.CapturedRecordTimeTracking{
			{
				Time:      200,
				StartTime: utils.ParseTime("2024-03-12T08:15:00.000Z"),
				EndTime:   utils.ParseTime("2024-03-12T08:18:20.000Z"),
			},
			{
				Time:      150,
				StartTime: utils.ParseTime("2024-04-05T11:20:00.000Z"),
				EndTime:   utils.ParseTime("2024-04-05T11:25:00.000Z"),
			},
		},
		TotalFocusedTime: 350,
		ID:               "captured_record_id_2",
	},
}
