package seeds

import (
	"encoding/json"
	"io"
	"os"

	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func ListKabupatenSeeder(db *gorm.DB) error {
	jsonFile, err := os.Open("./migrations/json/kabupaten.json")
	if err != nil {
		return err
	}

	jsonData, _ := io.ReadAll(jsonFile)

	var listKabupaten []entity.Kabupaten
	if err := json.Unmarshal(jsonData, &listKabupaten); err != nil {
		return err
	}

	hasTable := db.Migrator().HasTable(&entity.Kabupaten{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.Kabupaten{}); err != nil {
			return err
		}
	}

	for _, data := range listKabupaten {
		// var kabupaten entity.Kabupaten
		// err := db.Where(&entity.Kabupaten{Email: data.Email}).First(&kabupaten).Error
		// if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return err
		// }

		if err := db.Create(&data).Error; err != nil {
			return err
		}

	}

	return nil
}
