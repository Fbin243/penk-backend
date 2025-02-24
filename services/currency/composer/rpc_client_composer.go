package composer

import (
	"log"
	"os"

	"tenkhours/proto/pb/core"
	"tenkhours/services/currency/repo/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ComposeCoreClient() (*rpc.CoreClient, *grpc.ClientConn) {
	port, found := os.LookupEnv("CORE_GPRC_PORT")
	if !found {
		port = "50051"
	}

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient("localhost"+":"+port, opts)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return rpc.NewCoreClient(core.NewCoreClient(conn)), conn
}
