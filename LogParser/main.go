// Package main is the entry point of the LogParser application.
// It initializes the necessary configurations, sets up the database connection,
// and starts the application server.
package main

import (
	_ "LogParser/connection"
	"LogParser/helpers"
	"LogParser/logger"
	_ "log"
)

// main is the entry point for the application. It performs the following tasks:
// 1. Initializes the database connection.
// 2. Loads the configuration settings.
// 3. Sets up and starts the application server.
// 4. Logs the service start and failure messages appropriately.
func main() {
	logger.InitLogger("debug")
	logger.LogInfo("Starting Log Parser service...")
	conf := &helpers.Configs{}
	server := &helpers.Servers{}
	app := helpers.NewApplication(server, conf)

	// Call the SetUp method of the Application instance to configure the server
	// and load any necessary configurations or settings.
	// If an error occurs during setup, log the error and indicate a failure.
	if err := app.SetUp(); err != nil {
		logger.LogError("Setup Failed! Some internal Issues")
	}
	logger.LogInfo("Service stopped!")
}
