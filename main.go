package main

import (
	"log"     // for logging errors and other information
	"strconv" // for converting integers to strings
	"time"    // for working with time and dates

	// for building the API server
	"github.com/gofiber/fiber/v3" // Using Fiber v3

	// for validating API requests
	//"github.com/go-playground/validator/v10"

	// for working with databases
	//"github.com/jmoiron/sqlx"
	// sqlx is an extension of the standard
	// database/sql package that provides
	// additional features for working with databases in Go

	// for working with PostgreSQL databases
	//"github.com/lib/pq"
	// pq is a pure Go Postgres driver for the
	// database/sql package

	//load .env file in development
	//"github.com/joho/godotenv"
	// godotenv is a Go port of the Ruby dotenv library,
	// which loads environment variables from a
	// .env file into the process's environment

	// local packages
	"gobasics/config"     // for loading configuration from environment variables
	"gobasics/db"         // for connecting to the database
	"gobasics/handlers"   // for handling API requests
	"gobasics/repository" // for interacting with the database through a repository pattern
)

//---------------------------------------------------------------
// MIDDLEWARE
//---------------------------------------------------------------

// logger is custom logging middleware

// How middleware works in Fiber:
// - A middleware is just a handler that calls c.Next() to pass control
//   to the next handler in the chain.
// - Code Before c.Next() runs when the request comes in
// - Code after c.Next() runs after the response has been written

// - This 'sandwich' pattern lets us measure latency, log status codes,
//  capture errors, and more, without touching any route handlers.

// Lifecycle of a request in Fiber with middleware:
// 1. Request comes in and hits the first middleware (e.g., logger)
// 2. Logger records start time + logs
// 3. Logger calls c.Next() to pass control to the next handler (e.g., getAlbums)
// 4. Route Handler processes the request and sends a response
// 5. Control returns back to logger after c.Next()
// 6. Logger calculates latency, logs status code, and finishes

// .Handler() is used to create a custom middleware function in Fiber.
func Logger() fiber.Handler {
	return func(c fiber.Ctx) error {
		//--------------Before REST handler-------------------------
		start := time.Now() // Record the start time of the request

		// c.Method() returns the HTTP method of the request (e.g., GET, POST)
		// c.Path() returns the path of the request
		// (e.g., /albums , /albums/:id)
		// c.IP() returns the client's IP address, for audit logs
		log.Printf("Started %s %s for %s", c.Method(), c.Path(), c.IP())

		// Call the next handler in the chain
		// If handler returns an error, we capture it below
		err := c.Next() // This will execute the actual route handler

		//--------------After REST handler-------------------------
		// Calculate the latency of the request
		latency := time.Since(start)

		// c.Response().StatusCode() returns the
		// HTTP status code of the response
		// which is only meaningful after the route handler has executed
		statusCode := c.Response().StatusCode()

		//Color code the log output based on status code
		// for better visibility
		statusLabel := colorStatus(statusCode)

		log.Printf("Completed %s %s for %s with status %s and latency %v",
			c.Method(), c.Path(), c.IP(),
			statusLabel, latency.Round(time.Microsecond))

		// Return any error from downstream handlers to be handled by Fiber's error handling mechanism
		// If err is not nil, Fiber will handle it according to its error handling rules
		// (e.g., returning a 500 Internal Server Error if the error is not handled)
		return err
	}

}

// colorStatus is a helper function to color code status codes in logs
func colorStatus(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "\033[32m" + strconv.Itoa(statusCode) + "\033[0m" // Green
	case statusCode >= 400 && statusCode < 500:
		return "\033[33m" + strconv.Itoa(statusCode) + "\033[0m" // Yellow
	case statusCode >= 500:
		return "\033[31m" + strconv.Itoa(statusCode) + "\033[0m" // Red
	default:
		return strconv.Itoa(statusCode)
	}
}

//---------------------------------------------------------------
// Application entry point
//---------------------------------------------------------------

// This is the only place in the application where
// concrete types are created and dependencies are wired together.
// Everything else depends on abstractions (interfaces),
// which promotes loose coupling and makes the
// code easier to test and maintain.
// Dependency Injection in action
func main() {

	// 1. Load configuration from environment variables
	cfg := config.Load()

	// 2. Open a connection pool to the database
	database := db.Connect(cfg.DatabaseURL)
	defer database.Close() // Ensure the database connection is closed when main exits

	// 3. Build the dependency graph (innermost -> outermost)
	// In a real application, you would instantiate your
	// repository and service layers here,
	// injecting the database connection as a dependency
	// repository depends on the database connection
	// handler depends on the repository
	// routes depend on the handler

	// Create the album repository, which will interact with the database
	albumRepo := repository.NewAlbumRepository(database)
	albumHandler := handlers.NewAlbumHandler(albumRepo)

	// 4. Initializes Fiber app
	app := fiber.New(fiber.Config{
		// Return JSON responses with indentation for better readability
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if ok := false; !ok {
				_ = e
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// 5. Register the custom logger middleware
	// Every request will go through this middleware,
	// allowing us to log details about the request and response
	// Middleware is executed in the order it is registered,
	// so it will wrap around all route handlers defined below
	app.Use(Logger())

	// 6. Routes
	app.Get("/albums", albumHandler.GetAllAlbums)
	app.Get("/albums/:id", albumHandler.GetAlbumByID)
	app.Post("/albums", albumHandler.CreateAlbum)
	app.Put("/albums/:id", albumHandler.UpdateAlbum)
	app.Delete("/albums/:id", albumHandler.DeleteAlbum)

	// 7. Start the server on port 3000
	log.Println("Starting server on port 3000...")
	// Start the server on port 3000
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
