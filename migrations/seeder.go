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

	// if err := seeds.ListProvinsiSeeder(db); err != nil {
	// 	return err
	// }

	// if err := seeds.ListKabupatenSeeder(db); err != nil {
	// 	return err
	// }

	// if err := seeds.ListChurchSeeder(db); err != nil {
	// 	return err
	// }

	if err := seeds.ListPersonSeeder(db); err != nil {
		return err
	}

	// Seed LifeGroup related data
	// First seed persons (needed for users foreign key)

	// Temporarily disabled user seeders due to foreign key constraint issues
	// Users and lifegroups can be created manually for testing
	// 
	// if err := seeds.LifeGroupUserSeeder(db); err != nil {
	//	return err
	// }
	// if err := seeds.LifeGroupSeeder(db); err != nil {
	//	return err
	// }
	// if err := seeds.LifeGroupMemberSeeder(db); err != nil {
	//	return err
	// }

	if err := seeds.DepartmentSeeder(db); err != nil {
		return err
	}

	// Seed pelayanan (requires departments to exist first)
	if err := seeds.PelayananSeeder(db); err != nil {
		return err
	}

	// Seed pelayanan assignments (requires persons, pelayanan, and churches to exist)
	if err := seeds.PelayananAssignmentSeeder(db); err != nil {
		return err
	}

	// if err := seeds.NotificationSeeder(db); err != nil {
	// 	return err
	// }

	return nil
}
