package migrations

import (
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Person{},
		&entity.Church{},
		&entity.Department{},
		&entity.LifeGroup{},
		&entity.Notification{},
	); err != nil {
		return err
	}

	return nil
}
