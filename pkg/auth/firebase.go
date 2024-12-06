package auth

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type FirebaseManager struct {
	Client *auth.Client
	App    *firebase.App
}

var firebaseManager *FirebaseManager

func GetFirebaseManager() *FirebaseManager {
	if firebaseManager == nil {
		firebaseManager = &FirebaseManager{}
		opt := option.WithCredentialsFile(os.Getenv("FIREBASE_ADMIN"))
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			log.Fatal(err)
		}

		firebaseManager.App = app

		client, err := app.Auth(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		firebaseManager.Client = client
	}

	return firebaseManager
}

func GetProfileByIDToken(idToken string) (*FirebaseProfile, error) {
	token, err := GetFirebaseManager().Client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return nil, err
	}

	authProfile := FirebaseProfile{
		UID: token.UID,
	}

	if token.Claims["email"] != nil {
		authProfile.Email = token.Claims["email"].(string)
	}

	if token.Claims["name"] != nil {
		authProfile.Name = token.Claims["name"].(string)
	}

	if token.Claims["picture"] != nil {
		authProfile.Picture = token.Claims["picture"].(string)
	}

	return &authProfile, nil
}

func DeleteProfileOnFirebase(uid string) error {
	return GetFirebaseManager().Client.DeleteUser(context.Background(), uid)
}
