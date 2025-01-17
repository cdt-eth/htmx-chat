# Real-time Chat with Go, WebSocket, and HTMX

A real-time chat application built with Go, WebSocket for real-time messaging, and HTMX for dynamic UI interactions.

<img src="https://miro.medium.com/v2/resize:fit:1024/1*mr6lwBOE6xkRGOb0KTABaQ.png" width="400" alt="go-htmx">

## Features

- Real-time messaging using WebSocket
- Dynamic UI interactions with HTMX
- User authentication with JWT
- In-memory message storage
- Message deletion
- User identification

## Tech Stack

- Go for backend
- WebSocket for real-time communication
- HTMX for dynamic UI updates
- JWT for authentication
- Native CSS for styling
- Air for live reload during development

## Project Structure

```
htmx-chat/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── auth/
│   │   └── jwt.go          # JWT utilities
│   ├── handlers/
│   │   ├── auth.go         # Authentication handlers
│   │   └── chat.go         # WebSocket and chat handlers
│   ├── middleware/
│   │   └── auth.go         # Auth middleware
│   └── models/
│       ├── message.go      # Message model
│       └── user.go         # User model
├── templates/
│   ├── index.html          # Chat interface
│   └── login.html          # Auth interface
├── .air.toml               # Air configuration
├── .env                    # Environment variables
├── .env.example           # Example environment variables
├── .gitignore
├── go.mod                 # Go module definition
└── README.md
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

4. Set up environment variables:

```bash
cp .env.example .env
# Generate a JWT secret:
openssl rand -base64 32
# Add it to .env as JWT_SECRET
```

5. Run with Air (live reload):

```bash
air
```

Or run directly with Go:

```bash
go run cmd/server/main.go
```

6. Visit `http://localhost:8080` in your browser

## License

MIT
