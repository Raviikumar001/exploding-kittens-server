package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	if port == "" {
		port = "8080"
	}
	app.Get("/", func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{"message": "hello kitten"})
	})

	// Graceful shutdown support
	go func() {
		if err := app.Listen("0.0.0.0:" + port); err != nil {
			log.Printf("Fiber server stopped: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")
	_ = app.Shutdown()

}
