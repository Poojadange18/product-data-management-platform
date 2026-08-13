package service

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"product-data-management-platform/model"
	"product-data-management-platform/repository"
)

var ErrEmailTaken = errors.New("email already registered")
var ErrInvalidCredentials = errors.New("invalid email or password")

type AuthService struct {
	users  *repository.UserRepository
	secret []byte
}

func NewAuthService(users *repository.UserRepository, secret string) *AuthService {
	return &AuthService{users: users, secret: []byte(secret)}
}

func (s *AuthService) Register(input model.RegisterRequest) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{Email: strings.ToLower(strings.TrimSpace(input.Email)), PasswordHash: string(hash)}
	if err := s.users.Create(user); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, ErrEmailTaken
		}
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(input model.LoginRequest) (string, *model.User, error) {
	user, err := s.users.FindByEmail(strings.ToLower(strings.TrimSpace(input.Email)))
	if err == sql.ErrNoRows {
		return "", nil, ErrInvalidCredentials
	}
	if err != nil {
		return "", nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)) != nil {
		return "", nil, ErrInvalidCredentials
	}
	claims := jwt.MapClaims{"sub": user.ID, "email": user.Email, "exp": time.Now().Add(24 * time.Hour).Unix(), "iat": time.Now().Unix()}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *AuthService) Parse(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
}
