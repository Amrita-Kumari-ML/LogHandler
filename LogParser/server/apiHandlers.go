// Package server defines the server logic, including handling the HTTP requests and sending logs to a log processor.
package server

import (
	"LogParser/interfaces"
	"LogParser/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// count is a global variable to track the number of logs processed.
var count int

// ServerHandler struct represents the handler for processing requests and sending responses.
type ServerHandler struct {
	HandlerEndpoint interfaces.Handler     // The handler responsible for mapping HTTP requests to specific handler functions.
	ResponseW       interfaces.ResponseWrite // The response writer that formats and sends HTTP responses.
}

// HandlerMapper struct is responsible for mapping the handlers to the appropriate HTTP endpoint.
type HandlerMapper struct {
	ServerHandler *ServerHandler // Holds the reference to the ServerHandler to map handlers for different routes.
}

// NewHandlerMapper initializes and returns a new HandlerMapper instance.
func NewHandlerMapper(serverHandler *ServerHandler) *HandlerMapper {
	return &HandlerMapper{ServerHandler: serverHandler}
}

// sendLogToProcessor sends logs to a log processor service using an HTTP POST request.
func sendLogToProcessor(logs []models.Log) {
	// Marshal logs into JSON format
	logJson, err := json.Marshal(logs)
	if err != nil {
		log.Printf("Error marshalling log data: %v", err) // Log an error if marshalling fails
		return
	}
	
	fmt.Println("---------------------------------") // Placeholder for logic to store data to PostgreSQL
	// Send the logs to a log processing service
	resp, err := http.Post("http://localhost:8083/addlogs", "application/json", bytes.NewBuffer(logJson))
	if err != nil {
		log.Printf("Error sending log to LogProcessor: %v", err) // Log an error if sending fails
		return
	}
	defer resp.Body.Close()

	// Check the response status code and log success or failure
	if resp.StatusCode == http.StatusOK {
		log.Println("Log successfully sent to LogProcessor")
	} else {
		log.Printf("Failed to send log to LogProcessor, Status Code: %d", resp.StatusCode)
	}
}

// printToConsole prints the logs to the console in a formatted JSON structure.
func printToConsole(logs []models.Log) {
	// Marshal logs into a pretty-printed JSON format
	logJson, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		log.Printf("Error marshalling logs to JSON: %v", err) // Log an error if marshalling fails
		return
	}
	
	// Print the formatted JSON logs to the console
	fmt.Println(string(logJson))
}

