package repository

import (
	"wallet-backend/internal/model"

	"gorm.io/gorm"
)

type TxRepository struct {
	db *gorm.DB
}

func NewTxRepository(db *gorm.DB) *TxRepository {
	return &TxRepository{db: db}
}

func (r *TxRepository) Create(tx *model.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *TxRepository) FindByID(id uint) (*model.Transaction, error) {
	var tx model.Transaction
	err := r.db.First(&tx, id).Error
	return &tx, err
}

func (r *TxRepository) FindByHash(hash string) (*model.Transaction, error) {
	var tx model.Transaction
	err := r.db.Where("tx_hash = ?", hash).First(&tx).Error
	return &tx, err
}

func (r *TxRepository) FindByUserID(userID uint) ([]*model.Transaction, error) {
	var txs []*model.Transaction
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&txs).Error
	return txs, err
}

func (r *TxRepository) FindByAddress(address string) ([]*model.Transaction, error) {
	var txs []*model.Transaction
	err := r.db.Where("\"from\" = ? OR \"to\" = ?", address, address).Order("created_at DESC").Find(&txs).Error
	return txs, err
}

func (r *TxRepository) Update(tx *model.Transaction) error {
	return r.db.Save(tx).Error
}

func (r *TxRepository) Delete(id uint) error {
	return r.db.Delete(&model.Transaction{}, id).Error
}

// AutoMigrate 自动迁移数据库
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Wallet{},
		&model.Transaction{},
	)
}
