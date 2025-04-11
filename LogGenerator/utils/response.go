package utils

import (
	"LogGenerator/logger"
	"LogGenerator/models"
	"encoding/json"
	_ "log"
	"net/http"
)

// ResponseHandler is a utility struct that provides methods for
// sending standardized responses in the form of JSON
type ResponseHandler struct{}

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
func (r *ResponseHandler) SendResponse(w http.ResponseWriter, statusCode int, success bool, message string, data interface{}){
	var jsonData json.RawMessage
	if data != nil {
		var err error
		jsonData, err = json.Marshal(data)
		if err != nil {
			logger.LogError("Internal Server Error")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	resp := models.Response{
		Status:  success,
		Message: message,
		Data:    jsonData,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.LogError("Json decode failed!")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
