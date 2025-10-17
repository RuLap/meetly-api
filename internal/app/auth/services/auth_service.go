package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/RuLap/meetly-api/meetly/internal/app/auth/dto"
	mapper "github.com/RuLap/meetly-api/meetly/internal/app/auth/mapper/custom"
	"github.com/RuLap/meetly-api/meetly/internal/app/auth/models"
	"github.com/RuLap/meetly-api/meetly/internal/app/auth/repository"
	"github.com/RuLap/meetly-api/meetly/internal/pkg/jwt_helper"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log/slog"
	"net/http"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	GoogleAuth(ctx context.Context, req dto.GoogleAuthRequest) (*dto.AuthResponse, error)
	GenerateGoogleOAuthURL() (string, string, error)
}

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type authService struct {
	log          *slog.Logger
	jwtHelper    *jwt_helper.JWTHelper
	googleConfig *GoogleOAuthConfig
	authRepo     repository.AuthRepository
}

func NewAuthService(
	log *slog.Logger,
	jwtHelper *jwt_helper.JWTHelper,
	googleConfig *GoogleOAuthConfig,
	authRepo repository.AuthRepository,
) AuthService {
	return &authService{
		log:          log,
		jwtHelper:    jwtHelper,
		googleConfig: googleConfig,
		authRepo:     authRepo,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to hash password", "error", err)
		return nil, fmt.Errorf("произошла ошибка")
	}

	hashedPasswordStr := string(hashedPassword)

	user := mapper.RegisterRequestToUser(&req, hashedPasswordStr)

	userID, err := s.authRepo.CreateUser(ctx, user)
	if err != nil {
		s.log.Error("failed to create user", "error", err, "email", user.Email)
		return nil, err
	}

	token, err := s.jwtHelper.GenerateDefaultToken(*userID, req.Email)
	if err != nil {
		s.log.Error("failed to generate JWT token", "error", err)
		return nil, fmt.Errorf("произошла ошибка")
	}

	s.log.Info("user registered successfully", "user_id", *userID, "email", req.Email)

	return &dto.AuthResponse{
		AccessToken: token,
		UserID:      *userID,
		Email:       req.Email,
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.authRepo.GetByEmailProvider(ctx, req.Email, models.LocalProvider)
	if err != nil {
		s.log.Warn("user not found", "email", req.Email)
		return nil, fmt.Errorf("неверный email или пароль")
	}

	if user.Password == nil {
		s.log.Error("user entered empty password", "email", req.Email)
		return nil, fmt.Errorf("неверный email или пароль")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(req.Password)); err != nil {
		s.log.Error("user entered invalid password", "email", req.Email)
		return nil, fmt.Errorf("неверный email или пароль")
	}

	token, err := s.jwtHelper.GenerateDefaultToken(user.ID.String(), user.Email)
	if err != nil {
		s.log.Error("failed to generate JWT token", "error", err)
		return nil, fmt.Errorf("произошла ошибка")
	}

	s.log.Info("user logged in successfully", "user_id", user.ID, "email", req.Email)

	return &dto.AuthResponse{
		AccessToken: token,
		UserID:      user.ID.String(),
		Email:       user.Email,
	}, nil
}

func (s *authService) GoogleAuth(ctx context.Context, req dto.GoogleAuthRequest) (*dto.AuthResponse, error) {
	token, err := s.exchangeCodeForToken(req.Code)
	if err != nil {
		s.log.Error("failed to exchange code for token", "error", err)
		return nil, fmt.Errorf("ошибка авторизации через Google")
	}

	userInfo, err := s.getGoogleUserInfo(token)
	if err != nil {
		s.log.Error("failed to get user info from Google", "error", err)
		return nil, fmt.Errorf("ошибка получения данных от Google")
	}

	user, err := s.authRepo.GetByEmailProvider(ctx, userInfo.Email, models.GoogleProvider)
	if err != nil {
		s.log.Info("creating new google user", "email", userInfo.Email)
	}

	providerID := userInfo.ID
	user = &models.User{
		Email:          userInfo.Email,
		Provider:       models.GoogleProvider,
		ProviderID:     &providerID,
		EmailConfirmed: true,
	}

	userID, err := s.authRepo.CreateUser(ctx, user)
	if err != nil {
		s.log.Error("failed to create google user", "error", err, "email", userInfo.Email)
		return nil, err
	}

	user.ID, err = uuid.Parse(*userID)
	if err != nil {
		s.log.Error("failed to parse user id", "error", err, "email", userInfo.Email)
		return nil, fmt.Errorf("произошла ошибка")
	}

	jwtToken, err := s.jwtHelper.GenerateDefaultToken(user.ID.String(), user.Email)
	if err != nil {
		s.log.Error("failed to generate JWT token", "error", err)
		return nil, fmt.Errorf("произошла ошибка")
	}

	s.log.Info("google auth successful", "user_id", user.ID, "email", user.Email)

	return &dto.AuthResponse{
		AccessToken: jwtToken,
		UserID:      user.ID.String(),
		Email:       user.Email,
	}, nil
}

func (s *authService) GenerateGoogleOAuthURL() (string, string, error) {
	state, err := generateState()
	if err != nil {
		return "", "", fmt.Errorf("не удалось сгенерировать параметр безопасности")
	}

	url := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email profile&state=%s",
		s.googleConfig.ClientID,
		s.googleConfig.RedirectURL,
		state,
	)

	return url, state, nil
}

func (s *authService) ValidateToken(token string) (bool, error) {
	valid, err := s.jwtHelper.ValidateToken(token)
	if err != nil {
		s.log.Warn("token validation failed", "error", err)
		return false, fmt.Errorf("неверный токен")
	}
	return valid, nil
}

func (s *authService) exchangeCodeForToken(code string) (string, error) {
	url := "https://oauth2.googleapis.com/token"

	data := fmt.Sprintf(
		"code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code",
		code, s.googleConfig.ClientID, s.googleConfig.ClientSecret, s.googleConfig.RedirectURL,
	)

	resp, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.Error != "" {
		return "", fmt.Errorf("google oauth error: %s", result.Error)
	}

	return result.AccessToken, nil
}

func (s *authService) getGoogleUserInfo(accessToken string) (*dto.GoogleUserInfo, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo dto.GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(bytes.NewReader([]byte("random-seed-for-now")), b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
