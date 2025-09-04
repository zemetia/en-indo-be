package seeds

import (
	"errors"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func PelayananAssignmentSeeder(db *gorm.DB) error {
	// Define the assignment for yoelsit@gmail.com as Lifegroup PIC
	assignments := []entity.PersonPelayananGereja{
		{
			ID:          uuid.MustParse("22222222-2222-2222-2222-222222222001"),
			PersonID:    uuid.MustParse("77285a15-dcd1-4c50-bd45-46efe6aeea3a"), // yoelsit@gmail.com person ID
			PelayananID: uuid.MustParse("11111111-1111-1111-1111-111111111001"), // Lifegroup PIC pelayanan ID
			ChurchID:    uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"), // Church ID from lifegroups.json
		},
	}

	hasTable := db.Migrator().HasTable(&entity.PersonPelayananGereja{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.PersonPelayananGereja{}); err != nil {
			return err
		}
	}

	for _, data := range assignments {
		var existingAssignment entity.PersonPelayananGereja
		err := db.Where(&entity.PersonPelayananGereja{ID: data.ID}).First(&existingAssignment).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// Check if assignment already exists
		isData := db.Find(&existingAssignment, "id = ?", data.ID).RowsAffected
		if isData == 0 {
			// Also check if the same person-pelayanan-church combination exists
			var duplicateCheck entity.PersonPelayananGereja
			duplicateExists := db.Where("person_id = ? AND pelayanan_id = ? AND church_id = ?", 
				data.PersonID, data.PelayananID, data.ChurchID).First(&duplicateCheck).Error
			
			if errors.Is(duplicateExists, gorm.ErrRecordNotFound) {
				// No duplicate found, safe to create
				if err := db.Create(&data).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}