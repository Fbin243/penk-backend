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
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updated_at,omitempty"`
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
	InsertOne(m *M) (*M, error)
	FindByID(id primitive.ObjectID) (*M, error)
	UpdateByID(id primitive.ObjectID, m *M) (*M, error)
	DeleteByID(id primitive.ObjectID) (*M, error)
}

type BaseRepo[M IBaseModel] struct {
	*mongo.Collection
}

func NewBaseRepo[M IBaseModel](collection *mongo.Collection) *BaseRepo[M] {
	return &BaseRepo[M]{collection}
}

func (r *BaseRepo[M]) InsertOne(m *M) (*M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_m := *m
	_m.SetID(primitive.NewObjectID())
	_m.SetCreatedAtByNow()
	_m.SetUpdatedAtByNow()
	_, err := r.Collection.InsertOne(ctx, m)
	return m, err
}

func (r *BaseRepo[M]) FindByID(id primitive.ObjectID) (*M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var m M
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&m)
	return &m, err
}

func (r *BaseRepo[M]) UpdateByID(id primitive.ObjectID, m *M) (*M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_m := *m
	_m.SetUpdatedAtByNow()

	err := r.Collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": _m}, FindOneAndUpdateOptions).Decode(&m)
	return m, err
}

func (r *BaseRepo[M]) DeleteByID(id primitive.ObjectID) (*M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var m M
	err := r.Collection.FindOneAndDelete(ctx, bson.M{"_id": id}).Decode(&m)
	return &m, err
}
