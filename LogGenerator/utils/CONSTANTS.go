package utils

// Constants representing configuration keys and values for the log generator
// These constants are used to configure the log generator server's host, port, and other parameters
// such as the rate of log generation, unit of time, and URLs for the generator's API and parser.

const (
	// KEY_HOST represents the environment variable key for the log generator's host address.
	// It should be set to the hostname of the machine running the log generator service.
	// Example: "GENERATOR_HOST=loggenerate"
	KEY_HOST string = "GENERATOR_HOST"

	// KEY_PORT represents the environment variable key for the log generator's port number.
	// It should be set to the port on which the log generator service listens.
	// Example: "GENERATOR_PORT=:8080"
	KEY_PORT string = "GENERATOR_PORT"

	// KEY_ALIVE_URL represents the environment variable key for the URL path to check if the log generator is alive.
	// Example: "GENERATOR_ALIVE_URL=/"
	KEY_ALIVE_URL string = "GENERATOR_ALIVE_URL"

	// KEY_START_URL represents the environment variable key for the URL path to start log generation.
	// Example: "GENERATOR_START_URL=/logs"
	KEY_START_URL string = "GENERATOR_START_URL"

	// KEY_PARSER_API represents the environment variable key for the API endpoint to send logs for parsing.
	// Example: "PARSER_API=http://localhost:8083/logs"
	KEY_PARSER_API string = "PARSER_API"

	// KEY_RATE represents the environment variable key for the rate of log generation.
	// This value indicates how many logs to generate per unit time (in seconds, minutes, or hours).
	// Example: "GENERATOR_RATE=10"
	KEY_RATE string = "GENERATOR_RATE"

	// KEY_UNIT represents the environment variable key for the unit of time used in log generation.
	// The valid values are "s" for seconds, "m" for minutes, and "h" for hours.
	// Example: "GENERATOR_UNIT=s"
	KEY_UNIT string = "GENERATOR_UNIT"
)

// Constants representing default values for the log generator configuration.
// These values are used if the respective environment variables are not set or if they are invalid.

const (
	// GENERATOR_HOST represents the default hostname for the log generator service.
	// The service is expected to run on this hostname.
	// Default value: "loggenerate"
	GENERATOR_HOST string = "loggenerate"

	// GENERATOR_PORT represents the default port for the log generator service.
	// Default value: ":8080"
	GENERATOR_PORT string = ":8080"

	// GENERATOR_ALIVE_URL represents the default URL for checking if the log generator is alive.
	// Default value: "/"
	GENERATOR_ALIVE_URL string = "/"

	// GENERATOR_START_URL represents the default URL for starting log generation.
	// Default value: "/logs"
	GENERATOR_START_URL string = "/logs"

	// PARSER_API represents the default API endpoint for sending logs to be parsed.
	// Default value: "http://localhost:8083/logs"
	PARSER_API string = "http://localhost:8083/logs"

	// GENERATOR_RATE represents the default rate of log generation in logs per unit time.
	// Default value: 10 logs per unit time
	GENERATOR_RATE int = 10

	// GENERATOR_UNIT represents the default unit of time for log generation.
	// Default value: "s" for seconds
	GENERATOR_UNIT string = "s"
)


const FILE_NAME string = "config.yaml"
