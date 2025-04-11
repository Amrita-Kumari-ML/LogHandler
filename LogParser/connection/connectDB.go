// Package connection manages the connection to the database, including initialization,
// database pinging, and ensuring the necessary database tables (e.g., logs table) exist.
package connection

import (
	"LogParser/logger"
	"LogParser/models"
	"LogParser/utils"
	"database/sql"
	"fmt"
	_ "log"

	_ "github.com/lib/pq" // Importing the Postgres driver
)

var DB *sql.DB            // Global variable holding the database connection
var Config *models.DB_Config // Global variable holding the configuration data for the database

// InitDB initializes the database connection using the configuration data.
// It first loads the configuration, then attempts to connect to the database
// using the provided credentials and connection details. If the connection is successful,
// it checks the database connection with a ping and ensures the necessary logs table exists.
func InitDB() *sql.DB {
	// Load configuration settings
	err1 := FirstLoad()
	if err1 != nil {
		logger.LogError("Configuration not loaded. Exiting...\n")
		return nil
	}

	// Use the global ConfigData loaded from configuration
	Config = &ConfigData
	var err error

	// Database connection string using values from the loaded config
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%s",
		Config.Database.DBUsername,
		Config.Database.DBPassword,
		Config.Database.DBName,
		Config.Database.DBSslMode,
		Config.Database.DBHost,
		Config.Database.DBPort,
	)

	// Open the database connection
	DB, err = sql.Open(utils.DB_USERNAME, connStr)
	if err != nil {
		logger.LogError(fmt.Sprintf("Error connecting to the database: %v\n", err))
	}

	// Check if the connection to the database is successful
	PingDB()

	// Ensure the logs table exists, if not, create it
	createLogsTableIfNotExist(*Config)

	return DB
}

// PingDB checks the database connection by attempting to ping it.
// It returns a boolean indicating if the connection is successful or not,
// and the database connection object.
func PingDB() (bool, *sql.DB) {
	if DB == nil {
		logger.LogError("Database connection is nil.")
		return false, nil
	}

	// Ping the database to check if it's reachable
	err := DB.Ping()
	if err != nil {
		logger.LogError(fmt.Sprintf("Error pinging the database: %v\n", err))
		return false, nil
	}

	logger.LogInfo("Successfully connected to the database!")
	return true, DB
}

// createLogsTableIfNotExist ensures that the logs table exists in the database.
// If the table doesn't exist, it creates the table using the SQL query provided in the config.
func createLogsTableIfNotExist(config models.DB_Config) {
	var tableName string
	// Check if the logs table exists in the database
	err := DB.QueryRow(`SELECT table_name FROM information_schema.tables WHERE table_name = $1`, config.Logs.TableName).Scan(&tableName)
	if err == sql.ErrNoRows {
		// Table doesn't exist, so create it
		logger.LogDebug("Logs table doesn't exist, creating it...")
		_, err = DB.Exec(config.Logs.CreateTableQuery)
		if err != nil {
			logger.LogError(fmt.Sprintf("Error creating the logs table: %v\n", err))
		}
		indexExists("idx_time_local")
		logger.LogDebug("Logs table created successfully!")
	} else if err != nil {
		logger.LogDebug(fmt.Sprintf("Error checking if logs table exists: %v\n", err))
	} else {
		logger.LogDebug("Logs table already exists.")
	}
}

func indexExists(indexName string) bool {
    var index string
    err := DB.QueryRow(`SELECT indexname FROM pg_indexes WHERE indexname = $1`, indexName).Scan(&index)
    if err == sql.ErrNoRows {
        // Index does not exist
        return false
    } else if err != nil {
        logger.LogDebug(fmt.Sprintf("Error checking if index exists: %v\n", err))
    }
    // Index exists
    return true
}


