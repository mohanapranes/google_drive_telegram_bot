# Use an official Golang image as the base image
FROM golang:1.23-alpine AS builder

# Set the working directory in the container
WORKDIR /app

# Copy the Go modules and go.sum files
COPY go.mod go.sum .env ./

# Download the Go dependencies
RUN go mod tidy

# Copy the rest of the application code into the container
COPY . .

# Build the Go binary
RUN go build -o bot .

# Use a smaller image to run the app
FROM alpine:latest

# Install dependencies (for example, for glibc or other libraries if needed)
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/bot .

# Expose port (optional, depending on how your bot interacts)
EXPOSE 8080

# Command to run the bot
CMD ["./bot"]
