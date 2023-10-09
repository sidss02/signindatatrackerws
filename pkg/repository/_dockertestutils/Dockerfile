# Use the official Golang image as the base image and name the build stage as 'build'
FROM golang:1.17 AS build

# Set the working directory
WORKDIR /app

# Copy the Go application files to the container
COPY . .

# Build the Go application
RUN go mod download

RUN go build -o main .

# Start a new stage using debian:bullseye-slim as the base image
FROM debian:bullseye-slim

WORKDIR /app

# Create a non-root user and switch to that user
RUN useradd -u 10001 myuser
USER myuser

# Copy the built Go binary from the 'build' stage
COPY --from=build /app/main /app/main

# Expose port 8000
EXPOSE 8000

# Command to run the Test Go application
CMD ["/app/main"]