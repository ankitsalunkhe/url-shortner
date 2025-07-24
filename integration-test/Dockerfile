FROM golang:latest AS service
WORKDIR /app
COPY . .    
WORKDIR /app

RUN go build -o url-shortner ./internal/cmd/main.go