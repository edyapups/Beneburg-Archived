package main

import (
	"beneburg/pkg/database"
	"beneburg/pkg/database/model"
	"beneburg/pkg/site"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
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

func run(logger *zap.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbName := os.Getenv("MYSQL_DATABASE")
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	onlyMakeMigrations := os.Getenv("ONLY_MAKE_MIGRATIONS") == "true"

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "3306"
	}

	dataSourceName := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true" + "&" + "multiStatements=true"
	db, err := database.NewDatabase(ctx, dataSourceName, logger)
	if err != nil {
		return err
	}

	models := []interface{}{&model.User{}, &model.Token{}}

	// Making migrations
	err = db.AutoMigrate(models...)
	if err != nil {
		return err
	}

	// Generating query schema
	db.GenerateCode(models...)

	if onlyMakeMigrations {
		logger.Info("Migrations and code generation were made, exiting...")
		return nil
	}

	router := gin.Default()
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	index := site.NewIndexConfig(logger, ctx, db)

	router.GET("/", index.Index)

	errChan := make(chan error)

	logger.Info("Starting server...")
	go func() {
		errChan <- router.Run(":8080")
	}()
	logger.Info("Server started")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("All ready")
	select {
	case err := <-errChan:
		return err
	case sig := <-sigs:
		logger.Info("Received signal", zap.String("signal", sig.String()))
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
