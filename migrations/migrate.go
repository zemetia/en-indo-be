package migrations

import (
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// First run the auto migration for all entities
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Person{},
		&entity.Church{},
		&entity.Department{},
		&entity.LifeGroup{},
		&entity.LifeGroupPersonMember{},
		&entity.LifeGroupVisitorMember{},
		&entity.Notification{},
		&entity.Kabupaten{},
		&entity.Provinsi{},
		&entity.PersonPelayananGereja{},
		&entity.Pelayanan{},
		&entity.RecurrenceRule{},
		&entity.RecurrenceException{},
		&entity.Event{},
		&entity.DiscipleshipJourney{},
		&entity.Lagu{},
		&entity.Visitor{},
		&entity.VisitorInformation{},
	); err != nil {
		return err
	}

	// Run custom migration to drop is_verified column
	if err := DropIsVerifiedColumn(db); err != nil {
		return err
	}

	// Run custom migration to handle PIC field changes
	// Temporarily disabled due to is_pic column issue
	// if err := MigratePicField(db); err != nil {
	//	return err
	// }

	// Run custom migration to add church fields
	if err := AddChurchFields(db); err != nil {
		return err
	}

	// Remove deleted_at column from visitor_informations table (hard delete)
	if err := RemoveVisitorInformationDeletedAt(db); err != nil {
		return err
	}

	return nil
}
