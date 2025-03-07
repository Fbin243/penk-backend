package composer

import (
	"log"
	"os"

	"tenkhours/proto/pb/core"
	"tenkhours/proto/pb/currency"
	"tenkhours/services/timetracking/repo/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ComposeCoreClient() (*rpc.CoreClient, *grpc.ClientConn) {
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

	return rpc.NewCoreClient(core.NewCoreClient(conn)), conn
}

func ComposeCurrencyClient() (*rpc.CurrencyClient, *grpc.ClientConn) {
	port, found := os.LookupEnv("CURRENCY_GRPC_PORT")
	if !found {
		port = "50055"
	}

	host, found := os.LookupEnv("CURRENCY_GRPC_HOST")
	if !found {
		host = "localhost"
	}

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(host+":"+port, opts)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return rpc.NewCurrencyClient(currency.NewCurrencyClient(conn)), conn
}
