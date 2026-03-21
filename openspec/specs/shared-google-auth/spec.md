## ADDED Requirements

### Requirement: Shared Google ADC credential setup

The `server/google` package SHALL provide exported functions to create authenticated Google API clients using
Application Default Credentials (ADC).

#### Scenario: Create Sheets service

- **WHEN** `NewSheetsService(ctx)` is called with valid ADC configured
- **THEN** it returns a `*sheets.Service` scoped to `SpreadsheetsReadonlyScope`

#### Scenario: Create Drive service

- **WHEN** `NewDriveService(ctx)` is called with valid ADC configured
- **THEN** it returns a `*drive.Service` scoped to `DriveReadonlyScope`

#### Scenario: Missing credentials

- **WHEN** `NewSheetsService(ctx)` or `NewDriveService(ctx)` is called without ADC configured
- **THEN** it returns a descriptive error wrapping the underlying credential failure

### Requirement: Auth functions are reusable across server packages

Any server package under `server/google/` SHALL be able to import and use the auth functions without duplicating
credential logic.

#### Scenario: Multiple server packages share auth

- **WHEN** a new server package `server/google/calendar/` is added
- **THEN** it can call `google.NewCalendarService(ctx)` (or similar) following the same pattern, with the auth function
  defined in `server/google/auth.go`
