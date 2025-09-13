package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/helpers"
	"gorm.io/gorm"
)

type User struct {
	ID                        uuid.UUID  `gorm:"type:char(36);primary_key;unique" json:"id"`
	Email                     string     `gorm:"type:varchar(100);unique;not null"`
	Password                  string     `gorm:"type:varchar(100);not null"`
	ImageUrl                  string     `gorm:"type:text" json:"image_url"`
	IsActive                  bool       `gorm:"type:boolean;not null;default:true" json:"is_active"`
	HasChangedDefaultPassword bool       `gorm:"type:boolean;not null;default:false" json:"has_changed_default_password"`
	LastLoginAt               *time.Time `gorm:"type:timestamp" json:"last_login_at"`
	PersonID                  uuid.UUID  `gorm:"type:char(36);not null;unique"`
	Person                    Person     `gorm:"foreignKey:PersonID"`

	Timestamp
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var err error
	u.Password, err = helpers.HashPassword(u.Password)
	if err != nil {
		return err
	}
	return nil
}
