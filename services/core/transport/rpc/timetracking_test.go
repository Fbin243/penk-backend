package rpc_test

import (
	"log"
	"testing"

	"tenkhours/pkg/db/base"
	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/utils"
	"tenkhours/proto/pb/core"
	"tenkhours/services/core/entity"
	"tenkhours/services/core/transport/rpc"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

// Entity
var timeTracking = entity.TimeTracking{
	BaseEntity: &base.BaseEntity{
		ID:        mongodb.GenObjectID(),
		CreatedAt: utils.Now(),
		UpdatedAt: utils.Now(),
	},
	CharacterID:   mongodb.GenObjectID(),
	CategoryID:    lo.ToPtr(mongodb.GenObjectID()),
	ReferenceID:   lo.ToPtr(mongodb.GenObjectID()),
	ReferenceType: lo.ToPtr(entity.EntityTypeHabit),
	Timestamp:     utils.Now(),
}

func TestMapTimeTracking(t *testing.T) {
	// Act
	rpcTimeTracking, err := rpc.Map[entity.TimeTracking, core.TimeTracking](&timeTracking, append(rpc.UnixTimeConverter, rpc.EntityTypeConverter...))
	assert.NoError(t, err)

	log.Printf("timeTracking: %+v", utils.PrettyJSON(timeTracking))
	log.Printf("rpcTimeTracking: %+v", utils.PrettyJSON(rpcTimeTracking))

	// Assert
	assert.Equal(t, timeTracking.ID, rpcTimeTracking.Id)
	assert.Equal(t, timeTracking.CharacterID, rpcTimeTracking.CharacterId)
	assert.Equal(t, *timeTracking.CategoryID, *rpcTimeTracking.CategoryId)
	assert.Equal(t, *timeTracking.ReferenceID, *rpcTimeTracking.ReferenceId)
	assert.Equal(t, string(*timeTracking.ReferenceType), rpcTimeTracking.ReferenceType.String())
	assert.Equal(t, timeTracking.Timestamp.Unix(), rpcTimeTracking.Timestamp)
}

var rpcTimeTrackingInput = &core.TimeTrackingInput{
	ReferenceId:   lo.ToPtr(mongodb.GenObjectID()),
	ReferenceType: lo.ToPtr(core.EntityType_TaskType),
	Timestamp:     utils.Now().Unix(),
	Duration:      3600,
}

func TestMapTimeTrackingInput(t *testing.T) {
	// Act
	timeTrackingInput, err := rpc.Map[core.TimeTrackingInput, entity.TimeTrackingInput](rpcTimeTrackingInput, append(rpc.UnixTimeConverter, rpc.EntityTypeConverter...))
	assert.NoError(t, err)

	log.Printf("rpcTimeTrackingInput: %+v", utils.PrettyJSON(rpcTimeTrackingInput))
	log.Printf("timeTrackingInput: %+v", utils.PrettyJSON(timeTrackingInput))

	// Assert
	assert.Equal(t, *rpcTimeTrackingInput.ReferenceId, timeTrackingInput.ReferenceID)
	assert.Equal(t, rpcTimeTrackingInput.ReferenceType.String(), string(timeTrackingInput.ReferenceType))
	assert.Equal(t, rpcTimeTrackingInput.Timestamp, timeTrackingInput.Timestamp.Unix())
	assert.Equal(t, rpcTimeTrackingInput.Duration, int64(timeTrackingInput.Duration))
}
