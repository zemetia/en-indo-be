package migrations

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AddEventTypesTable creates the event_types table
func AddEventTypesTable(db *gorm.DB) error {
	// Check if table already exists
	if db.Migrator().HasTable("event_types") {
		fmt.Println("event_types table already exists, skipping creation")
		return nil
	}

	// Create the table
	if err := db.Migrator().CreateTable(&EventTypeForMigration{}); err != nil {
		return fmt.Errorf("failed to create event_types table: %v", err)
	}

	// Create indexes for better query performance
	if !db.Migrator().HasIndex(&EventTypeForMigration{}, "idx_event_types_is_active") {
		if err := db.Migrator().CreateIndex(&EventTypeForMigration{}, "idx_event_types_is_active"); err != nil {
			return fmt.Errorf("failed to create is_active index: %v", err)
		}
	}

	if !db.Migrator().HasIndex(&EventTypeForMigration{}, "idx_event_types_name") {
		if err := db.Migrator().CreateIndex(&EventTypeForMigration{}, "idx_event_types_name"); err != nil {
			return fmt.Errorf("failed to create name index: %v", err)
		}
	}

	fmt.Println("event_types table created successfully")
	return nil
}

// EventTypeForMigration represents the event_types table structure for migration
type EventTypeForMigration struct {
	ID          uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	Name        string    `gorm:"type:varchar(255);unique;not null;index:idx_event_types_name" json:"name"`
	DisplayName string    `gorm:"type:varchar(255);not null" json:"display_name"`
	Description string    `gorm:"type:text" json:"description"`
	Color       string    `gorm:"type:varchar(7)" json:"color"`
	Icon        string    `gorm:"type:varchar(50)" json:"icon"`
	IsActive    bool      `gorm:"default:true;not null;index:idx_event_types_is_active" json:"is_active"`
	SortOrder   int       `gorm:"default:0;not null" json:"sort_order"`
}

func (EventTypeForMigration) TableName() string {
	return "event_types"
}
