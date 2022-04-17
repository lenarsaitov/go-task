package main

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"go-task/internals"
	"go-task/internals/config"
	"go-task/internals/db"
	"go-task/internals/services/cards"
	"go-task/internals/services/users"
	"go-task/pkg/logging"
	"go-task/pkg/shutdown"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"
)

func main() {
	logging.Init()

	logger := logging.GetLogger()
	logger.Println("logger initialized")

	logger.Println("config initializing")
	cfg := config.GetConfig()

	logger.Println("database initializing")
	postgres := db.NewPostgresDB(cfg)

	logger.Println("router initializing")
	router := internals.NewServer()

	logger.Println("create and register service user's storage, service and handlers")
	userStorage := users.NewUserStorage(postgres)
	userService := users.NewUserService(userStorage)
	userHandlers := users.NewUserHandler(userService)
	userRoot := router.Group("", internals.DefaultJsonContentTypeMiddleware())

	logger.Println("create and register service card's storage, service and handlers")
	cardStorage := cards.NewCardStorage(postgres)
	cardService := cards.NewCardService(cardStorage)
	cardHandlers := cards.NewCardHandler(cardService)
	cardRoot := router.Group("", internals.DefaultJsonContentTypeMiddleware())

	userHandlers.Setup(userRoot)
	cardHandlers.Setup(cardRoot)

	start(router, logger, cfg)
}

func start(router *echo.Echo, logger logging.Logger, cfg *config.Config) {
	logger.Infof("bind application to host: %s and port: %s", cfg.Listen.BindIP, cfg.Listen.Port)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	if err != nil {
		logger.Fatal(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}, server)

	logger.Println("application initialized and started")

	if err = server.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warn("server shutdown")
		default:
			logger.Fatal(err)
		}
	}
}
