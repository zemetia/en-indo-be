package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

// FixPolymorphicFKConstraints drops incorrect foreign key constraints for polymorphic relationships
// Polymorphic relationships should NOT have database-level foreign keys because a single column
// cannot reference multiple tables simultaneously
func FixPolymorphicFKConstraints(db *gorm.DB) error {
	// Check if table exists
	if !db.Migrator().HasTable("event_participants") {
		return nil
	}

	fmt.Println("Fixing polymorphic foreign key constraints on event_participants...")

	// Drop the incorrect foreign key constraints that GORM created
	// These constraints cause errors because participant_id can point to either people OR visitors
	db.Exec(`ALTER TABLE event_participants DROP FOREIGN KEY IF EXISTS fk_event_participants_person`)
	db.Exec(`ALTER TABLE event_participants DROP FOREIGN KEY IF EXISTS fk_event_participants_visitor`)

	// Keep the Event foreign key (this is valid)
	// Keep the RecordedBy foreign key (this is valid)

	fmt.Println("✓ Dropped polymorphic foreign key constraints (person/visitor)")
	fmt.Println("✓ Application-level validation will be used instead")

	return nil
}
