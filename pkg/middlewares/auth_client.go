package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/pb"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authClient struct {
	CoreClient  pb.CoreClient
	RedisClient *redis.Client
}

func ComposeRPCClient() (*authClient, *grpc.ClientConn) {
	port, found := os.LookupEnv("CORE_RPC_PORT")
	if !found {
		port = "50051"
	}

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient("localhost"+":"+port, opts)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &authClient{CoreClient: pb.NewCoreClient(conn), RedisClient: db.GetRedisClient()}, conn
}

func (ac *authClient) IntrospectToken(ctx context.Context, token string) (*db.AuthSession, error) {
	firebaseProfile, err := auth.GetProfileByIDToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	var authSession db.AuthSession
	// Check if there is any active session in Redis
	keyFound, err := ac.RedisClient.Exists(context.Background(), db.GetAuthSessionKey(firebaseProfile.UID)).Result()
	if err != nil {
		return nil, err
	}

	// Cache hit, return the profile from redis
	if keyFound == 1 {
		profileJSON, err := ac.RedisClient.Get(context.Background(), db.GetAuthSessionKey(firebaseProfile.UID)).Result()
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(profileJSON), &authSession)
		if err != nil {
			return nil, err
		}
	}

	// Cache miss, create a new session from the profile in DB
	req := &pb.IntrospectReq{
		FirebaseUID: firebaseProfile.UID,
		Name:        firebaseProfile.Name,
		Email:       firebaseProfile.Email,
		Picture:     firebaseProfile.Picture,
	}
	res, err := ac.CoreClient.IntrospectProfile(ctx, req)
	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, fmt.Errorf("failed to introspect profile")
	}

	oid, err := primitive.ObjectIDFromHex(res.ProfileID)
	if err != nil {
		return nil, err
	}

	authSession.ProfileID = oid

	// Save profile in redis
	authSessionJSON, err := json.Marshal(authSession)
	if err != nil {
		return nil, err
	}

	err = ac.RedisClient.Set(context.Background(), db.GetAuthSessionKey(firebaseProfile.UID), authSessionJSON, time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return &authSession, nil
}
