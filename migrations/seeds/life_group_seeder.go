package seeds

import (
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func LifeGroupSeeder(db *gorm.DB) error {
	// Ambil beberapa user untuk dijadikan leader
	var leaders []entity.User
	if err := db.Find(&leaders).Error; err != nil {
		return err
	}

	// Ambil beberapa congregation untuk dijadikan anggota
	var congregations []entity.Person
	if err := db.Find(&congregations).Error; err != nil {
		return err
	}

	lifeGroups := []entity.LifeGroup{
		{
			Name:         "Life Group Jakarta Pusat",
			Location:     "Jakarta Pusat",
			WhatsAppLink: "https://wa.me/group/abc123",
			LeaderID:     leaders[0].ID,
		},
		{
			Name:         "Life Group Jakarta Selatan",
			Location:     "Jakarta Selatan",
			WhatsAppLink: "https://wa.me/group/def456",
			LeaderID:     leaders[1].ID,
		},
		{
			Name:         "Life Group Jakarta Timur",
			Location:     "Jakarta Timur",
			WhatsAppLink: "https://wa.me/group/ghi789",
			LeaderID:     leaders[2].ID,
		},
	}

	for _, lifeGroup := range lifeGroups {
		if err := db.Create(&lifeGroup).Error; err != nil {
			return err
		}

		// Tambahkan beberapa anggota (user dan congregation)
		if len(leaders) > 3 {
			lifeGroup.Members = leaders[3:5] // Tambahkan 2 user sebagai anggota
		}
		if len(congregations) > 0 {
			lifeGroup.CongregationMembers = congregations[:2] // Tambahkan 2 congregation sebagai anggota
		}

		if err := db.Save(&lifeGroup).Error; err != nil {
			return err
		}
	}

	return nil
}
