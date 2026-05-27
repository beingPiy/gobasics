package main

import (
	"log" // for logging errors and other information
	"strconv" // for converting integers to strings
	"time" // for working with time and dates

	// for building the API server
	"github.com/gofiber/fiber/v3" // Using Fiber v3

	// for validating API requests
	"github.com/go-playground/validator/v10"

	// for logging HTTP requests (middleware)
	//"github.com/gofiber/fiber/v3/middleware/logger"
)

// album represents data about a record album

type album struct {
	// validation tags are added to struct fields to specify 
	// validation rules
	ID     string  `json:"id"`
	Title  string  `json:"title" validate:"required,min=3"`
	Artist string  `json:"artist" validate:"required,max=100"`
	Price  float64 `json:"price" validate:"required,gt=0"`
}

// albums slice to seed record album data
// swap for a real repository once we move to a database
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

var validate = validator.New() // initializes a new validator instance

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
func Logger () fiber.Handler {
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

// main function
func main() {

	// # Initializes Fiber app
	app := fiber.New()

	// # Register the custom logger middleware
	// Every request will go through this middleware, 
	// allowing us to log details about the request and response
	// Middleware is executed in the order it is registered, 
	// so it will wrap around all route handlers defined below
	app.Use(Logger())

	// # Routes

	// getAlbums responds with the list of all albums as JSON,
	// return type is error because error is taken for granted in Go,
	// and it is a common practice to return an error type in case 
	// something goes wrong
	app.Get("/albums", getAlbums)

	// getAlbumByID locates the album whose ID value matches the 
	// id parameter sent by the client, then returns that album as response
	app.Get("/albums/:id", getAlbumByID)

	// postAlbums adds an album from JSON received in the request body,
	// and it also validates the incoming data using the validator package
	app.Post("/albums", postAlbums)

	// updateAlbum updates an existing album based on the 
	// ID provided in the URL path,
	// and it also validates the incoming data using the validator package
	app.Put("/albums/:id", updateAlbum)

	// deleteAlbum deletes an existing album based on the ID provided in the URL path
	// Note: In a real application, you would typically use a database to store and manage albums,
	// and the delete operation would involve removing the album from the database rather than a slice.
	app.Delete("/albums/:id", deleteAlbum)

	log.Println("Starting server on port 3000...")
	// Start the server on port 3000
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}


// Request handlers for the API endpoints defined in the main function

// Each handler function takes a fiber.Ctx parameter, 
// which provides methods for handling the HTTP request and response

// Each handler function returns an error, 
// which allows for proper error handling in case 
// something goes wrong during the request processing

func getAlbums(c fiber.Ctx) error {
	// return all albums as JSON response
	return c.JSON(albums)
}

func getAlbumByID(c fiber.Ctx) error {

	// Get the id parameter from the URL path
	id := c.Params("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter
	// reange return index and value, we only need value here,
	// so we use _ to ignore the index
	for _, a := range albums {
		if a.ID == id {
			// Status code to notify successful creation of a resource, and return the album as JSON response
			return c.Status(201).JSON(a)
		}
	}
	// If no album is found with the given ID,
	// return a 404 status code and an error message as JSON response
	return c.Status(404).JSON(fiber.Map{"error": "Album not found"})
}

func postAlbums(c fiber.Ctx) error {

	// Create a new album struct to hold the incoming data
	newAlbum := new(album)

	// Call Bind to bind the received JSON from request body to newAlbum
	if err := c.Bind().Body(newAlbum); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate the newly created newAlbum struct as per the validation tags
	// defined in the album struct
	if err := validate.Struct(newAlbum); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	// Add the new album to the slice
	// The * operator is used to dereference the pointer newAlbum,
	// so that we can append the actual album value to the albums slice
	albums = append(albums, *newAlbum)
	return c.Status(201).JSON(newAlbum)
}

func updateAlbum(c fiber.Ctx) error {

	// Get the id parameter from the URL path
	id := c.Params("id")

	// This function would handle updating an existing album based on the ID provided in the URL path
	// Create a new album struct to hold the incoming data
	updatedAlbum := new(album)

	// Call Bind to bind the received JSON from request body to updatedAlbum
	if err := c.Bind().Body(updatedAlbum); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate the newly created updatedAlbum struct as per the validation tags
	// defined in the album struct
	if err := validate.Struct(updatedAlbum); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	// Add the new album to the slice
	// The * operator is used to dereference the pointer updatedAlbum,
	// so that we can replace the actual album value to the albums slice
	for i, a := range albums {
		if a.ID == id {
			albums[i] = *updatedAlbum
			return c.Status(201).JSON(updatedAlbum)
		}
	}
	return c.Status(404).JSON(fiber.Map{"error": "Album not found"})
}

func deleteAlbum(c fiber.Ctx) error {

	// Get the id parameter from the URL path
	id := c.Params("id")

	// This function would handle deleting an existing album 
	// based on the ID provided in the URL path
	for i, a := range albums {
		if a.ID == id {
			// Remove the album from the slice by slicing out 
			// the album at index i
			albums = append(albums[:i], albums[i+1:]...)
			return c.Status(204).Send(nil) // No content response
		}
	}
	return c.Status(404).JSON(fiber.Map{"error": "Album not found"})
}
