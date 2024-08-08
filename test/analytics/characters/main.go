package characters

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
	user, ok := (*ctx).Value(common.User).(string)
	if !ok {
		return common.ErrNotFound(common.User)
	}

	character, ok := (*ctx).Value(s.CharacterKey).(string)
	if !ok {
		return common.ErrNotFound(s.CharacterKey)
	}

	characterId := gjson.Get(character, "id").Value()

	variables := map[string]interface{}{
		"characterID": characterId,
	}

	assertion := jsonpath.Chain().NotPresent("$.errors").
		Equal("$.data.createSnapshot.character", gjson.Get(user, fmt.Sprintf(`characters.#(id==%s)`, characterId)).Value())

	if s.ExpectError {
		assertion = common.AssertionError
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Url:       common.AnalyticsUrl,
			Query:     CreateSnapshotQuery,
			Variables: variables,
			Assertion: assertion.End(),
		})
}

type GetCharacterSnapshotsStage struct {
	common.Metadata
	SnapshotKey    common.ContextKey
	HasOneSnapshot bool
}

func (s GetCharacterSnapshotsStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	testingT, ok := (*ctx).Value(common.TestingT).(apitest.TestingT)
	if !ok {
		return common.ErrNotFound(common.TestingT)
	}

	snapshot, ok := (*ctx).Value(s.SnapshotKey).(string)
	if !ok {
		return common.ErrNotFound(s.SnapshotKey)
	}

	variables := map[string]interface{}{
		"characterID": gjson.Get(snapshot, "character.id").Value(),
	}

	assertion := common.AssertionSuccess.End()

	if s.HasOneSnapshot {
		assertion = func(res *http.Response, req *http.Request) error {
			json := common.ReadResponseJson(res)

			assert.Empty(testingT, gjson.Get(json, "errors").Value())
			assert.Equal(testingT, gjson.Get(json, "data.characterSnapshots.#").Int(), int64(1))
			assert.Equal(testingT, gjson.Get(json, "data.characterSnapshots.0").Value(), gjson.Parse(snapshot).Value())

			return nil
		}
	}

	if s.ExpectError {
		assertion = common.AssertionError.End()
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Url:       common.AnalyticsUrl,
			Query:     CharacterSnapshotsQuery,
			Variables: variables,
			Assertion: assertion,
		})
}

type GetUserSnapshotsStage struct {
	common.Metadata
	HasTwoSnapshots bool
}

func (s GetUserSnapshotsStage) Exec(ctx *context.Context) error {
	log.Println("--> Stage: ", s.Describe)
	testingT, ok := (*ctx).Value(common.TestingT).(apitest.TestingT)
	if !ok {
		return common.ErrNotFound(common.TestingT)
	}

	user, ok := (*ctx).Value(common.User).(string)
	if !ok {
		return common.ErrNotFound(common.User)
	}

	variables := map[string]interface{}{
		"characterID": gjson.Get(user, "id").Value(),
	}

	assertion := common.AssertionSuccess.End()

	if s.HasTwoSnapshots {
		assertion = func(res *http.Response, req *http.Request) error {
			json := common.ReadResponseJson(res)

			assert.Empty(testingT, gjson.Get(json, "errors").Value())
			assert.Equal(testingT, gjson.Get(json, "data.userSnapshots.#").Int(), int64(2))

			return nil
		}
	}

	if s.ExpectError {
		assertion = common.AssertionError.End()
	}

	return common.QueryGraphQL(ctx,
		&common.QueryParams{
			Url:       common.AnalyticsUrl,
			Query:     UserSnapshotsQuery,
			Variables: variables,
			Assertion: assertion,
		})
}
