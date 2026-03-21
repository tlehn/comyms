# comyms

A CLI that hosts multiple MCP servers as subcommands.

## Available Servers

- `comyms google sheets` — Google Sheets MCP server

## Prerequisites

- Go 1.21+
- A Google Cloud project with the Google Sheets API enabled

## Google Cloud Setup

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project (or select an existing one)
3. Enable the **Google Sheets API**:
   - Navigate to **APIs & Services > Library**
   - Search for "Google Drive API" and click **Enable**
   - Search for "Google Sheets API" and click **Enable**
4. Create a service account:
   - Navigate to **APIs & Services > Credentials**
   - Click **Create Credentials > Service account**
   - Give it a name and click **Done**
5. Create a key for the service account:
   - Click on the service account you just created
   - Go to the **Keys** tab
   - Click **Add Key > Create new key > JSON**
   - Save the downloaded JSON file somewhere safe (e.g., `~/.config/comyms/credentials.json`)
6. Share your spreadsheet(s) with the service account's email address (found in the JSON file as `client_email`)

## MCP Configuration

Add the following to your MCP config (e.g., `claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "google-sheets": {
      "command": "/path/to/comyms",
      "args": ["google", "sheets"],
      "env": {
        "GOOGLE_APPLICATION_CREDENTIALS": "/path/to/credentials.json"
      }
    }
  }
}
```

## Building

```sh
make build
```

## Development

```sh
make run    # run the CLI
make test   # run tests
make lint   # vet the code
make clean  # remove the binary
```
