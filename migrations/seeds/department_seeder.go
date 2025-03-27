package seeds

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func DepartmentSeeder(db *gorm.DB) error {
	departments := []entity.Department{
		{
			Name:        "Pemusik",
			Description: "Departemen yang menangani pelayanan musik",
			ChurchID:    uuid.UUID{}, // Every Nation Church Jakarta
		},
		{
			Name:        "Multimedia",
			Description: "Departemen yang menangani pelayanan multimedia",
			ChurchID:    uuid.UUID{},
		},
		{
			Name:        "Dokumentasi",
			Description: "Departemen yang menangani pelayanan dokumentasi",
			ChurchID:    uuid.UUID{},
		},
		{
			Name:        "Usher",
			Description: "Departemen yang menangani pelayanan usher",
			ChurchID:    uuid.UUID{},
		},
		{
			Name:        "Pendidikan",
			Description: "Departemen yang menangani pelayanan pendidikan",
			ChurchID:    uuid.UUID{},
		},
		{
			Name:        "Konseling",
			Description: "Departemen yang menangani pelayanan konseling",
			ChurchID:    uuid.UUID{},
		},
	}

	for _, department := range departments {
		if err := db.Create(&department).Error; err != nil {
			return err
		}
	}

	return nil
}
