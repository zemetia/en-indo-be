package migrations

import (
	"gorm.io/gorm"
)

func RemoveVisitorInformationDeletedAt(db *gorm.DB) error {
	// Check if column exists before trying to drop it
	if db.Migrator().HasColumn("visitor_informations", "deleted_at") {
		// Drop the deleted_at column from visitor_informations table
		if err := db.Migrator().DropColumn("visitor_informations", "deleted_at"); err != nil {
			return err
		}
	}
	
	return nil
}