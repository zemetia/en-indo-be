package entity

import (
	"time"

	"gorm.io/gorm"
)

type Timestamp struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type TimestampHardDelete struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Authorization struct {
	Token string `json:"token"`
	Yoel  string `json:"yoel"`
}
