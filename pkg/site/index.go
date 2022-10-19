package site

import (
	"beneburg/pkg/database"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type IndexConfig struct {
	logger *zap.Logger
	ctx    context.Context
	db     database.Database
}

func NewIndexConfig(ctx context.Context, db database.Database, logger *zap.Logger) *IndexConfig {
	return &IndexConfig{
		logger: logger,
		ctx:    ctx,
		db:     db,
	}
}

func (i IndexConfig) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Home Page",
	})
}
