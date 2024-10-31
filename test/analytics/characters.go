package analytics

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"tenkhours/test/common"

	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

type CreateSnapshotStage struct {
	common.Metadata
	CharacterKey common.ContextKey
}

func (s CreateSnapshotStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	profile, ok := (*ctx).Value(common.Profile).(string)
	if !ok {
		return common.ErrNotFoundInContext("Profile")
	}

	character, ok := (*ctx).Value(s.CharacterKey).(string)
	if !ok {
		return common.ErrNotFoundInContext("CharacterKey")
	}

	characterId := gjson.Get(character, "id").Value()

	variables := map[string]interface{}{
		"characterID": characterId,
		"description": "Test snapshot description",
	}

	characterMap := gjson.Get(profile, fmt.Sprintf(`characters.#(id==%s)`, characterId)).Value().(map[string]interface{})

	assertion := jsonpath.Chain().NotPresent("$.errors").
		Equal("$.data.createSnapshot.character.name", characterMap["name"]).
		Equal("$.data.createSnapshot.character.tags", characterMap["tags"]).
		Equal("$.data.createSnapshot.character.gender", characterMap["gender"]).
		Equal("$.data.createSnapshot.character.profileID", characterMap["profileID"]).
		Equal("$.data.createSnapshot.character.totalFocusedTime", characterMap["totalFocusedTime"]).
		Equal("$.data.createSnapshot.description", variables["description"]).
		NotEqual("$.data.createSnapshot.timestamp", nil)

	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     CreateSnapshotQuery,
			Variables: variables,
			Assertion: assertion.End(),
		})
}

type GetSnapshotsStage struct {
	common.Metadata
	SnapshotKey     common.ContextKey
	HasOneSnapshot  bool
	HasTwoSnapshots bool
}

func (s GetSnapshotsStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	testingT, ok := (*ctx).Value(common.TestingT).(apitest.TestingT)
	if !ok {
		return common.ErrNotFoundInContext("TestingT")
	}

	variables := map[string]interface{}{}

	assertion := common.AssertionSuccess.End()

	if s.HasOneSnapshot {
		snapshot, ok := (*ctx).Value(s.SnapshotKey).(string)
		if !ok {
			return common.ErrNotFoundInContext("SnapshotKey")
		}

		variables["characterID"] = gjson.Get(snapshot, "character.id").Value()
		assertion = func(res *http.Response, req *http.Request) error {
			json := common.ReadResponseJson(res)

			passed := assert.Empty(testingT, gjson.Get(json, "errors").Value()) &&
				assert.Equal(testingT, gjson.Get(json, "data.snapshots.#").Int(), int64(1)) &&
				assert.Equal(testingT, gjson.Get(json, "data.snapshots.0").Value(), gjson.Parse(snapshot).Value())

			if !passed {
				return fmt.Errorf("assertion failed")
			}

			return nil
		}
	}

	if s.HasTwoSnapshots {
		assertion = func(res *http.Response, req *http.Request) error {
			json := common.ReadResponseJson(res)

			passed := assert.Empty(testingT, gjson.Get(json, "errors").Value()) &&
				assert.Equal(testingT, gjson.Get(json, "data.snapshots.#").Int(), int64(2))

			if !passed {
				return fmt.Errorf("assertion failed")
			}

			return nil
		}
	}

	if s.ExpectError {
		assertion = common.AssertionError.End()
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Query:     SnapshotsQuery,
			Variables: variables,
			Assertion: assertion,
		})
}
