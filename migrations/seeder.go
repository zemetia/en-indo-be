package migrations

import (
	"github.com/zemetia/en-indo-be/migrations/seeds"
	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	if err := seeds.ChurchSeeder(db); err != nil {
		return err
	}

	if err := seeds.ListUserSeeder(db); err != nil {
		return err
	}

	if err := seeds.RolePermissionSeeder(db); err != nil {
		return err
	}

	if err := seeds.DepartmentSeeder(db); err != nil {
		return err
	}

	if err := seeds.LifeGroupSeeder(db); err != nil {
		return err
	}

	if err := seeds.NotificationSeeder(db); err != nil {
		return err
	}

	return nil
}
