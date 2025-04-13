package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
	"github.com/smthjapanese/avito_pvz/internal/domain/usecase"
	"github.com/smthjapanese/avito_pvz/internal/pkg/errors"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "user"
)

type AuthMiddleware struct {
	userUseCase usecase.UserUseCase
}

func NewAuthMiddleware(userUseCase usecase.UserUseCase) *AuthMiddleware {
	return &AuthMiddleware{
		userUseCase: userUseCase,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(authorizationHeader)
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "empty auth header"})
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid auth header"})
			return
		}

		token := headerParts[1]
		user, err := m.userUseCase.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return
		}

		c.Set(userCtx, user)
		c.Next()
	}
}

func (m *AuthMiddleware) CheckRole(roles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userValue, exists := c.Get(userCtx)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "user not found in context"})
			return
		}

		user, ok := userValue.(*models.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "user is of invalid type"})
			return
		}

		for _, role := range roles {
			if user.Role == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "access denied"})
	}
}

func GetUser(c *gin.Context) (*models.User, error) {
	userValue, exists := c.Get(userCtx)
	if !exists {
		return nil, errors.ErrUnauthorized
	}

	user, ok := userValue.(*models.User)
	if !ok {
		return nil, errors.ErrInternal
	}

	return user, nil
}
