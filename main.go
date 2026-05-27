package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3" // Using Fiber v3
	//"github.com/gofiber/fiber/v3/middleware/logger"
)

// album represents data about a record album

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title" validate:"required,min=3"`
	Artist string  `json:"artist" validate:"required,max=10"`
	Price  float64 `json:"price" validate:"required,gt=0"`
}

// albums slice to seed record album data
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

var validate = validator.New()

// main function
func main() {

	// Initialize Fiber app
	app := fiber.New()

	app.Get("/albums", getAlbums)
	app.Get("/albums/:id", getAlbumByID)
	app.Post("/albums", postAlbums)
	// router.PUT("/albums/:id", updateAlbum)
	// router.DELETE("/albums/:id", deleteAlbum)

	// Start the server on port 3000
	app.Listen(":3000")

}

// getAlbums responds with the list of all albums as JSON
func getAlbums(c fiber.Ctx) error {
	return c.JSON(albums)
}

// postAlbums adds an album from JSON received in the request body
func postAlbums(c fiber.Ctx) error {
	newAlbum := new(album)

	// Call Bind to bind the received JSON to newAlbum
	if err := c.Bind().Body(newAlbum); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate the newAlbum struct
	if err := validate.Struct(newAlbum); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	// Add the new album to the slice
	albums = append(albums, *newAlbum)
	return c.Status(201).JSON(newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as response

func getAlbumByID(c fiber.Ctx) error {
	id := c.Params("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter
	for _, a := range albums {
		if a.ID == id {
			return c.Status(201).JSON(a)
		}
	}
	return c.Status(404).JSON(fiber.Map{"error": "Album not found"})
}
