package routes

import (
	"github.com/Raviikumar001/exploding-kittens-server/controller"
	"github.com/gofiber/fiber/v2"
	"github.com/Raviikumar001/exploding-kittens-server/middleware"
)

func GameRoutes(a *fiber.App) {
	
	route := a.Group("/api/v1")

	//Game routes
	route.Patch("/update-game-result",middleware.JWTMiddleware(),  controller.UpdateGameResults)
	route.Get("/get-result", middleware.JWTMiddleware(), controller.GetGameResult)

}
