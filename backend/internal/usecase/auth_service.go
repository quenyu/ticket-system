package usecase

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"ticket-system/backend/internal/domain"
	"ticket-system/backend/internal/model"
)

type AuthService struct {
	UserRepo  UserRepository
	JWTSecret string
}

func NewAuthService(userRepo UserRepository, jwtSecret string) *AuthService {
	return &AuthService{UserRepo: userRepo, JWTSecret: jwtSecret}
}

func (s *AuthService) Login(username, password string) (model.LoginResponse, error) {
	user, err := s.UserRepo.FindByUsername(username)
	if err != nil || user == nil {
		return model.LoginResponse{}, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return model.LoginResponse{}, errors.New("invalid credentials")
	}
	token, err := s.generateJWT(user)
	if err != nil {
		return model.LoginResponse{}, err
	}
	return model.LoginResponse{
		Token: token,
		User: model.LoginUserDTO{
			UserID:       user.ID.String(),
			Username:     user.Username,
			RoleID:       user.RoleID.String(),
			DepartmentID: user.DepartmentID,
		},
	}, nil
}

func (s *AuthService) generateJWT(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"role_id": user.RoleID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.JWTSecret))
}

type UserRegisterInput struct {
	Username     string
	Password     string
	Email        string
	RoleID       string
	DepartmentID int16
}

type UserRegisterOutput struct {
	UserID string
}

func (s *AuthService) Register(req model.UserRegisterRequest) (model.UserRegisterResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.UserRegisterResponse{}, errors.New("failed to hash password")
	}
	user := &domain.User{
		ID:           uuid.New(),
		Username:     req.Username,
		PasswordHash: string(hash),
		Email:        req.Email,
		RoleID:       uuid.MustParse(req.RoleID),
		DepartmentID: req.DepartmentID,
	}
	err = s.UserRepo.Create(user)
	if err != nil {
		return model.UserRegisterResponse{}, err
	}
	return model.UserRegisterResponse{UserID: user.ID.String()}, nil
}
