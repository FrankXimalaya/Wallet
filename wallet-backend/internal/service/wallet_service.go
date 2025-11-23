package service

import (
	"context"
	"wallet-backend/internal/eth"
	"wallet-backend/internal/model"
	"wallet-backend/internal/repository"

	"github.com/ethereum/go-ethereum/common"
)

type WalletService struct {
	walletRepo *repository.WalletRepository
	ethClient  *eth.Client
}

func NewWalletService(walletRepo *repository.WalletRepository, ethClient *eth.Client) *WalletService {
	return &WalletService{
		walletRepo: walletRepo,
		ethClient:  ethClient,
	}
}

// CreateWallet 创建新钱包
func (s *WalletService) CreateWallet(userID uint, label string, isMain bool) (*model.Wallet, error) {
	// 生成新钱包
	address, privateKey, err := eth.GenerateWallet()
	if err != nil {
		return nil, err
	}

	// TODO: 加密私钥后再存储
	wallet := &model.Wallet{
		UserID:     userID,
		Address:    address,
		PrivateKey: privateKey,
		Label:      label,
		IsMain:     isMain,
	}

	if err := s.walletRepo.Create(wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

// GetBalance 获取钱包余额
func (s *WalletService) GetBalance(userID uint) (string, error) {
	wallet, err := s.walletRepo.FindMainWalletByUserID(userID)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	balance, err := s.ethClient.GetBalance(ctx, common.HexToAddress(wallet.Address))
	if err != nil {
		return "", err
	}

	return balance.String(), nil
}

// GetWalletByAddress 根据地址获取钱包
func (s *WalletService) GetWalletByAddress(address string) (*model.Wallet, error) {
	return s.walletRepo.FindByAddress(address)
}

// GetUserWallets 获取用户的所有钱包
func (s *WalletService) GetUserWallets(userID uint) ([]*model.Wallet, error) {
	return s.walletRepo.FindByUserID(userID)
}

// GetMainWallet 获取用户的主钱包
func (s *WalletService) GetMainWallet(userID uint) (*model.Wallet, error) {
	return s.walletRepo.FindMainWalletByUserID(userID)
}
