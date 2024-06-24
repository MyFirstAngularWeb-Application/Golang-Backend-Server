// main.go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"backend-server/config"
	"backend-server/routes"
)

func main() {
	// Initialize MongoDB
	config.ConnectMongoDB()

	// Create Fiber app
	app := fiber.New()

	// Enable CORS
	// app.Use(cors.New())
	app.Use(func(c *fiber.Ctx) error {
		log.Printf("Incoming request: %s %s", c.Method(), c.OriginalURL())
		return c.Next()
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*", // Replace with your frontend URL
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: false,
	}))

	// Setup routes
	routes.SetupRoutes(app)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
