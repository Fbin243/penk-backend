package repo_test

import (
	"testing"

	"tenkhours/pkg/utils"
	"tenkhours/services/core/repo"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var snapshot = &repo.Snapshot{
	ID:        primitive.NewObjectID(),
	Timestamp: utils.Now(),
	Metadata: repo.Metadata{
		ProfileID:   primitive.NewObjectID(),
		CharacterID: primitive.NewObjectID(),
	},
	Character: repo.SnapshotCharacter{},
	// Asset: :,
}

func TestCreateNewSnapshot(t *testing.T) {
	_, err := snapshotsRepo.CreateSnapshot(snapshot)

	assert.Nil(t, err)
}

func TestGetSnapshotsByProfileID(t *testing.T) {
	profileID := primitive.NewObjectID()
	snapshot.Metadata.ProfileID = profileID

	_, err := snapshotsRepo.CreateSnapshot(snapshot)
	assert.Nil(t, err)

	// queriedSnapshots, err := repo.GetSnapshotsByProfileID(profileID)
	// assert.Nil(t, err)
	// assert.Len(t, queriedSnapshots, 1)
	// assert.Equal(t, *snapshot, queriedSnapshots[0])
}

func TestGetSnapshotsByCharacterID(t *testing.T) {
	characterID := primitive.NewObjectID()
	snapshot.Metadata.CharacterID = characterID

	_, err := snapshotsRepo.CreateSnapshot(snapshot)
	assert.Nil(t, err)

	// queriedSnapshots, err := snapshotsRepo.GetSnapshotsByCharacterID(characterID)
	// assert.Nil(t, err)
	// assert.Len(t, queriedSnapshots, 1)
	// assert.Equal(t, *snapshot, queriedSnapshots[0])
}

func TestGetLatestSnapshotByCharacterID(t *testing.T) {
	// Insert with the past time
	_, err := snapshotsRepo.CreateSnapshot(snapshot)
	assert.Nil(t, err)

	// Insert with the latest time
	latestTime := utils.Now()
	snapshot.Timestamp = latestTime
	_, err = snapshotsRepo.CreateSnapshot(snapshot)
	assert.Nil(t, err)

	latestSnapshot, err := snapshotsRepo.GetLatestSnapshotByCharacterID(snapshot.Metadata.CharacterID)

	assert.Nil(t, err)
	assert.Equal(t, latestSnapshot.Timestamp, latestTime)
}
