package middlewares

import (
	"context"

	"tenkhours/pkg/auth"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type authInterceptor struct {
	authClient *AuthClient
}

func NewAuthInterceptor(authClient *AuthClient) *authInterceptor {
	return &authInterceptor{authClient: authClient}
}

func (a *authInterceptor) UnaryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	service := md.Get("service-name")
	if len(service) == 0 || len(service) > 0 && service[0] != "penk" {
		return handler(ctx, req)
	}

	userID := md.Get("x-user-id")
	if len(userID) == 0 {
		return nil, errors.New("missing user id")
	}

	authSession, err := a.authClient.IntrospectUser(ctx, "", userID[0], "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to introspect user")
	}

	// Save auth session to context
	ctx = context.WithValue(ctx, auth.AuthSessionKey, *authSession)

	return handler(ctx, req)
}
