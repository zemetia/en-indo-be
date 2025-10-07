package migrations

import (
	"fmt"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

// AddEventChurchesTable creates the junction table for event-church many-to-many relationship
func AddEventChurchesTable(db *gorm.DB) error {
	// Check if table already exists
	if db.Migrator().HasTable("event_churches") {
		fmt.Println("event_churches table already exists, skipping creation")
		return nil
	}

	// Create the junction table
	if err := db.Migrator().CreateTable(&EventChurchForMigration{}); err != nil {
		return fmt.Errorf("failed to create event_churches table: %v", err)
	}

	// Create indexes for better query performance
	if !db.Migrator().HasIndex(&EventChurchForMigration{}, "idx_event_churches_event_id") {
		if err := db.Migrator().CreateIndex(&EventChurchForMigration{}, "idx_event_churches_event_id"); err != nil {
			return fmt.Errorf("failed to create event_id index: %v", err)
		}
	}

	if !db.Migrator().HasIndex(&EventChurchForMigration{}, "idx_event_churches_church_id") {
		if err := db.Migrator().CreateIndex(&EventChurchForMigration{}, "idx_event_churches_church_id"); err != nil {
			return fmt.Errorf("failed to create church_id index: %v", err)
		}
	}

	fmt.Println("event_churches table created successfully")
	return nil
}

// EventChurchForMigration represents the junction table structure for migration
type EventChurchForMigration struct {
	EventID  uuid.UUID `gorm:"type:char(36);not null;index:idx_event_churches_event_id;primaryKey" json:"event_id"`
	ChurchID uuid.UUID `gorm:"type:char(36);not null;index:idx_event_churches_church_id;primaryKey" json:"church_id"`
}

func (EventChurchForMigration) TableName() string {
	return "event_churches"
}