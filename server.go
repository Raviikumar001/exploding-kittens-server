package main

import (
	"fmt"
	// "os"
	"github.com/Raviikumar001/exploding-kittens-server/db"
	"github.com/Raviikumar001/exploding-kittens-server/routes"


	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)


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

	routes.AuthRoutes(app)
	routes.GameRoutes(app)
	// app.Get("/private", func(c *fiber.Ctx) error {
	// 	return c.JSON(fiber.Map{"success": true, "path": "private"})

	// })

	// app.Get("/public", func(c *fiber.Ctx) error {
	// 	return c.JSON(fiber.Map{"success": true, "path": "public"})

	// })

	app.Get("/", func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{"message": "hello kitten"})
	})

	if err := app.Listen(":5000"); err != nil {
		panic(err)
	}

}
