package helpers

import (
	"LogGenerator/interfaces"
	"LogGenerator/loggenerator"
	"LogGenerator/logger"
	"LogGenerator/server"
	"LogGenerator/utils"
	"fmt"
	_ "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Servers struct responsible for start and stop of the server
type Servers struct{}

// StartServer is responsible for starting the server where it has listen and serve
// and the handlers are also aattached to handle the api end point
// Example usage:
//
//	// Initialize and start the server
//	server := &Servers{}
//	server.startServer()
func (s *Servers) StartServer() error {
	serv := &server.ServerHandler{
		ResponseW: &utils.ResponseHandler{},
		LogGen:    &loggenerator.Generator{},
	}
	http.HandleFunc(utils.GloablMetaData.IsAliveUrl, serv.IsAlive)
	http.HandleFunc(utils.GloablMetaData.GenerateUrl, serv.LogHandler)
	http.HandleFunc("/logs/stop", serv.StopHandler)
	http.HandleFunc("/logs/status", serv.StatusHandler)

	//http.HandleFunc("/gen", serv.LogTestHandler)

	logger.LogInfo("Starting log generator server on port " + utils.GloablMetaData.Port + "...")
	logger.LogDebug(utils.ConfigData)
	if err := http.ListenAndServe(utils.GloablMetaData.Port, nil); err != nil {
		logger.LogError(fmt.Sprintf("Error starting server: %v", err))
		os.Exit(1)
	}
	return nil
}

// StopServer stops the HTTP server gracefully. It listens for signals to shut down the server.
// Example usage:
//
//	// Initialize and stop the server
//	server := &Servers{}
//	server.stopServer()
func (s *Servers) StopServer() error {
	<-done
	logger.LogInfo("Server Stopped......")
	os.Exit(1)
	return nil
}

// Configs struct is reponsible for handling the refresh of the config data
type Configs struct{}

// RefreshServer refreshes the server's configuration. It loads the config data.
// Example usage:
//
//	// Initialize and refresh the server configuration
//	configs := &Configs{}
//	configs.refreshServer()
func (c *Configs) RefreshServer() error {
	if err := utils.FirstLoad(); err != nil {
		return fmt.Errorf("error loading configuration: %v", err)
	}
	logger.LogDebug(fmt.Sprintf("Updated Data : %v", utils.ConfigData))
	return nil
}

// RefreshConfigura calls the Refresh server periodically to refresh for the configuration
// Example usage:
//
//	// Initialize configuration and refresh it every 1 minute.
//	configs := &Configs{}
//	RefreshConfigura(configs, time.Minute)
func RefreshConfigura(configs interfaces.ConfigurationLoader, t time.Duration) {
	ticker := time.NewTicker(1 * t)
	defer ticker.Stop()

	for range ticker.C {
		if err := configs.RefreshServer(); err != nil {
			logger.LogError(err)
		}
	}
}

// Application stuct is used to bind the config and server structs together
type Application struct {
	Server        interfaces.ServerLoader
	Configuration interfaces.ConfigurationLoader
}

// NewApplication creates a new Application instance with the given server and configuration loaders.
// Example usage:
//
//	// Create a new application with server and configuration loaders
//	app := NewApplication(&Servers{}, &Configs{})
//	app.SetUp() // Set up the application
func NewApplication(servers interfaces.ServerLoader, configs interfaces.ConfigurationLoader) *Application {
	return &Application{
		Server:        servers,
		Configuration: configs,
	}
}

// done is the channel used to carry the stop signal of program shutdown
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
//
//	// Set up the application
//	app := NewApplication(&Servers{}, &Configs{})
//	app.SetUp()
func (app *Application) SetUp() error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done = make(chan bool, 1)

	go func() {
		sig := <-sigs
		logger.LogWarn(sig)
		done <- true
	}()

	if err := app.Configuration.RefreshServer(); err != nil {
		logger.LogError(err)
		return err
	}

	go RefreshConfigura(app.Configuration, time.Minute)
	go app.Server.StopServer()
	app.Server.StartServer()

	return nil
}
