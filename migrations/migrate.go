package migrations

import (
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// Fix visitor kabupaten foreign key constraint before AutoMigrate
	if err := FixVisitorKabupatenFK(db); err != nil {
		return err
	}

	// Clean up duplicate and invalid participant records before AutoMigrate
	if err := CleanupDuplicateParticipants(db); err != nil {
		return err
	}

	// Fix polymorphic foreign key constraints (must run BEFORE AutoMigrate)
	if err := FixPolymorphicFKConstraints(db); err != nil {
		return err
	}

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
		&entity.EventType{},
		&entity.EventPIC{},
		&entity.EventPICRole{},
		&entity.EventPICHistory{},
		&entity.DiscipleshipJourney{},
		&entity.Lagu{},
		&entity.Visitor{},
		&entity.VisitorInformation{},
		&entity.Ketersediaan{},
		&entity.EventDepartment{},
		&entity.EventParticipant{},
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

	// Drop leader_id and co_leader_id columns from life_groups table
	if err := DropLifeGroupLeaderColumns(db); err != nil {
		return err
	}

	// Add expected participant fields to events table
	if err := AddEventParticipantFields(db); err != nil {
		return err
	}

	// Add event_churches junction table for event-church many-to-many relationship
	if err := AddEventChurchesTable(db); err != nil {
		return err
	}

	// Add event_types table
	if err := AddEventTypesTable(db); err != nil {
		return err
	}

	// Drop icon column from event_types table
	if err := DropEventTypesIconColumn(db); err != nil {
		return err
	}

	// Add event_departments junction table for event-department many-to-many relationship
	if err := AddEventDepartmentsTable(db); err != nil {
		return err
	}

	// Add event_participants table with proper indexes
	// Note: CleanupDuplicateParticipants is now called before AutoMigrate
	if err := AddEventParticipantsTable(db); err != nil {
		return err
	}

	// Add qr_scan_timestamp column to event_participants table
	migration := &AddQRTimestampToParticipants{}
	if err := migration.Up(db); err != nil {
		return err
	}

	// Add is_active column to person_pelayanan_gereja table
	isActiveMigration := &AddIsActiveToPersonPelayananGereja{}
	if err := isActiveMigration.Up(db); err != nil {
		return err
	}

	return nil
}
