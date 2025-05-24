package services

import (
	"errors"
	"fmt"
	"time"

	"auth-service/internal/config"
	"auth-service/internal/models"
	"auth-service/internal/repositories"
	"auth-service/internal/security/jwt"
	"auth-service/internal/security/password"
	"auth-service/internal/utils"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound                     = errors.New("user not found")
	ErrInvalidPassword                  = errors.New("invalid password")
	ErrEmailExists                      = errors.New("email already exists")
	ErrInvalidToken                     = errors.New("invalid token")
	ErrTokenExpired                     = errors.New("token expired")
	ErrInvalidCode                      = errors.New("invalid verification code")
	ErrCodeExpired                      = errors.New("verification code expired")
	ErrInvalidRefreshToken              = errors.New("invalid refresh token")
	ErrRefreshTokenExpired              = errors.New("refresh token expired")
	ErrEmailNotConfirmed                = errors.New("email not confirmed")
	ErrVerificationSent                 = errors.New("verification code sent")
	ErrTokenInvalidatedByPasswordChange = errors.New("token invalidated by password change")
)

type AuthService struct {
	userRepo         *repositories.UserRepository
	refreshRepo      *repositories.RefreshRepository
	verificationRepo *repositories.VerificationRepository
	emailService     *EmailService
	tokenManager     *jwt.JWTManager
	cfg              *config.Config
	passwordHasher   *password.Hasher
}

func NewAuthService(
	userRepo *repositories.UserRepository,
	refreshRepo *repositories.RefreshRepository,
	verificationRepo *repositories.VerificationRepository,
	emailService *EmailService,
	tokenManager *jwt.JWTManager,
	passwordHasher *password.Hasher,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshRepo:      refreshRepo,
		verificationRepo: verificationRepo,
		emailService:     emailService,
		tokenManager:     tokenManager,
		cfg:              cfg,
		passwordHasher:   passwordHasher,
	}
}

func (s *AuthService) ValidateToken(tokenString string) (*models.TokenClaims, error) {
	claims, err := s.tokenManager.ParseAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Получаем актуальные данные пользователя из БД
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	// Проверяем, не был ли изменен пароль после создания токена
	if user.PasswordChangedAt.Unix() > claims.PasswordChangedAt {
		return nil, ErrTokenInvalidatedByPasswordChange
	}

	return claims, nil
}

func (s *AuthService) Register(input *models.UserCreate) (*models.User, error) {
	exists, err := s.userRepo.CheckEmailExists(input.Email)
	if err != nil {
		return nil, fmt.Errorf("error checking email existence: %w", err)
	}
	if exists {
		return nil, ErrEmailExists
	}

	hashedPassword, err := s.passwordHasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	now := time.Now()
	user := &models.User{
		ID:                uuid.New(),
		Email:             input.Email,
		PasswordHash:      hashedPassword,
		Role:              "student",
		Confirmed:         false,
		CreatedAt:         now,
		PasswordChangedAt: now,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	// Отправляем код подтверждения
	if err := s.sendVerificationCode(user.ID, user.Email, models.VerificationTypeRegistration); err != nil {
		return nil, fmt.Errorf("error sending verification code: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(input *models.UserLogin) (*models.TokenPair, error) {
	user, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if err := s.passwordHasher.Compare(input.Password, user.PasswordHash); err != nil {
		return nil, ErrInvalidPassword
	}

	if !user.Confirmed {
		return nil, ErrEmailNotConfirmed
	}

	// Отправляем код подтверждения для входа
	if err := s.sendVerificationCode(user.ID, user.Email, models.VerificationTypeLogin); err != nil {
		return nil, fmt.Errorf("error sending verification code: %w", err)
	}

	return nil, ErrVerificationSent
}

func (s *AuthService) VerifyEmail(email, code string) (*models.TokenPair, error) {
	return s.verifyCode(email, code, models.VerificationTypeRegistration)
}

func (s *AuthService) VerifyLogin(email, code string) (*models.TokenPair, error) {
	return s.verifyCode(email, code, models.VerificationTypeLogin)
}

func (s *AuthService) verifyCode(email, code string, verificationType models.VerificationType) (*models.TokenPair, error) {
	verificationCode, err := s.verificationRepo.GetActiveCode(email, verificationType)
	if err != nil {
		return nil, fmt.Errorf("error getting verification code: %w", err)
	}
	if verificationCode == nil {
		return nil, ErrCodeExpired
	}

	if verificationCode.Code != code {
		return nil, ErrInvalidCode
	}

	// Помечаем код как использованный
	if err := s.verificationRepo.MarkAsUsed(verificationCode.ID); err != nil {
		return nil, fmt.Errorf("error marking code as used: %w", err)
	}

	user, err := s.userRepo.GetByID(verificationCode.UserID)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	// Если это подтверждение регистрации, подтверждаем email
	if verificationType == models.VerificationTypeRegistration {
		if err := s.userRepo.UpdateConfirmation(user.ID); err != nil {
			return nil, fmt.Errorf("error confirming email: %w", err)
		}
	}

	// Генерируем токены
	accessToken, err := s.tokenManager.GenerateAccessToken(user.ID, user.Role, user.PasswordChangedAt)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("error generating refresh token: %w", err)
	}

	// Используем конфигурацию для TTL
	session := &models.RefreshSession{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.cfg.Token.RefreshTTL).Unix(),
		CreatedAt:    time.Now().Unix(),
	}

	if err := s.refreshRepo.Create(session); err != nil {
		return nil, fmt.Errorf("error saving refresh session: %w", err)
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) sendVerificationCode(userID uuid.UUID, email string, verificationType models.VerificationType) error {
	code, err := utils.GenerateVerificationCode()
	if err != nil {
		return fmt.Errorf("error generating verification code: %w", err)
	}

	verificationCode := &models.VerificationCode{
		ID:        uuid.New(),
		UserID:    userID,
		Email:     email,
		Code:      code,
		Type:      verificationType,
		Used:      false,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		CreatedAt: time.Now(),
	}

	if err := s.verificationRepo.Create(verificationCode); err != nil {
		return fmt.Errorf("error saving verification code: %w", err)
	}

	if err := s.emailService.SendVerificationCode(email, code, verificationType); err != nil {
		return fmt.Errorf("error sending verification email: %w", err)
	}

	return nil
}

// RefreshTokens обновляет пару токенов с помощью refresh token
func (s *AuthService) RefreshTokens(refreshToken string) (*models.TokenPair, error) {
	// Проверяем refresh token в базе данных
	session, err := s.refreshRepo.GetByToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	// Проверяем срок действия refresh token
	if time.Now().Unix() > session.ExpiresAt {
		// Удаляем просроченный токен
		_ = s.refreshRepo.Delete(session.ID)
		return nil, ErrRefreshTokenExpired
	}

	// Получаем пользователя
	user, err := s.userRepo.GetByID(session.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Создаем новую пару токенов
	accessToken, err := s.createAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	newRefreshToken := uuid.New().String()
	expiresAt := time.Now().Add(s.cfg.Token.RefreshTTL).Unix()

	// Обновляем refresh token в базе данных
	err = s.refreshRepo.Update(session.ID, newRefreshToken, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update refresh token: %w", err)
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) createAccessToken(user *models.User) (string, error) {
	return s.tokenManager.GenerateAccessToken(user.ID, user.Role, user.PasswordChangedAt)
}

func (s *AuthService) InitiatePasswordReset(email string) error {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Отправляем код подтверждения для сброса пароля
	if err := s.sendVerificationCode(user.ID, user.Email, models.VerificationTypePassword); err != nil {
		return fmt.Errorf("error sending verification code: %w", err)
	}

	return nil
}

func (s *AuthService) ResetPassword(email, code, newPassword string) error {
	verificationCode, err := s.verificationRepo.GetActiveCode(email, models.VerificationTypePassword)
	if err != nil {
		return fmt.Errorf("error getting verification code: %w", err)
	}
	if verificationCode == nil {
		return ErrCodeExpired
	}

	if verificationCode.Code != code {
		return ErrInvalidCode
	}

	// Помечаем код как использованный
	if err := s.verificationRepo.MarkAsUsed(verificationCode.ID); err != nil {
		return fmt.Errorf("error marking code as used: %w", err)
	}

	// Хешируем новый пароль
	hashedPassword, err := s.passwordHasher.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// Обновляем пароль пользователя
	if err := s.userRepo.UpdatePassword(verificationCode.UserID, hashedPassword); err != nil {
		return fmt.Errorf("error updating password: %w", err)
	}

	// Удаляем все активные сессии пользователя
	if err := s.refreshRepo.DeleteAllUserSessions(verificationCode.UserID); err != nil {
		return fmt.Errorf("error deleting user sessions: %w", err)
	}

	return nil
}
