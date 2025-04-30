package core

var GoalsQuery = `
query Goals($characterID: ID!) {
    goals(characterID: $characterID) {
        id
        createdAt
        updatedAt
        characterID
        name
        description
        startTime
        endTime
        status
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
                value
                unit
                condition
                targetValue
                rangeValue {
                    min
                    max
                }
            }
        }
        metrics {
            id
            name
            value
            unit
            condition
            targetValue
            rangeValue {
                min
                max
            }
        }
        checkboxes {
            id
            name
            value
        }
    }
}`

var UpsertGoalQuery = `
mutation UpsertGoal($input: GoalInput!) {
    upsertGoal(input: $input) {
        id
        createdAt
        updatedAt
        characterID
        name
        description
        startTime
        endTime
        status
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
                value
                unit
                condition
                targetValue
                rangeValue {
                    min
                    max
                }
            }
        }
        metrics {
            id
            name
            value
            unit
            condition
            targetValue
            rangeValue {
                min
                max
            }
        }
        checkboxes {
            id
            name
            value
        }
    }
}`

var DeleteGoalQuery = `
mutation DeleteGoal($id: ID!) {
    deleteGoal(id: $id) {
        id
        createdAt
        updatedAt
        characterID
        name
        description
        startTime
        endTime
        status
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
                value
                unit
                condition
                targetValue
                rangeValue {
                    min
                    max
                }
            }
        }
        metrics {
            id
            name
            value
            unit
            condition
            targetValue
            rangeValue {
                min
                max
            }
        }
        checkboxes {
            id
            name
            value
        }
    }
}`
