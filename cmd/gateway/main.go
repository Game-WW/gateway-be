package main

import (
	"log"
	"strconv"
	"time"

	"github.com/Game-WW/gateway-be/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("config: .env not found, using defaults")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("fatal: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName:      "gateway",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	if err := app.Listen(":" + strconv.Itoa(cfg.Port)); err != nil {
		log.Fatalf("server: failed to listen on %s: %v", cfg.Port, err)
	}
}
