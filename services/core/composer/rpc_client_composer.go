package composer

import (
	"log"
	"os"

	"tenkhours/proto/pb/currency"
	"tenkhours/proto/pb/notification"
	"tenkhours/services/core/repo/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

func ComposeNotificationClient() (*rpc.NotificationClient, *grpc.ClientConn) {
	port, found := os.LookupEnv("NOTIFICATION_GRPC_PORT")
	if !found {
		port = "50054"
	}

	host, found := os.LookupEnv("NOTIFICATION_GRPC_HOST")
	if !found {
		host = "localhost"
	}

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(host+":"+port, opts)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return rpc.NewNotificationClient(notification.NewNotificationClient(conn)), conn
}
