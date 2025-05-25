# 1. Base image
FROM golang:1.24-alpine

# 2. Set working directory
WORKDIR /app

# 3. Copy go.mod & go.sum
COPY go.mod ./
COPY go.sum ./

# 4. Download dependencies
RUN go mod download

# 5. Copy the entire backend
COPY . .

# 6. Build the Go app
RUN go build -o main ./cmd/server

# 7. Expose port and run
EXPOSE 8080
CMD ["./main"]
