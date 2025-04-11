// Package helpers manages server lifecycle, configuration loading, and request handling.
// It defines the core components for starting and stopping the server, refreshing configurations,
// and mapping handler functions to the appropriate URLs. The application also handles graceful
// shutdowns upon receiving termination signals.
package helpers

import (
	"LogParser/connection"
	"LogParser/handlers"
	_"LogParser/interfaces"
	"LogParser/logger"
	_ "LogParser/server"
	"LogParser/utils"
	"fmt"
	_ "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ServerLoader interface defines methods for starting and stopping the server.
type ServerLoader interface{
	// startServer starts the server and listens on the specified port.
	startServer() error

	// stopServer stops the server gracefully when triggered (e.g., during shutdown).
	stopServer() error
}

// ConfigurationLoader interface defines a method to refresh the server configuration.
type ConfigurationLoader interface{
	// refreshServer refreshes the server configuration by reloading environment variables 
	// and reloading database configurations.
	refreshServer() error
}

// Servers struct implements the ServerLoader interface. It contains methods for starting 
// and stopping the HTTP server. It is responsible for managing the server lifecycle.
type Servers struct{}

// EndPointHandler struct is used to map handler names (from the config) to corresponding HTTP 
// handler functions. It allows dynamic routing of requests based on handler names.
type EndPointHandler struct{}

// startServer starts the HTTP server, which listens for incoming requests on the port 
// defined in the configuration. The server handles requests for specific paths and endpoints.
func (s *Servers) startServer() error{
	// Logs the port number the server is starting on.
	fmt.Println("Starting log generator server on port", utils.ConfigData.PORT)
	
	// Register HTTP handler functions for various paths.
	http.HandleFunc(utils.PARSER_ALIVE_URL, handlers.IsAlive)            // Handler for /alive
	http.HandleFunc(utils.PARSER_MAIN_URL, handlers.HandleType)          // Handler for /parse
	http.HandleFunc(utils.PARSER_GET_COUNT_URL, handlers.GetLogsCountHandler) // Handler for /logs/count
	
	// Log the current configuration data (for debugging and verification).
	fmt.Println("Current Configuration Data:", utils.ConfigData)
	
	// Start the HTTP server and listen on the configured port.
	serverPort := utils.ConfigData.PORT
	if err := http.ListenAndServe(fmt.Sprintf("%s", serverPort), nil); err != nil {
		logger.LogError(fmt.Sprintf("Error starting server: %v", err))
		os.Exit(1)
	}

	return nil
}
/*
// MapHandlerToFunc maps a handler name to a corresponding HTTP handler function.
// This function is used to dynamically assign the correct handler based on configuration.
func MapHandlerToFunc(handlerName string, handler interfaces.Handler) http.HandlerFunc {
	return handler.MapHandler(handlerName)
}

// MapHandler maps handler names to corresponding HTTP handler functions.
// It returns the appropriate handler for each API endpoint (based on the handler name).
func (url *EndPointHandler) MapHandler(handlerName string) http.HandlerFunc{
	switch handlerName {
	case utils.API_HANDLER_IS_ALIVE: // Handler for the /alive endpoint
		return handlers.IsAlive
	case utils.API_HANDLER_ADD_LOGS: // Handler for the /logs/add endpoint
		return handlers.AddLogsHandler
	case utils.API_HANDLER_GET_LOGS: // Handler for the /logs endpoint
		return handlers.GetLogsHandler
	case utils.API_HANDLER_GET_COUNT_LOGS: // Handler for the /logs/count endpoint
		return handlers.GetLogsCountHandler
	case utils.API_HANDLER_DELETE_LOGS: // Handler for the /logs/delete endpoint
		return handlers.DeleteLogsHandler
	default:
		// If the handler is not recognized, return a 404 (Not Found).
		logger.LogWarn(fmt.Sprintf("No handler found for %s, returning a 404", handlerName))
		return http.NotFound
	}
}
*/
// stopServer gracefully shuts down the server when a termination signal is received.
func (s *Servers) stopServer() error{
	// Wait for a signal (e.g., SIGINT or SIGTERM) to stop the server.
	<-done
	fmt.Println("Server Stopped......")
	os.Exit(1)
	return nil
}

// Configs struct implements the ConfigurationLoader interface, which is responsible for 
// refreshing the configuration (including environment variables and database configurations).
type Configs struct{}

// refreshServer refreshes the configuration of the server by reloading the environment 
// variables and reloading the database configuration (through the connection package).
func (c *Configs) refreshServer() error {
	// Load configuration settings (e.g., from environment variables or config files).
	if err := utils.FirstLoad(); err != nil {
		return fmt.Errorf("error loading configuration: %v", err)
	}

	// Reload the database configuration.
	if err := connection.FirstLoad(); err != nil {
		return fmt.Errorf("error loading Database configuration: %v", err)
	}

	// Log the updated configuration data (for debugging and verification).
	fmt.Println(utils.ConfigData)
	logger.LogDebug("Configuration Updated!")
	return nil
}

// RefreshConfigura refreshes the server's configuration at regular intervals using a ticker.
// The ticker triggers configuration refresh every `t` duration.
func RefreshConfigura(configs ConfigurationLoader, t time.Duration){
	// Create a ticker to trigger configuration refresh at regular intervals.
	ticker := time.NewTicker(1 * t)
	defer ticker.Stop()

	// Continuously refresh the configuration as long as the ticker is active.
	for range ticker.C {
		//log.SetFlags(log.LstdFlags | log.Lshortfile)
		if err := configs.refreshServer(); err != nil{
			// Log any errors encountered while refreshing the configuration.
			logger.LogError(err)
		}
	}
}

// Application struct encapsulates the server and configuration loader, managing the application's 
// lifecycle, including starting the server, refreshing configurations, and handling graceful shutdowns.
type Application struct{
	server       ServerLoader     // ServerLoader interface instance to manage server lifecycle.
	configuration ConfigurationLoader // ConfigurationLoader interface instance to manage configuration updates.
}

// NewApplication creates a new Application instance, initializing it with the provided ServerLoader 
// and ConfigurationLoader implementations.
func NewApplication(servers ServerLoader, configs ConfigurationLoader) *Application{
	return &Application{
		server:       servers,
		configuration: configs,
	}
}

// done channel is used to signal the termination of the server when a shutdown signal is received.
var done chan bool

// SetUp initializes and sets up the application. It starts the server, begins periodic configuration 
// refreshes, and listens for termination signals to gracefully stop the server.
func (app *Application) SetUp() error{
	// Create a channel to capture OS signals (e.g., SIGINT, SIGTERM).
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Create a channel for notifying when the server should stop.
	done = make(chan bool, 1)

	// Goroutine to listen for termination signals (e.g., SIGINT or SIGTERM).
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	// Refresh the server configuration at startup.
	if err := app.configuration.refreshServer(); err != nil {
		//log.SetFlags(log.LstdFlags | log.Lshortfile)
    	logger.LogError(err)
		return nil
	}

	// Start refreshing the configuration periodically.
	go RefreshConfigura(app.configuration, time.Minute)
	// Start the server and listen for incoming requests.
	go app.server.stopServer()
	app.server.startServer()

	return nil
}
