# 🔧 Workspace Channel Cleaner CLI

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ashfaqahmed/workspace-channel-cleaner-cli)](https://goreportcard.com/report/github.com/ashfaqahmed/workspace-channel-cleaner-cli)

> A simple, fast CLI tool I built to help clean up stale channels in Slack workspaces. Perfect for automation, batch processing, and quick cleanup tasks.
> 💡 **Looking for a beautiful interactive TUI version?** Check out [workspace-channels-cleaner](https://github.com/ashfaqahmed/workspace-channels-cleaner) for a full-featured terminal user interface with pagination, configuration management, and skip list editing.


## ✨ Features

### Fast & Efficient
- **Concurrent Processing**: Uses goroutines for efficient API calls
- **Rate Limit Handling**: Automatic retry with exponential backoff
- **Batch Operations**: Process multiple channels at once
- **Lightweight**: No dependencies on UI frameworks

### 🔍 Smart Filtering
- **Time-based**: Find channels inactive for N days
- **Keyword-based**: Filter channels by name patterns
- **Type Selection**: Target public, private, or both channel types
- **Skip List Protection**: Important channels are automatically protected

### Safety First
- **Interactive Confirmation**: Always asks before leaving channels
- **Hardcoded Skip List**: Protects critical channels from accidental removal
- **Error Handling**: Graceful failure without data loss
- **Rate Limit Respect**: Never overwhelms Slack's API

## 🚀 Quick Start

### Prerequisites
- Go 1.24 or higher
- Slack API token with required scopes

### Installation

```bash
# Clone the repository
git clone https://github.com/ashfaqahmed/workspace-channel-cleaner-cli
cd workspace-channel-cleaner-cli

# Install dependencies
go mod tidy

# Build the application
go build -o workspace-cleaner-cli main.go

# Run the application
./workspace-cleaner-cli --days 30
```

## 📖 Usage Examples

### Basic Usage

```bash
# Find channels inactive for 30 days
./workspace-cleaner-cli --days 30

# Find channels with "test" in the name
./workspace-cleaner-cli --keyword test

# Find channels inactive for 60 days with verbose output
./workspace-cleaner-cli --days 60 --verbose

# Find only public channels inactive for 30 days
./workspace-cleaner-cli --days 30 --types public

# Find only private channels with "old" in the name
./workspace-cleaner-cli --keyword old --types private

# Find both public and private channels (default)
./workspace-cleaner-cli --days 45 --types public,private
```

### Command Line Arguments

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--days` | int | 0 | Filter channels inactive for N days |
| `--keyword` | string | "" | Filter channels containing keyword in name |
| `--types` | string | "public,private" | Channel types: public, private, or both |
| `--verbose` | bool | false | Enable detailed logging |

### Output Example
```
[1/5] #old-project-2023 ID: C1234567890 — Last message: Mon, 15 Jan 2024 10:30:00 UTC
[2/5] #test-channel ID: C0987654321 — Last message: Tue, 20 Dec 2023 15:45:00 UTC
[3/5] #archived-discussions ID: C1122334455 — Last message: Wed, 10 Jan 2024 09:15:00 UTC
[4/5] #deprecated-feature ID: C2233445566 — Last message: Thu, 05 Jan 2024 14:20:00 UTC
[5/5] #old-team-chat ID: C3344556677 — Last message: Fri, 29 Dec 2023 11:30:00 UTC
❓ Do you want to leave these channels now? (y/N):
```

### 💡 Best Practices for Rate Limit Management

To avoid hitting Slack's API rate limits, use specific search criteria:

**✅ Recommended (Narrow searches):**
```bash
# Specific keyword with time filter
./workspace-cleaner-cli --keyword "test-2023" --days 30

# Only public channels with specific term
./workspace-cleaner-cli --keyword "deprecated" --types public

# Short time range for initial testing
./workspace-cleaner-cli --days 7 --keyword "old"
```

**❌ Avoid (Broad searches):**
```bash
# Too broad - may hit rate limits
./workspace-cleaner-cli --days 365

# Very generic keyword
./workspace-cleaner-cli --keyword "test"
```

**Pro Tips:**
- Start with smaller time ranges (`--days 7` or `--days 30`)
- Use specific keywords that match your channel naming patterns
- Combine filters to narrow down results
- Use `--verbose` to monitor API call progress

## ⚙️ Configuration

### Environment Setup

1. Copy the example environment file:
```bash
cp example.env .env
```

2. Edit the `.env` file and add your Slack API token:
```env
# Required: Your Slack API token
SLACK_API_TOKEN=xoxb-your-slack-bot-token-here
```

**Important**: Never commit your `.env` file to version control. The `.gitignore` file is configured to exclude it.

### API Token Requirements

Your Slack API token must have the following scopes:

- `channels:history` - Read channel message history
- `groups:history` - Read private channel history
- `conversations.list` - List all channels
- `conversations.leave` - Leave channels

### Getting Your API Token

1. Go to [Slack API Apps](https://api.slack.com/apps)
2. Create a new app or select an existing one
3. Go to "OAuth & Permissions"
4. Add the required scopes listed above
5. Install the app to your workspace
6. Copy the "Bot User OAuth Token" (starts with `xoxb-`) or "User OAuth Token" (starts with `xoxp-`)

## 🛡️ Protected Channels

The application includes a configurable skip list to protect important channels. By default, it protects common critical channels like:

- `general`, `company-announcements`
- Team channels: `team-leads`, `support-example`, `dev-team`
- Security channels: `security`, `infrastructure`, `admin-only`
- And many more critical channels

### Customizing the Skip List

You can customize which channels to skip by editing the `skip-channels.json` file:

```json
{
  "skip_channels": [
    "general",
    "company-announcements",
    "team-leads",
    "support-example",
    "dev-team",
    "infrastructure",
    "security",
    "hr-announcements",
    "product-updates"
  ]
}
```

**Note**: If the JSON file is missing or invalid, the application will use a built-in default skip list for safety.

## 🔧 Development

### Project Structure
```workspace-channel-cleaner-cli/
├── main.go # Main CLI application
├── go.mod # Go module file
├── go.sum # Go dependencies checksum
├── .env # Environment variables (create from example.env)
├── example.env # Example environment file
├── skip-channels.json # Configurable channel skip list
├── .gitignore # Git ignore rules
└── README.md # This file
```


### Building from Source

```bash
# Clone the repository
git clone https://github.com/ashfaqahmed/workspace-channel-cleaner-cli
cd workspace-channel-cleaner-cli

# Install dependencies
go mod tidy

# Build the application
go build -o workspace-cleaner-cli main.go

# Run the application
./workspace-cleaner-cli --days 30
```

### Dependencies

- `github.com/slack-go/slack` - Slack API client
- `github.com/joho/godotenv` - Environment variable loading

## 🚀 Automation & Scripting

This CLI tool is perfect for automation:

### Cron Job Example

```bash
# Add to crontab to run weekly cleanup
0 9 * * 1 /path/to/workspace-cleaner-cli --days 30 --verbose >> /var/log/channel-cleanup.log 2>&1
```

### CI/CD Integration

```yaml
# GitHub Actions example
- name: Clean up stale channels
  run: |
    ./workspace-cleaner-cli --days 60 --types public
  env:
    SLACK_API_TOKEN: ${{ secrets.SLACK_API_TOKEN }}
```

### Docker Usage

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o workspace-cleaner-cli main.go

FROM alpine:latest
COPY --from=builder /app/workspace-cleaner-cli /usr/local/bin/
ENTRYPOINT ["workspace-cleaner-cli"]
```

## 🔒 Security

- **Token Protection**: API tokens are never logged or displayed
- **Skip List**: Protects important channels from accidental removal
- **Confirmation**: Double-confirmation for destructive actions
- **Error Handling**: Graceful failure without data loss
- **Rate Limiting**: Respects Slack's API rate limits

## 🐛 Troubleshooting

### Common Issues

**"SLACK_API_TOKEN not set"**
- Ensure your `.env` file exists and contains the token
- Check that the token has the required scopes

**"Rate limit hit"**
- The application automatically handles rate limits
- Wait for the retry mechanism to complete
- **💡 Tip**: Make your search criteria more specific to reduce API calls:
  - Use `--keyword "specific-term"` instead of broad searches
  - Combine `--days` with `--keyword` for targeted filtering
  - Use `--types public` or `--types private` to limit scope
  - Start with smaller time ranges (e.g., `--days 7`) before larger ones

**"No matching channels found"**
- Check your filter settings
- Verify you're a member of channels
- Review the skip list for protected channels

**"You must provide either --days or --keyword"**
- Provide at least one filtering criteria
- Use `--days 30` for time-based filtering
- Use `--keyword test` for keyword-based filtering

### Debug Mode

Enable verbose output for detailed logging:

```bash
./workspace-cleaner-cli --days 30 --verbose
```

## Contributing

Feel free to contribute! If you want to add something big, open an issue first so we can talk about it.

1. Fork the repo
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go coding standards
- Add tests for new features
- Update docs when needed
- Keep the CLI simple and focused

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## ⚖️ Legal Disclaimers

**Important**: Please read our Legal Disclaimer for important information about third-party trademarks, copyrights, and terms of use.

### Copyright Disclaimer

This project is **NOT** affiliated with, endorsed by, or sponsored by Slack Technologies, Inc. or any of its subsidiaries.

### Third-Party Trademarks and Copyrights

- **Slack** is a registered trademark of Slack Technologies, Inc.
- **Slack API** and related services are owned by Slack Technologies, Inc.
- This project uses the official Slack Go SDK which is subject to its own license terms.

### Fair Use

This project is developed for educational and productivity purposes, using Slack's publicly available API in accordance with their API Terms of Service. The use of Slack's API and SDK is subject to Slack's own terms and conditions.

### No Warranty

This project is provided "as is" without any warranties. Users are responsible for ensuring their use complies with Slack's terms of service and applicable laws.

## Acknowledgments

- Slack Go SDK - Slack API client
- The open-source community for inspiration and tools

---

## ☕ Support My Work

If this tool saved you time or effort, consider buying me a coffee.
Your support helps me keep building and maintaining open-source projects like this!

You can either scan the QR code below or click the link to tip me:

👉 [buymeacoffee.com/ashfaqueali](https://buymeacoffee.com/ashfaqueali)

**Buy Me a Coffee QR**

<img src="https://ashfaqsolangi.com/images/bmc_qr.png" alt="Buy Me a Coffee QR" width="220" height="220" />

---

**Happy channel cleaning! 🧹✨**

*Made with ❤️ by [Ashfaque Ali](https://github.com/ashfaqahmed)* 
