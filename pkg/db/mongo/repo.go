package mongodb

import (
	"time"

	"tenkhours/pkg/db/base"

	"go.mongodb.org/mongo-driver/mongo"
)

// Mongo base repo implements IBaseRepo
type BaseRepo[M base.IBaseEntity, N any] struct {
	Collection    *mongo.Collection
	WithTimestamp bool
	Timeout       time.Duration
}

func NewBaseRepo[M base.IBaseEntity, N any](collection *mongo.Collection, withTimestamp bool) *BaseRepo[M, N] {
	return &BaseRepo[M, N]{collection, withTimestamp, 5 * time.Second}
}
