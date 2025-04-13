package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smthjapanese/avito_pvz/internal/domain/models"
	"github.com/smthjapanese/avito_pvz/internal/domain/usecase/mock"
	"github.com/smthjapanese/avito_pvz/internal/pkg/errors"
	"github.com/smthjapanese/avito_pvz/internal/pkg/logger"
)

func TestUserHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	handler := NewUserHandler(mockUserUseCase, mockLogger)

	req := registerRequest{
		Email:    "test@example.com",
		Password: "password123",
		Role:     models.EmployeeRole,
	}
	reqBody, _ := json.Marshal(req)

	user := &models.User{
		ID:    uuid.New(),
		Email: req.Email,
		Role:  req.Role,
	}
	mockUserUseCase.EXPECT().Register(gomock.Any(), req.Email, req.Password, req.Role).Return(user, nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/register", handler.Register)

	c.Request, _ = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, user.Email, response.Email)
	assert.Equal(t, user.Role, response.Role)
}

func TestUserHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	handler := NewUserHandler(mockUserUseCase, mockLogger)

	req := loginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	reqBody, _ := json.Marshal(req)

	token := "jwt-token"
	mockUserUseCase.EXPECT().Login(gomock.Any(), req.Email, req.Password).Return(token, nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/login", handler.Login)

	c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), token)
}

func TestUserHandler_DummyLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	handler := NewUserHandler(mockUserUseCase, mockLogger)

	req := dummyLoginRequest{
		Role: models.EmployeeRole,
	}
	reqBody, _ := json.Marshal(req)

	token := "dummy-jwt-token"
	mockUserUseCase.EXPECT().DummyLogin(gomock.Any(), req.Role).Return(token, nil)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/dummyLogin", handler.DummyLogin)

	c.Request, _ = http.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), token)
}

func TestUserHandler_Register_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	handler := NewUserHandler(mockUserUseCase, mockLogger)

	req := registerRequest{
		Email:    "invalid-email",
		Password: "123",
		Role:     "invalid-role",
	}
	reqBody, _ := json.Marshal(req)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/register", handler.Register)

	c.Request, _ = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "message")
}

func TestUserHandler_Register_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	handler := NewUserHandler(mockUserUseCase, mockLogger)

	req := registerRequest{
		Email:    "test@example.com",
		Password: "password123",
		Role:     models.EmployeeRole,
	}
	reqBody, _ := json.Marshal(req)

	mockUserUseCase.EXPECT().Register(gomock.Any(), req.Email, req.Password, req.Role).Return(nil, errors.ErrUserAlreadyExists)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/register", handler.Register)

	c.Request, _ = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "already exists")
}

func TestUserHandler_Login_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mock.NewMockUserUseCase(ctrl)
	mockLogger, _ := logger.NewLogger("debug")
	handler := NewUserHandler(mockUserUseCase, mockLogger)

	req := loginRequest{
		Email:    "test@example.com",
		Password: "wrong-password",
	}
	reqBody, _ := json.Marshal(req)

	mockUserUseCase.EXPECT().Login(gomock.Any(), req.Email, req.Password).Return("", errors.ErrInvalidCredentials)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/login", handler.Login)

	c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, c.Request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid credentials")
}
