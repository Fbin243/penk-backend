package middlewares

import (
	"context"
	"fmt"
	"log"
	"os"

	rdb "tenkhours/pkg/db/redis"
	"tenkhours/proto/pb/core"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	core.CoreClient
}

func ComposeAuthClient() (*AuthClient, *grpc.ClientConn) {
	port, found := os.LookupEnv("CORE_GRPC_PORT")
	if !found {
		port = "50051"
	}

	host, found := os.LookupEnv("CORE_GRPC_HOST")
	if !found {
		host = "localhost"
	}

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(host+":"+port, opts)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return &AuthClient{CoreClient: core.NewCoreClient(conn)}, conn
}

func (ac *AuthClient) IntrospectUser(ctx context.Context, token, userID, deviceID string) (*rdb.AuthSession, error) {
	req := &core.IntrospectReq{
		Token:    token,
		UserId:   userID,
		DeviceId: deviceID,
	}

	res, err := ac.CoreClient.IntrospectUser(ctx, req)
	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, fmt.Errorf("failed to introspect token")
	}

	return &rdb.AuthSession{
		ProfileID:          res.ProfileId,
		DeviceID:           res.DeviceId,
		FirebaseUID:        res.FirebaseUid,
		CurrentCharacterID: res.CurrentCharacterId,
	}, nil
}
