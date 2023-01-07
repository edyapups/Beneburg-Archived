package main

import (
	"beneburg/pkg/api"
	"beneburg/pkg/database"
	"beneburg/pkg/middleware"
	"beneburg/pkg/telegram"
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"time"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	err := run(logger)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

type Config struct {
	Database struct {
		User     string
		Password string
		Host     string
		Port     string
		Name     string

		DataSourceName     string
		OnlyMakeMigrations bool
	}
	Telegram struct {
		Token string
	}
	noAuth bool
}

func run(logger *zap.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := loadConfig()

	// Creating database connection
	db, err := database.NewDatabase(config.Database.DataSourceName, logger.Named("database"))
	if err != nil {
		return err
	}

	// Making migrations
	err = db.AutoMigrate(database.Models...)
	if err != nil {
		return err
	}

	// Stop if only migrations are needed
	if config.Database.OnlyMakeMigrations {
		logger.Info("Migrations and code generation were made, exiting...")
		return nil
	}

	// Configuring bot
	if token := config.Telegram.Token; token != "" {
		botAPI, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			return err
		}
		bot := telegram.NewBot(ctx, botAPI, db, logger.Named("telegram"))
		bot.Start()
	}

	// Configuring gin
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./frontend/dist", true)))
	router.Use(cors.Default())

	// Configuring API
	apiGroup := router.Group("/api")
	var tokenAuthMiddleware middleware.TokenAuth
	if config.noAuth {
		tokenAuthMiddleware = middleware.NewDevTokenAuth()
	} else {
		tokenAuthMiddleware = middleware.NewTokenAuth(db, logger.Named("TokenAuthMiddleware"))
	}
	apiGroup.Use(tokenAuthMiddleware.Auth)

	// User API
	usersAPI := api.NewUsersAPI(ctx, db, logger.Named("api/users"))
	usersAPI.RegisterRoutes(apiGroup)

	// Starting server
	logger.Info("Starting server...")
	errChan := make(chan error)
	go func() {
		errChan <- router.Run(":8080")
	}()
	logger.Info("Server started")

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, os.Interrupt)

	logger.Info("All ready")

	// Waiting for signal
	select {
	case err := <-errChan:
		return err
	case sig := <-sigs:
		logger.Info("Received signal", zap.String("signal", sig.String()))
		cancel()
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(time.Second * 10):
			return fmt.Errorf("timeout while waiting for context to be done")
		}
	case <-ctx.Done():
		return ctx.Err()
	}
}

func loadConfig() (*Config, error) {
	dbName := os.Getenv("MYSQL_DATABASE")
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	onlyMakeMigrations := os.Getenv("ONLY_MAKE_MIGRATIONS") == "true"
	botToken := os.Getenv("BOT_TOKEN")
	noAuth := os.Getenv("NO_AUTH") == "true"

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "3306"
	}

	dataSourceName := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true" + "&" + "multiStatements=true"
	return &Config{
		Database: struct {
			User               string
			Password           string
			Host               string
			Port               string
			Name               string
			DataSourceName     string
			OnlyMakeMigrations bool
		}{
			User:               dbUser,
			Password:           dbPassword,
			Host:               dbHost,
			Port:               dbPort,
			Name:               dbName,
			DataSourceName:     dataSourceName,
			OnlyMakeMigrations: onlyMakeMigrations,
		},
		Telegram: struct {
			Token string
		}{
			Token: botToken,
		},
		noAuth: noAuth,
	}, nil
}
