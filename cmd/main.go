package main

import (
	"log"

	_ "learngo/docs" // Import generated docs
	"learngo/internal/common"

	httpapp "learngo/internal/delivery/http"
	mysql "learngo/internal/repository/mysql"

	"learngo/internal/usecase"
	"learngo/pkg/cache"
	"learngo/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title My Tasks API
// @version 1.0
// @description A task management API with authentication
// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using default values")
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: common.ErrorHandler,
	})
	app.Use(recover.New())

	app.Use(common.Logger)

	// Swagger route

	db := database.InitDB()
	cache.InitRedis()

	userRepo := mysql.NewUserMySQLRepo(db)
	userUC := usecase.NewUserUsecase(userRepo)
	httpapp.NewUserHandler(app, userUC)

	taskRepo := mysql.NewTaskMySQLRepo(db)
	taskUC := usecase.NewTaskUsecase(taskRepo)
	httpapp.NewTaskHandler(app, taskUC)
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	if err := app.Listen(":3000"); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
