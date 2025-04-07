package helpers

import (
	"LogGenerator/loggenerator"
	"LogGenerator/server"
	"LogGenerator/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

)


// ServerLoader defines the interface for starting and stopping the server.
type ServerLoader interface{

	// startServer starts the HTTP server to handle incoming requests.
	startServer() error

	// stopServer stops the running HTTP server gracefully.
	stopServer() error
}

// ConfigurationLoader defines the interface for refreshing and updating the server configuration.
type ConfigurationLoader interface {
	// refreshServer refreshes the server's configuration by reloading the settings.
	refreshServer() error
}

// Servers struct handles server-related operations like starting and stopping the server.
type Servers struct{}

// startServer starts the log generator HTTP server and listens for incoming requests on the configured port.
//
// The function sets up HTTP handlers for different endpoints like "IsAlive" and "GenerateUrl".
// These handlers define the behavior for specific URLs that clients can hit to interact with the server.
//
// It starts an HTTP server on the configured port, and if the server fails to start, it logs the error and exits the application.
//
// Example usage:
//   // Initialize and start the server
//   server := &Servers{}
//   server.startServer()
func (s *Servers) startServer() error{
	// Initialize the ServerHandler, which handles the server responses and log generation.
	serv := &server.ServerHandler{
		ResponseW: &utils.ResponseHandler{},
		LogGen: &loggenerator.Generator{},
	}
	//server start logic
	http.HandleFunc(utils.GloablMetaData.IsAliveUrl, serv.IsAlive)
	http.HandleFunc(utils.GloablMetaData.GenerateUrl, serv.LogHandler)
	http.HandleFunc("/metrics", serv.MetricsHandler)
	//http.HandleFunc("/gen", serv.LogTestHandler)

	log.Println("Starting log generator server on port ",utils.GloablMetaData.Port,"...")
	log.Println(utils.ConfigData)
	if err := http.ListenAndServe(utils.GloablMetaData.Port, nil); err != nil {
		log.Println("Error starting server:", err)
		os.Exit(1)
	}
	return nil
}

// stopServer stops the HTTP server gracefully. It listens for signals to shut down the server.
//
// The function waits for the "done" channel to receive a signal, indicating that the server should stop.
//
// Example usage:
//   // Initialize and stop the server
//   server := &Servers{}
//   server.stopServer()
func (s *Servers) stopServer() error{
	//server stop logic
	<-done
	log.Println("Server Stopped......")
	os.Exit(1)
	return nil
}

// Configs struct handles the configuration-related operations, like refreshing the configuration periodically.
type Configs struct{}

// refreshServer refreshes the server's configuration. It loads the configuration settings anew from the source (e.g., YAML file or environment variables).
//
// If there is an error loading the configuration, it returns an error message.
//
// Example usage:
//   // Initialize and refresh the server configuration
//   configs := &Configs{}
//   configs.refreshServer()
func (c *Configs) refreshServer() error{
	if err := utils.FirstLoad(); err != nil {
		return fmt.Errorf("error loading configuration: %v", err)
	}
	log.Println(utils.ConfigData)
	log.Println("Configuration Updated")
	return nil
}

// RefreshConfigura refreshes the server configuration at regular intervals (every "t" duration).
//
// The function sets up a ticker that triggers every "t" duration and calls the refreshServer method
// to reload the configuration periodically. This ensures that the application is always using the latest configuration.
//
// Example usage:
//   // Initialize configuration and refresh it every 1 minute.
//   configs := &Configs{}
//   RefreshConfigura(configs, time.Minute)
func RefreshConfigura(configs ConfigurationLoader, t time.Duration){
	ticker := time.NewTicker(1 * t)
	defer ticker.Stop()

	for range ticker.C {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		if err := configs.refreshServer(); err != nil{
			log.Println(err)
		}
	}
}

// Application struct combines server and configuration functionalities to initialize and run the application.
type Application struct{
	server ServerLoader
	configuration ConfigurationLoader
}

// NewApplication creates a new Application instance with the given server and configuration loaders.
//
// The function returns a pointer to the created Application object, which will be used to start the server and handle configuration updates.
//
// Example usage:
//   // Create a new application with server and configuration loaders
//   app := NewApplication(&Servers{}, &Configs{})
//   app.SetUp() // Set up the application
func NewApplication(servers ServerLoader, configs ConfigurationLoader) *Application{
	return &Application{
		server: servers,
		configuration: configs,
	}
}
var done chan bool

// SetUp sets up the application environment, loading the configuration and starting the server.
//
// It starts a background goroutine that listens for system signals (e.g., SIGINT or SIGTERM).
// When such a signal is received, the server is gracefully stopped, and the application exits.
//
// Additionally, it loads the configuration for the first time and sets up a periodic refresh
// of the configuration every minute.
//
// Example usage:
//   // Set up the application
//   app := NewApplication(&Servers{}, &Configs{})
//   app.SetUp()
func (app *Application) SetUp() error{
//todo updatye
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done = make(chan bool, 1)

	go func() {
		sig := <-sigs
		log.Println()
		log.Println(sig)
		done <- true
	}()
	
	if err := app.configuration.refreshServer(); err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
    	log.Println(err)
		return nil
	}

	go RefreshConfigura(app.configuration,time.Minute)
	go app.server.stopServer()
	app.server.startServer()

	return nil
}