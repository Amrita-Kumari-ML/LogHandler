FROM golang:1.24-alpine AS build

WORKDIR /app

# Environment variables
#ENV GENERATOR_START_URL="/logs"
#ENV GENERATOR_ALIVE_URL="/"
#ENV GENERATOR_PORT=":8081"
#ENV PARSER_API="http://localhost:8082/addlogs"
#"http://localhost:8082/logs"
#ENV RATE=10
#ENV UNIT="s"

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go binary
RUN go build -o loggenerate .

# Give execution permissions to the binary
RUN chmod +x loggenerate

# Copy the config.yaml if it's present in your project directory
COPY config.yaml ./config.yaml

# Expose the port the app will run on
EXPOSE 8081

# Start the app
CMD ["./loggenerate"]

