// Package utils provides utility functions for managing configuration settings 
// by loading them either from environment variables or from a YAML configuration file.
// It ensures that global configuration settings are loaded and provides functions
// for retrieving configuration values with fallback options.
package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"LogParser/models"
	"gopkg.in/yaml.v2"
)

var ConfigData models.Config // Global variable to hold the application configuration

// FirstLoad handles the creation and updating of configuration data. 
// It first attempts to load global configuration from environment variables. 
// If environment variables are not present, it falls back to loading configuration from a YAML file.
func FirstLoad() (error) {
	// Retrieve the server port from environment variables, falling back to the default value
	port := getEnvString(KEY_PORT, PARSER_PORT)

	// Set the global ConfigData object with the retrieved port value
	ConfigData = models.Config{
		PORT: port, 
	}

	// If the port is still set to the default value (meaning the environment variable was not set),
	// fall back to loading the configuration from the YAML file
	if port == PARSER_PORT {
		log.Println("Using config.yaml values or default settings.")

		// Attempt to load the YAML file
		if err := LoadConfigFromYaml(); err != nil {
			return fmt.Errorf("error loading config from YAML: %v", err)
		}
	}

	return nil
}

// LoadConfigFromYaml loads configuration data from a YAML file (config.yaml).
// This is called when essential environment variables are missing or have default values.
// It unmarshals the YAML data into the global ConfigData variable.
func LoadConfigFromYaml() error {
	// Read the YAML file
	yamlFile, err := os.ReadFile(CONFIG_FILE_NAME)
	if err != nil {
		log.Printf("error reading YAML file: %v\n", err)
		return fmt.Errorf("error reading YAML file: %v\n", err)
	}

	// Unmarshal the YAML content into ConfigData
	err = yaml.Unmarshal(yamlFile, &ConfigData)
	if err != nil {
		log.Printf("error unmarshalling YAML file: %v\n", err)
		return fmt.Errorf("error unmarshalling YAML file: %v", err)
	}

	return nil
}

/*
This commented-out function is an older approach to load configuration from environment variables directly,
but it is not used in the current implementation.

func LoadEnvironmentVariables() models.Config {
	// Load environment variables and set defaults if not provided
	envData := models.Config{
		Server: struct {
			Port string `yaml:"port"`
		} {
			Port: getEnvString("CURRENT_SERVICE_PORT", ":8083"),
		},
		API: struct {
			Endpoints []struct {
				Path    string `yaml:"path"`
				Method  string `yaml:"method"`
				Handler string `yaml:"handler"`
			} `yaml:"endpoints"`
		}{
			Endpoints: []struct {
				Path    string `yaml:"path"`
				Method  string `yaml:"method"`
				Handler string `yaml:"handler"`
			}{
				{
					Path:    getEnvString("API_PATH_ALIVE", "/"),
					Method:  getEnvString("API_METHOD_LOGS", "GET"),
					Handler: getEnvString("API_HANDLER_LOGS", "isAlive"),
				},
				{
					Path:    getEnvString("API_PATH_ADDLOGS", "/addlogs"),
					Method:  getEnvString("API_METHOD_ADDLOGS", "POST"),
					Handler: getEnvString("API_HANDLER_ADDLOGS", "AddLogsHandler"),
				},
				{
					Path:    getEnvString("API_PATH_GETLOGS", "/getlogs"),
					Method:  getEnvString("API_METHOD_GETLOGS", "GET"),
					Handler: getEnvString("API_HANDLER_GETLOGS", "GetLogsHandler"),
				},
				{
					Path:    getEnvString("API_PATH_GETLOGSCOUNT", "/getlogsCount"),
					Method:  getEnvString("API_METHOD_GETLOGSCOUNT", "GET"),
					Handler: getEnvString("API_HANDLER_GETLOGSCOUNT", "GetLogsCountHandler"),
				},
				{
					Path:    getEnvString("API_PATH_DELETELLOGS", "/deletelogs"),
					Method:  getEnvString("API_METHOD_DELETELOGS", "DELETE"),
					Handler: getEnvString("API_HANDLER_DELETELOGS", "DeleteLogsHandler"),
				},
			},
		},
	}
	return envData
}
*/

 // getEnvString retrieves a string value from an environment variable or returns a default value if the environment variable is not set.
func getEnvString(key string, defaultValue string) string {
	// Attempt to fetch the environment variable
	value := os.Getenv(key)
	// If the value is empty (environment variable not set), return the default value
	if value == "" {
		return defaultValue
	}
	// Return the value of the environment variable
	return value
}

// getEnvInt retrieves an integer value from an environment variable or returns a default value if the environment variable is not set.
// It also handles any errors that occur during the conversion from string to int.
func getEnvInt(key string, defaultValue int) int {
	// Attempt to fetch the environment variable
	value := os.Getenv(key)
	// If the value is empty (environment variable not set), return the default value
	if value == "" {
		return defaultValue
	}

	// Attempt to parse the value as an integer
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		// Log an error if the value cannot be converted to an integer
		log.Printf("Error parsing int value for key %s, defaulting to %d", key, defaultValue)
		return defaultValue
	}
	// Return the parsed integer value
	return parsedValue
}
