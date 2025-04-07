// Package utils implements a simple utility package
// It consists of fetching and storing configuration data,
// either from environment variables if present or from configuration file
package utils

import (
	"fmt"
	"os"
	"log"
	"strconv"
	"LogGenerator/models"
	"github.com/go-yaml/yaml"
)

var ConfigData models.AllConfigModel
var RateData models.RequestPayload

var GloablMetaData models.GlobalConstantvariables

// FirstLoad handles the creation and updating of configuration data.
// It loads global data from environment variables, and if they are not set,
// it loads the data from a configuration file (config.yaml).
// If any configuration is missing or invalid, it will fall back to defaults.
func FirstLoad() (error){
	// Load values from environment variables or use default values
	port := getEnvString(KEY_PORT, GENERATOR_PORT)
	alive_url := getEnvString(KEY_ALIVE_URL, GENERATOR_ALIVE_URL)
	generate_url := getEnvString(KEY_START_URL, GENERATOR_START_URL)
	parser_api := getEnvString(KEY_PARSER_API, PARSER_API)

	// Initialize GlobalMetaData with values
	GloablMetaData = models.GlobalConstantvariables{
		Port:        port,
		IsAliveUrl:  alive_url,
		GenerateUrl: generate_url,
		ProcessorApi:parser_api,
	}

	RateData = models.RequestPayload{
		NumLogs : int64(getEnvInt(KEY_RATE, GENERATOR_RATE)),
		Unit: getEnvString(KEY_UNIT, GENERATOR_UNIT),
	}

	// If any essential environment variable is missing, fall back to loading from config.yaml
	if port == GENERATOR_PORT {
		log.Println("Using config.yaml values or default settings.")
		err := LoadConfigFromYaml()
		if err != nil {
			return fmt.Errorf("error loading config from YAML: %v", err)
		}
	}

	return nil
}

// GetEnvString this function is reponsible for fetching
// string type environment variables anf if not present then 
// sets default value
func getEnvString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Println("Environment variable not set, using default value for'", key, "':", defaultValue)
		return defaultValue
	}
	return value
}

// getEnvInt this function is reponsible for fetching
// integer type environment variables anf if not present then 
// sets default value
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("\nError parsing int value for key '%s', defaulting to '%d'", key, defaultValue)
		return defaultValue
	}
	return parsedValue
}

// LoadConfigFromYaml is responsible for setting the data to global
// variables based on the configuration file
func LoadConfigFromYaml() error {
	fileData, err := os.ReadFile("config.yaml")
	if err != nil {
		return fmt.Errorf("failed to read config.yaml: %v", err)
	}

	if err := yaml.Unmarshal(fileData, &ConfigData); err != nil {
		return fmt.Errorf("failed to parse config.yaml: %v", err)
	}

	// Update global variables with data from config.yaml if necessary
	GloablMetaData.Port = ConfigData.CurrentService.KEY_PORT
	GloablMetaData.IsAliveUrl = ConfigData.CurrentService.KEY_ALIVE_URL
	GloablMetaData.GenerateUrl = ConfigData.CurrentService.KEY_START_URL
	GloablMetaData.ProcessorApi = ConfigData.ParserService.KEY_PARSER_API

	if RateData.NumLogs <= 0 {
		RateData.NumLogs = int64(ConfigData.KEY_RATE)
	}
	if !(RateData.Unit == "s" || RateData.Unit == "m" || RateData.Unit == "h") {
		RateData.Unit = ConfigData.KEY_UNIT
	}

	return nil
}

// ReloadRateData this functions reloads the data every time 
// when rate changes and sets the global rate data which
// consists of unit and rate as parameters
func ReloadRateData(rd models.RequestPayload) error{
	if (rd.NumLogs <= 0) || !(rd.Unit == "s" || rd.Unit == "m" || rd.Unit == "h"){
		return fmt.Errorf("invalid rate or unit found")
	}

	RateData.NumLogs = rd.NumLogs
	RateData.Unit = rd.Unit

	return nil
}