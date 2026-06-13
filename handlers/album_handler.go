package handlers

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"gobasics/models"
	"gobasics/repository"
)

type AlbumHandler struct {
	repo repository.AlbumRepository
	validate *validator.Validate
}

func NewAlbumHandler(repo repository.AlbumRepository) *AlbumHandler {
	return &AlbumHandler{
		repo: repo,
		validate: validator.New(),
	}
}

// GetAllAlbums handles GET /albums requests,
// retrieves all albums from the repository and returns them as JSON.
func (h *AlbumHandler) GetAllAlbums(c fiber.Ctx) error {	
	albums, err := h.repo.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve albums",
		})
	}
	return c.JSON(albums)
}

// GetAlbumByID handles GET /albums/:id requests,
// retrieves a single album by its ID and returns it as JSON.
func (h *AlbumHandler) GetAlbumByID(c fiber.Ctx) error {	
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid album ID",
		})
	}
	album, err := h.repo.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Album not found",
		})
	}
	return c.JSON(album)
}

// CreateAlbum handles POST /albums requests,
// validates the incoming JSON payload and creates a new album in the repository.
func (h *AlbumHandler) CreateAlbum(c fiber.Ctx) error {	
	var req models.CreateAlbumRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}
	album, err := h.repo.Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create album",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(album)
}

// UpdateAlbum handles PUT /albums/:id requests,
// validates the incoming JSON payload and 
// updates an existing album in the repository.
func (h *AlbumHandler) UpdateAlbum(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid album ID",
		})
	}
	var req models.CreateAlbumRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request data",
		})
	}
	album, err := h.repo.Update(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(album)
}

// DeleteAlbum handles DELETE /albums/:id requests,
// deletes an album from the repository by its ID.
func (h *AlbumHandler) DeleteAlbum(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid album ID",
		})
	}
	err = h.repo.Delete(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}