package site

import (
	"beneburg/pkg/database"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type indexConfig struct {
	logger *zap.Logger
	ctx    context.Context
	db     database.Database
}

func NewIndexConfig(logger *zap.Logger, ctx context.Context, db database.Database) *indexConfig {
	return &indexConfig{
		logger: logger,
		ctx:    ctx,
		db:     db,
	}
}

func (i indexConfig) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Home Page",
	})
}
