package composer

import (
	"log"
	"os"

	"tenkhours/proto/pb/analytic"
	"tenkhours/proto/pb/currency"
	"tenkhours/services/core/repo/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ComposeAnalyticClient() (*rpc.AnalyticClient, *grpc.ClientConn) {
	port, found := os.LookupEnv("ANALYTIC_GPRC_PORT")
	if !found {
		port = "50052"
	}

	host, found := os.LookupEnv("ANALYTIC_GRPC_HOST")
	if !found {
		host = "localhost"
	}

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(host+":"+port, opts)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return rpc.NewAnalyticClient(analytic.NewAnalyticClient(conn)), conn
}

func ComposeCurrencyClient() (*rpc.CurrencyClient, *grpc.ClientConn) {
	port, found := os.LookupEnv("CURRENCY_GPRC_PORT")
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
