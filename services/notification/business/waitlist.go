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
