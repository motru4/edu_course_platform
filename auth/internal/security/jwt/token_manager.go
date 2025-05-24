package jwt

import (
	"auth-service/internal/config"
	"auth-service/internal/models"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManager interface {
	GenerateAccessToken(userID uuid.UUID, role models.Role, passwordChangedAt time.Time) (string, error)
	GenerateRefreshToken() (string, error)
	ParseAccessToken(token string) (*models.TokenClaims, error)
	ParseRefreshToken(token string) (uuid.UUID, error)
}

type JWTManager struct {
	config config.TokenConfig
}

func NewJWTManager(config config.TokenConfig) *JWTManager {
	return &JWTManager{
		config: config,
	}
}

func (m *JWTManager) GenerateAccessToken(userID uuid.UUID, role models.Role, passwordChangedAt time.Time) (string, error) {
	claims := jwt.MapClaims{
		"user_id":     userID.String(),
		"role":        role,
		"exp":         time.Now().Add(m.config.AccessTTL).Unix(),
		"pepper":      m.config.PepperStr,
		"pwd_changed": passwordChangedAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.Secret))
}

func (m *JWTManager) GenerateRefreshToken() (string, error) {
	return uuid.NewString(), nil
}

func (m *JWTManager) ParseAccessToken(tokenString string) (*models.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.config.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	if pepper, ok := claims["pepper"].(string); !ok || pepper != m.config.PepperStr {
		return nil, fmt.Errorf("invalid token: pepper mismatch")
	}

	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return nil, fmt.Errorf("invalid user_id claim")
	}

	pwdChanged, ok := claims["pwd_changed"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid password_changed claim")
	}

	return &models.TokenClaims{
		UserID:            userID,
		Role:              models.Role(claims["role"].(string)),
		PasswordChangedAt: int64(pwdChanged),
	}, nil
}

func (m *JWTManager) ParseRefreshToken(token string) (uuid.UUID, error) {
	return uuid.Parse(token)
}
