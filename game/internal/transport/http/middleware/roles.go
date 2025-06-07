package middleware

import (
	"context"
	"net/http"
	"strings"

	pb "game/internal/transport/grpc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// RolesMiddleware предоставляет функционал для проверки ролей пользователя
type RolesMiddleware struct {
	authClient pb.AuthServiceClient
}

// NewRolesMiddleware создает новый экземпляр RolesMiddleware
func NewRolesMiddleware(authConn *grpc.ClientConn) *RolesMiddleware {
	return &RolesMiddleware{
		authClient: pb.NewAuthServiceClient(authConn),
	}
}

// RequireRoles проверяет, имеет ли пользователь указанные роли
func (m *RolesMiddleware) RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "отсутствует токен авторизации"})
			return
		}

		// Убираем префикс "Bearer " если он есть
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Создаем запрос к сервису авторизации
		req := &pb.CheckAccessRequest{
			Token:         token,
			RequiredRoles: roles,
		}

		// Отправляем запрос
		resp, err := m.authClient.CheckAccess(context.Background(), req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "ошибка при проверке прав доступа"})
			return
		}

		// Проверяем ответ
		if !resp.Allowed {
			statusCode := http.StatusForbidden
			if resp.Error == "invalid or expired token" {
				statusCode = http.StatusUnauthorized
			}
			c.AbortWithStatusJSON(statusCode, gin.H{"message": resp.Error})
			return
		}

		// Преобразуем строковый ID в UUID
		userID, err := uuid.Parse(resp.UserId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "некорректный формат ID пользователя"})
			return
		}

		// Сохраняем ID пользователя и роль в контексте
		c.Set("user_id", userID)
		//c.Set("user_role", resp.Role)
		c.Next()
	}
}
