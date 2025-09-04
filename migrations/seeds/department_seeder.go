package seeds

import (
	"errors"

	"github.com/google/uuid"
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
		{
			Name:        "Lifegroup",
			Description: "Departemen yang menangani pelayanan lifegroup",
		},
	}

	hasTable := db.Migrator().HasTable(&entity.Department{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.Department{}); err != nil {
			return err
		}
	}

	for _, department := range departments {
		var existingDepartment entity.Department
		err := db.Where("name = ?", department.Name).First(&existingDepartment).Error
		
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		
		// Only create if department doesn't exist
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Generate UUID for the new department
			department.ID = uuid.New()
			if err := db.Create(&department).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
