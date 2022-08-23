package main

import (
	"go-api-redis/order"
	"log"

	_ "go-api-redis/docs"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
)

// @title           Go API Redis
// @version         1.0
// @description     Go API Redis
// @BasePath /
// @schemes http https
func main() {
	var configuration struct {
		Listen   string `envconfig:"LISTEN" required:"true"`
		Redis    string `envconfig:"REDIS" required:"true"`
		Database int    `envconfig:"DATABASE" required:"true"`
		Password string `envconfig:"PASSWORD" required:"true"`
	}

	if err := envconfig.Process("", &configuration); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault)

	rdb := redis.NewClient(&redis.Options{
		Addr:     configuration.Redis,
		Password: configuration.Password,
		DB:       configuration.Database,
	})

	order.NewServer(app.Group("/order"), rdb).Route()

	if err := app.Listen(configuration.Listen); err != nil {
		log.Fatal(err)
	}
}
