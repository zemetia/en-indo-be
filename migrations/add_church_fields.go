package migrations

import (
	"fmt"
	"gorm.io/gorm"
)

// AddChurchFields adds church_code, latitude, and longitude fields to churches table
func AddChurchFields(db *gorm.DB) error {
	// Add church_code column
	if !db.Migrator().HasColumn(&ChurchForMigration{}, "church_code") {
		if err := db.Migrator().AddColumn(&ChurchForMigration{}, "church_code"); err != nil {
			return fmt.Errorf("failed to add church_code column: %v", err)
		}
	}

	// Add latitude column
	if !db.Migrator().HasColumn(&ChurchForMigration{}, "latitude") {
		if err := db.Migrator().AddColumn(&ChurchForMigration{}, "latitude"); err != nil {
			return fmt.Errorf("failed to add latitude column: %v", err)
		}
	}

	// Add longitude column
	if !db.Migrator().HasColumn(&ChurchForMigration{}, "longitude") {
		if err := db.Migrator().AddColumn(&ChurchForMigration{}, "longitude"); err != nil {
			return fmt.Errorf("failed to add longitude column: %v", err)
		}
	}

	// Add unique index for church_code
	if !db.Migrator().HasIndex(&ChurchForMigration{}, "idx_church_code") {
		if err := db.Migrator().CreateIndex(&ChurchForMigration{}, "idx_church_code"); err != nil {
			return fmt.Errorf("failed to create church_code index: %v", err)
		}
	}

	// Add index for coordinates for location-based queries
	if !db.Migrator().HasIndex(&ChurchForMigration{}, "idx_coordinates") {
		if err := db.Migrator().CreateIndex(&ChurchForMigration{}, "idx_coordinates"); err != nil {
			return fmt.Errorf("failed to create coordinates index: %v", err)
		}
	}

	return nil
}

// ChurchForMigration is a temporary struct for migration purposes
type ChurchForMigration struct {
	ChurchCode string  `gorm:"type:varchar(10);uniqueIndex:idx_church_code;null" json:"church_code"`
	Latitude   float64 `gorm:"type:decimal(10,8);index:idx_coordinates;null" json:"latitude"`
	Longitude  float64 `gorm:"type:decimal(11,8);index:idx_coordinates;null" json:"longitude"`
}

func (ChurchForMigration) TableName() string {
	return "churches"
}
