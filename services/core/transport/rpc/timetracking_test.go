package rpc_test

import (
	"log"
	"testing"
	"time"

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
	StartTime:     utils.Now(),
	EndTime:       utils.Now().Add(1 * time.Hour),
}

func TestMapTimeTracking(t *testing.T) {
	// Act
	rpcTimeTracking, err := rpc.MapEntityToRPC[entity.TimeTracking, core.TimeTracking](&timeTracking, append(rpc.UnixTimeConverter, rpc.EntityTypeConverter...))
	assert.NoError(t, err)

	log.Printf("timeTracking: %+v", utils.PrettyJSON(timeTracking))
	log.Printf("rpcTimeTracking: %+v", utils.PrettyJSON(rpcTimeTracking))

	// Assert
	assert.Equal(t, timeTracking.ID, rpcTimeTracking.Id)
	assert.Equal(t, timeTracking.CharacterID, rpcTimeTracking.CharacterId)
	assert.Equal(t, *timeTracking.CategoryID, *rpcTimeTracking.CategoryId)
	assert.Equal(t, *timeTracking.ReferenceID, *rpcTimeTracking.ReferenceId)
	assert.Equal(t, string(*timeTracking.ReferenceType), rpcTimeTracking.ReferenceType.String())
	assert.Equal(t, timeTracking.StartTime.Unix(), rpcTimeTracking.StartTime)
	assert.Equal(t, timeTracking.EndTime.Unix(), *rpcTimeTracking.EndTime)
}

var rpcTimeTrackingInput = &core.CreateTimeTrackingReq{
	CharacterId:   mongodb.GenObjectID(),
	CategoryId:    lo.ToPtr(mongodb.GenObjectID()),
	ReferenceId:   lo.ToPtr(mongodb.GenObjectID()),
	ReferenceType: lo.ToPtr(core.EntityType_Task),
	StartTime:     utils.Now().Unix(),
}

func TestMapTimeTrackingInput(t *testing.T) {
	// Act
	timeTrackingInput, err := rpc.MapRPCInputToEntityInput[core.CreateTimeTrackingReq, entity.TimeTrackingInput](rpcTimeTrackingInput, append(rpc.UnixTimeConverter, rpc.EntityTypeConverter...))
	assert.NoError(t, err)

	log.Printf("rpcTimeTrackingInput: %+v", utils.PrettyJSON(rpcTimeTrackingInput))
	log.Printf("timeTrackingInput: %+v", utils.PrettyJSON(timeTrackingInput))

	// Assert
	assert.Equal(t, *rpcTimeTrackingInput.CategoryId, *timeTrackingInput.CategoryID)
	assert.Equal(t, *rpcTimeTrackingInput.ReferenceId, *timeTrackingInput.ReferenceID)
	assert.Equal(t, rpcTimeTrackingInput.ReferenceType.String(), string(*timeTrackingInput.ReferenceType))
	assert.Equal(t, rpcTimeTrackingInput.StartTime, timeTrackingInput.StartTime.Unix())
}
