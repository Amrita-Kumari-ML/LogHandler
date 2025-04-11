package interfaces

// ServerLoader defines the interface for starting and stopping the server.
type ServerLoader interface{

	// startServer starts the HTTP server to handle incoming requests.
	StartServer() error

	// stopServer stops the running HTTP server gracefully.
	StopServer() error
}

// ConfigurationLoader defines the interface for refreshing and updating the server configuration.
type ConfigurationLoader interface {
	// refreshServer refreshes the server's configuration by reloading the settings.
	RefreshServer() error
}