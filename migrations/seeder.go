package migrations

import (
	"github.com/zemetia/en-indo-be/migrations/seeds"
	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {

	// if err := seeds.ListUserSeeder(db); err != nil {
	// 	return err
	// }

	// db.Exec("SET FOREIGN_KEY_CHECKS = 0;")

	if err := seeds.ListProvinsiSeeder(db); err != nil {
		return err
	}

	if err := seeds.ListKabupatenSeeder(db); err != nil {
		return err
	}

	if err := seeds.ListChurchSeeder(db); err != nil {
		return err
	}

	if err := seeds.ListPersonSeeder(db); err != nil {
		return err
	}

	// if err := seeds.ListUserSeeder(db); err != nil {
	// 	return err
	// }

	// if err := seeds.DepartmentSeeder(db); err != nil {
	// 	return err
	// }

	// if err := seeds.NotificationSeeder(db); err != nil {
	// 	return err
	// }

	return nil
}
