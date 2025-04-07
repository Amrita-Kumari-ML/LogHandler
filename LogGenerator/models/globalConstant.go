package models

// GlobalConstantvariables contains configuration constants used throughout the application.
// These constants define various essential parameters such as server port, API endpoints, and processor API URL.
// These values are typically loaded from configuration files (e.g., YAML, JSON) during application startup.
//
// Fields:
//   - Port: A string representing the port on which the server should run.
//     This is typically a port number (e.g., "8080") that the HTTP server will listen on for incoming requests.
//
//   - IsAliveUrl: A string representing the URL path to check if the service is alive and running.
//     This endpoint is typically used for health-checking purposes to monitor if the service is operational.
//
//   - GenerateUrl: A string representing the URL path used to trigger the log generation process.
//     When this endpoint is hit, the service will generate logs based on the defined parameters (e.g., rate, log format).
//
//   - ProcessorApi: A string representing the URL to the log processor API that the generated logs will be sent to.
//     The logs will be forwarded to this API for processing or further handling.
//
// Example YAML configuration (as an example of how these constants might be set in a config file):
//   KEY_PORT: "8080"
//   KEY_ALIVE_URL: "/"
//   KEY_START_URL: "/logs"
//   KEY_PARSER_API: "http://localhost:8082/logs"
//
type GlobalConstantvariables struct {
	Port        string `yaml:"KEY_PORT"`        // The port on which the application server listens for requests.
	IsAliveUrl  string `yaml:"KEY_ALIVE_URL"`    // The URL path for checking if the service is alive.
	GenerateUrl string `yaml:"KEY_START_URL"`    // The URL path to trigger log generation.
	ProcessorApi string `yaml:"KEY_PARSER_API"`   // The API endpoint to which logs are sent for processing.
}
