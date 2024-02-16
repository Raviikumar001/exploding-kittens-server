package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Raviikumar001/exploding-kittens-server/db"
	"github.com/Raviikumar001/exploding-kittens-server/models"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func UpdateGameResults(c *fiber.Ctx) error {
	// 1. Extract Data & Validate
	fmt.Println("hi")
	updateData := models.UpdateData{}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Invalid request body",
		})
	}

	// 2. Fetch User from Redis
	key := fmt.Sprintf("user:%s", updateData.ID)
	redisClient := db.GetRedisClient() // Assuming you have this set up
	userJSON, err := redisClient.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "User not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Redis error",
		})
	}

	// 3. Modify the User Data
	user := &models.User{}
	err = json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to deserialize user data",
		})
	}

	if updateData.GameResult {
		user.TotalPoints += 1
	} else {
		user.TotalGamesLost += 1
	}

	// 4. Store Back in Redis
	newJSONUser, err := json.Marshal(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Failed to serialize user data",
		})
	}

	err = redisClient.Set(context.Background(), key, newJSONUser, 0).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Redis error",
		})
	}

	return c.JSON(fiber.Map{
		"msg": "Game results updated",
	})
}

func GetGameResult(c *fiber.Ctx) error {
	fmt.Println("hello")

	// Retrieve the ID from query parameters
	userID := c.Query("id")
	fmt.Println(userID, "hee")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Missing ID in query parameters",
		})
	}

	key := fmt.Sprintf("user:%s", userID)
	redisClient := db.GetRedisClient()

	userJSON, err := redisClient.Get(context.Background(), key).Result()
	// ... (rest of your error handling and response logic remains the same) ...

	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "User not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   "Redis error",
		})
	}

	// Send the JSON response directly (no need to deserialize with json.Unmarshal)
	return c.JSON(fiber.Map{
		"user": userJSON,
	})
}

// func GetGameResult(c *fiber.Ctx) error {

// 	fmt.Println("hello")
//      idData := models.IDData{}
//     if err := c.BodyParser(&idData); err != nil {
//         return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
//             "error": true,
//             "msg":   "Invalid request body",
//         })
//     }

//     key := fmt.Sprintf("user:%s", idData.ID)
//     redisClient := db.GetRedisClient() // Assuming you have this

//     userJSON, err := redisClient.Get(context.Background(), key).Result()
//     if err == redis.Nil {
//         return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
//             "error": true,
//             "msg":   "User not found",
//         })
//     } else if err != nil {
//         return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//             "error": true,
//             "msg":   "Redis error",
//         })
//     }

//     // Send the JSON response directly (no need to deserialize with json.Unmarshal)
//     return c.JSON(fiber.Map{
//         "user": userJSON,
//     })
// }
