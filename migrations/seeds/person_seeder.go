package seeds

import (
	"encoding/json"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func ListPersonSeeder(db *gorm.DB) error {
	jsonFile, err := os.Open("./migrations/json/person.json")
	if err != nil {
		return err
	}

	jsonData, _ := io.ReadAll(jsonFile)

	var listPerson []entity.Person
	if err := json.Unmarshal(jsonData, &listPerson); err != nil {
		return err
	}

	hasTable := db.Migrator().HasTable(&entity.Person{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.Person{}); err != nil {
			return err
		}
	}

	for _, data := range listPerson {
		// var data entity.Person
		// err := db.Where(&entity.Person{Email: data.Email}).First(&Person).Error
		// if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return err
		// }

		// isData := db.Find(&Person, "email = ?", data.Email).RowsAffected
		// if isData == 0 {
		data.ID = uuid.New()

		if err := db.Create(&data).Error; err != nil {
			return err
		}
		// }
	}

	return nil
}
