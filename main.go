package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/XEFF09/fulldev-be/inits"
	"github.com/gofiber/swagger"

	"github.com/XEFF09/fulldev-be/cogs/adapters/repositories"
	"github.com/XEFF09/fulldev-be/cogs/adapters/rest"
	"github.com/XEFF09/fulldev-be/cogs/usecases"
)

// @title Book API
// @description This is a sample server for a book API.
// @version 1.0
// @host localhost:8088
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	//TODO: init database connection
	dbInstance := inits.DatabaseConnection()

	//TODO: init clean
	gormUserRepo := repositories.NewGormUserRepository(dbInstance)
	userService := usecases.NewUserService(gormUserRepo)
	userRest := rest.NewUserRest(userService)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, HEAD, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	app.Use("/hello", authentication)

	app.Post("/register", userRest.Register)
	app.Post("/login", userRest.Login)
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":8088")
}
