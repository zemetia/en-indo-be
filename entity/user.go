package entity

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/helpers"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID `gorm:"type:char(36);primary_key;unique" json:"id"`
	Email      string    `gorm:"type:varchar(100);unique;not null"`
	Password   string    `gorm:"type:varchar(100);not null"`
	ImageUrl   string    `gorm:"type:text" json:"image_url"`
	IsVerified bool      `gorm:"type:boolean;not null;default:false" json:"is_verified"`
	PersonID   uuid.UUID `gorm:"type:char(36);not null;unique"`
	Person     Person    `gorm:"foreignKey:PersonID"`

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
