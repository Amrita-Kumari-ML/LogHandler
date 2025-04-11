package models

// AllConfigModel represents the configuration settings for the log generation application.
// It holds various configuration parameters for both the current service and the parser service,
// as well as settings related to log generation, such as rate and unit. This struct is designed to 
// be loaded from a YAML file and used throughout the application for consistent configuration access.
//
// Fields:
//   - KEY_RATE: This represents the rate at which logs should be generated. 
//     The value is typically an integer that defines how many logs to generate per unit of time.
//     Example: `100` logs per minute.
//
//   - KEY_UNIT: This specifies the time unit for the rate, indicating how often the logs should 
//     be generated based on the rate. It is typically a string and can take values like `"second"`, 
//     `"minute"`, `"hour"`, etc.
//     Example: `"minute"`
//
//   - CurrentService: This struct holds the configuration details for the current service (the log generator service).
//     - KEY_START_URL: The URL endpoint used to start the log generation process.
//       Example: `"/start"`
//     - KEY_ALIVE_URL: The URL endpoint used to check if the service is alive and running.
//       Example: `"/alive"`
//     - KEY_PORT: The port number where the current service will run (e.g., `8080`).
//
//   - ParserService: This struct holds the configuration details for the parser service, which processes the logs.
//     - KEY_PARSER_API: The API endpoint where the generated logs are sent for parsing and processing.
//       Example: `"http://localhost:5000/processLogs"`

type AllConfigModel struct {
	// KEY_RATE defines the rate of log generation, indicating how many logs to generate.
	// This rate is used in conjunction with KEY_UNIT to determine the frequency of log generation.
	KEY_RATE int `yaml:"KEY_RATE"`

	// KEY_UNIT specifies the time unit for the rate of log generation.
	// Common units include "second", "minute", "hour", etc.
	KEY_UNIT string `yaml:"KEY_UNIT"`

	// CurrentService holds the configuration for the log generation service.
	// This includes the URL endpoints and port number where the service is running.
	CurrentService struct {
		// KEY_START_URL is the endpoint that triggers the start of log generation.
		KEY_START_URL string `yaml:"KEY_START_URL"`

		// KEY_ALIVE_URL is the endpoint used to check if the service is alive and running.
		KEY_ALIVE_URL string `yaml:"KEY_ALIVE_URL"`

		// KEY_PORT is the port number where the log generator service listens.
		KEY_PORT string `yaml:"KEY_PORT"`
	} `yaml:"currentService"`

	// ParserService holds the configuration for the log parser service.
	// It includes the API endpoint for processing logs.
	ParserService struct {
		// KEY_PARSER_API is the API endpoint where the generated logs are sent for parsing and processing.
		KEY_PARSER_API string `yaml:"KEY_PARSER_API"`
	} `yaml:"parserService"`
}
