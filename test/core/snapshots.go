package core

import (
	"context"
	"log"

	"tenkhours/test/common"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/tidwall/gjson"
)

type CreateSnapshotStage struct {
	common.Metadata
	CharacterKey common.ContextKey
}

func (s CreateSnapshotStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	character, ok := (*ctx).Value(s.CharacterKey).(string)
	if !ok {
		return common.ErrNotFoundInContext("CharacterKey")
	}

	characterId := gjson.Get(character, "id").Value()

	variables := map[string]interface{}{
		"characterID": characterId,
		"description": "Test snapshot description",
	}

	assertion := jsonpath.Chain().NotPresent("$.errors").
		Equal("$.data.createSnapshot.character.name", gjson.Get(character, "name").Value()).
		Equal("$.data.createSnapshot.character.tags", gjson.Get(character, "tags").Value()).
		Equal("$.data.createSnapshot.character.gender", gjson.Get(character, "gender").Value()).
		Equal("$.data.createSnapshot.character.profileID", gjson.Get(character, "profileID").Value()).
		NotEqual("$.data.createSnapshot.timestamp", nil)

	if s.ExpectError {
		assertion = jsonpath.Chain().Present("$.errors")
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     CreateSnapshotQuery,
			Variables: variables,
			Assertion: []common.Assertion{assertion.End()},
		})
}

type GetSnapshotsStage struct {
	common.Metadata
	SnapshotKey       common.ContextKey
	NumberOfSnapshots int
}

func (s GetSnapshotsStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	snapshot, ok := (*ctx).Value(s.SnapshotKey).(string)
	if !ok {
		return common.ErrNotFoundInContext("SnapshotKey")
	}

	variables := map[string]interface{}{}

	assertion := jsonpath.Chain().NotPresent("$.errors").
		Equal("$.data.snapshots[0].description", gjson.Get(snapshot, "description").Value()).
		Equal("$.data.snapshots[0].character.name", gjson.Get(snapshot, "character.name").Value()).
		Equal("$.data.snapshots[0].character.tags", gjson.Get(snapshot, "character.tags").Value()).
		NotEqual("$.data.snapshots[0].timestamp", nil).
		NotEqual("$.data.snapshots[0].id", nil)

	assertions := []common.Assertion{}
	assertions = append(assertions, jsonpath.Len("$.data.snapshots", s.NumberOfSnapshots))

	if s.ExpectError {
		assertion = jsonpath.Chain().Present("$.errors")
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     SnapshotsQuery,
			Variables: variables,
			Assertion: append(assertions, assertion.End()),
		})
}
