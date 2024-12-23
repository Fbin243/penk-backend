package core

var SnapshotsQuery = `
query Snapshots ($characterID: ObjectID) {
    snapshots(characterID: $characterID) {
        id
        timestamp
        description
        character {
            id
            profileID
            name
            gender
            tags
            totalFocusedTime
            customMetrics {
                id
                name
                description
                time
                style {
                    color
                    icon
                }
                properties {
                    id
                    name
                    type
                    value
                    unit
                }
            }
        }
    }
}
`

var CreateSnapshotQuery = `
mutation CreateSnapshot ($characterID: ObjectID!, $description: String) {
    createSnapshot(characterID: $characterID, description: $description) {
        id
        timestamp
        description
        character {
            gender
            id
            name
            tags
            totalFocusedTime
            profileID
            customMetrics {
                description
                id
                name
                time
                properties {
                    id
                    name
                    type
                    unit
                    value
                }
                style {
                    color
                    icon
                }
            }
        }
    }
}
`
