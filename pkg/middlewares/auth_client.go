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

func (ac *AuthClient) IntrospectToken(ctx context.Context, token, deviceID string) (*rdb.AuthSession, error) {
	req := &core.IntrospectReq{
		Token:    token,
		DeviceId: deviceID,
	}
	// TODO: @Fbin243 refactor auth later
	res, err := ac.CoreClient.IntrospectToken(ctx, req)
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
