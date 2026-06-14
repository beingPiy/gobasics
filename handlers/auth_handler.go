package handlers

import (
	"os"      // for accessing environment variables
	"time"    // for working with time and dates	
	
	"golang.org/x/crypto/bcrypt" // for hashing passwords securely
	"github.com/golang-jwt/jwt/v4" // for generating and validating JWT tokens
	"github.com/gofiber/fiber/v3" // for building the API server
	"github.com/go-playground/validator/v10" // for validating API requests

	"gobasics/models"     // for the User model and related request/response structs
	"gobasics/repository" // for interacting with the database through a repository pattern
)

type AuthHandler struct {
	userRepo repository.UserRepository // repository for user-related database operations
	validate *validator.Validate      // validator for validating API requests
}

func NewAuthHandler(userRepo repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		validate: validator.New(),
	}
}	

// Register handles user registration requests
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req models.RegisterRequest	
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Validate the request payload
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to hash password"})
	}

	user, err := h.userRepo.Create(c.Context(), req.Email, string(hashedPassword))
	if err != nil {
		// email already registered
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}
	
	token, err := generateJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}

	return c.Status(fiber.StatusCreated).JSON(models.AuthResponse{
		Token: token,
		User:  *user,
	})
}

// Login handles user login requests POST /login
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.userRepo.GetByEmail(c.Context(), req.Email)
	// return unauthorized if user not found or password does not match
	// A different error message for "user not found" vs "invalid password" can help attackers determine which emails are registered, so we return the same error for both cases
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid email or password"})
	}

	// Compare the provided password with the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid email or password"})
	}

	token, err := generateJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}

	return c.Status(fiber.StatusOK).JSON(models.AuthResponse{
		Token: token,
		User:  *user,
	})

}

// generateJWT creates a JWT token for the authenticated user	
// Extracted into its own  function so both Register and Login can use it
func generateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,	
		"email":   user.Email,
		"role":    user.Role,
		// exp : expiration time is a standard claim
		//  that indicates when the token should 
		// expire. Setting it to 24 hours means 
		// the token will be valid for 24 hours 
		// after it's issued. 
		// After that, the client will need to 
		// log in again to get a new token.
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // token expires in 24 hours
		"iat":     time.Now().Unix(), // issued at time : useful for auditing and debugging
	}

	// Create the token using the HS256 signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the secret key from environment variable
	secretKey := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secretKey))
}