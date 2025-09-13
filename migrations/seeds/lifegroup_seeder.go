package seeds

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func LifeGroupSeeder(db *gorm.DB) error {
	jsonFile, err := os.Open("./migrations/json/lifegroups.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, _ := io.ReadAll(jsonFile)

	var listLifeGroup []entity.LifeGroup
	if err := json.Unmarshal(jsonData, &listLifeGroup); err != nil {
		return err
	}

	hasTable := db.Migrator().HasTable(&entity.LifeGroup{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.LifeGroup{}); err != nil {
			return err
		}
	}

	for _, data := range listLifeGroup {
		var lifeGroup entity.LifeGroup
		err := db.Where(&entity.LifeGroup{ID: data.ID}).First(&lifeGroup).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		isData := db.Find(&lifeGroup, "id = ?", data.ID).RowsAffected
		if isData == 0 {
			if err := db.Create(&data).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
