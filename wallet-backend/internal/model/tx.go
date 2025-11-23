package model

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	TxHash    string         `gorm:"uniqueIndex;not null" json:"tx_hash"`
	From      string         `gorm:"index;not null" json:"from"`
	To        string         `gorm:"index;not null" json:"to"`
	Amount    string         `gorm:"not null" json:"amount"` // 使用字符串存储，避免精度问题
	GasPrice  string         `json:"gas_price"`
	GasUsed   uint64         `json:"gas_used"`
	Status    string         `gorm:"index;not null" json:"status"` // pending, confirmed, failed
	TxType    string         `gorm:"not null" json:"tx_type"` // deposit, withdraw, transfer
	BlockNum  uint64         `json:"block_num"`
	Timestamp int64          `gorm:"index" json:"timestamp"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
