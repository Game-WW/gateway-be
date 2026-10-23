package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Game-WW/gateway-be/docs"
	"github.com/Game-WW/gateway-be/internal/client"
	"github.com/Game-WW/gateway-be/internal/config"
	"github.com/Game-WW/gateway-be/internal/middleware"
	"github.com/Game-WW/gateway-be/internal/proxy"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	usersConn, err := client.NewConn(cfg.UsersGRPCAddr)
	if err != nil {
		log.Fatalf("fatal: %v", err)
	}
	defer usersConn.Close()

	mux, err := proxy.NewMux(context.Background(), usersConn)
	if err != nil {
		log.Fatalf("fatal: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName:      "gateway",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowHeaders:     "Authorization, Content-Type",
		AllowCredentials: true,
	}))

	docs.Register(app)

	proxyHandler := adaptor.HTTPHandler(mux)

	api := app.Group("/api/v1", middleware.Headers())
	api.All("/auth/*", proxyHandler)
	api.All("/*", middleware.Auth(cfg.JWTSecret), proxyHandler)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	addr := ":" + strconv.Itoa(cfg.Port)
	go func() {
		if err := app.Listen(addr); err != nil {
			log.Fatalf("server: failed to listen on %s: %v", addr, err)
		}
	}()

	<-ctx.Done()
	log.Println("server: shutting down")

	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Printf("server: shutdown: %v", err)
	}
}
