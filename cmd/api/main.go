package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/High-la/greenlight/internal/data"
	_ "github.com/lib/pq"
)

// application version number
const version = "1.0.0"

// Define a config struct to hold all configuration settings for the app
// port number and
// name of the current operating envt for the application
// (development, staging, production, etc...)

// Add a db struct field to hold the config settings for our database connection
// pool. For now this only holds the DSN(Data source name), which we will read in
// from the command line flag.

// Add maxOpenConns, maxIdleConns and maxIdleTime fields to hold the configuration
// settings for the connection pool.
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}

	// Add a new limiter struct containing fields for the requests per second and burst
	// values, and a boolean field which we can use to enable/disable rate limiting
	//
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

// Define application struct to hold the dependencies for HTTP handlers, helpers,
// and middleware.

// Add a models field to hold our new Models struct.
type application struct {
	config config
	logger *slog.Logger
	models data.Models
}

func main() {

	// Initialize a new structured logger which writes log entries to the standard out
	// stream.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// 1. Load env first
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	if env != "production" {

		envFileName := fmt.Sprintf(".env.%s", env)
		fmt.Println(envFileName)
		if err := godotenv.Load(envFileName); err != nil {
			logger.Error("error loading env file", "file", envFileName)
			os.Exit(1)
		}
	}

	// 2. Declare config
	// Declare an instance of the config struct
	var cfg config

	// 3. Define all flags

	// Read the value of the port and env from command-line flags into config struct.
	// we default to using the port number 4000 and the environment "development" if no
	// corresponding flags are provided.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// Create command line flags to read the setting values into the config struct.
	// Notice that we use true as the default for the 'enabled' setting?
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	// Use the value of the GREENLIGHT_DB_DSN environment var as the default value
	// for our db-dsn command line flag.
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")

	// Read the connection pool settings from command-line flags into the config struct.
	// Notice that the default values we're using are the ones we discussed above?
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	// 4. parse ONCE
	flag.Parse()

	// Call the openDB() helper function(see below) to create the connection pool,
	// passing in the config struct. If this returns an error, we log it and exit the
	// application immediately.
	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// Defer a call to db.Close() so that the connection pool is closed before the
	// main() function exits.
	defer db.Close()

	// Also log a message to say that the connection pools has been successfully
	// established.
	logger.Info("database connection pool established")

	// Declare an instance of the application struct, containing th econfig struct and
	// the logger.

	// Use the data.NewModels() function to initialize a Model struct, passing in the
	// connection pool as a parameter.
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	// Call app.serve() to start the server
	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)

	}

}

// The openDB() function returns a sql.DB connection pool.
func openDB(cfg config) (*sql.DB, error) {

	// Use sql.Open() to crearte an empty connection pool, using the DSN from the config
	// struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open (in-use + idle) connections in the pool. Note that
	// passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	// Set the maximum number of idle connections in the pool. Again, passing a value
	// less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	// Set the maximum idle timeout for connections in the pool. Passing a duration less
	// than or equal to 0 will mean that connections are not closed due to their idle time.
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use PingContext() to establish a new connection to the database, passing in the
	// context we create above as a parameter. If the connection couldn't be
	// established successfully within 5 second deadline, then this will return an error.
	// If we get this error, or any other, we close the connection pool and
	// return the error.
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	// Return the sql.DB connection pool.
	return db, nil
}
