package main

import (
	"fmt"
	"os"
	// "os"
	"github.com/Raviikumar001/exploding-kittens-server/db"
	"github.com/Raviikumar001/exploding-kittens-server/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {

	app := fiber.New()
	err := godotenv.Load("app.env")

	if err != nil {
		fmt.Println("Error loading .env files")
	}
	db.StartRedis()

	app.Use(cors.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST,PATCH",
		AllowHeaders: "Origins, Content-Type, Accept",
	}))

	routes.AuthRoutes(app)
	routes.GameRoutes(app)

	port := os.Getenv("PORT")
	app.Get("/", func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{"message": "hello kitten"})
	})

	if err := app.Listen("0.0.0.0:" + port); err != nil {
		panic(err)
	}

}
