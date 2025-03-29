package seeds

import (
	"github.com/google/uuid"
	"github.com/zemetia/en-indo-be/entity"
	"gorm.io/gorm"
)

func NotificationSeeder(db *gorm.DB) error {
	churchID := uuid.New()
	notifications := []entity.Notification{
		{
			Title:    "Selamat Datang di Every Nation Church",
			Message:  "Terima kasih telah bergabung dengan kami. Kami berharap Anda dapat terhubung dengan baik dalam komunitas kami.",
			Type:     "info",
			IsRead:   false,
			UserID:   uuid.New(), // User pertama
			ChurchID: &churchID,  // Every Nation Church Jakarta
		},
		{
			Title:    "Pengingat Ibadah Minggu",
			Message:  "Jangan lupa untuk hadir di ibadah minggu besok pukul 09:00 WIB.",
			Type:     "info",
			IsRead:   false,
			UserID:   uuid.New(),
			ChurchID: &churchID,
		},
		{
			Title:    "Pendaftaran Life Group",
			Message:  "Pendaftaran Life Group periode baru telah dibuka. Silakan daftar melalui aplikasi.",
			Type:     "success",
			IsRead:   false,
			UserID:   uuid.New(),
			ChurchID: &churchID,
		},
		{
			Title:    "Perubahan Jadwal Ibadah",
			Message:  "Ada perubahan jadwal ibadah minggu depan karena ada acara khusus.",
			Type:     "warning",
			IsRead:   false,
			UserID:   uuid.New(),
			ChurchID: &churchID,
		},
	}

	for _, notification := range notifications {
		if err := db.Create(&notification).Error; err != nil {
			return err
		}
	}

	return nil
}
