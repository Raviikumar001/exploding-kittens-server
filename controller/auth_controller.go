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

	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func UserSignUp(c *fiber.Ctx) error {

	siginUp := models.SignUp{}
	redisClient := db.GetRedisClient()

	if err := c.BodyParser(&siginUp); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	fmt.Println(siginUp.Name)

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

	usernameExists, err := redisClient.Exists(context.Background(), fmt.Sprintf("username:%s", siginUp.Username)).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Redis error",
		})
	}
	if usernameExists == 1 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": true,
			"msg":   "Username already exists",
		})
	}

	key := fmt.Sprintf("user:%s", user.ID.String())

	err = redisClient.Set(context.Background(), key, jsonUser, 0).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Redis error",
		})
	}

	err = redisClient.Set(context.Background(), fmt.Sprintf("username:%s", siginUp.Username), user.ID.String(), 0).Err()
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Error storing username mapping",
		})
	}

	token, exp, err := createJWTToken(*user)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"token": token, "exp": exp, "user": string(jsonUser), "msg": "User Created!"})

}

func UserSignIn(c *fiber.Ctx) error {
	signIn := models.SignIn{}
	redisClient := db.GetRedisClient()

	if err := c.BodyParser(&signIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	userID, err := redisClient.Get(context.Background(), fmt.Sprintf("username:%s", signIn.Username)).Result()
	if err == redis.Nil {

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "Username does not Exist",
		})
	} else if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Redis error",
		})
	}

	key := fmt.Sprintf("user:%s", userID)
	userJSON, err := redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to retrieve user data",
		})
	}

	var user models.User
	err = json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to deserialize user data",
		})
	}

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
		"user":  userJSON,
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
