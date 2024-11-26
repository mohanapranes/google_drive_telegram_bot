# google_drive_telegram_bot

This is a Go-based Telegram bot that sends a file to a specific chat once every 24 hours. The bot is containerized using Docker for easy deployment.

## Features

- Responds to the `/updateId` command to set the main chat ID.
- Sends a file (e.g., a document) to the specified chat ID every 24 hours.
- Uses Docker for containerization.

## Prerequisites

Before running the bot, make sure you have the following:

- A [Telegram Bot API token](https://core.telegram.org/bots#botfather).
- Go (version 1.20 or higher) installed on your system for building the project (if you are not using Docker).
- Docker (if you prefer to use Docker for running the bot).

## Installation

### 1. Clone the Repository

Clone this repository to your local machine:

```bash
git clone https://github.com/mohanapranes/google_drive_telegram_bot.git
cd google_drive_telegram_bot
```

2. Set Up the Go Project
Option A: Build Without Docker

Install Go and dependencies:

```bash
go mod tidy
```

Build the bot:

```bash
go build -o bot .
```

Run the bot:
```bash
./bot
```
Option B: Build and Run With Docker

Build the Docker image:

```bash
docker build -t telegram-bot-file-sender .
```

Run the Docker container:
```bash
docker run -d --name telegram-bot-container telegram-bot-file-sender
```

