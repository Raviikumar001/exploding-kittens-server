package main

import (
	"fmt"
	// "os"
	"github.com/Raviikumar001/exploding-kittens-server/db"
	// // "github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

type SignupRequest struct {
	name     string
	username string
}

func main() {

	app := fiber.New()
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env files")
	}
	db.StartRedis()


	app.Use(cors.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST",
		AllowHeaders: "Origins, Content-Type, Accept",
	}))

	app.Post("/signup", func(c *fiber.Ctx) error {

		req := new(SignupRequest)
		if err := c.BodyParser(req); err != nil {
			return err
		}

		if req.name == "" || req.username == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Invalied Signup Credentials")
		}

		return nil

	})

	app.Post("/login", func(c *fiber.Ctx) error {
		return nil

	})

	app.Get("/private", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "path": "private"})

	})

	app.Get("/public", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "path": "public"})

	})

	app.Get("/", func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{"message": "hello api"})
	})

	if err := app.Listen(":5000"); err != nil {
		panic(err)
	}

}
