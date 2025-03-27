package seeds

import (
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func ChurchSeeder(db *gorm.DB) error {
	churches := []entity.Church{
		{
			Name:     "Every Nation Church Jakarta",
			Address:  "Jl. Sudirman No. 1",
			City:     "Jakarta Pusat",
			Province: "DKI Jakarta",
		},
		{
			Name:     "Every Nation Church Bandung",
			Address:  "Jl. Merdeka No. 2",
			City:     "Bandung",
			Province: "Jawa Barat",
		},
		{
			Name:     "Every Nation Church Surabaya",
			Address:  "Jl. Pemuda No. 3",
			City:     "Surabaya",
			Province: "Jawa Timur",
		},
	}

	for _, church := range churches {
		if err := db.Create(&church).Error; err != nil {
			return err
		}
	}

	return nil
}
