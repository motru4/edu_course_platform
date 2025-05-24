package handler

import (
	"fmt"
	"net/http"

	"auth-service/internal/config"
	"auth-service/internal/models"
	"auth-service/internal/services"
	"auth-service/internal/transport/http/middleware"
	"auth-service/internal/utils"

	_ "auth-service/docs" // импорт сгенерированной документации

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	authService *services.AuthService
	cfg         *config.Config
}

func NewHandler(router *gin.Engine, authService *services.AuthService, cfg *config.Config) *Handler {
	h := &Handler{
		authService: authService,
		cfg:         cfg,
	}

	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit)

	// Public routes
	v1 := router.Group("/api/v1/auth")
	{
		v1.POST("/register", rateLimiter.RateLimit(), h.register)
		v1.POST("/login", rateLimiter.RateLimit(), h.login)
		v1.POST("/verify-email", rateLimiter.RateLimit(), h.verifyEmail)
		v1.POST("/verify-login", rateLimiter.RateLimit(), h.verifyLogin)
		v1.POST("/refresh", rateLimiter.RateLimit(), h.refresh)
		v1.GET("/oauth/google", rateLimiter.RateLimit(), h.googleOAuth)
		v1.POST("/reset-password/request", rateLimiter.RateLimit(), h.requestPasswordReset)
		v1.POST("/reset-password/confirm", rateLimiter.RateLimit(), h.confirmPasswordReset)
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return h
}

// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя и отправляет код подтверждения на email
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.UserCreate true "Данные для регистрации"
// @Success 200 {object} string "Код подтверждения отправлен"
// @Failure 400 {object} string "Некорректные входные данные"
// @Failure 409 {object} string "Email уже существует"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /register [post]
func (h *Handler) register(c *gin.Context) {
	var input models.UserCreate

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !utils.IsValidEmail(input.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		return
	}

	if !utils.IsValidPassword(input.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 8 characters and contain numbers and letters in different cases"})
		return
	}

	_, err := h.authService.Register(&input)
	if err != nil {
		if err == services.ErrEmailExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verification code sent to your email"})
}

// @Summary Аутентификация пользователя
// @Description Проверяет учетные данные и отправляет код подтверждения на email
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.UserLogin true "Данные для входа"
// @Success 200 {object} string "Код подтверждения отправлен"
// @Failure 400 {object} string "Некорректные входные данные"
// @Failure 401 {object} string "Неверные учетные данные"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /login [post]
func (h *Handler) login(c *gin.Context) {
	var input models.UserLogin

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.authService.Login(&input)
	if err != nil {
		switch err {
		case services.ErrVerificationSent:
			c.JSON(http.StatusOK, gin.H{"message": "verification code sent to your email"})
		case services.ErrUserNotFound, services.ErrInvalidPassword:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		case services.ErrEmailNotConfirmed:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "email not confirmed"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
}

// @Summary Подтверждение email при регистрации
// @Description Подтверждает email с помощью кода подтверждения
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.VerificationRequest true "Код подтверждения"
// @Success 200 {object} models.TokenPair "Токены доступа"
// @Failure 400 {object} string "Некорректный код"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /verify-email [post]
func (h *Handler) verifyEmail(c *gin.Context) {
	var input models.VerificationRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.authService.VerifyEmail(input.Email, input.Code)
	if err != nil {
		switch err {
		case services.ErrInvalidCode:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification code"})
		case services.ErrCodeExpired:
			c.JSON(http.StatusBadRequest, gin.H{"error": "verification code expired"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// @Summary Подтверждение входа
// @Description Подтверждает вход с помощью кода подтверждения
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.VerificationRequest true "Код подтверждения"
// @Success 200 {object} models.TokenPair "Токены доступа"
// @Failure 400 {object} string "Некорректный код"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /verify-login [post]
func (h *Handler) verifyLogin(c *gin.Context) {
	var input models.VerificationRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.authService.VerifyLogin(input.Email, input.Code)
	if err != nil {
		switch err {
		case services.ErrInvalidCode:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification code"})
		case services.ErrCodeExpired:
			c.JSON(http.StatusBadRequest, gin.H{"error": "verification code expired"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// @Summary OAuth авторизация через Google
// @Description Выполняет OAuth 2.0 авторизацию через Google
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} models.TokenPair "Токены доступа"
// @Failure 501 {object} string "Функционал не реализован"
// @Router /oauth/google [get]
func (h *Handler) googleOAuth(c *gin.Context) {
	// TODO: Implement Google OAuth
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// @Summary Обновление токена доступа
// @Description Обновляет access token с помощью refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.RefreshInput true "Refresh токен"
// @Success 200 {object} models.TokenPair "Новая пара токенов"
// @Failure 400 {object} string "Некорректные входные данные"
// @Failure 401 {object} string "Невалидный refresh token"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /refresh [post]
func (h *Handler) refresh(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token is required"})
		return
	}

	tokens, err := h.authService.RefreshTokens(input.RefreshToken)
	if err != nil {
		switch err {
		case services.ErrInvalidRefreshToken:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		case services.ErrRefreshTokenExpired:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expired"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// @Summary Запрос на сброс пароля
// @Description Отправляет код подтверждения на email для сброса пароля
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.PasswordResetRequest true "Email пользователя"
// @Success 200 {object} string "Код подтверждения отправлен"
// @Failure 400 {object} string "Некорректные входные данные"
// @Failure 404 {object} string "Пользователь не найден"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /reset-password/request [post]
func (h *Handler) requestPasswordReset(c *gin.Context) {
	var input models.PasswordResetRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.InitiatePasswordReset(input.Email); err != nil {
		switch err {
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verification code sent to your email"})
}

// @Summary Подтверждение сброса пароля
// @Description Сбрасывает пароль с помощью кода подтверждения
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.PasswordResetConfirm true "Данные для сброса пароля"
// @Success 200 {object} string "Пароль успешно изменен"
// @Failure 400 {object} string "Некорректные входные данные"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /reset-password/confirm [post]
func (h *Handler) confirmPasswordReset(c *gin.Context) {
	var input models.PasswordResetConfirm

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !utils.IsValidPassword(input.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 8 characters and contain numbers and letters in different cases"})
		return
	}

	if err := h.authService.ResetPassword(input.Email, input.Code, input.NewPassword); err != nil {
		switch err {
		case services.ErrInvalidCode:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification code"})
		case services.ErrCodeExpired:
			c.JSON(http.StatusBadRequest, gin.H{"error": "verification code expired"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "password successfully changed, all active sessions have been terminated",
		"details": "you will need to log in again on all devices",
	})
}
