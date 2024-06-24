// routes/userRoutes.go
package routes

import (
	"github.com/gofiber/fiber/v2"

	"backend-server/handlers"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/users", handlers.CreateUser)
	api.Get("/users/:email", handlers.GetUser)
	api.Put("/users/:id", handlers.UpdateUser)
	api.Delete("/users/:id", handlers.DeleteUser)
	api.Get("/users", handlers.GetUsers)
}
