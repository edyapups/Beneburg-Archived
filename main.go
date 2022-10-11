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
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("MYSQL_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	forceMigrationStr := os.Getenv("FORCE_MIGRATION")
	forceMigration, err := strconv.Atoi(forceMigrationStr)
	if err != nil && forceMigrationStr != "" {
		return err
	}

	dataSourceName := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true" + "&" + "multiStatements=true"
	db, err := database.NewDatabase(ctx, logger, dataSourceName)
	if err != nil {
		return err
	}
	defer func(db database.Database) {
		_ = db.Close()
	}(db)

	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	err = db.MakeMigrations(forceMigration)
	if err != nil {
		return err
	}

	router := gin.Default()
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	index := site.NewIndexConfig(logger, ctx, db)

	router.GET("/", index.Index)

	errChan := make(chan error)
	go func() {
		errChan <- router.Run(":8080")
	}()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

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
