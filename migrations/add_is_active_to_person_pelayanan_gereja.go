package migrations

import (
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

// AddIsActiveToPersonPelayananGereja adds is_active column to person_pelayanan_gereja table
type AddIsActiveToPersonPelayananGereja struct{}

func (m *AddIsActiveToPersonPelayananGereja) GetName() string {
	return "add_is_active_to_person_pelayanan_gereja"
}

func (m *AddIsActiveToPersonPelayananGereja) Up(db *gorm.DB) error {
	// Add is_active column to person_pelayanan_gereja table
	if db.Migrator().HasColumn(&entity.PersonPelayananGereja{}, "IsActive") {
		return nil
	}
	return db.Migrator().AddColumn(&entity.PersonPelayananGereja{}, "IsActive")
}

func (m *AddIsActiveToPersonPelayananGereja) Down(db *gorm.DB) error {
	// Remove is_active column
	if !db.Migrator().HasColumn(&entity.PersonPelayananGereja{}, "IsActive") {
		return nil
	}
	return db.Migrator().DropColumn(&entity.PersonPelayananGereja{}, "IsActive")
}
