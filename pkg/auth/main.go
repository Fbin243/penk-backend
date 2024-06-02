package auth

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

func GetFirebaseApp() (*firebase.App, error) {
	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_ADMIN"))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	return app, nil
}

func GetAuthClient() (*auth.Client, error) {
	app, err := GetFirebaseApp()
	if err != nil {
		return nil, err
	}

	return app.Auth(context.Background())
}

type AuthProfile struct {
	UID           string
	Email         string
	EmailVerified bool
	Name          string
}

func VerifyClientIDToken(idToken string) (*AuthProfile, error) {
	authClient, err := GetAuthClient()
	if err != nil {
		return nil, err
	}

	token, err := authClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return nil, err
	}

	authProfile := AuthProfile{
		UID:           token.UID,
		Email:         token.Claims["email"].(string),
		EmailVerified: token.Claims["email_verified"].(bool),
		Name:          token.Claims["name"].(string),
	}

	return &authProfile, nil
}
