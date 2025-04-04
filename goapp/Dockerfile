# Use a newer Go version
FROM golang:1.20 AS build

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (better caching)
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the source code
COPY . .

# Build the Go application
RUN go build -o goapp .

# Create a smaller final image
FROM golang:1.20
WORKDIR /app
COPY --from=build /app/goapp .

# Expose port 8888
EXPOSE 8888

# Run the application
CMD ["/app/goapp"]
