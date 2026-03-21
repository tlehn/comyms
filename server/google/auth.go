package google

import (
	"context"
	"fmt"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// NewSheetsService creates an authenticated Google Sheets API client using Application Default Credentials.
// The service is scoped to read-only access (SpreadsheetsReadonlyScope).
//
// Returns an error if GOOGLE_APPLICATION_CREDENTIALS is not set or points to an invalid/expired
// service account key, or if the Sheets API is not enabled in the GCP project.
func NewSheetsService(ctx context.Context) (*sheets.Service, error) {
	creds, err := google.FindDefaultCredentials(ctx, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("finding credentials: %w", err)
	}

	srv, err := sheets.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("creating sheets service: %w", err)
	}

	return srv, nil
}

// NewDriveService creates an authenticated Google Drive API client using Application Default Credentials.
// The service is scoped to read-only access (DriveReadonlyScope).
//
// Returns an error if GOOGLE_APPLICATION_CREDENTIALS is not set or points to an invalid/expired
// service account key, or if the Drive API is not enabled in the GCP project.
func NewDriveService(ctx context.Context) (*drive.Service, error) {
	creds, err := google.FindDefaultCredentials(ctx, drive.DriveReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("finding credentials: %w", err)
	}

	srv, err := drive.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("creating drive service: %w", err)
	}

	return srv, nil
}
