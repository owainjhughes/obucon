package services

import (
	"context"
	"errors"
	"fmt"
	"obucon/internal/models"
	"obucon/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication operations
type AuthService interface {
	Register(ctx context.Context, email, username, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	LoginWithUserID(ctx context.Context, email, password string) (string, uint, error)
	ValidateToken(tokenString string) (uint, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user account
func (s *authService) Register(ctx context.Context, email, username, password string) (*models.User, error) {
	// Validate email
	if email == "" {
		return nil, errors.New("email is required")
	}

	if len(username) < 3 || len(username) > 50 {
		return nil, errors.New("username must be between 3 and 50 characters")
	}

	if len(password) < 4 {
		return nil, errors.New("password must be at least 4 characters")
	}

	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already registered")
	}

	existingUser, err = s.userRepo.GetByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, errors.New("username already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Email:        email,
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	// Retrieve user
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if user == nil {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (s *authService) LoginWithUserID(ctx context.Context, email, password string) (string, uint, error) {
	// Retrieve user
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", 0, errors.New("invalid email or password")
	}

	if user == nil {
		return "", 0, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", 0, errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", 0, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, user.ID, nil
}

func (s *authService) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token claims")
	}

	// Check token expiration
	if exp, ok := (*claims)["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return 0, errors.New("token expired")
		}
	}

	userID, ok := (*claims)["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user_id in token")
	}

	return uint(userID), nil
}
