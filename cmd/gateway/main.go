package main

import (
	"log"
	"strconv"

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
        ReadTimeout:  cfg.ReadTimeout,                                                                                                                                             
        WriteTimeout: cfg.WriteTimeout,                                                                                                                                            
        IdleTimeout:  cfg.IdleTimeout,                                                                                                                                             
  })   

	if err := app.Listen(":" + strconv.Itoa(cfg.Port)); err != nil {                                                                                                                 
        log.Fatalf("server: failed to listen on %s: %v", cfg.Port, err)                                                                                                            
  }  
}
