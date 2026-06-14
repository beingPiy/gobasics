package middleware

import (
	"fmt" // for formatting error messages
	"os" // for accessing environment variables	
	
	"github.com/gofiber/fiber/v3" // for building the API server
	"github.com/golang-jwt/jwt/v4" // for generating and validating JWT tokens
)

func Protected() fiber.Handler {
	return func(c fiber.Ctx) error {

		// Step 1: Get the Authorization header
		// HTTP Convention : "Authorization: Bearer <token>	"
		authHeader := c.Get("Authorization")
		if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
				"hint": "Authorization header must be in the format 'Bearer <token>'",
			})
		}
		rawToken := authHeader[7:] // Extract the token part after "Bearer "

		// Step 2: Parse and validate the JWT token
		secretKey := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(rawToken, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method : %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		// Step 3: Extract claims and store them in Fiber's request scoped locals.
		// c.Locals() is a convenient way to store data that can be accessed by subsequent handlers in the request chain.
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token claims",
			})
		}

		c.Locals("user_id", int(claims["user_id"].(float64))) // JWT claims are float64 by default, convert to int
		c.Locals("email", claims["email"].(string))
		c.Locals("role", claims["role"].(string))	
		
		// Step 4: Call the next handler in the chain
		return c.Next()
	}
}

// AdminOnly is a middleware that checks if the user has an admin role
func AdminOnly() fiber.Handler {
	return func(c fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "forbidden - admin access required",
			})
		}	
		return c.Next()
	}
}

