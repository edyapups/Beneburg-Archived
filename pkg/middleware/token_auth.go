package middleware

import (
	"beneburg/pkg/database"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TokenAuth interface {
	Auth(ctx *gin.Context)
}
type tokenAuth struct {
	db     database.Database
	logger *zap.Logger
}

func (t tokenAuth) Auth(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")
	if token == "" {
		ctx.AbortWithStatusJSON(401, gin.H{"error": "token is required"})
		return
	}
	userID, err := t.db.GetUserIDByToken(ctx, token)
	if err != nil {
		ctx.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
		return
	}
	ctx.Set("currentUserID", userID)
}

func NewTokenAuth(db database.Database, logger *zap.Logger) TokenAuth {
	return &tokenAuth{
		db:     db,
		logger: logger,
	}
}

// devTokenAuth is a TokenAuth implementation that always returns the same user id.
type devTokenAuth struct {
}

func (t devTokenAuth) Auth(ctx *gin.Context) {
	ctx.Set("currentUserID", 1)
}

func NewDevTokenAuth() TokenAuth {
	return &devTokenAuth{}
}
