package routes

import (
	"github.com/Raviikumar001/exploding-kittens-server/controller"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(a *fiber.App) {
	
	route := a.Group("/auth/v1")

	//Routes for POST Method
	route.Post("/register", controller.UserSignUp)
	route.Post("/login", controller.UserSignIn)

}
