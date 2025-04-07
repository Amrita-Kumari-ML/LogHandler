package models

import (
	"encoding/json"
)

// Response represents a standardized structure for HTTP responses
// typically used to send back results of an operation (e.g., a log generation request).
// It is structured in a way that it includes a status, a message, and optional data
// that is encoded in JSON format.
//
// Fields:
//   - Status: A boolean indicating whether the operation was successful or not.
//     - `true` means the operation was successful.
//     - `false` indicates an error or failure in the operation.
//
//   - Message: A string message that provides additional information about the response.
//     It can describe the status or give further details on what happened, such as error descriptions.
//
//   - Data: This field holds the actual data associated with the response, if any.
//     It is represented as `json.RawMessage` to handle arbitrary JSON data. It could be 
//     the result of a log generation process, for example, a list of logs, or `null` 
//     if no data is included.
//
// Example usage:
//   // Successful response with data
//   response := models.Response{
//       Status:  true,
//       Message: "Logs generated successfully",
//       Data:    json.RawMessage(`[{"log": "data"}]`),
//   }
//
//   // Failed response with an error message
//   response := models.Response{
//       Status:  false,
//       Message: "Failed to generate logs",
//       Data:    nil,
//   }

type Response struct {
	// Status indicates whether the operation was successful or not.
	// - `true` means success
	// - `false` means failure
	Status bool `json:"status"`

	// Message provides a textual description of the response.
	// It could be a success message or an error message depending on the operation's outcome.
	Message string `json:"message"`

	// Data holds the actual data associated with the response, encoded as JSON.
	// It is represented as `json.RawMessage` to allow flexibility in handling different types of data.
	// It could be the result of an operation or `null` if no data is provided.
	Data json.RawMessage `json:"data"`
}