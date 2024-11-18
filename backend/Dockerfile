# Use the official Golang image as the base image
FROM golang:1.22 as builder

# Set the working directory
WORKDIR /app

# Copy the Go modules and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o collaborative-gis-app .

# Start a new stage to reduce the image size
FROM golang:1.19
WORKDIR /app
COPY --from=builder /app/collaborative-gis-app .

# Expose the port
EXPOSE 8000

# Run the executable
CMD ["./collaborative-gis-app"]