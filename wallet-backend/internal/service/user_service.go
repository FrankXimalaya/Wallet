package service

import (
	"errors"
	"wallet-backend/internal/model"
	"wallet-backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo      *repository.UserRepository
	walletService *WalletService
}

func NewUserService(userRepo *repository.UserRepository, walletService *WalletService) *UserService {
	return &UserService{
		userRepo:      userRepo,
		walletService: walletService,
	}
}

// Register 注册新用户并自动创建钱包
func (s *UserService) Register(username, email, password string) (*model.User, error) {
	// 检查用户是否已存在
	if _, err := s.userRepo.FindByUsername(username); err == nil {
		return nil, errors.New("username already exists")
	}

	if _, err := s.userRepo.FindByEmail(email); err == nil {
		return nil, errors.New("email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &model.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 自动创建主钱包
	wallet, err := s.walletService.CreateWallet(user.ID, "Main Wallet", true)
	if err != nil {
		return nil, err
	}

	// 更新用户的主地址
	user.Address = wallet.Address
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(username, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// TODO: 生成 JWT token
	token := "dummy-jwt-token"

	return token, nil
}

// GetUserByID 根据 ID 获取用户
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.FindByID(id)
}
