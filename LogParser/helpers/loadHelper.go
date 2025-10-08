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
	fmt.Println("Starting log generator server on port", utils.ConfigData.PORT)
		
	http.HandleFunc(utils.PARSER_ALIVE_URL, handlers.IsAlive)            // Handler for /alive
	http.HandleFunc(utils.PARSER_MAIN_URL, handlers.HandleType)          // Handler for /parse
	http.HandleFunc(utils.PARSER_GET_COUNT_URL, handlers.GetLogsCountHandler) // Handler for /logs/count

	// Statistics endpoints
	http.HandleFunc("/stats/status", handlers.GetStatusStatsHandler)     // Handler for /stats/status
	http.HandleFunc("/stats/ip", handlers.GetIPStatsHandler)             // Handler for /stats/ip
	http.HandleFunc("/stats/time", handlers.GetTimeStatsHandler)         // Handler for /stats/time
	http.HandleFunc("/stats/dashboard", handlers.GetDashboardStatsHandler) // Handler for /stats/dashboard

	// ML/AI endpoints
	http.HandleFunc("/ml/insights", handlers.GetMLInsightsHandler)       // Handler for comprehensive ML insights
	http.HandleFunc("/ml/anomalies", handlers.GetAnomalyDetectionHandler) // Handler for anomaly detection
	http.HandleFunc("/ml/predictions", handlers.GetPredictionsHandler)   // Handler for traffic predictions
	http.HandleFunc("/ml/security", handlers.GetSecurityThreatsHandler)  // Handler for security threat analysis
	http.HandleFunc("/ml/clusters", handlers.GetUserClustersHandler)     // Handler for user behavior clustering
	http.HandleFunc("/ml/realtime-anomaly", handlers.GetRealTimeAnomalyHandler) // Handler for real-time anomaly detection
	http.HandleFunc("/ml/config", handlers.GetMLConfigHandler)           // Handler for ML configuration
	http.HandleFunc("/ml/config/update", handlers.UpdateMLConfigHandler) // Handler for updating ML configuration

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
	<-Done
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
	if err := utils.FirstLoad(); err != nil {
		return fmt.Errorf("error loading configuration: %v", err)
	}

	db := connection.InitDB()
	if db == nil {
		logger.LogDebug("Database not configured!")
	}
	
	if err := connection.FirstLoad(); err != nil {
		return fmt.Errorf("error loading Database configuration: %v", err)
	}

	fmt.Println(utils.ConfigData)
	logger.LogDebug("Configuration Updated!")
	return nil
}

// RefreshConfigura refreshes the server's configuration at regular intervals using a ticker.
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
type Application struct{
	server       ServerLoader     // ServerLoader interface instance to manage server lifecycle.
	configuration ConfigurationLoader // ConfigurationLoader interface instance to manage configuration updates.
}

// NewApplication creates a new Application instance, initializing it with the provided ServerLoader 
func NewApplication(servers ServerLoader, configs ConfigurationLoader) *Application{
	return &Application{
		server:       servers,
		configuration: configs,
	}
}

// done channel is used to signal the termination of the server when a shutdown signal is received.
var Done chan bool

// SetUp initializes and sets up the application. It starts the server, begins periodic config refresh 
func (app *Application) SetUp() error{
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	Done = make(chan bool, 1)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		Done <- true
	}()

	if err := app.configuration.refreshServer(); err != nil {
		//log.SetFlags(log.LstdFlags | log.Lshortfile)
    	logger.LogError(err)
		return nil
	}

	// Initialize ML service
	if err := handlers.InitializeMLService(); err != nil {
		logger.LogWarn(fmt.Sprintf("ML service initialization failed: %v", err))
		// Continue without ML features
	} else {
		logger.LogInfo("ML service initialized successfully")
	}

	go RefreshConfigura(app.configuration, time.Minute)
	go app.server.stopServer()
	app.server.startServer()

	return nil
}
