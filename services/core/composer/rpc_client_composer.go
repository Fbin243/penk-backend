package composer

import (
	"log"
	"os"

	"tenkhours/pkg/pb"
	"tenkhours/services/core/repo/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ComposeAnalyticClient() (*rpc.AnalyticClient, *grpc.ClientConn) {
	port, found := os.LookupEnv("ANALYTIC_GPRC_PORT")
	if !found {
		port = "50052"
	}

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient("localhost"+":"+port, opts)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return rpc.NewAnalyticClient(pb.NewAnalyticClient(conn)), conn
}

func ComposeCurrencyClient() (*rpc.CurrencyClient, *grpc.ClientConn) {
	port, found := os.LookupEnv("CURRENCY_GPRC_PORT")
	if !found {
		port = "50055"
	}

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient("localhost"+":"+port, opts)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return rpc.NewCurrencyClient(pb.NewCurrencyClient(conn)), conn
}
