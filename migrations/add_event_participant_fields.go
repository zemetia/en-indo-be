package migrations

import (
	"fmt"
	"gorm.io/gorm"
)

// AddEventParticipantFields adds expected participant count fields to events table
func AddEventParticipantFields(db *gorm.DB) error {
	// Add expected_participants column
	if !db.Migrator().HasColumn(&EventForParticipantMigration{}, "expected_participants") {
		if err := db.Migrator().AddColumn(&EventForParticipantMigration{}, "expected_participants"); err != nil {
			return fmt.Errorf("failed to add expected_participants column: %v", err)
		}
	}

	// Add expected_adults column
	if !db.Migrator().HasColumn(&EventForParticipantMigration{}, "expected_adults") {
		if err := db.Migrator().AddColumn(&EventForParticipantMigration{}, "expected_adults"); err != nil {
			return fmt.Errorf("failed to add expected_adults column: %v", err)
		}
	}

	// Add expected_youth column
	if !db.Migrator().HasColumn(&EventForParticipantMigration{}, "expected_youth") {
		if err := db.Migrator().AddColumn(&EventForParticipantMigration{}, "expected_youth"); err != nil {
			return fmt.Errorf("failed to add expected_youth column: %v", err)
		}
	}

	// Add expected_kids column
	if !db.Migrator().HasColumn(&EventForParticipantMigration{}, "expected_kids") {
		if err := db.Migrator().AddColumn(&EventForParticipantMigration{}, "expected_kids"); err != nil {
			return fmt.Errorf("failed to add expected_kids column: %v", err)
		}
	}

	return nil
}

// EventForParticipantMigration is a temporary struct for migration purposes
type EventForParticipantMigration struct {
	ExpectedParticipants int `gorm:"default:0"`
	ExpectedAdults      int `gorm:"default:0"`
	ExpectedYouth       int `gorm:"default:0"`
	ExpectedKids        int `gorm:"default:0"`
}

func (EventForParticipantMigration) TableName() string {
	return "events"
}