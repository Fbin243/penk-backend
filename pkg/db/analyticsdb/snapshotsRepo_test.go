package analyticsdb

import (
	"testing"

	"tenkhours/pkg/db/coredb"
	"tenkhours/pkg/utils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var snapshot = &Snapshot{
	ID:        primitive.NewObjectID(),
	Timestamp: utils.Now(),
	Metadata: Metadata{
		UserID:      primitive.NewObjectID(),
		CharacterID: primitive.NewObjectID(),
	},
	Character: coredb.Character{},
	// Asset: :,
}

func TestCreateNewSnapshot(t *testing.T) {
	_, err := snapshotsRepo.CreateSnapshot(snapshot)

	assert.Nil(t, err)
}

func TestGetSnapshotsByUserID(t *testing.T) {
	userID := primitive.NewObjectID()
	snapshot.Metadata.UserID = userID

	_, err := snapshotsRepo.CreateSnapshot(snapshot)
	assert.Nil(t, err)

	queriedSnapshots, err := snapshotsRepo.GetSnapshotsByUserID(userID)
	assert.Nil(t, err)
	assert.Len(t, queriedSnapshots, 1)
	assert.Equal(t, *snapshot, queriedSnapshots[0])
}

func TestGetSnapshotsByCharacterID(t *testing.T) {
	characterID := primitive.NewObjectID()
	snapshot.Metadata.CharacterID = characterID

	_, err := snapshotsRepo.CreateSnapshot(snapshot)
	assert.Nil(t, err)

	queriedSnapshots, err := snapshotsRepo.GetSnapshotsByCharacterID(characterID)
	assert.Nil(t, err)
	assert.Len(t, queriedSnapshots, 1)
	assert.Equal(t, *snapshot, queriedSnapshots[0])
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
