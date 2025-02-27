package timetracking

var GetCurrentTimeTrackingQuery = `
query CurrentTimeTracking {
    currentTimeTracking {
        id
        characterID
        categoryID
        startTime
        endTime
    }
}`

var CreateTimeTrackingQuery = `
	mutation CreateTimeTracking ($characterID: ID!, $categoryID: ID, $startTime: Time!) {
		createTimeTracking(characterID: $characterID, startTime: $startTime, categoryID: $categoryID) {
			characterID
			categoryID
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
				categoryID
				startTime
				endTime
			}	
			gold
			normal
		}
	}
`
