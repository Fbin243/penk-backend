package analytics

var UserSnapshotsQuery = `
query UserSnapshots {
    userSnapshots {
        id
        timestamp
        character {
            gender
            id
            limitedMetricNumber
            name
            tags
            totalFocusedTime
            profileID
            customMetrics {
                description
                id
                limitedPropertyNumber
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

var CharacterSnapshotsQuery = `
query CharacterSnapshots ($characterID: ID!) {
    characterSnapshots(characterID: $characterID) {
        id
        timestamp
        character {
            gender
            id
            limitedMetricNumber
            name
            tags
            totalFocusedTime
            profileID
            customMetrics {
                description
                id
                limitedPropertyNumber
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

var CreateSnapshotQuery = `
mutation CreateSnapshot ($characterID: ID!) {
    createSnapshot(characterID: $characterID) {
        id
        timestamp
        character {
            gender
            id
            limitedMetricNumber
            name
            tags
            totalFocusedTime
            profileID
            customMetrics {
                description
                id
                limitedPropertyNumber
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
