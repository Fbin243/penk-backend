package timetrackings

var CreateTimeTrackingQuery = `
	mutation CreateTimeTracking ($characterID: ID!, $metricID: ID, $startTime: DateTime!) {
		createTimeTracking(characterID: $characterID, startTime: $startTime, metricID: $metricID) {
			characterID
			metricID
			endTime
			id
			startTime
		}
	}
`

var UpdateTimeTrackingQuery = `
	mutation UpdateTimeTracking ($id: ID!) {
		updateTimeTracking(id: $id) {
			characterID
			metricID
			endTime
			id
			startTime
		}
	}
`
