package coredb

import (
	"context"
	"fmt"
	"time"

	"tenkhours/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CharactersRepo struct {
	*mongo.Collection
}

func NewCharactersRepo(mongodb *mongo.Database) *CharactersRepo {
	return &CharactersRepo{mongodb.Collection(db.CharacterCollection)}
}

func (r *CharactersRepo) GetCharacterByID(id primitive.ObjectID) (*Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	character := Character{}
	err := r.FindOne(ctx, bson.M{"_id": id}).Decode(&character)

	return &character, err
}

func (r *CharactersRepo) GetCharactersByProfileID(profileID primitive.ObjectID) ([]Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.Find(ctx, bson.M{"profile_id": profileID})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var characters []Character
	err = cursor.All(ctx, &characters)
	if err != nil {
		return nil, err
	}

	return characters, nil
}

func (r *CharactersRepo) GetAllCharacters() ([]Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var characters []Character
	cursor, err := r.Find(ctx, primitive.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find characters: %v", err)
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &characters)
	if err != nil {
		return nil, err
	}

	return characters, nil
}

func (r *CharactersRepo) CreateCharacter(character *Character) (*Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.InsertOne(ctx, character)

	return character, err
}

func (r *CharactersRepo) UpdateCharacter(character *Character) (*Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.FindOneAndUpdate(ctx, bson.M{"_id": character.ID}, bson.M{"$set": character}, db.FindOneAndUpdateOptions).Decode(character)

	return character, err
}

func (r *CharactersRepo) DeleteCharacter(id primitive.ObjectID) (*Character, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	character := &Character{}
	err := r.FindOneAndDelete(ctx, bson.M{"_id": id}).Decode(character)

	return character, err
}

func (r *CharactersRepo) CreateCustomMetric(characterID primitive.ObjectID, metric *CustomMetric) (*CustomMetric, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.UpdateOne(ctx, bson.M{"_id": characterID}, bson.M{"$push": bson.M{"custom_metrics": *metric}})

	return metric, err
}

func (r *CharactersRepo) UpdateCustomMetric(characterID primitive.ObjectID, metric *CustomMetric) (*CustomMetric, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.UpdateOne(ctx, bson.M{"_id": characterID, "custom_metrics._id": metric.ID}, bson.M{"$set": bson.M{
		"custom_metrics.$.name":        metric.Name,
		"custom_metrics.$.description": metric.Description,
		"custom_metrics.$.style":       metric.Style,
		"custom_metrics.$.properties":  metric.Properties,
	}})

	return metric, err
}

func (r *CharactersRepo) DeleteCustomMetric(characterID primitive.ObjectID, metricID primitive.ObjectID) (*CustomMetric, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	metric := &CustomMetric{}
	err := r.FindOneAndUpdate(ctx, bson.M{"_id": characterID}, bson.M{
		"$pull": bson.M{
			"custom_metrics": bson.M{"_id": metricID},
		},
	}).Decode(metric)
	return metric, err
}

func (r *CharactersRepo) CreateMetricProperty(characterID primitive.ObjectID, metricID primitive.ObjectID, property *MetricProperty) (*MetricProperty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.UpdateOne(ctx, bson.M{"_id": characterID, "custom_metrics._id": metricID}, bson.M{"$push": bson.M{"custom_metrics.$.properties": *property}})

	return property, err
}

func (r *CharactersRepo) UpdateMetricProperty(characterID primitive.ObjectID, metricID primitive.ObjectID, property *MetricProperty) (*MetricProperty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.UpdateOne(ctx, bson.M{"_id": characterID, "custom_metrics._id": metricID, "custom_metrics.properties._id": property.ID}, bson.M{
		"$set": bson.M{
			"custom_metrics.$.properties.$[property].name":  property.Name,
			"custom_metrics.$.properties.$[property].type":  property.Type,
			"custom_metrics.$.properties.$[property].value": property.Value,
			"custom_metrics.$.properties.$[property].unit":  property.Unit,
		},
	}, &options.UpdateOptions{
		ArrayFilters: &options.ArrayFilters{
			Filters: []interface{}{bson.M{"property._id": property.ID}},
		},
	})

	return property, err
}

func (r *CharactersRepo) DeleteMetricProperty(characterID primitive.ObjectID, metricID primitive.ObjectID, propertyID primitive.ObjectID) (*MetricProperty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	property := &MetricProperty{}
	err := r.FindOneAndUpdate(ctx, bson.M{"_id": characterID, "custom_metrics._id": metricID}, bson.M{
		"$pull": bson.M{
			"custom_metrics.$.properties": bson.M{"_id": propertyID},
		},
	}).Decode(property)

	return property, err
}
