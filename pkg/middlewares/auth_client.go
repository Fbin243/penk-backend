package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/pb"

	rdb "tenkhours/pkg/db/redis"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authClient struct {
	CoreClient  pb.CoreClient
	RedisClient *redis.Client
}

func ComposeAuthClient() (*authClient, *grpc.ClientConn) {
	port, found := os.LookupEnv("CORE_GRPC_PORT")
	if !found {
		port = "50051"
	}

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient("localhost"+":"+port, opts)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &authClient{CoreClient: pb.NewCoreClient(conn), RedisClient: rdb.GetRedisClient()}, conn
}

func (ac *authClient) IntrospectToken(ctx context.Context, token string) (*rdb.AuthSession, error) {
	firebaseProfile, err := auth.GetProfileByIDToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	var authSession rdb.AuthSession
	// Check if there is any active session in Redis
	keyFound, err := ac.RedisClient.Exists(context.Background(), rdb.GetAuthSessionKey(firebaseProfile.UID)).Result()
	if err != nil {
		return nil, err
	}

	// Cache hit, return the profile from redis
	if keyFound == 1 {
		profileJSON, err := ac.RedisClient.Get(context.Background(), rdb.GetAuthSessionKey(firebaseProfile.UID)).Result()
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

	authSession.ProfileID = res.ProfileID

	// Save profile in redis
	authSessionJSON, err := json.Marshal(authSession)
	if err != nil {
		return nil, err
	}

	err = ac.RedisClient.Set(context.Background(), rdb.GetAuthSessionKey(firebaseProfile.UID), authSessionJSON, time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return &authSession, nil
}
