# Use golang image to build the app
FROM golang:1.24-alpine AS build

# Set the working directory inside the container
WORKDIR /app
#psql -U postgres -h 172.17.0.4 -d logsdb
# "172.17.0.2"

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Copy the config files to the container
COPY config.yaml ./config.yaml
COPY connection/dbConfig.yaml ./connection/dbConfig.yaml

# Build the Go binary
RUN go build -o logparser .

# Expose the port the app will run on
EXPOSE 8082

# Start the app
CMD ["./logparser"]
