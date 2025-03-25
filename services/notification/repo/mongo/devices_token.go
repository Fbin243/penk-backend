package mongorepo

import (
	"context"
	"log"
	"time"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/utils"
	"tenkhours/services/notification/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DeviceTokenRepo struct {
	*mongodb.BaseRepo[entity.DevicesToken, DevicesToken]
}

func NewDevicesTokenRepo(db *mongo.Database) *DeviceTokenRepo {
	devicesTokenCollection := db.Collection(mongodb.DevicesTokensCollection)
	_, err := devicesTokenCollection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "profile_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		log.Println("failed to create indexes for devices token collection")
		return nil
	}

	return &DeviceTokenRepo{
		mongodb.NewBaseRepo(devicesTokenCollection,
			&mongodb.Mapper[entity.DevicesToken, DevicesToken]{},
			true),
	}
}

// UpsertDeviceToken adds a new device token to the user's profile or updates the existing token for the given deviceID
func (r *DeviceTokenRepo) UpsertDeviceToken(ctx context.Context, profileID, token, deviceID, platform string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"profile_id": profileID, "tokens.device_id": deviceID}
	update := bson.M{
		"$set": bson.M{"tokens.$.token": token, "tokens.$.platform": platform, "tokens.$.create_at": time.Now().Format(time.RFC3339)},
	}

	opts := options.Update().SetUpsert(true)

	_, err := r.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		// If the token for the deviceID does not exist, add it
		if err == mongo.ErrNoDocuments {
			filter = bson.M{"profile_id": profileID}
			update = bson.M{
				"$addToSet": bson.M{"tokens": bson.M{
					"device_id": deviceID,
					"token":     token,
					"platform":  platform,
					"create_at": utils.Now().Format(time.RFC3339),
				}},
			}
			_, err = r.UpdateOne(ctx, filter, update, opts)
		}
	}
	return err
}

// RemoveDeviceToken removes a device token from the user's profile
func (r *DeviceTokenRepo) RemoveDeviceToken(ctx context.Context, profileID, token string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"profile_id": profileID}
	update := bson.M{
		"$pull": bson.M{"tokens": bson.M{"token": token}},
	}

	opts := options.Update()
	_, err := r.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *DeviceTokenRepo) GetDeviceTokenByDeviceID(ctx context.Context, deviceID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"tokens.device_id": deviceID}
	devicesToken := new(entity.DevicesToken)

	err := r.FindOne(ctx, filter).Decode(devicesToken)
	if err != nil {
		return "", err
	}

	for _, token := range devicesToken.Tokens {
		if token.DeviceID == deviceID {
			return token.Token, nil
		}
	}
	return "", nil
}
