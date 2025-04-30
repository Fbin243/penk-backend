package auth

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type FirebaseManager struct {
	Client          *auth.Client
	MessagingClient *messaging.Client
	App             *firebase.App
}

var firebaseManager *FirebaseManager

func GetFirebaseManager() *FirebaseManager {
	if firebaseManager == nil {
		opt := option.WithCredentialsFile(os.Getenv("FIREBASE_ADMIN"))
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			log.Fatal(err)
		}

		client, err := app.Auth(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		messagingClient, err := app.Messaging(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		firebaseManager = &FirebaseManager{
			Client:          client,
			App:             app,
			MessagingClient: messagingClient,
		}
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

	nameClaim, ok := token.Claims["name"].(string)
	if !ok || nameClaim == "" {
		authProfile.Name = "Anonymous"
	} else {
		authProfile.Name = nameClaim
	}

	if token.Claims["picture"] != nil {
		authProfile.Picture = token.Claims["picture"].(string)
	}

	return &authProfile, nil
}

func DeleteProfileOnFirebase(uid string) error {
	return GetFirebaseManager().Client.DeleteUser(context.Background(), uid)
}
