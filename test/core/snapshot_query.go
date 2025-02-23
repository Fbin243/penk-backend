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
            categories {
                id
                name
                description
                style {
                    color
                    icon
                }
                metrics {
                    id
                    name
                    type
                    value
                    unit
                }
            }
        }
    }
}`

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
            profileID
            categories {
                description
                id
                name
                metrics {
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
}`
