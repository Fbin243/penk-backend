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
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	service := md.Get("service-name")
	if len(service) == 0 || len(service) > 0 && service[0] != "penk" {
		return handler(ctx, req)
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return nil, errors.New("missing authorization token")
	}

	token := authHeaders[0]
	authSession, err := a.authClient.IntrospectToken(ctx, token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to introspect token")
	}

	// Save auth session to context
	ctx = context.WithValue(ctx, auth.AuthSessionKey, *authSession)

	return handler(ctx, req)
}
