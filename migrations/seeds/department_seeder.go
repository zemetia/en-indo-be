package seeds

import (
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func DepartmentSeeder(db *gorm.DB) error {
	departments := []entity.Department{
		{
			Name:        "Pemusik",
			Description: "Departemen yang menangani pelayanan musik",
		},
		{
			Name:        "Multimedia",
			Description: "Departemen yang menangani pelayanan multimedia",
		},
		{
			Name:        "Dokumentasi",
			Description: "Departemen yang menangani pelayanan dokumentasi",
		},
		{
			Name:        "Usher",
			Description: "Departemen yang menangani pelayanan usher",
		},
		{
			Name:        "Pendidikan",
			Description: "Departemen yang menangani pelayanan pendidikan",
		},
		{
			Name:        "Konseling",
			Description: "Departemen yang menangani pelayanan konseling",
		},
	}

	for _, department := range departments {
		if err := db.Create(&department).Error; err != nil {
			return err
		}
	}

	return nil
}
