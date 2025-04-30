package business_test

import "tenkhours/services/analytic/entity"

var capturedRecords = []entity.CapturedRecord{
	{
		CharacterID: "character_id_1",
		CategoryID:  "category_id_1",
		Year:        2025,
		Date:        "2025-01-10",
		Month:       1,
		Week:        2,
		Day:         10,
		Time:        120,
	},
	{
		CharacterID: "character_id_1",
		CategoryID:  "category_id_2",
		Year:        2025,
		Date:        "2025-01-10",
		Month:       1,
		Week:        2,
		Day:         10,
		Time:        180,
	},
}
