package interfaces

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// ResponseWrite defines an interface for handling HTTP responses with a standardized structure in JSON format.
type ResponseWrite interface {
	
	// SendResponse sends a standardized HTTP response in JSON format. The response includes
	// the status, message, and data. It is structured according to the models.Response format.
	//
	// Parameters:
	//   - w: The http.ResponseWriter used to write the response.
	//   - statusCode: The HTTP status code (e.g., http.StatusOK, http.StatusBadRequest).
	//   - success: A boolean indicating whether the operation was successful or not.
	//   - message: A string message that provides additional information about the response.
	//   - data: An interface{} that contains the actual data to be sent in the response (e.g., a user object, a list of records, etc.).
	//
	// If `data` is not nil, it will be marshaled into a JSON format and included in the response.
	// If `data` is nil, no data field will be included in the response.
	//
	// This method automatically sets the Content-Type to "application/json" and writes the provided
	// statusCode to the response header. In case of any issues with marshaling or writing the response,
	// appropriate error messages will be logged and a generic internal server error (HTTP 500) will be returned.
	//
	// Example usage:
	//   // Initialize a ResponseHandler instance
	//   handler := utils.ResponseHandler{}
	//
	//   // Send a successful response with data
	//   handler.SendResponse(w, http.StatusOK, true, "Request successful", data)
	//
	//   // Send a failed response without data
	//   handler.SendResponse(w, http.StatusBadRequest, false, "Invalid input", nil)
	SendResponse(w http.ResponseWriter, statusCode int, success bool, message string, data interface{})
}

// LogGenerator defines an interface for generating logs concurrently.
type LogGenerator interface {
	
	// GenerateLogsConcurrently generates logs in parallel based on the specified rate and duration.
	//
	// Parameters:
	//   - ctx: A context object that can be used for cancelling or timing out the log generation.
	//   - rate: The rate at which logs should be generated (e.g., number of logs per second).
	//   - duration: The duration for which the log generation should occur (e.g., 5 minutes, 1 hour).
	//   - wg: A sync.WaitGroup that helps manage concurrent operations, ensuring that all log generation tasks complete before continuing.
	//
	// This method performs log generation concurrently using goroutines, ensuring that logs are generated efficiently
	// and that the application can continue processing other tasks without waiting for each log generation task to finish.
	//
	// Example usage:
	//   // Initialize a log generator instance
	//   logGen := loggenerator.Generator{}
	//
	//   // Create a WaitGroup to wait for all concurrent log generation tasks to complete
	//   var wg sync.WaitGroup
	//
	//   // Start generating logs concurrently with a rate of 10 logs per second for 5 minutes
	//   logGen.GenerateLogsConcurrently(ctx, 10, 5*time.Minute, &wg)
	GenerateLogsConcurrently(ctx context.Context, rate int, duration time.Duration, wg *sync.WaitGroup, statusChan chan<- string)
}
