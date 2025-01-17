FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with environment variable support
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/

EXPOSE 8080

# Use environment variables from Railway
CMD ["./main"] 