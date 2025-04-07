// Package connection handles the loading of database configuration settings
// from environment variables or a configuration YAML file. It also sets up the
// global configuration needed for database connections in the LogParser application.
package connection

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"LogParser/models"
	"LogParser/utils"
	"gopkg.in/yaml.v2"
)

var ConfigData models.DB_Config // Global variable for storing the loaded database configuration

// FirstLoad initializes the configuration for the database connection by:
// 1. Loading values from environment variables if available.
// 2. Falling back to loading the configuration from a YAML file if environment variables are missing.
func FirstLoad() error {
	// Load database connection settings from environment variables or defaults
	dbPort := getEnvString(utils.KEY_DB_PORT, utils.DB_PORT)
	dbHost := getEnvString(utils.KEY_DB_HOST, utils.DB_HOST)
	dbUsername := getEnvString(utils.KEY_DB_USERNAME, utils.DB_USERNAME)
	dbPassword := getEnvString(utils.KEY_DB_PASSWORD, utils.DB_PASSWORD)
	dbName := getEnvString(utils.KEY_DB_NAME, utils.DB_NAME)
	dbSslMode := getEnvString(utils.KEY_DB_SSLMODE, utils.DB_SSLMODE)

	// Set the database configuration
	ConfigData.Database = struct {
		DBPort     string `yaml:"DB_PORT"`
		DBHost     string `yaml:"DB_HOST"`
		DBUsername string `yaml:"DB_USERNAME"`
		DBPassword string `yaml:"DB_PASSWORD"`
		DBName     string `yaml:"DB_NAME"`
		DBSslMode  string `yaml:"DB_SSLMODE"`
	}{
		DBPort:     dbPort,
		DBHost:     dbHost,
		DBUsername: dbUsername,
		DBPassword: dbPassword,
		DBName:     dbName,
		DBSslMode:  dbSslMode,
	}

	// Set the log table configuration
	ConfigData.Logs = struct {
		TableName       string `yaml:"table_name"`
		CreateTableQuery string `yaml:"create_table_query"`
	}{
		TableName:       getEnvString(utils.KEY_DB_TABLE_NAME, utils.DB_TABLE_NAME),
		CreateTableQuery: getEnvString(utils.KEY_DB_CREATE_TABLE_QUERY, utils.DB_CREATE_TABLE_QUERY),
	}

	// If essential environment variables are missing, fall back to loading from the YAML file
	if dbHost == utils.DB_HOST {
		log.Println("Using config.yaml values or default settings.")
		err := LoadConfigFromYaml(utils.CONFIG_DB_FILE_NAME)
		if err != nil {
			return fmt.Errorf("error loading config from YAML: %v", err)
		}
	}

	return nil
}

// LoadConfigFromYaml loads the configuration data from a specified YAML file.
// This function unmarshals the YAML file contents into the global `ConfigData` variable.
func LoadConfigFromYaml(filePath string) error {
	// Read the YAML file into memory
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading YAML file: %v", err)
	}

	// Unmarshal the YAML contents into the DB_Config struct
	var config models.DB_Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return fmt.Errorf("error unmarshalling YAML file: %v", err)
	}

	// Update global ConfigData with the contents from the YAML file
	ConfigData = config
	return nil
}

// getEnvString retrieves an environment variable as a string. If the variable is not set,
// it returns the specified default value.
func getEnvString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvInt retrieves an environment variable as an integer. If the variable is not set,
// it returns the specified default value. If the conversion fails, it logs an error and
// returns the default value.
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Error parsing environment variable %s, using default value %d\n", key, defaultValue)
		return defaultValue
	}
	return parsedValue
}
