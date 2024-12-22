package composer

import (
	"log"
	"os"
	"tenkhours/pkg/pb"
	"tenkhours/services/timetrackings/repo/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ComposeRPCClient() (*rpc.RPCClient, *grpc.ClientConn) {
	port, found := os.LookupEnv("CORE_RPC_PORT")
	if !found {
		port = "50051"
	}

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient("localhost"+":"+port, opts)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return rpc.NewRPCClient(pb.NewCoreClient(conn)), conn
}
