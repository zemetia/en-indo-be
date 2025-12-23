package migrations

import (
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

// AddQRTimestampToParticipants adds qr_scan_timestamp column to event_participants table
type AddQRTimestampToParticipants struct{}

func (m *AddQRTimestampToParticipants) GetName() string {
	return "add_qr_timestamp_to_participants"
}

func (m *AddQRTimestampToParticipants) Up(db *gorm.DB) error {
	// Check if column already exists
	if db.Migrator().HasColumn(&entity.EventParticipant{}, "qr_scan_timestamp") {
		return nil
	}

	// Add qr_scan_timestamp column to event_participants table
	return db.Exec(`
		ALTER TABLE event_participants
		ADD COLUMN qr_scan_timestamp TIMESTAMP NULL
		COMMENT 'Timestamp from QR code when participant scanned for attendance'
	`).Error
}

func (m *AddQRTimestampToParticipants) Down(db *gorm.DB) error {
	// Remove qr_scan_timestamp column
	return db.Exec(`
		ALTER TABLE event_participants
		DROP COLUMN qr_scan_timestamp
	`).Error
}
