// Package models defines the data structure used for storing configuration settings.
package models

// Config struct represents the configuration settings of the application.
// It holds the necessary configuration data, such as the server port.
type Config struct {
	// PORT holds the port number on which the server should listen.
	// It is fetched from a YAML configuration file and passed as a string.
	// Example: "8080"
	PORT string `yaml:"PORT"`
}
