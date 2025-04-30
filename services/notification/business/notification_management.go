package business

import (
	"context"
	"fmt"
	"log"
	"os"

	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type NotificationBusiness struct {
	MessagingClient  *messaging.Client
	DevicesTokenRepo IDeviceTokenRepo
}

func NewNotificationBusiness(messagingClient *messaging.Client, devicesTokenRepo IDeviceTokenRepo) *NotificationBusiness {
	return &NotificationBusiness{
		MessagingClient:  messagingClient,
		DevicesTokenRepo: devicesTokenRepo,
	}
}

func (n *NotificationBusiness) AddEmailToWaitlist(ctx context.Context, email string) error {
	googleSheetsCredentials := os.Getenv("GOOGLE_SHEETS_CREDENTIALS")
	spreadsheetID := os.Getenv("GOOGLE_SHEET_ID")
	sheetName := os.Getenv("GOOGLE_SHEET_NAME")

	creds := option.WithCredentialsJSON([]byte(googleSheetsCredentials))

	srv, err := sheets.NewService(ctx, creds)
	if err != nil {
		return fmt.Errorf("failed to create sheets service: %w", err)
	}

	vr := &sheets.ValueRange{
		Values: [][]interface{}{{email}},
	}

	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, sheetName, vr).
		ValueInputOption("USER_ENTERED").
		Do()
	if err != nil {
		return fmt.Errorf("failed to append to sheet: %w", err)
	}

	log.Printf("Successfully added %s to Google Sheet.", email)
	return nil
}

func (biz *NotificationBusiness) RegisterDeviceToken(ctx context.Context, profileID, token, deviceID, platform string) (bool, error) {
	err := biz.DevicesTokenRepo.UpsertDeviceToken(ctx, profileID, token, deviceID, platform)
	if err != nil {
		return false, fmt.Errorf("failed to add Devices Token: %v", err)
	}

	return true, nil
}

// RemoveDeviceToken removes a device token from the user's profile
func (biz *NotificationBusiness) RemoveDeviceToken(ctx context.Context, profileID, token string) (bool, error) {
	err := biz.DevicesTokenRepo.RemoveDeviceToken(ctx, profileID, token)
	if err != nil {
		return false, fmt.Errorf("failed to remove Devices Token: %v", err)
	}

	return true, nil
}
