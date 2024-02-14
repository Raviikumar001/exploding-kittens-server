package controller

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Raviikumar001/exploding-kittens-server/db"
	"github.com/Raviikumar001/exploding-kittens-server/models"

	"context"

	"github.com/gofiber/fiber/v2"

	// jwtware "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func UserSignUp(c *fiber.Ctx) error {

	siginUp := models.SignUp{}
	redisClient := db.GetRedisClient()

	if err := c.BodyParser(&siginUp); err != nil {

		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	fmt.Println(siginUp.Name)

	// Create a new user struct.
	user := &models.User{}

	user.ID = uuid.New()
	user.Name = siginUp.Name
	user.Username = siginUp.Username
	user.TotalPoints = 0
	user.TotalGamesLost = 0

	fmt.Println(user, "user ")

	jsonUser, err := json.Marshal(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to serialize user data",
		})
	}

	// fmt.Println(string(jsonUser))

	// Check for existing username in Redis
	usernameExists, err := redisClient.Exists(context.Background(), fmt.Sprintf("username:%s", siginUp.Username)).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Redis error",
		})
	}
	if usernameExists == 1 { // Username exists in Redis
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   "Username already exists",
		})
	}


	key := fmt.Sprintf("user:%s", user.ID.String())

	// Set the User data (you have previously checked for existing usernames)
	err = redisClient.Set(context.Background(), key, jsonUser, 0).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Redis error", 
		})
	}
	
	// Store a mapping of username to userID (no need for SetNX since the check is done)
	err = redisClient.Set(context.Background(), fmt.Sprintf("username:%s", siginUp.Username), user.ID.String(), 0).Err()
	if err != nil {
		// Handle potential errors in setting the username mapping
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Error storing username mapping", 
		})
	}
	
	token, exp, err := createJWTToken(*user)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"token": token, "exp": exp, "user": string(jsonUser),"msg": "User Created!" })

}


func UserSignIn(c *fiber.Ctx) error {
    signIn := models.SignIn{} // Assuming this struct only has a username field
    redisClient := db.GetRedisClient()

    if err := c.BodyParser(&signIn); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": true,
            "msg":   err.Error(),
        })
    }

    // 1. Check if the username exists in Redis
    userID, err := redisClient.Get(context.Background(), fmt.Sprintf("username:%s", signIn.Username)).Result()
    if err == redis.Nil {
        // Username not found
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": true,
            "msg":   "Username does not Exist", 
        })
    } else if err != nil {
        // Redis error
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": true,
            "msg":   "Redis error",
        })
    }

    // 2. Retrieve user data from Redis 
    key := fmt.Sprintf("user:%s", userID)
    userJSON, err := redisClient.Get(context.Background(), key).Result()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": true,
            "msg":   "Failed to retrieve user data",
        })
    }

    // 3. Unmarshal the user data (deserialize)
    var user models.User
    err = json.Unmarshal([]byte(userJSON), &user)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": true,
            "msg":   "Failed to deserialize user data",
        })
    }

    // 4. Generate JWT token on successful login
    token, exp, err := createJWTToken(user)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": true,
            "msg":   "Error generating token",
        })
    }

    return c.JSON(fiber.Map{
        "token": token,
        "exp":   exp,
        "msg":   "Login successful",
        "user":  userJSON, // Send back the user data (optional)
    })
}


func createJWTToken(user models.User) (string, int64, error) {

	secret := os.Getenv("JWT_SECRET")
    exp := time.Now().Add(time.Hour * 24).Unix()

    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["user_id"] = user.ID
    claims["exp"] = exp

    t, err := token.SignedString([]byte(secret))
    if err != nil {
        return "", 0, err
    }

    return t, exp, nil
}