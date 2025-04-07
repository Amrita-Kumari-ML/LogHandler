// Package models defines the configuration structure used for database connection
// settings and logging configuration. This package contains a structure that holds
// the necessary configuration parameters for connecting to a database and managing logs.
package models

// DB_Config struct represents the configuration for the database connection and logs.
// It holds the details for connecting to a database and managing logs within that database.
type DB_Config struct {

	// Database struct holds the connection details for the database, including
	// the port, host, username, password, database name, and SSL mode.
	Database struct {
		// DBPort is the port number on which the database server is running.
		// This value is used to establish the connection to the database.
		DBPort string `yaml:"DB_PORT"`

		// DBHost is the hostname or IP address of the database server.
		// This value is used to connect to the database.
		DBHost string `yaml:"DB_HOST"`

		// DBUsername is the username used to authenticate to the database.
		// This is a required value for connecting securely to the database.
		DBUsername string `yaml:"DB_USERNAME"`

		// DBPassword is the password corresponding to the DBUsername for authentication.
		// It is required to ensure secure access to the database.
		DBPassword string `yaml:"DB_PASSWORD"`

		// DBName is the name of the database to connect to.
		// This value identifies which database within the database server to use.
		DBName string `yaml:"DB_NAME"`

		// DBSslMode determines the SSL mode for the connection.
		// This can be values like "disable", "require", "verify-full", etc., depending on 
		// the security requirements of the database server.
		DBSslMode string `yaml:"DB_SSLMODE"`
	} `yaml:"database"`

	// Logs struct defines the log table settings, including the table name and 
	// the SQL query to create the table if it does not exist.
	Logs struct {
		// TableName is the name of the table where logs will be stored.
		// This value is used to check the existence of the logs table and 
		// for performing operations like insert, update, delete, etc.
		TableName string `yaml:"table_name"`

		// CreateTableQuery holds the SQL query to create the log table.
		// This query will be executed if the log table does not already exist in the database.
		CreateTableQuery string `yaml:"create_table_query"`
	} `yaml:"logs"`
}
