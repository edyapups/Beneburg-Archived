package api

import (
	"beneburg/pkg/database"
	"beneburg/pkg/database/model"
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
	router.GET("/users/:user_id", u.GetUser)
	router.GET("/users", u.ListUsers)
	router.GET("/me", u.GetMe)
	router.PUT("/me", u.UpdateMe)
}

func (u usersApi) GetUser(g *gin.Context) {
	var userId uint64
	userIdStr := g.Param("user_id")
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
	users, err := u.db.GetAllUsers(u.ctx)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
		return
	}
	g.JSON(200, users)
}

func (u usersApi) GetMe(c *gin.Context) {
	currentUserId := c.GetUint("currentUserID")
	user, err := u.db.GetUserByID(c, currentUserId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}

func (u usersApi) UpdateMe(c *gin.Context) {
	var err error
	currentUserId := c.GetUint("currentUserID")
	updatedUser := model.User{}
	err = c.BindJSON(&updatedUser)
	if err != nil {
		return
	}
	result, err := u.db.UpdateUserByID(c, currentUserId, &updatedUser)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "user updated"})
}
