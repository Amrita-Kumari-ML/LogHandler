package main

import (
	"LogGenerator/helpers"
	"log"
)

// main is the entry point of the application, where the server and application setup is initialized and started.
//
// The main function performs the following tasks:
// 1. Initializes configuration and server objects to handle configurations and server logic.
// 2. Creates a new instance of the application using the configuration and server.
// 3. Attempts to set up the application by calling the SetUp() method.
// 4. If setup fails, logs the failure and exits.
// 5. Once setup is successful, the service is started, and the server is kept running until it is stopped.
//
// Example usage:
//   // Initialize the main entry point
//   main()
func main() {
	
	// Initialize a configuration object to handle application settings.
	// This will load configuration data for the application.
	conf := &helpers.Configs{}

	// Initialize the server object to handle the server setup and logic.
	// This will manage the server's lifecycle, including starting and stopping the server.
	server := &helpers.Servers{}

	// Create a new application instance with the server and configuration objects.
	// The NewApplication function will wire together the server and configuration.
	app := helpers.NewApplication(server, conf)

	// Call SetUp() to set up the application.
	// This function will configure the server, load environment variables, and start the application.
	// If there is an error during setup, it will be logged and the service will exit.
	if err := app.SetUp(); err != nil{ 
		// If there is a failure in setting up the application, log the error message.
		// This will help in debugging any issues during initialization.
		log.Println("Setup Failed! Some internal Issues")
	}

	// Log a message indicating that the service has been stopped.
	// This indicates that the server has either been shut down or encountered an issue.
	log.Println("Service stopped!")
}

