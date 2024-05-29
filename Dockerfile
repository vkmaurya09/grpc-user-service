# Use the official Golang image to create a build artifact.
FROM golang:1.20 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build main.go

# Expose port 50051 to the outside world
EXPOSE 50051

# Command to run the executable
CMD ["./main"]
