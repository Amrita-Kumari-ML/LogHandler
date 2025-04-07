// Package models defines the data structure and utility function used for sending
// structured JSON responses in HTTP API endpoints.
package models

import (
	"encoding/json"
	"net/http"
)

// Response struct is used to format the response sent to the client.
// It contains the status, message, and optional data related to the API response.
type Response struct {
	// Status indicates whether the operation was successful or not.
	// A value of `true` indicates success, while `false` indicates failure.
	Status bool `json:"status"`

	// Message provides additional information regarding the operation.
	// It can be a success message or an error message depending on the outcome.
	Message string `json:"message"`

	// Data contains the actual data returned by the operation, if any.
	// It is serialized as `json.RawMessage` to handle any type of data.
	// If no data is to be sent, this field can be `null` or omitted.
	Data json.RawMessage `json:"data"`
}

// SendResponse is a utility function used to send a structured JSON response to the client.
// It sets the correct HTTP status code, formats the response, and encodes it as JSON.
// If the `data` parameter is not `nil`, it will be included in the response body as JSON data.
// If an error occurs while encoding the response or marshaling data, an error message is sent to the client.
func SendResponse(w http.ResponseWriter, statusCode int, success bool, message string, data interface{}) {

	// If the data is not nil, attempt to marshal it into a JSON object.
	var jsonData json.RawMessage
	if data != nil {
		var err error
		// Marshal the data into JSON
		jsonData, err = json.Marshal(data)
		if err != nil {
			// If there is an error marshaling the data, return a 500 Internal Server Error.
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// Create a Response object that contains the status, message, and data.
	resp := Response{
		Status:  success,
		Message: message,
		Data:    jsonData,
	}

	// Set the response header to indicate that the response is in JSON format.
	w.Header().Set("Content-Type", "application/json")
	// Set the HTTP status code as passed in the function argument.
	w.WriteHeader(statusCode)

	// Encode the response struct into JSON and write it to the HTTP response.
	// If an error occurs while encoding, return a 500 Internal Server Error.
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
