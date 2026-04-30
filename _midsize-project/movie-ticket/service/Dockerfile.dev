FROM golang:alpine

# Install air for hot-reload
RUN go install github.com/air-verse/air@latest

# Set working directory
WORKDIR /app

# Run air
CMD ["air"]