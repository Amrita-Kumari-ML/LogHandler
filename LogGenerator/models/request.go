package models

// RequestPayload represents the request payload structure used to define the parameters
// required for log generation. It is typically used to specify how many logs to generate 
// and the time duration (in seconds) for generating those logs.
//
// Fields:
//   - NumLogs: The number of logs to be generated. This is an integer value representing 
//     the total number of logs that should be produced during the specified time period.
//     Example: `1000` logs.
//
//   - Unit: The time unit (in seconds) over which the logs should be generated. It specifies 
//     the duration for log generation. This field represents the length of time for which 
//     the `NumLogs` should be distributed for generation.
//     Example: `"60"` (i.e., generate `NumLogs` logs over 60 seconds).
//
// Example usage:
//   // Example of a RequestPayload struct in a log generation request
//   requestPayload := models.RequestPayload{
//       NumLogs: 1000,
//       Unit: "60", // Generate 1000 logs over 60 seconds
//   }
type RequestPayload struct{
	// NumLogs defines the total number of logs that should be generated.
	// This number will be distributed over the time period specified by the `Unit`.
	NumLogs int64 `json:"num_logs"`

	// Unit defines the time period in seconds over which the logs will be generated.
	// Example: "60" means logs will be generated in the span of 60 seconds.
	Unit string `json:"time"` // in seconds
}
