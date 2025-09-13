package seeds

import (
	"encoding/json"
	"io"
	"os"

	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func ListProvinsiSeeder(db *gorm.DB) error {
	jsonFile, err := os.Open("./migrations/json/provinsi.json")
	if err != nil {
		return err
	}

	jsonData, _ := io.ReadAll(jsonFile)

	var listProvinsi []entity.Provinsi
	if err := json.Unmarshal(jsonData, &listProvinsi); err != nil {
		return err
	}

	hasTable := db.Migrator().HasTable(&entity.Provinsi{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.Provinsi{}); err != nil {
			return err
		}
	}

	for _, data := range listProvinsi {
		// var kabupaten entity.Provinsi
		// err := db.Where(&entity.Provinsi{Email: data.Email}).First(&kabupaten).Error
		// if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return err
		// }

		if err := db.Create(&data).Error; err != nil {
			return err
		}

	}

	return nil
}
