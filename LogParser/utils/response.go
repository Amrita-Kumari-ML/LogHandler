// Package utils contains utilities for handling HTTP responses in a structured way.
// The ResponseHandler provides a method for formatting and sending consistent responses
// to HTTP requests, including status, message, and data in JSON format.
package utils

import (
	"LogParser/models"
	"encoding/json"
	"net/http"
)

// ResponseHandler is a struct used to handle HTTP responses.
// It provides methods to send a standardized JSON response.
type ResponseHandler struct{}

// SendResponse sends an HTTP response with a JSON body containing a status, message, and data.
// The response is structured according to the `models.Response` format.
// Parameters:
//   - w: The HTTP ResponseWriter to send the response to.
//   - statusCode: The HTTP status code to return, e.g., 200, 404, etc.
//   - success: A boolean indicating whether the request was successful or not.
//   - message: A string message providing additional information about the result.
//   - data: The actual data to include in the response body, if any. Can be any type.
func (r *ResponseHandler) SendResponse(w http.ResponseWriter, statusCode int, success bool, message string, data interface{}) {
	// If data is provided, marshal it into JSON. If data is nil, it will be omitted from the response.
	var jsonData json.RawMessage
	if data != nil {
		var err error
		jsonData, err = json.Marshal(data)
		if err != nil {
			// If there is an error marshalling the data, return an internal server error response.
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// Create the response struct with status, message, and data.
	resp := models.Response{
		Status:  success,
		Message: message,
		Data:    jsonData,
	}

	// Set the content type to "application/json" for the response.
	w.Header().Set("Content-Type", "application/json")
	// Write the response status code to the ResponseWriter.
	w.WriteHeader(statusCode)
	// Encode the response struct into JSON and send it as the response body.
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		// If there is an error encoding the response, return an internal server error response.
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
