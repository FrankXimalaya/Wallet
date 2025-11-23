package model

import (
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	Address    string         `gorm:"uniqueIndex;not null" json:"address"`
	PrivateKey string         `gorm:"not null" json:"-"` // 加密存储，不在 JSON 中显示
	Label      string         `json:"label"` // 钱包标签
	IsMain     bool           `gorm:"default:false" json:"is_main"` // 是否为主钱包
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
