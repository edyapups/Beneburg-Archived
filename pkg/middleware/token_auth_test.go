package middleware

import (
	mock_database "beneburg/pkg/database/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func Test_TokenAuth(t *testing.T) {
	t.Run("Auth", func(t *testing.T) {
		w := httptest.NewRecorder()
		controller := gomock.NewController(t)
		defer controller.Finish()
		dbMock := mock_database.NewMockDatabase(controller)
		r := gin.New()
		authMiddleware := NewTokenAuth(dbMock, nil)
		r.Use(authMiddleware.Auth)
		r.GET("/test", func(ctx *gin.Context) {
			ctx.String(200, "%v", ctx.GetUint("currentUserID"))
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "test")

		dbMock.EXPECT().GetUserIDByToken(gomock.Any(), "test").Return(uint(10), nil).Times(1)

		r.ServeHTTP(w, req)
		assert.Equal(t, "10", w.Body.String())
	})
}
