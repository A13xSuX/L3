package service

import (
	"context"
	"errors"
	"l3/WarehouseControl/internal/auth"
	"l3/WarehouseControl/internal/customErrs"
	"l3/WarehouseControl/internal/models"
	"l3/WarehouseControl/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	userRepo *repository.UserRepo
	jwt      *auth.JWT
}

func NewLoginService(userRepo *repository.UserRepo, jwt *auth.JWT) *LoginService {
	return &LoginService{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

func (s *LoginService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, customErrs.ErrInvalidCredentials
	}
	err = auth.CheckPasswordHash(req.Password, user.PasswordHash)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, customErrs.ErrInvalidCredentials
		}
		return nil, err
	}
	token, err := s.jwt.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	resp := models.LoginResponse{
		Username: user.Username,
		Role:     user.Role,
		Token:    token,
	}
	return &resp, nil
}
