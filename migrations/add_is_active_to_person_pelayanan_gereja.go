package migrations

import (
	"gorm.io/gorm"
)

// AddIsActiveToPersonPelayananGereja adds is_active column to person_pelayanan_gereja table
type AddIsActiveToPersonPelayananGereja struct{}

func (m *AddIsActiveToPersonPelayananGereja) GetName() string {
	return "add_is_active_to_person_pelayanan_gereja"
}

func (m *AddIsActiveToPersonPelayananGereja) Up(db *gorm.DB) error {
	// Add is_active column to person_pelayanan_gereja table
	return db.Exec(`
		ALTER TABLE person_pelayanan_gereja
		ADD COLUMN is_active BOOLEAN DEFAULT true
		COMMENT 'Indicates if this pelayanan assignment is currently active'
	`).Error
}

func (m *AddIsActiveToPersonPelayananGereja) Down(db *gorm.DB) error {
	// Remove is_active column
	return db.Exec(`
		ALTER TABLE person_pelayanan_gereja
		DROP COLUMN is_active
	`).Error
}
