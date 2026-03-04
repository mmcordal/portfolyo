package service

import (
	"context"
	"errors"
	"fmt"
	"portfolyo/internal/infrastructure/config"
	"portfolyo/internal/model"
	"portfolyo/internal/repository"
	"portfolyo/internal/viewmodel"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, vm *viewmodel.RegisterRequest) error
	Login(ctx context.Context, vm viewmodel.LoginRequest) (*viewmodel.LoginResponse, error)
	GetUserProfile(ctx context.Context, email string) (*viewmodel.UserVM, error)
	UpdateUser(ctx context.Context, email string, vm *viewmodel.UpdateRequest) error
	DeleteUser(ctx context.Context, email string) error
}

type authService struct {
	ur repository.UserRepository
}

func NewAuthService(ur repository.UserRepository) AuthService {
	return &authService{ur: ur}
}

func (s *authService) Register(ctx context.Context, vm *viewmodel.RegisterRequest) error {
	user := &model.User{
		Name:     vm.Name,
		Surname:  vm.Surname,
		Email:    vm.Email,
		Password: vm.Password,
	}

	if user.Name == "" || user.Surname == "" || user.Email == "" || user.Password == "" {
		return errors.New("all fields are required")
	}

	user.Email = strings.ToLower(user.Email)
	exist, err := s.ur.ExistEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("email already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing errorsx: %v", err)
	}
	user.Password = string(hashed)

	return s.ur.Create(ctx, user)
}

type accessToken struct {
	jwt.RegisteredClaims
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
	UserID  int64  `json:"user_id"`
}

func (s *authService) Login(ctx context.Context, vm viewmodel.LoginRequest) (*viewmodel.LoginResponse, error) {
	if vm.Email == "" || vm.Password == "" {
		return nil, errors.New("email and password are required")
	}

	vm.Email = strings.ToLower(vm.Email)
	user, err := s.ur.GetByEmail(ctx, vm.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if user.DeletedAt != nil {
		return nil, errors.New("user account is deleted")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(vm.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	jwtSecret := []byte(config.Get().Secret.JWTSecret)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		Name:    user.Name,
		Surname: user.Surname,
		Email:   user.Email,
		UserID:  user.ID,
	}).SignedString(jwtSecret)

	if err != nil {
		return nil, fmt.Errorf("token generation errorsx: %v", err)
	}
	return &viewmodel.LoginResponse{
		Token: token,
	}, nil
}

func (s *authService) GetUserProfile(ctx context.Context, email string) (*viewmodel.UserVM, error) {
	user, err := s.ur.GetUserProfile(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	c := &viewmodel.CurrentByAssetType{
		Asset: "",
		Price: 0,
	}
	vm := viewmodel.ToUserVM(user)
	vm.UserAssets = viewmodel.ToUserAssetVMs(user.Assets, c)
	vm.Reminders = viewmodel.ToReminderVMs(user.Reminders)

	return vm, nil
}
func (s *authService) UpdateUser(ctx context.Context, email string, vm *viewmodel.UpdateRequest) error {
	if vm == nil {
		return errors.New("update request is required")
	}

	email = strings.ToLower(email)
	user, err := s.ur.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	if vm.Name != "" {
		user.Name = vm.Name
	}
	if vm.Surname != "" {
		user.Surname = vm.Surname
	}
	if vm.Email != "" {
		vm.Email = strings.ToLower(vm.Email)
		exist, err := s.ur.ExistEmail(ctx, vm.Email)
		if err != nil {
			return err
		}
		if exist && vm.Email != email {
			return errors.New("email already exists")
		}
		user.Email = vm.Email
	}
	if vm.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(vm.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hashing errorsx: %v", err)
		}
		user.Password = string(hashed)
	}

	user.UpdatedAt = time.Now()
	return s.ur.Update(ctx, user)
}

func (s *authService) DeleteUser(ctx context.Context, email string) error {
	email = strings.ToLower(email)
	user, err := s.ur.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	return s.ur.Delete(ctx, user)
}
