package migrations

import (
	"gorm.io/gorm"
)

// AddEventParticipantsTable creates the event_participants table with proper indexes
func AddEventParticipantsTable(db *gorm.DB) error {
	// Check if table exists, if not, nothing to do (AutoMigrate will create it)
	if !db.Migrator().HasTable("event_participants") {
		return nil
	}

	// Check if unique index already exists
	var indexExists bool
	err := db.Raw(`
		SELECT COUNT(*) > 0
		FROM information_schema.statistics
		WHERE table_schema = DATABASE()
		  AND table_name = 'event_participants'
		  AND index_name = 'idx_event_participant_unique'
	`).Scan(&indexExists).Error

	if err != nil {
		return err
	}

	// Create unique composite index if it doesn't exist
	// Index: event_id + occurrence_date + participant_type + participant_id must be unique
	if !indexExists {
		if err := db.Exec(`
			CREATE UNIQUE INDEX idx_event_participant_unique
			ON event_participants(event_id, occurrence_date, participant_type, participant_id)
		`).Error; err != nil {
			return err
		}
	}

	// Create other indexes (ignore errors if they already exist)
	// Create index for querying by event and occurrence
	db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_event_occurrence
		ON event_participants(event_id, occurrence_date)
	`)

	// Create index for querying by participant
	db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_participant
		ON event_participants(participant_type, participant_id)
	`)

	// Create index for attendance status queries
	db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_attendance_status
		ON event_participants(event_id, occurrence_date, attendance_status)
	`)

	return nil
}
