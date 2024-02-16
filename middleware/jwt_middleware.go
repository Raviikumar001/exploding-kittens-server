package middleware

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)
  
  func JWTMiddleware() fiber.Handler {
	  secret := os.Getenv("JWT_SECRET")
  
	  return func(c *fiber.Ctx) error {
		  authHeader := c.Get("Authorization")
		  fmt.Println(authHeader)
		  if !strings.Contains(authHeader, "Bearer") {
			  return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				  "error": true,
				  "msg":   "Unauthorized: Missing or malformed token",
			  })
		  }
  
		  // Extract token 
		  tokenString := strings.Replace(authHeader, "Bearer ", "", -1)
  
		  // Parse and validate
		  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method") // Lowercase
			}
			return []byte(secret), nil
		})
  
		  if err != nil {
			  return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				  "error": true,
				  "msg":   "Unauthorized: Invalid token",
			  })
		  }
  
		  claims, ok := token.Claims.(jwt.MapClaims)
		  if !ok || !token.Valid {
			  return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				  "error": true,
				  "msg":   "Unauthorized: Invalid token",
			  })
		  }
  
		  // Check for expiration 
		  exp := int64(claims["exp"].(float64))
		  if exp < time.Now().Unix() {
			  return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				  "error": true,
				  "msg":   "Unauthorized: Token expired", 
			  })
		  }
  
		  // If we reach here, token is valid. You might want to extract the user_id
		  // and make it available in context:
		  userID := claims["user_id"].(string)
		  c.Locals("userID", userID) 
  
		  return c.Next()
	  }
  }
  