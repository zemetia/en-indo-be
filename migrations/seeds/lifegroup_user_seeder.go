package seeds

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func LifeGroupUserSeeder(db *gorm.DB) error {
	jsonFile, err := os.Open("./migrations/json/lifegroup_leader_users.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, _ := io.ReadAll(jsonFile)

	var listUser []entity.User
	if err := json.Unmarshal(jsonData, &listUser); err != nil {
		return err
	}

	hasTable := db.Migrator().HasTable(&entity.User{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.User{}); err != nil {
			return err
		}
	}

	for _, data := range listUser {
		var user entity.User
		err := db.Where(&entity.User{Email: data.Email}).First(&user).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		isData := db.Find(&user, "email = ?", data.Email).RowsAffected
		if isData == 0 {
			if err := db.Create(&data).Error; err != nil {
				return err
			}
		}
	}

	return nil
}