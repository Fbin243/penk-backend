package auth

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func getFirebaseApp() (*firebase.App, error) {
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_ADMIN"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	return app, nil
}

func getAuthClient() (*auth.Client, error) {
	app, err := getFirebaseApp()
	if err != nil {
		return nil, err
	}

	return app.Auth(context.Background())
}

func GetProfileByIDToken(idToken string) (*Profile, error) {
	authClient, err := getAuthClient()
	if err != nil {
		return nil, err
	}

	token, err := authClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return nil, err
	}

	authProfile := Profile{
		UID:           token.UID,
		Email:         token.Claims["email"].(string),
		EmailVerified: token.Claims["email_verified"].(bool),
		Name:          token.Claims["name"].(string),
	}

	return &authProfile, nil
}

func GetProfileByContext(ctx context.Context) (*Profile, error) {
	untypedProfile := ctx.Value(ProfileContextKey)
	if untypedProfile == nil {
		return nil, ErrorProfileNotFound
	}

	profile, ok := untypedProfile.(Profile)
	if !ok {
		return nil, ErrorCannotParseProfile
	}

	return &profile, nil
}
