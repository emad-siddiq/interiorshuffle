# Use Go 1.23.4 base image
FROM golang:1.23.4-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download the Go dependencies
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o myapp .

# Expose the port (modify this based on your app's port)
EXPOSE 8080

# Run the application
CMD ["./myapp"]
