package api

import (
	"beneburg/pkg/database"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type UsersAPI interface {
	RegisterRoutes(router *gin.RouterGroup)
	GetUser(*gin.Context)
	ListUsers(*gin.Context)
}

var _ UsersAPI = &usersApi{}

type usersApi struct {
	db     database.Database
	ctx    context.Context
	logger *zap.Logger
}

func NewUsersAPI(ctx context.Context, db database.Database, logger *zap.Logger) UsersAPI {
	return &usersApi{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (u usersApi) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/users/:userId", u.GetUser)
}

func (u usersApi) GetUser(g *gin.Context) {
	var userId uint64
	userIdStr := g.Param("userId")
	userId, err := strconv.ParseUint(userIdStr, 10, 32)
	if err != nil {
		g.JSON(400, gin.H{"error": "invalid user id"})
		return
	}
	user, err := u.db.GetUserByID(u.ctx, uint(userId))
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	g.JSON(200, user)
}

func (u usersApi) ListUsers(g *gin.Context) {
	//TODO implement me
	panic("implement me")
}
