package seeds

import (
	"errors"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func PelayananSeeder(db *gorm.DB) error {
	// First, we need to get the department IDs
	var lifegroupDept entity.Department
	if err := db.Where("name = ?", "Lifegroup").First(&lifegroupDept).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Lifegroup department not found. Please ensure DepartmentSeeder is run first")
		}
		return err
	}

	pelayanan := []entity.Pelayanan{
		{
			ID:           uuid.MustParse("11111111-1111-1111-1111-111111111001"),
			Pelayanan:    "Lifegroup",
			Description:  "Person In Charge untuk mengelola lifegroup",
			DepartmentID: lifegroupDept.ID,
			IsPic:        true,
		},
		{
			ID:           uuid.MustParse("11111111-1111-1111-1111-111111111002"),
			Pelayanan:    "Lifegroup Leader",
			Description:  "Pemimpin lifegroup",
			DepartmentID: lifegroupDept.ID,
			IsPic:        false,
		},
		{
			ID:           uuid.MustParse("11111111-1111-1111-1111-111111111003"),
			Pelayanan:    "Lifegroup Co-Leader",
			Description:  "Wakil pemimpin lifegroup",
			DepartmentID: lifegroupDept.ID,
			IsPic:        false,
		},
	}

	hasTable := db.Migrator().HasTable(&entity.Pelayanan{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.Pelayanan{}); err != nil {
			return err
		}
	}

	for _, data := range pelayanan {
		var existingPelayanan entity.Pelayanan
		err := db.Where(&entity.Pelayanan{ID: data.ID}).First(&existingPelayanan).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		isData := db.Find(&existingPelayanan, "id = ?", data.ID).RowsAffected
		if isData == 0 {
			if err := db.Create(&data).Error; err != nil {
				return err
			}
		}
	}

	return nil
}