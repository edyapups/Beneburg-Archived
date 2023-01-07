package api

import (
	mock_database "beneburg/pkg/database/mocks"
	"beneburg/pkg/database/model"
	"beneburg/pkg/middleware"
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gen"
	"gorm.io/gorm"
	"net/http/httptest"
	"testing"
)

func Test_usersApi(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	defer controller.Finish()
	dbMock := mock_database.NewMockDatabase(controller)
	api := NewUsersAPI(ctx, dbMock, nil)

	t.Run("GetMe", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := gin.New()
		r.Use(middleware.NewTokenAuth(dbMock, nil).Auth)
		api.RegisterRoutes(r.Group("/"))

		req := httptest.NewRequest("GET", "/me", nil)
		req.Header.Set("Authorization", "test")

		expectedUser := model.User{
			Model: gorm.Model{
				ID: 10,
			},
			TelegramID: 11,
		}
		dbMock.EXPECT().GetUserIDByToken(gomock.Any(), "test").Return(uint(10), nil).Times(1)
		dbMock.EXPECT().GetUserByID(gomock.Any(), uint(10)).Return(&expectedUser, nil).Times(1)

		r.ServeHTTP(w, req)

		actualUser := model.User{}
		err := json.Unmarshal(w.Body.Bytes(), &actualUser)
		assert.NoError(t, err)
		assert.True(t, assert.ObjectsAreEqualValues(expectedUser, actualUser))
	})

	t.Run("updateMe", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := gin.New()
		r.Use(middleware.NewTokenAuth(dbMock, nil).Auth)
		api.RegisterRoutes(r.Group("/"))

		newUser := model.User{
			TelegramID: 11,
			Model: gorm.Model{
				ID: 10,
			},
		}
		jsonBody, err := json.Marshal(newUser)
		assert.NoError(t, err)
		req := httptest.NewRequest("PUT", "/me", bytes.NewReader(jsonBody))
		req.Header.Set("Authorization", "test")

		dbMock.EXPECT().GetUserIDByToken(gomock.Any(), "test").Return(uint(10), nil).Times(1)
		dbMock.EXPECT().UpdateUserByID(gomock.Any(), uint(10), &newUser).Return(&gen.ResultInfo{
			RowsAffected: 0,
			Error:        nil,
		}, nil).Times(1)

		r.ServeHTTP(w, req)

		resultStruct := struct {
			Message string `json:"message"`
		}{}
		err = json.Unmarshal(w.Body.Bytes(), &resultStruct)
		assert.NoError(t, err)
		assert.Equal(t, "user updated", resultStruct.Message)
	})
}
