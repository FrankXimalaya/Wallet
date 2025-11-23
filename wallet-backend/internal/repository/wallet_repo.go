package repository

import (
	"wallet-backend/internal/model"

	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(wallet *model.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *WalletRepository) FindByID(id uint) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.First(&wallet, id).Error
	return &wallet, err
}

func (r *WalletRepository) FindByAddress(address string) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.Where("address = ?", address).First(&wallet).Error
	return &wallet, err
}

func (r *WalletRepository) FindByUserID(userID uint) ([]*model.Wallet, error) {
	var wallets []*model.Wallet
	err := r.db.Where("user_id = ?", userID).Find(&wallets).Error
	return wallets, err
}

func (r *WalletRepository) FindMainWalletByUserID(userID uint) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.Where("user_id = ? AND is_main = ?", userID, true).First(&wallet).Error
	return &wallet, err
}

func (r *WalletRepository) FindAll() ([]*model.Wallet, error) {
	var wallets []*model.Wallet
	err := r.db.Find(&wallets).Error
	return wallets, err
}

func (r *WalletRepository) Update(wallet *model.Wallet) error {
	return r.db.Save(wallet).Error
}

func (r *WalletRepository) Delete(id uint) error {
	return r.db.Delete(&model.Wallet{}, id).Error
}
