package main

import (
	"beneburg/pkg/database"
	"beneburg/pkg/middleware"
	"beneburg/pkg/telegram"
	"beneburg/pkg/views"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
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
		logger.Info("Migrations were made, exiting...")
		return nil
	}



	// TLS
	var tlsConfig *tls.Config
	var manager *autocert.Manager
	if config.tls {
		manager = &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("api.telegram.org"),
			Cache:      autocert.DirCache("/root/.cache/golang-autocert"),
		}
		tlsConfig = &tls.Config{
			GetCertificate: manager.GetCertificate,
		}
	}



	// Configuring bot
	var SendFunc telegram.TelegramBotSendFunc
	if token := config.Telegram.Token; token != "" {
		botAPI, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			return err
		}
		bot := telegram.NewBot(ctx, botAPI, db)
		bot.Start()
		SendFunc = bot.GetSendFunc()
		if config.Telegram.AdminID != nil {
			logger = logger.WithOptions(zap.Hooks(func(entry zapcore.Entry) error {
				if entry.Level < zapcore.WarnLevel {
					return nil
				}
				msg := tgbotapi.NewMessage(*config.Telegram.AdminID, fmt.Sprintf("%s: %s (%s)", entry.Level, entry.Message, entry.Caller))
				SendFunc(msg)
				return nil
			}))
		}
		bot.SetLogger(logger.Named("telegram"))
	}

	// Log panic
	defer func() {
		if r := recover(); r != nil {
			logger.Panic("panic", zap.Any("panic", r))
		}
	}()



	// Configuring gin
	router := gin.Default()
	if config.trustedProxy != "" {
		err := router.SetTrustedProxies(strings.Split(config.trustedProxy, ","))
		if err != nil {
			return err
		}
	}
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*")
	router.Use(cors.Default())

	// Configuring API
	var tokenAuthMiddleware middleware.TokenAuth
	if config.noAuth {
		tokenAuthMiddleware = middleware.NewDevTokenAuth()
	} else {
		tokenAuthMiddleware = middleware.NewTokenAuth(db, logger.Named("TokenAuthMiddleware"))
	}

	// Configuring groups
	loginGroup := router.Group("/login")
	profileGroup := router.Group("/profile")
	mainGroup := router.Group("/")

	// TokenAuthMiddleware
	mainGroup.Use(tokenAuthMiddleware.Auth)
	profileGroup.Use(tokenAuthMiddleware.Auth)

	// ProfileRedirectMiddleware
	mainGroup.Use(middleware.ProfileRedirectMiddleware())

	// Views
	viewsModule := views.NewViews(db, logger.Named("views"), SendFunc)
	viewsModule.RegisterRoutes(mainGroup)
	viewsModule.RegisterLogin(loginGroup)
	viewsModule.RegisterProfile(profileGroup)



	// Starting server
	logger.Info("Starting server...")
	if config.tls {
		logger.Info("Starting HTTPS server...")
		s1 := &http.Server{
			Addr:    ":http",
			Handler: http.HandlerFunc(redirect),
		}
		s2 := &http.Server{
			TLSConfig: tlsConfig,
			Handler:   router,
		}

		// TODO: add goroutine waitgroup
		go func() {
			err := s1.ListenAndServe()
			if err != nil {
				logger.Error("ListenAndServe", zap.Error(err))
			}
		}()
		go func() {
			err := s2.Serve(manager.Listener())
			if err != nil {
				logger.Error("Serve", zap.Error(err))
			}
		}()
	} else {
		logger.Info("Starting HTTP server...")
		s1 := &http.Server{
			Addr:    ":8080",
			Handler: router,
		}
		// TODO: add goroutine waitgroup
		go func() {
			err := s1.ListenAndServe()
			if err != nil {
				logger.Error("ListenAndServe", zap.Error(err))
			}
		}()
	}

	logger.Info("Server started")
	logger.Info("All ready")



	// Waiting for signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case sig := <-sigs:
		logger.Info("Received signal", zap.String("signal", sig.String()))
		cancel()
		select {
		// TODO: wait for all goroutines to finish
		case <-time.After(time.Second * 4):
			return fmt.Errorf("timeout while waiting for context to be done")
		}
	case <-ctx.Done():
		return ctx.Err()
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
		Token   string
		AdminID *int64
		GroupID *int64
	}
	tls          bool
	trustedProxy string
	noAuth       bool
}

func loadConfig() (*Config, error) {
	err := os.Setenv("XDG_CACHE_HOME", "/root/.cache")
	if err != nil {
		return nil, err
	}
	dir, _ := os.UserHomeDir()
	fmt.Printf("os.UserHomeDir() = %s\n", dir)
	home := os.Getenv("HOME")
	fmt.Printf("os.Getenv(\"HOME\") = %s\n", home)

	// Loading config
	dbName := os.Getenv("MYSQL_DATABASE")
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	onlyMakeMigrations := os.Getenv("ONLY_MAKE_MIGRATIONS") == "true"
	botToken := os.Getenv("BOT_TOKEN")
	noAuth := os.Getenv("NO_AUTH") == "true"
	tlsEnv := os.Getenv("TLS") == "true"
	trustedProxy := os.Getenv("TRUSTED_PROXY")

	var adminID, groupID *int64
	if strVal := os.Getenv("ADMIN_ID"); strVal != "" {
		val, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			return nil, err
		}
		adminID = &val
	}
	if strVal := os.Getenv("GROUP_ID"); strVal != "" {
		val, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			return nil, err
		}
		groupID = &val
	}

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
			Token   string
			AdminID *int64
			GroupID *int64
		}{
			Token:   botToken,
			AdminID: adminID,
			GroupID: groupID,
		},
		tls:          tlsEnv,
		trustedProxy: trustedProxy,
		noAuth:       noAuth,
	}, nil
}

func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.RequestURI

	http.Redirect(w, req, target, http.StatusMovedPermanently)
}
