package notification

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type NotificationService interface {
	AddEmailToWaitlist(ctx context.Context, email string) error
}

type notificationService struct{}

func NewNotificationService() NotificationService {
	return &notificationService{}
}

func (n *notificationService) AddEmailToWaitlist(ctx context.Context, email string) error {
	googleSheetsCredentials := os.Getenv("GOOGLE_SHEETS_CREDENTIALS")
	spreadsheetID := os.Getenv("GOOGLE_SHEET_ID")
	sheetName := os.Getenv("GOOGLE_SHEET_NAME")

	creds := option.WithCredentialsJSON([]byte(googleSheetsCredentials))

	srv, err := sheets.NewService(ctx, creds)
	if err != nil {
		return fmt.Errorf("Failed to create sheets service: %w", err)
	}

	vr := &sheets.ValueRange{
		Values: [][]interface{}{{email}},
	}

	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, sheetName, vr).
		ValueInputOption("USER_ENTERED").
		Do()
	if err != nil {
		return fmt.Errorf("Failed to append to sheet: %w", err)
	}

	log.Printf("Successfully added %s to Google Sheet.", email)
	return nil
}
