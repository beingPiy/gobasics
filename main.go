package main

import (
	// for building the API server
	"github.com/gofiber/fiber/v3" // Using Fiber v3

	// for validating API requests
	"github.com/go-playground/validator/v10"
	// for logging HTTP requests (middleware)
	//"github.com/gofiber/fiber/v3/middleware/logger"
)

// album represents data about a record album

type album struct {
	// validation tags are added to struct fields to specify validation rules
	ID     string  `json:"id"`
	Title  string  `json:"title" validate:"required,min=3"`
	Artist string  `json:"artist" validate:"required,max=100"`
	Price  float64 `json:"price" validate:"required,gt=0"`
}

// albums slice to seed record album data
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

var validate = validator.New() // initializes a new validator instance

// main function
func main() {

	// Initializes Fiber app
	app := fiber.New()

	// getAlbums responds with the list of all albums as JSON,
	// return type is error because error is taken for granted in Go,
	// and it is a common practice to return an error type in case something goes wrong
	app.Get("/albums", getAlbums)

	// getAlbumByID locates the album whose ID value matches the id parameter
	// sent by the client, then returns that album as response
	app.Get("/albums/:id", getAlbumByID)

	// postAlbums adds an album from JSON received in the request body,
	// and it also validates the incoming data using the validator package
	app.Post("/albums", postAlbums)

	// updateAlbum updates an existing album based on the ID provided in the URL path,
	// and it also validates the incoming data using the validator package
	app.Put("/albums/:id", updateAlbum)

	// deleteAlbum deletes an existing album based on the ID provided in the URL path
	// Note: In a real application, you would typically use a database to store and manage albums,
	// and the delete operation would involve removing the album from the database rather than a slice.
	app.Delete("/albums/:id", deleteAlbum)

	// Start the server on port 3000
	app.Listen(":3000")

}

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

	// This function would handle deleting an existing album based on the ID provided in the URL path
	for i, a := range albums {
		if a.ID == id {
			// Remove the album from the slice by slicing out the album at index i
			albums = append(albums[:i], albums[i+1:]...)
			return c.Status(204).Send(nil) // No content response
		}
	}
	return c.Status(404).JSON(fiber.Map{"error": "Album not found"})
}
