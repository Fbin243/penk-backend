package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IBaseModel interface {
	SetID(id primitive.ObjectID)
	SetCreatedAtByNow()
	SetUpdatedAtByNow()
}

type BaseModel struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updated_at,omitempty"`
}

func (m *BaseModel) SetID(id primitive.ObjectID) {
	m.ID = id
}

func (m *BaseModel) SetCreatedAtByNow() {
	m.CreatedAt = time.Now()
}

func (m *BaseModel) SetUpdatedAtByNow() {
	m.UpdatedAt = time.Now()
}

type IBaseRepo[M IBaseModel] interface {
	InsertOne(m M) (*M, error)
	FindById(id primitive.ObjectID) (*M, error)
	UpdateById(id primitive.ObjectID, m M) (*M, error)
}

type BaseRepo[M IBaseModel] struct {
	*mongo.Collection
}

func NewBaseRepo[M IBaseModel](collection *mongo.Collection) *BaseRepo[M] {
	return &BaseRepo[M]{collection}
}

func (r *BaseRepo[M]) InsertOne(m M) (*M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m.SetCreatedAtByNow()
	m.SetUpdatedAtByNow()
	_, err := r.Collection.InsertOne(ctx, m)
	return &m, err
}

func (r *BaseRepo[M]) FindById(id primitive.ObjectID) (*M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var m M
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&m)
	return &m, err
}

func (r *BaseRepo[M]) UpdateById(id primitive.ObjectID, m M) (*M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m.SetUpdatedAtByNow()
	_, err := r.Collection.ReplaceOne(ctx, bson.M{"_id": id}, m)
	return &m, err
}
