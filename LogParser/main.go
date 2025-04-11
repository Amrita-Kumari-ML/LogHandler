// Package main is the entry point of the LogParser application.
// It initializes the necessary configurations, sets up the database connection,
// and starts the application server.
package main

import (
	"LogParser/connection"
	"LogParser/helpers"
	"LogParser/logger"
	_"log"
)

// main is the entry point for the application. It performs the following tasks:
// 1. Initializes the database connection.
// 2. Loads the configuration settings.
// 3. Sets up and starts the application server.
// 4. Logs the service start and failure messages appropriately.
func main() {
	// Initialize the database connection by calling the InitDB function.
	// This will set up the connection to the database for the application.
	logger.InitializeLogger("debug")

	// Example of using the logger in your code
	logger.LogInfo("Starting Log Parser service...")

	connection.InitDB()

	// Create a new instance of the Configs struct for handling configuration.
	// The Configs struct is expected to manage configuration loading and settings.
	conf := &helpers.Configs{}

	// Create a new instance of the Servers struct which is responsible for
	// managing server-specific settings and operations.
	server := &helpers.Servers{}

	// Create a new Application instance using the server and configuration instances.
	// The NewApplication function initializes the application with server and configuration settings.
	app := helpers.NewApplication(server, conf)

	// Call the SetUp method of the Application instance to configure the server
	// and load any necessary configurations or settings.
	// If an error occurs during setup, log the error and indicate a failure.
	if err := app.SetUp(); err != nil {
		logger.LogError("Setup Failed! Some internal Issues")
	}

	// If the application setup succeeds, log that the service has stopped.
	logger.LogInfo("Service stopped!")
}
