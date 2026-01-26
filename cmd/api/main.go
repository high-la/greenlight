package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// application version number
const version = "1.0.0"

// Define a config struct to hold all configuration settings for the app
// port number and
// name of the current operating envt for the application
// (development, staging, production, etc...)
type config struct {
	port int
	env  string
}

// Define application struct to hold the dependencies for HTTP handlers, helpers,
// and middleware.
type application struct {
	config config
	logger *slog.Logger
}

func main() {

	// Declare an instance of the config struct
	var cfg config

	// Read the value of the port and env from command-line flags into config struct.
	// we default to using the port number 4000 and the environment "development" if no
	// corresponding flags are provided.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Initialize a new structured logger which writes log entries to the standard out
	// stream.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Declare an instance of the application struct, containing th econfig struct and
	// the logger.
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Use the httprouter instance returned by app.routes() as the server handler.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// Start the HTTP server.
	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
