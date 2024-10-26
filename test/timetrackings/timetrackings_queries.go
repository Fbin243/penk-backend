package timetrackings

var CreateTimeTrackingQuery = `
	mutation CreateTimeTracking ($characterID: ObjectID!, $customMetricID: ObjectID, $startTime: Time!) {
		createTimeTracking(characterID: $characterID, startTime: $startTime, customMetricID: $customMetricID) {
			characterID
			customMetricID
			endTime
			id
			startTime
		}
	}
`

var UpdateTimeTrackingQuery = `
	mutation UpdateTimeTracking ($id: ObjectID!) {
		updateTimeTracking(id: $id) {
			characterID
			customMetricID
			endTime
			id
			startTime
		}
	}
`
