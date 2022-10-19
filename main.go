package main

import (
	"beneburg/pkg/database"
	"beneburg/pkg/site"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strconv"
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
	forceMigrationStr := os.Getenv("FORCE_MIGRATION")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	forceMigration, err := strconv.Atoi(forceMigrationStr)
	if err != nil && forceMigrationStr != "" {
		return err
	}

	dataSourceName := constructDBSourceName(dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := database.NewDatabase(ctx, logger, dataSourceName)
	if err != nil {
		return err
	}
	defer func(db database.Database) { _ = db.Close() }(db)

	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	logger.Info("Making migrations...")
	err = db.MakeMigrations(forceMigration)
	if err != nil {
		return err
	}
	logger.Info("Migrations done")

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

func constructDBSourceName(dbUser string, dbPassword string, dbHost string, dbPort string, dbName string) string {
	return dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true" + "&" + "multiStatements=true"
}
