package timetracking

var CreateTimeTrackingQuery = `
	mutation CreateTimeTracking ($characterID: ID!, $customMetricID: ID, $startTime: Time!) {
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
	mutation UpdateTimeTracking {
		updateTimeTracking {
			timeTracking {
				id
				characterID
				customMetricID
				startTime
				endTime
			}	
			gold
			normal
		}
	}
`
