package middleware

import (
	"beneburg/pkg/database"
	"beneburg/pkg/database/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type TokenAuth interface {
	Auth(ctx *gin.Context)
}
type tokenAuth struct {
	db     database.Database
	logger *zap.Logger
}

func (t tokenAuth) Auth(ctx *gin.Context) {
	token, err := ctx.Cookie("token")
	if err != nil {
		t.logger.Named("Auth").Info("No token provided", zap.Any("headers", ctx.Request.Header))
		ctx.SetCookie("token", "", -1, "/", "", false, true)
		ctx.Redirect(http.StatusFound, "/login")
		ctx.Abort()
		return
	}

	user, err := t.db.GetUserByToken(ctx, token)
	if err != nil {
		t.logger.Named("Auth").Info("Error getting user id by token", zap.Error(err))
		ctx.SetCookie("token", "", -1, "/", "", false, true)
		ctx.Redirect(http.StatusFound, "/login")
		ctx.Abort()
		return
	}
	if user.Status == model.UserStatusBanned {
		t.logger.Named("Auth").Info("User is banned", zap.Any("user", user))
		ctx.Redirect(http.StatusFound, "/ban")
		ctx.Abort()
		return
	}
	ctx.Set("currentUser", user)
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
