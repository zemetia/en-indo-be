package script

import (
	"log"
	"time"

	"github.com/zemetia/en-indo-be/repository"
	"gorm.io/gorm"
)

type CleanupKetersediaan struct {
	db *gorm.DB
}

func NewCleanupKetersediaan(db *gorm.DB) *CleanupKetersediaan {
	return &CleanupKetersediaan{
		db: db,
	}
}

func (s *CleanupKetersediaan) Run() error {
	repo := repository.NewKetersediaanRepository(s.db)

	// Calculate date 2 months ago
	retentionDate := time.Now().AddDate(0, -2, 0)

	log.Printf("Cleaning up availability records older than %s...", retentionDate.Format("2006-01-02"))

	if err := repo.CleanupOldKetersediaan(retentionDate); err != nil {
		return err
	}

	log.Println("Cleanup completed successfully.")
	return nil
}
