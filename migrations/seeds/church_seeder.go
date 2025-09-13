package seeds

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func ListChurchSeeder(db *gorm.DB) error {
	jsonFile, err := os.Open("./migrations/json/church.json")
	if err != nil {
		return err
	}

	jsonData, _ := io.ReadAll(jsonFile)

	var listChurch []entity.Church
	if err := json.Unmarshal(jsonData, &listChurch); err != nil {
		return err
	}

	hasTable := db.Migrator().HasTable(&entity.Church{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.Church{}); err != nil {
			return err
		}
	}

	for _, data := range listChurch {
		// var Church entity.Church
		// err := db.Where(&entity.Church{Email: data.Email}).First(&Church).Error
		// if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return err
		// }

		if err := db.Create(&data).Error; err != nil {
			return fmt.Errorf("error creating Church: %s | Error: %v", data.Name, err)
		}

	}

	return nil
}
