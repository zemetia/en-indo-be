package migrations

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AddEventDepartmentsTable creates the junction table for event-department many-to-many relationship
func AddEventDepartmentsTable(db *gorm.DB) error {
	// Check if table already exists
	if db.Migrator().HasTable("event_departments") {
		fmt.Println("event_departments table already exists, skipping creation")
		return nil
	}

	// Create the junction table
	if err := db.Migrator().CreateTable(&EventDepartmentForMigration{}); err != nil {
		return fmt.Errorf("failed to create event_departments table: %v", err)
	}

	// Create indexes for better query performance
	if !db.Migrator().HasIndex(&EventDepartmentForMigration{}, "idx_event_departments_event_id") {
		if err := db.Migrator().CreateIndex(&EventDepartmentForMigration{}, "idx_event_departments_event_id"); err != nil {
			return fmt.Errorf("failed to create event_id index: %v", err)
		}
	}

	if !db.Migrator().HasIndex(&EventDepartmentForMigration{}, "idx_event_departments_department_id") {
		if err := db.Migrator().CreateIndex(&EventDepartmentForMigration{}, "idx_event_departments_department_id"); err != nil {
			return fmt.Errorf("failed to create department_id index: %v", err)
		}
	}

	fmt.Println("event_departments table created successfully")
	return nil
}

// EventDepartmentForMigration represents the junction table structure for migration
type EventDepartmentForMigration struct {
	ID           uuid.UUID  `gorm:"type:char(36);primary_key" json:"id"`
	EventID      uuid.UUID  `gorm:"type:char(36);not null;index:idx_event_departments_event_id" json:"event_id"`
	DepartmentID uuid.UUID  `gorm:"type:char(36);not null;index:idx_event_departments_department_id" json:"department_id"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (EventDepartmentForMigration) TableName() string {
	return "event_departments"
}
