# Real-time Chat with Go and HTMX

A simple real-time chat application built with Go and HTMX.

<img src="https://miro.medium.com/v2/resize:fit:1024/1*mr6lwBOE6xkRGOb0KTABaQ.png" width="400" alt="Chat App Screenshot">

## Features

- Real-time messaging using HTMX
- In-memory message storage
- Message deletion
- User identification
- Auto-refresh every 2 seconds

## Tech Stack

- Go for backend
- HTMX for dynamic updates
- Native CSS for styling
- Air for live reload during development

## Project Structure

```
htmx-chat/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── handlers/
│   │   └── chat.go         # HTTP handlers for chat functionality
│   └── models/
│       └── message.go      # Message data structure and storage
├── templates/
│   └── index.html          # Main chat interface
├── .air.toml               # Air configuration for live reload
├── .env                    # Environment variables
├── .env.example           # Example environment variables
├── .gitignore             # Git ignore file
├── go.mod                 # Go module definition
└── README.md              # Project documentation
```

## Running Locally

1. Clone the repository:

```bash
git clone https://github.com/cdt-eth/htmx-chat.git
```

2. Install dependencies:

```bash
go mod tidy
```

3. Install Air for live reload:

```bash
go install github.com/cosmtrek/air@latest
```

4. Run with Air (live reload):

```bash
air
```

Or run directly with Go:

```bash
go run cmd/server/main.go
```

5. Visit `http://localhost:8080` in your browser

## License

MIT
