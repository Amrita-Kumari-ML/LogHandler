// Package interfaces defines the abstractions for handling HTTP responses and mapping handler names 
// to corresponding HTTP handler functions. This package provides interfaces that can be used to 
// structure how the server responds to HTTP requests and how handlers are mapped dynamically.
package interfaces

import (
	"net/http"
)

// ResponseWrite interface defines the method for sending HTTP responses. It allows different 
// implementations to structure responses uniformly with status code, success flag, message, and data.
type ResponseWrite interface {
	// SendResponse sends a response to the HTTP client.
	// w: The HTTP response writer.
	// statusCode: The HTTP status code to return (e.g., 200, 404).
	// success: A boolean indicating whether the request was successful or not.
	// message: A string message providing additional information about the request result.
	// data: Any data to be included in the response (can be any type).
	SendResponse(w http.ResponseWriter, statusCode int, success bool, message string, data interface{})
}

// Handler interface defines a method for mapping handler names to their corresponding HTTP handler functions.
// This is useful for dynamic routing and flexible handler management.
type Handler interface {
	// MapHandler maps a handler name to a specific HTTP handler function.
	// handlerName: The name of the handler to map.
	// Returns: An HTTP handler function that corresponds to the provided handler name.
	MapHandler(handlerName string) http.HandlerFunc
}
