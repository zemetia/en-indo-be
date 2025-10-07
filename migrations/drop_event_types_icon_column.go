package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// DropEventTypesIconColumn drops the icon column from event_types table
func DropEventTypesIconColumn(db *gorm.DB) error {
	// Check if table exists
	if !db.Migrator().HasTable("event_types") {
		fmt.Println("event_types table does not exist, skipping icon column drop")
		return nil
	}

	// Check if icon column exists
	if !db.Migrator().HasColumn(&EventTypeForMigration{}, "icon") {
		fmt.Println("icon column does not exist in event_types table, skipping drop")
		return nil
	}

	// Drop the icon column
	if err := db.Migrator().DropColumn(&EventTypeForMigration{}, "icon"); err != nil {
		return fmt.Errorf("failed to drop icon column from event_types table: %v", err)
	}

	fmt.Println("icon column dropped successfully from event_types table")
	return nil
}
