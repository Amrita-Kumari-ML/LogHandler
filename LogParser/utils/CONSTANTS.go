// Package utils defines constant values for configuration keys and default values 
// related to the server and database. These constants are used across the application
// to manage environment variables and configuration settings.
package utils

// Constants for environment variable keys related to the parser host, port, and URLs.
const KEY_HOST string = "PARSER_HOST"               // The key for the host of the parser service.
const KEY_PORT string = "PARSER_PORT"               // The key for the port on which the parser service runs.
const KEY_ALIVE_URL string = "PARSER_ALIVE_URL"     // The key for the URL that checks the parser service's health.
const KEY_GET_COUNT_URL string = "PARSER_GET_COUNT_URL"  // The key for the URL to get the log count.
const KEY_MAIN_URL string = "PARSER_MAIN_URL"       // The key for the main URL endpoint for logs.


// Constants for database configuration keys.
const KEY_DB_PORT string = "DB_PORT"                // The key for the database port.
const KEY_DB_HOST string = "DB_HOST"                // The key for the database host.
const KEY_DB_USERNAME string = "DB_USERNAME"        // The key for the database username.
const KEY_DB_PASSWORD string = "DB_PASSWORD"        // The key for the database password.
const KEY_DB_NAME string = "DB_NAME"                // The key for the database name.
const KEY_DB_SSLMODE string = "DB_SSLMODE"          // The key for the database SSL mode.

// Constants for database table and query keys.
const KEY_DB_TABLE_NAME string = "TABLE_NAME"       // The key for the database table name.
const KEY_DB_CREATE_TABLE_QUERY string = "CREATE_TABLE_QUERY"  // The key for the SQL query to create the logs table.


// Default values for the parser service configuration.
const PARSER_HOST string = "logparser"              // Default host for the parser service.
const PARSER_PORT string = ":8083"                  // Default port for the parser service.
const PARSER_ALIVE_URL string = "/"                 // Default URL for checking the parser service's health.
const PARSER_MAIN_URL string = "/logs"              // Default main URL for the logs endpoint.
const PARSER_GET_COUNT_URL string = "/logs/count"   // Default URL for retrieving the log count.


// Default values for the database connection configuration.
const DB_PORT string = "5432"                       // Default port for the PostgreSQL database.
const DB_HOST string = "postgres"                   // Default host for the PostgreSQL database.
const DB_USERNAME string = "postgres"               // Default username for the PostgreSQL database.
const DB_PASSWORD string = "123456"                 // Default password for the PostgreSQL database.
const DB_NAME string = "logsdb"                     // Default name for the PostgreSQL database.
const DB_SSLMODE string = "disable"                 // Default SSL mode for the PostgreSQL database connection.

// Default values for the database table name and table creation query.
const DB_TABLE_NAME string = "logs"                 // Default table name for storing logs in the database.
const DB_CREATE_TABLE_QUERY string = "CREATE TABLE IF NOT EXISTS logs (id SERIAL PRIMARY KEY, remote_addr VARCHAR(255), remote_user VARCHAR(255), time_local TIMESTAMP, request VARCHAR(255), status INT, body_bytes_sent INT, http_referer VARCHAR(255), http_user_agent VARCHAR(255), http_x_forwarded_for VARCHAR(255));"  // SQL query for creating the logs table if it doesn't exist.


// Constants for the HTTP request methods.
const REQUEST_GET_METHOD string = "GET"             // HTTP GET method.
const REQUEST_POST_METHOD string = "POST"           // HTTP POST method.
const REQUEST_DELETE_METHOD string = "DELETE"       // HTTP DELETE method.


// Constants for the names of API handlers.
const API_HANDLER_IS_ALIVE string = "isAlive"       // The handler for checking if the service is alive.
const API_HANDLER_ADD_LOGS string = "AddLogsHandler" // The handler for adding logs.
const API_HANDLER_GET_LOGS string = "GetLogsHandler" // The handler for retrieving logs.
const API_HANDLER_GET_COUNT_LOGS string = "GetLogsCountHandler" // The handler for retrieving the log count.
const API_HANDLER_DELETE_LOGS string = "DeleteLogsHandler" // The handler for deleting logs.


// Constants for configuration file names.
const CONFIG_FILE_NAME string = "config.yaml"        // The name of the main configuration file.
const CONFIG_DB_FILE_NAME string = "connection/dbConfig.yaml" // The name of the database connection configuration file.
