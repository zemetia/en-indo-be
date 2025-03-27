package entity

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/helpers"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Email         string         `gorm:"type:varchar(100);unique;not null"`
	Password      string         `gorm:"type:varchar(100);not null"`
	ImageUrl      string         `json:"image_url"`
	IsVerified    bool           `json:"is_verified"`
	Churches      []Church       `gorm:"many2many:user_churches;"`
	Roles         []Role         `gorm:"many2many:user_has_roles;"`
	Departments   []Department   `gorm:"many2many:user_departments;"`
	Notifications []Notification `gorm:"foreignKey:UserID"`

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
